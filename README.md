# configen

[![License](https://img.shields.io/github/license/kukymbr/configen.svg)](https://github.com/kukymbr/configen/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/configen.svg)](https://github.com/kukymbr/configen/releases/latest)
[![GoDoc](https://godoc.org/github.com/kukymbr/configen?status.svg)](https://godoc.org/github.com/kukymbr/configen)
[![GoReport](https://goreportcard.com/badge/github.com/kukymbr/configen)](https://goreportcard.com/report/github.com/kukymbr/configen)

The `configen` is a config files generator for the [Golang](https://go.dev), converting this:

```go
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
    Env       string `env:"ENV" envDefault:"development" json:"env" yaml:"env"`
    Namespace string `env:"NAMESPACE" envDefault:"unknown" json:"namespace" yaml:"namespace"`
    Domain    string `env:"DOMAIN" json:"domain" yaml:"domain"`
}

type LoggerConfig struct {
    Level string `env:"LEVEL" envDefault:"debug" json:"level" yaml:"level"`
}

type APIConfig struct {
    Host string `env:"HOST" envDefault:"0.0.0.0" json:"host" yaml:"host"`
    Port int    `env:"PORT" envDefault:"8080" json:"port" yaml:"port"`
}
```

into this:

```yaml
# Config is an example application config.

# App is an application common settings.
app:
  env: development
  namespace: unknown
  domain: ""
# Logger is a logging setup values.
logger:
  level: debug
# API is an API server configuration.
api:
  host: 0.0.0.0
  port: 8080
```

and this:

```dotenv
# Config is an example application config.

APP_ENV=development
APP_NAMESPACE=unknown
APP_DOMAIN=
LOG_LEVEL=debug
API_HOST=0.0.0.0
API_PORT=8080
```

## Why?

1. To simplify the creation of config files, obviously;
2. and to easily keep up-to-date an example config files when config structure has changed.

## Installation

The go 1.24 is a minimal requirement for the `configen`.

The `go tool` is a preferred way to install:

```shell
go get -tool github.com/kukymbr/configen/cmd/configen
```

## Usage

1. Create the structure you want to generate files from, e.g. `Config`;
2. set the tags, see the available tags in the table below;
3. add the `//go:generate` directive (run `go tool configen --help` for available flags):
   ```go
    package config  

   //go:generate go tool configen --struct=Config --yaml=true --env=true
   ```
4. run the `go generate` command:
   ```shell
   go generate ./...
   ```
5. see the generated files.

| Tag          | Value                                                                               |
|--------------|-------------------------------------------------------------------------------------|
| `yaml`       | key for the value in YAML file, or `-` to skip                                      |
| `env`        | key for the value in dotenv file, fields without this tag are not added to env file |
| `default`    | default value to write to config files, prioritized for YAML                        |
| `envDefault` | default value to write to config files, prioritized for env                         |
| `example`    | default value to write to config files, to use with swaggo for example              |

See the [example](example) directory for usage and generated code example.

The `configen --help` output:

```text
Usage:
  configen [flags]

Flags:
      --fmt string            Formatter used to format generated go files (gofmt|noop) (default "gofmt")
  -h, --help                  help for configen
      --package string        Target package name of the generated code 
  -s, --silent                Silent mode
      --target string         Directory for the generated Go files (default ".")
  -v, --version               version for configen
```

## Contributing

Please refer the [CONTRIBUTING.md](CONTRIBUTING.md) doc.

## License

[MIT](LICENSE).