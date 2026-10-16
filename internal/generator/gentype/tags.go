package gentype

import "reflect"

const (
	TagYAML = "yaml"

	TagEnv        = "env"
	TagEnvPrefix  = "envPrefix"
	TagEnvDefault = "envDefault"

	TagDefault = "default"
	TagExample = "example"

	// TagEnvSeparator and TagEnvKeyValSeparator are the separators of the slice and map values,
	// as the env parsers (github.com/caarlos0/env and its forks) read them.
	TagEnvSeparator       = "envSeparator"
	TagEnvKeyValSeparator = "envKeyValSeparator"

	DefaultEnvSeparator       = ","
	DefaultEnvKeyValSeparator = ":"
)

var (
	valueTagsYAML = []string{TagDefault, TagExample, TagEnvDefault}
	valueTagsEnv  = []string{TagEnvDefault, TagDefault, TagExample}
)

// ValueSeparators returns the separators of slice items and map key-value pairs in a default value of the field.
func ValueSeparators(tag string) (separator string, keyValSeparator string) {
	st := reflect.StructTag(tag)

	separator = st.Get(TagEnvSeparator)
	if separator == "" {
		separator = DefaultEnvSeparator
	}

	keyValSeparator = st.Get(TagEnvKeyValSeparator)
	if keyValSeparator == "" {
		keyValSeparator = DefaultEnvKeyValSeparator
	}

	return separator, keyValSeparator
}

func ValueTagsYAML(override ...string) []string {
	return appendSlicesFiltered(override, valueTagsYAML)
}

func ValueTagsEnv(override ...string) []string {
	return appendSlicesFiltered(override, valueTagsEnv)
}
