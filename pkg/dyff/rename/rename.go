// Package rename contains modified code from go-git's rename detection logic.
//
// go-git is licensed under Apache License 2.0, and you may obtain a copy of their original code and license from:
// https://github.com/go-git/go-git
package rename

import (
	"errors"
	"io"
	"sort"
	"strings"
)

type DetectOptions struct {
	// RenameScore is the threshold to of similarity between files to consider
	// that a pair of delete and insert are a rename. The number must be
	// exactly between 0 and 100.
	RenameScore uint
	// RenameLimit is the maximum amount of files that can be compared when
	// detecting renames. The number of comparisons that have to be performed
	// is equal to the number of deleted files * the number of added files.
	// That means, that if 100 files were deleted and 50 files were added, 5000
	// file comparisons may be needed. So, if the rename limit is 50, the number
	// of both deleted and added needs to be equal or less than 50.
	// A value of 0 means no limit.
	RenameLimit uint
}

// DefaultDetectOptions are the default and recommended options.
var DefaultDetectOptions = &DetectOptions{
	RenameScore: 60,
	RenameLimit: 50,
}

type Changes interface {
	Deleted() []File
	Added() []File

	MarkAsRename(deleted, added File) error
}

type File interface {
	Name() string
	Reader() (io.ReadCloser, error)
	Size() (int64, error)
}

// DetectRenames detects the renames in the given changes on two trees with
// the given options. It will return the given changes grouping additions and
// deletions into modifications when possible.
// If options is nil, the default diff tree options will be used.
func DetectRenames(
	changes Changes,
	opts *DetectOptions,
) error {
	if opts == nil {
		opts = DefaultDetectOptions
	}

	detector := &renameDetector{
		c:           changes,
		deleted:     changes.Deleted(),
		added:       changes.Added(),
		renameScore: int(opts.RenameScore),
		renameLimit: int(opts.RenameLimit),
	}

	return detector.detect()
}

// renameDetector will detect and resolve renames in a set of changes.
// see: https://github.com/eclipse/jgit/blob/master/org.eclipse.jgit/src/org/eclipse/jgit/diff/RenameDetector.java
type renameDetector struct {
	c       Changes
	deleted []File
	added   []File

	renameScore int
	renameLimit int
}

func (d *renameDetector) detect() error {
	if len(d.added) > 0 && len(d.deleted) > 0 {
		return d.detectContentRenames()
	}
	return nil
}

// detectContentRenames detects renames based on the similarity of the content
// in the files by building a matrix of pairs between sources and destinations
// and matching by the highest score.
// see: https://github.com/eclipse/jgit/blob/master/org.eclipse.jgit/src/org/eclipse/jgit/diff/SimilarityRenameDetector.java
func (d *renameDetector) detectContentRenames() error {
	cnt := max(len(d.added), len(d.deleted))
	if d.renameLimit > 0 && cnt > d.renameLimit {
		return nil
	}

	srcs, dsts := d.deleted, d.added
	matrix, err := buildSimilarityMatrix(srcs, dsts, d.renameScore)
	if err != nil {
		return err
	}

	// Match rename pairs on a first-come-first-serve basis until
	// we have looked at everything that is above the minimum score.
	for i := len(matrix) - 1; i >= 0; i-- {
		pair := matrix[i]
		src := srcs[pair.deleted]
		dst := dsts[pair.added]

		if dst == nil || src == nil {
			// It was already matched before
			continue
		}

		if err = d.c.MarkAsRename(src, dst); err != nil {
			return err
		}

		// Mark as matched
		srcs[pair.deleted] = nil
		dsts[pair.added] = nil
	}
	return nil
}

