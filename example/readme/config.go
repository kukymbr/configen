package readme

//go:generate go run ../../cmd/configen/main.go --source=. --struct=config --yaml=true --env=true --go=true

// config is an example application config.
type config struct {
	InstanceID int `env:"INSTANCE_ID" default:"1" yaml:"instance_id"`

	// App is an application common settings.
	App struct {
		Env       string `env:"ENV" default:"development" yaml:"env"`
		Namespace string `env:"NAMESPACE" default:"unknown" yaml:"namespace"`
		Domain    string `env:"DOMAIN" yaml:"domain"`
	} `envPrefix:"APP_" yaml:"app"`

	// Logger is a logging setup values.
	Logger struct {
		Level string `env:"LEVEL" default:"debug" yaml:"level"`
	} `envPrefix:"LOG_" yaml:"logger"`
}
