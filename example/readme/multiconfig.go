package readme

//go:generate go run ../../cmd/configen/main.go --source=. --struct=MultiConfig --yaml=production.yaml
//go:generate go run ../../cmd/configen/main.go --source=. --struct=MultiConfig --yaml=local.yaml --yaml-tag=local --value-tag=localDefault

type MultiConfig struct {
	Env                 string `yaml:"env" default:"production" local:"env" localDefault:"development"`
	ProductionOnlyValue string `yaml:"production_only_value" local:"-" default:"very productional value"`
}
