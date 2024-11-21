package rename

import (
	"errors"
	"io"
	"sort"
)

const (
	keyShift      = 32
	maxCountValue = (1 << keyShift) - 1
)

var errIndexFull = errors.New("index is full")

// similarityIndex is an index structure of lines/blocks in one file.
// This structure can be used to compute an approximation of the similarity
// between two files.
// To save space in memory, this index uses a space efficient encoding which
// will not exceed 1MiB per instance. The index starts out at a smaller size
// (closer to 2KiB), but may grow as more distinct blocks within the scanned
// file are discovered.
// see: https://github.com/eclipse/jgit/blob/master/org.eclipse.jgit/src/org/eclipse/jgit/diff/SimilarityIndex.java
type similarityIndex struct {
	hashed uint64
	// number of non-zero entries in hashes
	numHashes int
	growAt    int
	hashes    []keyCountPair
	hashBits  int
}

func fileSimilarityIndex(f File) (*similarityIndex, error) {
	idx := newSimilarityIndex()
	if err := idx.hash(f); err != nil {
		return nil, err
	}

	sort.Stable(keyCountPairs(idx.hashes))

	return idx, nil
}

func newSimilarityIndex() *similarityIndex {
	return &similarityIndex{
		hashBits: 8,
		hashes:   make([]keyCountPair, 1<<8),
		growAt:   shouldGrowAt(8),
	}
}

func (i *similarityIndex) hash(f File) error {
	r, err := f.Reader()
	if err != nil {
		return err
	}

	defer checkClose(r, &err)

	size, err := f.Size()
	if err != nil {
		return err
	}
	return i.hashContent(r, size)
}

func (i *similarityIndex) hashContent(r io.Reader, size int64) error {
	var buf = make([]byte, 4096)
	var ptr, cnt int
	remaining := size

	for 0 < remaining {
		hash := 5381
		var blockHashedCnt uint64

		// Hash one line or block, whatever happens first
		n := int64(0)
		for {
			if ptr == cnt {
				ptr = 0
				var err error
				cnt, err = io.ReadFull(r, buf)
				if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
					return err
				}

				if cnt == 0 {
					return io.EOF
				}
			}
			n++
			c := buf[ptr] & 0xff
			ptr++

			// Ignore CR in CRLF sequence
			if c == '\r' && ptr < cnt && buf[ptr] == '\n' {
				continue
			}
			blockHashedCnt++

			if c == '\n' {
				break
			}

			hash = (hash << 5) + hash + int(c)

			if n >= 64 || n >= remaining {
				break
			}
		}
		i.hashed += blockHashedCnt
		if err := i.add(hash, blockHashedCnt); err != nil {
			return err
		}
		remaining -= n
	}

	return nil
}

// score computes the similarity score between this index and another one.
// A region of a file is defined as a line in a text file or a fixed-size
// block in a binary file. To prepare an index, each region in the file is
// hashed; the values and counts of hashes are retained in a sorted table.
// Define the similarity fraction F as the count of matching regions between
// the two files divided between the maximum count of regions in either file.
// The similarity score is F multiplied by the maxScore constant, yielding a
// range [0, maxScore]. It is defined as maxScore for the degenerate case of
// two empty files.
// The similarity score is symmetrical; i.e. a.score(b) == b.score(a).
func (i *similarityIndex) score(other *similarityIndex, maxScore int) int {
	var maxHashed = i.hashed
	if maxHashed < other.hashed {
		maxHashed = other.hashed
	}
	if maxHashed == 0 {
		return maxScore
	}

	return int(i.common(other) * uint64(maxScore) / maxHashed)
}

func (i *similarityIndex) common(dst *similarityIndex) uint64 {
	srcIdx, dstIdx := 0, 0
	if i.numHashes == 0 || dst.numHashes == 0 {
		return 0
	}

	var common uint64
	srcKey, dstKey := i.hashes[srcIdx].key(), dst.hashes[dstIdx].key()

	for {
		if srcKey == dstKey {
			srcCnt, dstCnt := i.hashes[srcIdx].count(), dst.hashes[dstIdx].count()
			if srcCnt < dstCnt {
				common += srcCnt
			} else {
				common += dstCnt
			}

			srcIdx++
			if srcIdx == len(i.hashes) {
				break
			}
			srcKey = i.hashes[srcIdx].key()

			dstIdx++
			if dstIdx == len(dst.hashes) {
				break
			}
			dstKey = dst.hashes[dstIdx].key()
		} else if srcKey < dstKey {
			// Region of src that is not in dst
			srcIdx++
			if srcIdx == len(i.hashes) {
				break
			}
			srcKey = i.hashes[srcIdx].key()
		} else {
			// Region of dst that is not in src
			dstIdx++
			if dstIdx == len(dst.hashes) {
				break
			}
			dstKey = dst.hashes[dstIdx].key()
		}
	}

	return common
}

func (i *similarityIndex) add(key int, cnt uint64) error {
	key = int(uint32(key) * 0x9e370001 >> 1)

	j := i.slot(key)
	for {
		v := i.hashes[j]
		if v == 0 {
			// It's an empty slot, so we can store it here.
			if i.growAt <= i.numHashes {
				if err := i.grow(); err != nil {
					return err
				}
				j = i.slot(key)
				continue
			}

			var err error
			i.hashes[j], err = newKeyCountPair(key, cnt)
			if err != nil {
				return err
			}
			i.numHashes++
			return nil
		} else if v.key() == key {
			// It's the same key, so increment the counter.
			var err error
			i.hashes[j], err = newKeyCountPair(key, v.count()+cnt)
			return err
		} else if j+1 >= len(i.hashes) {
			j = 0
		} else {
			j++
		}
	}
}

type keyCountPair uint64

func newKeyCountPair(key int, cnt uint64) (keyCountPair, error) {
	if cnt > maxCountValue {
		return 0, errIndexFull
	}

	return keyCountPair((uint64(key) << keyShift) | cnt), nil
}

func (p keyCountPair) key() int {
	return int(p >> keyShift)
}

func (p keyCountPair) count() uint64 {
	return uint64(p) & maxCountValue
}

func (i *similarityIndex) slot(key int) int {
	// We use 31 - hashBits because the upper bit was already forced
	// to be 0 and we want the remaining high bits to be used as the
	// table slot.
	return int(uint32(key) >> uint(31-i.hashBits))
}

func shouldGrowAt(hashBits int) int {
	return (1 << uint(hashBits)) * (hashBits - 3) / hashBits
}

func (i *similarityIndex) grow() error {
	if i.hashBits == 30 {
		return errIndexFull
	}

	old := i.hashes

	i.hashBits++
	i.growAt = shouldGrowAt(i.hashBits)

	// TODO: find a way to check if it will OOM and return errIndexFull instead.
	i.hashes = make([]keyCountPair, 1<<uint(i.hashBits))

	for _, v := range old {
		if v != 0 {
			j := i.slot(v.key())
			for i.hashes[j] != 0 {
				j++
				if j >= len(i.hashes) {
					j = 0
				}
			}
			i.hashes[j] = v
		}
	}

	return nil
}

type keyCountPairs []keyCountPair

func (p keyCountPairs) Len() int           { return len(p) }
func (p keyCountPairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p keyCountPairs) Less(i, j int) bool { return p[i] < p[j] }
