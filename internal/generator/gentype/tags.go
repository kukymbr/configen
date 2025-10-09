package gentype

const (
	TagEnvPrefix  = "envPrefix"
	TagEnvDefault = "envDefault"

	TagDefault = "default"
	TagExample = "example"
)

var (
	valueTagsYAML = []string{TagDefault, TagExample, TagEnvDefault}
	valueTagsEnv  = []string{TagEnvDefault, TagDefault, TagExample}
)

func ValueTagsYAML() []string {
	return valueTagsYAML
}

func ValueTagsEnv() []string {
	return valueTagsEnv
}
