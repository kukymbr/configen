package testdata

import "errors"

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
	Env       string `env:"ENV" default:"development" json:"env" yaml:"env"`
	Namespace string `env:"NAMESPACE" default:"unknown" json:"namespace" yaml:"namespace"`
	Domain    string `env:"DOMAIN" example:"example.com" json:"domain" yaml:"domain"`
}

type LoggerConfig struct {
	Level LogLevel `env:"LEVEL" envDefault:"debug" json:"level" yaml:"level"`
}

type APIConfig struct {
	Host string `env:"HOST" envDefault:"0.0.0.0" json:"host" yaml:"host"`
	Port int    `env:"PORT" envDefault:"8080" json:"port" yaml:"port"`
}

type LogLevel int

func (l *LogLevel) MarshalText() ([]byte, error) {
	switch *l {
	case 0:
		return []byte("debug"), nil
	case 1:
		return []byte("info"), nil
	case 2:
		return []byte("error"), nil
	}

	return []byte(""), errors.New("invalid log level")
}
