package generator

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

var interfaces struct {
	textMarshaler *types.Interface
	stringer      *types.Interface
}

//nolint:gochecknoinits
func init() {
	if err := LoadInterfaces(); err != nil {
		panic(err)
	}
}

func LoadInterfaces() error {
	var err error

	interfaces.textMarshaler, err = lookupInterface("encoding", "TextMarshaler")
	if err != nil {
		return err
	}

	interfaces.stringer, err = lookupInterface("fmt", "Stringer")
	if err != nil {
		return err
	}

	return nil
}

func lookupInterface(pkg string, name string) (*types.Interface, error) {
	conf := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(conf, pkg)
	if err != nil || len(pkgs) == 0 {
		return nil, fmt.Errorf("fetch package info for %s: %w", pkg, err)
	}

	scope := pkgs[0].Types.Scope()
	if tm := scope.Lookup(name); tm != nil {
		return tm.Type().Underlying().(*types.Interface), nil
	}

	return nil, fmt.Errorf("failed to find interface %s in %s", name, pkg)
}

func isTextMarshaler(t types.Type) bool {
	return implements(t, interfaces.textMarshaler)
}

func isStringer(t types.Type) bool {
	return implements(t, interfaces.stringer)
}

func implements(t types.Type, intf *types.Interface) bool {
	if types.Implements(t, intf) {
		return true
	}

	if types.Implements(types.NewPointer(t), intf) {
		return true
	}

	return false
}
