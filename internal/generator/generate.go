package generator

import (
	"go/types"
	"reflect"
	"strings"
)

const (
	tagEnvPrefix  = "envPrefix"
	tagEnvDefault = "envDefault"

	tagDefault = "default"
	tagExample = "example"
)

var (
	valueTagsYAML = []string{tagDefault, tagExample, tagEnvDefault}
	valueTagsEnv  = []string{tagEnvDefault, tagDefault, tagExample}
)

func parseDefaultValue(tagValue string, tags ...string) string {
	if tagValue == "" {
		return ""
	}

	st := reflect.StructTag(tagValue)

	for _, tag := range tags {
		if v := st.Get(tag); v != "" {
			return v
		}
	}

	return ""
}

func defaultValueForType(t types.Type, value string) string {
	if value != "" {
		return value
	}

	switch tt := t.(type) {
	case *types.Basic:
		switch tt.Kind() {
		case types.Bool:
			return "false"
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64:
			return "0"
		case types.String:
			return ""
		default:
			return ""
		}
	}

	return ""
}

func isStructLike(t types.Type) bool {
	switch t.(type) {
	case *types.Struct, *types.Named:
		return true
	}

	return false
}

func underlyingStruct(t types.Type) (*types.Struct, bool) {
	switch tt := t.(type) {
	case *types.Pointer:
		return underlyingStruct(tt.Elem())
	case *types.Named:
		if st, ok := tt.Underlying().(*types.Struct); ok {
			return st, true
		}
	case *types.Struct:
		return tt, true
	}

	return nil, false
}

func parseNameTag(tagContent string, tagName string, fallback string) string {
	if tagContent == "" {
		return fallback
	}

	st := reflect.StructTag(tagContent)

	nameTag := st.Get(tagName)
	if nameTag == "" {
		return fallback
	}

	parts := strings.Split(nameTag, ",")
	if parts[0] == "-" {
		return ""
	}

	if parts[0] == "" {
		return fallback
	}

	return parts[0]
}
