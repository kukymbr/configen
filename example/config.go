package example

import (
	"errors"
	"net/http"
	"time"
)

// Added as an example usage.
// To regenerate example files in the configen repository, use `make generate_example`.
//go:generate go tool configen --struct=config --yaml=true --env=example.env --go=config.gen.go

// Config godoc
//
// Main application config.
//
//nolint:unused // TODO: skip nolint comments (or comments without space after //)
type config struct {
	// App is an application common settings.
	App appConfig `envPrefix:"APP_" json:"app" yaml:"app"`

	// Logger is a logging setup values.
	Logger loggerConfig `envPrefix:"LOG_" json:"log" yaml:"logger"`

	// API is an API server configuration.
	API apiConfig `envPrefix:"API_" json:"api" yaml:"api"`
}

//nolint:unused
type appConfig struct {
	// Application environment mode: development|production
	Env string `env:"ENV" envDefault:"development" json:"env" yaml:"env"`

	// Environment namespace (e.g. "dev1")
	Namespace string `env:"NAMESPACE" envDefault:"unknown" json:"namespace" yaml:"namespace"`

	// Top-level domain for the cookies
	Domain string `json:"domain" yaml:"domain"`
}

//nolint:unused
type loggerConfig struct {
	Level LogLevel `env:"LEVEL" envDefault:"debug" json:"level" yaml:"level"`
}

//nolint:unused
type apiConfig struct {
	Host       string        `env:"HOST" envDefault:"0.0.0.0" json:"host" yaml:"host"`
	Port       int           `env:"PORT" envDefault:"8080" json:"port" yaml:"port"`
	Secret     string        `env:"SECRET,unset" envDefault:"secret" json:"secret" yaml:"secret"`
	ReqTTL     time.Duration `env:"REQ_TTL" envDefault:"1h" json:"req_ttl" yaml:"req_ttl"`
	RespTTL    time.Duration `env:"RESP_TTL" envDefault:"1h" json:"resp_ttl" yaml:"resp_ttl"`
	DefaultReq *http.Request `yaml:"-"`
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

func (l *LogLevel) UnmarshalText(text []byte) error {
	switch string(text) {
	case "debug":
		*l = 0
	case "info":
		*l = 1
	case "error":
		*l = 2
	default:
		return errors.New("invalid log level")
	}

	return nil
}
