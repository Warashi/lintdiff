package lintdiff

import (
	"io"
	"sort"

	"github.com/haya14busa/errorformat"
)

// LintResults fold results of lint
type LintResults []*errorformat.Entry

func (r LintResults) Len() int           { return len(r) }
func (r LintResults) Less(i, j int) bool { return less(r[i], r[j]) }
func (r LintResults) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func equal(i, j *errorformat.Entry) bool {
	if i.Lnum != j.Lnum {
		return false
	}
	if i.Col != j.Col {
		return false
	}
	if i.Type != j.Type {
		return false
	}
	if i.Nr != j.Nr {
		return false
	}
	return i.Text == j.Text
}

func less(i, j *errorformat.Entry) bool {
	if i.Lnum != j.Lnum {
		return i.Lnum < j.Lnum
	}
	if i.Col != j.Col {
		return i.Col < j.Col
	}
	if i.Type != j.Type {
		return i.Type < j.Type
	}
	if i.Nr != j.Nr {
		return i.Nr < j.Nr
	}
	return i.Text < j.Text
}

func scan(r io.Reader, format []string) (LintResults, error) {
	efm, err := errorformat.NewErrorformat(format)
	if err != nil {
		return nil, err
	}
	s := efm.NewScanner(r)

	var lr LintResults
	for s.Scan() {
		lr = append(lr, s.Entry())
	}
	sort.Sort(lr)
	return lr, nil
}
