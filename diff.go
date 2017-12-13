package lintdiff

import (
	"io"

	"github.com/haya14busa/errorformat"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Operation int

const (
	DiffDelete Operation = -1
	DiffEqual  Operation = 0
	DiffInsert Operation = 1
)

type Diff struct {
	*errorformat.Entry
	Type    Operation
	OldLnum int
	NewLnum int
}

func lineDiff(src1, src2 string) (*diffmatchpatch.DiffMatchPatch, []diffmatchpatch.Diff) {
	dmp := diffmatchpatch.New()
	a, b, _ := dmp.DiffLinesToChars(src1, src2)
	diffs := dmp.DiffMain(a, b, false)
	return dmp, diffs
}

// Additional returns additional warning, error
func diff(dfm *diffmatchpatch.DiffMatchPatch, diffs []diffmatchpatch.Diff, from, to []*errorformat.Entry) []Diff {
	d := make([]Diff, 0, len(from)+len(to))
	var i int
	for j := range from {
		f := *from[j] // copy
		f.Lnum = dfm.DiffXIndex(diffs, f.Lnum)
		for i < len(to)-1 && less(to[i], &f) {
			d = append(d, Diff{
				Entry:   to[i],
				Type:    DiffInsert,
				NewLnum: to[i].Lnum,
			})
			i++
		}
		if equal(to[i], &f) {
			d = append(d, Diff{
				Entry:   to[i],
				Type:    DiffEqual,
				OldLnum: from[j].Lnum,
				NewLnum: to[i].Lnum,
			})
		} else {
			d = append(d, Diff{
				Entry:   from[j],
				Type:    DiffDelete,
				OldLnum: from[j].Lnum,
			})
		}
	}
	return d
}

func DiffMain(oldCode, newCode string, oldLint, newLint io.Reader, format []string) ([]Diff, error) {
	dfm, diffs := lineDiff(oldCode, newCode)
	oldEfms, err := scan(oldLint, format)
	if err != nil {
		return nil, err
	}
	newEfms, err := scan(newLint, format)
	if err != nil {
		return nil, err
	}
	return diff(dfm, diffs, oldEfms, newEfms), nil
}
