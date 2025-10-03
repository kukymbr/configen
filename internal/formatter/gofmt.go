package formatter

import (
	"context"
	"go/format"
)

func NewGoFmtFormatter() Formatter {
	return &goFmt{
		executable: "go",
	}
}

type goFmt struct {
	executable string
}

func (f *goFmt) Format(_ context.Context, content []byte) ([]byte, error) {
	return format.Source(content)
}
