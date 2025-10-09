package gogetter

type StructInfo struct {
	Name             string
	SourceStructName string
	Doc              string
	Fields           []FieldInfo
}

type FieldInfo struct {
	Name       string
	ExportName string
	TypeName   string
	Comment    string
	IsStruct   bool
}
