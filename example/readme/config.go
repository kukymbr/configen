package readme

//go:generate go run ../../cmd/configen/main.go --source=. --struct=config --yaml=true --env=true --go=true

// config is an example application config.
type config struct {
	// App is an application common settings.
	App appConfig `envPrefix:"APP_" yaml:"app"`

	// Logger is a logging setup values.
	Logger loggerConfig `envPrefix:"LOG_" yaml:"logger"`

	// API is an API server configuration.
	API apiConfig `envPrefix:"API_" yaml:"api"`
}

type appConfig struct {
	Env       string `env:"ENV" envDefault:"development" yaml:"env"`
	Namespace string `env:"NAMESPACE" envDefault:"unknown" yaml:"namespace"`
	Domain    string `env:"DOMAIN" yaml:"domain"`
}

type loggerConfig struct {
	Level string `env:"LEVEL" envDefault:"debug" yaml:"level"`
}

type apiConfig struct {
	Host string `env:"HOST" envDefault:"0.0.0.0" yaml:"host"`
	Port int    `env:"PORT" envDefault:"8080" yaml:"port"`
}