func nameSimilarityScore(a, b string) int {
	aDirLen := strings.LastIndexByte(a, '/') + 1
	bDirLen := strings.LastIndexByte(b, '/') + 1

	dirMin := min(aDirLen, bDirLen)
	dirMax := max(aDirLen, bDirLen)

	var dirScoreLtr, dirScoreRtl int
	if dirMax == 0 {
		dirScoreLtr = 100
		dirScoreRtl = 100
	} else {
		var dirSim int

		for ; dirSim < dirMin; dirSim++ {
			if a[dirSim] != b[dirSim] {
				break
			}
		}

		dirScoreLtr = dirSim * 100 / dirMax

		if dirScoreLtr == 100 {
			dirScoreRtl = 100
		} else {
			for dirSim = 0; dirSim < dirMin; dirSim++ {
				if a[aDirLen-1-dirSim] != b[bDirLen-1-dirSim] {
					break
				}
			}
			dirScoreRtl = dirSim * 100 / dirMax
		}
	}

	fileMin := min(len(a)-aDirLen, len(b)-bDirLen)
	fileMax := max(len(a)-aDirLen, len(b)-bDirLen)

	fileSim := 0
	for ; fileSim < fileMin; fileSim++ {
		if a[len(a)-1-fileSim] != b[len(b)-1-fileSim] {
			break
		}
	}
	fileScore := fileSim * 100 / fileMax

	return (((dirScoreLtr + dirScoreRtl) * 25) + (fileScore * 50)) / 100
}

type similarityMatrix []similarityPair

func (m similarityMatrix) Len() int      { return len(m) }
func (m similarityMatrix) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m similarityMatrix) Less(i, j int) bool {
	if m[i].score == m[j].score {
		if m[i].added == m[j].added {
			return m[i].deleted < m[j].deleted
		}
		return m[i].added < m[j].added
	}
	return m[i].score < m[j].score
}

type similarityPair struct {
	// index of the added file
	added int
	// index of the deleted file
	deleted int
	// similarity score
	score int
}

const maxMatrixSize = 10000

func buildSimilarityMatrix(srcs, dsts []File, renameScore int) (similarityMatrix, error) {
	// Allocate for the worst-case scenario where every pair has a score
	// that we need to consider. We might not need that many.
	matrixSize := len(srcs) * len(dsts)
	if matrixSize > maxMatrixSize {
		matrixSize = maxMatrixSize
	}
	matrix := make(similarityMatrix, 0, matrixSize)
	srcSizes := make([]int64, len(srcs))
	dstSizes := make([]int64, len(dsts))
	dstTooLarge := make(map[int]bool)

	// Consider each pair of files, if the score is above the minimum
	// threshold we need to record that scoring in the matrix so we can
	// later find the best matches.
outerLoop:
	for srcIdx, src := range srcs {
		// Declare the from file and the similarity index here to be able to
		// reuse it inside the inner loop. The reason to not initialize them
		// here is so we can skip the initialization in case they happen to
		// not be needed later. They will be initialized inside the inner
		// loop if and only if they're needed and reused in subsequent passes.
		var s *similarityIndex
		var err error
		for dstIdx, dst := range dsts {
			if dstTooLarge[dstIdx] {
				continue
			}

			srcSize := srcSizes[srcIdx]
			if srcSize == 0 {
				srcSize, err = src.Size()
				if err != nil {
					return nil, err
				}
				srcSize += 1
				srcSizes[srcIdx] = srcSize
			}

			dstSize := dstSizes[dstIdx]
			if dstSize == 0 {
				dstSize, err = dst.Size()
				if err != nil {
					return nil, err
				}
				dstSize += 1
				dstSizes[dstIdx] = dstSize
			}

			minSize := min(srcSize, dstSize)
			maxSize := max(srcSize, dstSize)

			if int(minSize*100/maxSize) < renameScore {
				// File sizes are too different to be a match
				continue
			}

			if s == nil {
				s, err = fileSimilarityIndex(src)
				if err != nil {
					if errors.Is(err, errIndexFull) {
						continue outerLoop
					}
					return nil, err
				}
			}

			di, err := fileSimilarityIndex(dst)
			if err != nil {
				if errors.Is(err, errIndexFull) {
					dstTooLarge[dstIdx] = true
				}

				return nil, err
			}

			contentScore := s.score(di, 10000)
			// The name score returns a value between 0 and 100, so we need to
			// convert it to the same range as the content score.
			nameScore := nameSimilarityScore(src.Name(), dst.Name()) * 100
			score := (contentScore*99 + nameScore*1) / 10000

			if score < renameScore {
				continue
			}

			matrix = append(matrix, similarityPair{added: dstIdx, deleted: srcIdx, score: score})
		}
	}

	sort.Stable(matrix)

	return matrix, nil
}
