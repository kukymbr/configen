package gentype

const (
	TagYAML = "yaml"

	TagEnv        = "env"
	TagEnvPrefix  = "envPrefix"
	TagEnvDefault = "envDefault"

	TagDefault = "default"
	TagExample = "example"
)

var (
	valueTagsYAML = []string{TagDefault, TagExample, TagEnvDefault}
	valueTagsEnv  = []string{TagEnvDefault, TagDefault, TagExample}
)

func ValueTagsYAML(override ...string) []string {
	return appendSlicesFiltered(override, valueTagsYAML)
}

func ValueTagsEnv(override ...string) []string {
	return appendSlicesFiltered(override, valueTagsEnv)
}
