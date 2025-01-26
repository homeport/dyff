package dyff

import (
	"bytes"
	"errors"
	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff/rename"
	yamlv3 "gopkg.in/yaml.v3"
	"io"
)

func mapSlice[E any, S ~[]E, T any](slice S, fn func(e E) T) []T {
	ret := make([]T, len(slice))
	for i, e := range slice {
		ret[i] = fn(e)
	}
	return ret
}

func reject[E comparable, S ~[]E](slice S, elt E) (ret S, ok bool) {
	ret = make(S, 0, len(slice))
	for _, e := range slice {
		if elt == e {
			ok = true
		} else {
			ret = append(ret, e)
		}
	}
	return
}

type modifiedPair struct {
	from *renameCandidate
	to   *renameCandidate
}

type documentChanges struct {
	deleted []*renameCandidate
	added   []*renameCandidate

	modifiedPairs []modifiedPair
}

func newDocumentChanges(deleted []*renameCandidate, added []*renameCandidate) *documentChanges {
	return &documentChanges{
		deleted: deleted,
		added:   added,
	}
}

func (d *documentChanges) Deleted() []rename.File {
	return mapSlice(d.deleted, func(r *renameCandidate) rename.File { return r })
}

func (d *documentChanges) Added() []rename.File {
	return mapSlice(d.added, func(r *renameCandidate) rename.File { return r })
}

func (d *documentChanges) MarkAsRename(deleted, added rename.File) error {
	var ok bool
	d.deleted, ok = reject(d.deleted, deleted.(*renameCandidate))
	if !ok {
		return errors.New("deleted element not found")
	}
	d.added, ok = reject(d.added, added.(*renameCandidate))
	if !ok {
		return errors.New("added element not found")
	}
	d.modifiedPairs = append(d.modifiedPairs, modifiedPair{
		from: deleted.(*renameCandidate),
		to:   added.(*renameCandidate),
	})
	return nil
}

type renameCandidate struct {
	path *ytbx.Path
	doc  *yamlv3.Node

	content []byte
}

func (r *renameCandidate) Name() string {
	name, _ := k8sItem.Name(r.doc)
	return name
}

func (r *renameCandidate) Reader() (io.ReadCloser, error) {
	if r.content == nil {
		if err := r.marshal(); err != nil {
			return nil, err
		}
	}
	return io.NopCloser(bytes.NewReader(r.content)), nil
}

func (r *renameCandidate) Size() (int64, error) {
	if r.content == nil {
		if err := r.marshal(); err != nil {
			return 0, err
		}
	}
	return int64(len(r.content)), nil
}

func (r *renameCandidate) marshal() error {
	var err error
	r.content, err = yamlv3.Marshal(r.doc)
	return err
}
