package gentype

type OutputOptions struct {
	// Enable is a flag to enable an output.
	Enable bool

	// Path is a target file path.
	Path string

	// Tag is a field names tag.
	Tag string

	TargetStructName string

	TargetPackageName string
}

type GeneratorFunc func(src *SourceStruct, out OutputOptions) error
