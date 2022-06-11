package reload

import (
	"io"
	"strings"

	"github.com/mitranim/gg"
)

type FlagStrMultiline string

func (f *FlagStrMultiline) Parse(src string) error {
	*f += FlagStrMultiline(withNewline(RepSingleMulti(src)))
	return nil
}

func (f FlagStrMultiline) Dump(out io.Writer) {
	if len(f) > 0 && out != nil {
		gg.Nop2(out.Write(gg.ToBytes(f)))
	}
}

type FlagExtensions []string

func (self *FlagExtensions) Parse(src string) (err error) {
	defer gg.Rec(&err)
	vals := commaSplit(src)
	gg.Each(vals, validateExtension)
	gg.AppendVals(self, vals...)
	return
}

func (self FlagExtensions) Allow(path string) bool {
	return gg.IsEmpty(self) || gg.Has(self, cleanExtension(path))
}

type FlagWatch []string

func (self *FlagWatch) Parse(src string) error {
	gg.AppendVals(self, commaSplit(src)...)
	return nil
}

type FlagIgnoredPaths []string

func (f *FlagIgnoredPaths) Parse(src string) error {
	values := FlagIgnoredPaths(commaSplit(src))
	values.Norm()
	//gg.AppendVals(f, values...)
	return nil
}

func (f FlagIgnoredPaths) Norm() {
	gg.MapMut(f, toAbsDirPath)
}

func (f FlagIgnoredPaths) Allow(path string) bool {
	return !f.Ignore(path)
}

// Assumes that the input is an absolute path.
func (f FlagIgnoredPaths) Ignore(path string) bool {
	return gg.Some(f, func(val string) bool {
		return strings.HasPrefix(path, val)
	})
}
