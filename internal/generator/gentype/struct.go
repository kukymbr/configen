package gentype

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type SourceStruct struct {
	Package *packages.Package
	Struct  *types.Struct
	Named   *types.Named

	Name     string
	Doc      string
	Comments map[token.Pos]string
}
