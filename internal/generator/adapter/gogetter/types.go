package gogetter

type StructInfo struct {
	Name             string
	SourceStructName string
	Doc              string
	IsAnonymous      bool
	Fields           []FieldInfo
}

type FieldInfo struct {
	Name       string
	ExportName string
	TypeName   string
	Comment    string
	// CloneFunc is a function to copy the value in the getter (maps.Clone or slices.Clone), empty if no copy is needed.
	CloneFunc  string
	IsStruct   bool
	StructInfo *StructInfo
}
