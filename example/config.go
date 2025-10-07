package example

// Added as an example usage.
// To regenerate example files in the configen repository, use `make generate_example`.
//go:generate go tool configen --struct=Config --yaml=true --env=example.env

// Config is an example application config.
type Config struct {
	// App is an application common settings.
	App AppConfig `envPrefix:"APP_" json:"app" yaml:"app"`

	// Logger is a logging setup values.
	Logger LoggerConfig `envPrefix:"LOG_" json:"log" yaml:"logger"`

	// API is an API server configuration.
	API APIConfig `envPrefix:"API_" json:"api" yaml:"api"`
}

type AppConfig struct {
	// Application environment mode: development|production
	Env string `env:"ENV" envDefault:"development" json:"env" yaml:"env"`

	// Environment namespace (e.g. "dev1")
	Namespace string `env:"NAMESPACE" envDefault:"unknown" json:"namespace" yaml:"namespace"`

	// Top-level domain for the cookies
	Domain string `json:"domain" yaml:"domain"`
}

type LoggerConfig struct {
	Level string `env:"LEVEL" envDefault:"debug" json:"level" yaml:"level"`
}

type APIConfig struct {
	Host   string `env:"HOST" envDefault:"0.0.0.0" json:"host" yaml:"host"`
	Port   int    `env:"PORT" envDefault:"8080" json:"port" yaml:"port"`
	Secret string `env:"SECRET,unset" envDefault:"secret" json:"secret" yaml:"secret"`
}
