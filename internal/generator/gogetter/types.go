package gogetter

type StructInfo struct {
	Name             string
	SourceStructName string
	Doc              string
	Fields           []FieldInfo
}

type FieldInfo struct {
	Name       string // private field name
	ExportName string // exported getter name
	TypeName   string
	Comment    string
	IsStruct   bool
}
