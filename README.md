# configen

[![License](https://img.shields.io/github/license/kukymbr/configen.svg)](https://github.com/kukymbr/configen/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/configen.svg)](https://github.com/kukymbr/configen/releases/latest)
[![GoDoc](https://godoc.org/github.com/kukymbr/configen?status.svg)](https://godoc.org/github.com/kukymbr/configen)
[![GoReport](https://goreportcard.com/badge/github.com/kukymbr/configen)](https://goreportcard.com/report/github.com/kukymbr/configen)

The `configen` is a config files generator for the [Golang](https://go.dev), converting this:

```go
package readme

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
```

into this:

```yaml
instance_id: 1
# App is an application common settings.
app:
  env: development
  namespace: unknown
  domain: ""
# Logger is a logging setup values.
logger:
  level: debug
```

and this:

```dotenv
INSTANCE_ID=1

APP_ENV=development
APP_NAMESPACE=unknown
APP_DOMAIN=

LOG_LEVEL=debug
```

Also, the `configen` is able to generate a new go struct with read-only getters of the config value.

<details>
  <summary>See the example of generated struct</summary>

```golang
package readme

type Config struct {
	instanceID int
	app        struct {
		env       string
		namespace string
		domain    string
	}
	logger struct {
		level string
	}

	origin any
}

func (c Config) InstanceID() int {
	return c.instanceID
}

// App is an application common settings.
func (c Config) App() struct {
	env       string
	namespace string
	domain    string
} {
	return c.app
}

// Logger is a logging setup values.
func (c Config) Logger() struct {
	level string
} {
	return c.logger
}

// NewConfig is a constructor converting config into the Config.
func NewConfig(dto config) Config {
	return Config{
		instanceID: dto.InstanceID,
		app: struct {
			env       string
			namespace string
			domain    string
		}{
			env:       dto.App.Env,
			namespace: dto.App.Namespace,
			domain:    dto.App.Domain,
		},
		logger: struct {
			level string
		}{
			level: dto.Logger.Level,
		},

		origin: dto,
	}
}

type ConfigAppProvider struct {
	env       string
	namespace string
	domain    string

	origin any
}

func (c ConfigAppProvider) Env() string {
	return c.env
}

func (c ConfigAppProvider) Namespace() string {
	return c.namespace
}

func (c ConfigAppProvider) Domain() string {
	return c.domain
}

type ConfigLoggerProvider struct {
	level string

	origin any
}

func (c ConfigLoggerProvider) Level() string {
	return c.level
}
```

</details>

## Why?

1. To simplify the creation of config files, obviously;
2. and to easily keep up to date an example config files when the config structure has changed;
3. to see all the available config values in one place;
4. to hide config values behind the read-only struct.

## Installation

The go 1.24 is a minimal requirement for the `configen`.

The `go tool` is a preferred way to install:

```shell
go get -tool github.com/kukymbr/configen/cmd/configen
```

## Usage

1. Create the structure you want to generate files from, e.g. `Config`;
2. set the structure tags (see the available in the table below),
3. add the `//go:generate` directive (run `go tool configen --help` for available flags):
   ```go
    package config  

   //go:generate go tool configen --struct=Config --yaml=true --env=true
   ```
4. run the `go generate` command:
   ```shell
   go generate ./...
   ```
5. enjoy the generated files.

### Supported struct tags

| Tag          | Value                                                                                |
|--------------|--------------------------------------------------------------------------------------|
| `yaml`       | key for the value in YAML file, or `-` to skip                                       |
| `env`        | key for the value in dotenv file, fields without this tag are not added to env file  |
| `envPrefix`  | prefix for sub-structs in dotenv file                                                |
| `default`    | default value to write to config files, prioritized for YAML                         |
| `envDefault` | default value to write to config files, prioritized for env                          |
| `example`    | default value to write to config files, general use (to use with swaggo for example) |

See the [example](example) directory for usage and generated code example.

### Command arguments to generate things

| Argument                   | Required | Value                                                                      |
|----------------------------|----------|----------------------------------------------------------------------------|
| `--struct=<StructName>`    | ✅        | Name of the struct to generate config from                                 |
| `--source=<dir>`           |          | Directory of the source go files (default `.`)                             |
| `--yaml=<filepath/true>`   |          | Path to YAML config file, set `true` to enable with default path           |
| `--yaml-tag=<tag>`         |          | Tag name for a YAML field names (default `yaml`)                           |
| `--env=<filepath/true>`    |          | Path to dotenv config file, set `true` to enable with default path         |
| `--env-tag=<tag>`          |          | Tag name for a dotenv variables names (default `env`)                      |
| `--env-prefix-tag=<tag>`   |          | Tag name for a dotenv subs-struct variables prefixes (default `envPrefix`) |
| `--go=<filepath/true>`     |          | Path to Golang config getter file, set `true` to enable with default path  |
| `--go-pkg=<package>`       |          | Target package name (default is equal to source package)                   |
| `--go-struct=<StructName>` |          | Target struct name (default is exported variant of incoming struct name)   |
| `--value-tag=<tag>`        |          | Custom tag name for default values                                         |

<details>
<summary>
    The <code>configen --help</code> output
</summary>

```text
Usage:
  configen [flags]

Flags:
      --env string              Path to dotenv config file, set 'true' to enable with default path
      --env-prefix-tag string   Tag name for a dotenv variable prefixes (default "envPrefix")
      --env-tag string          Tag name for a dotenv variables names (default "env")
      --go string               Path to Golang config getter file, set 'true' to enable with default path
      --go-pkg string           Target package name
      --go-struct string        Target struct name (default is exported variant of incoming struct name)
  -h, --help                    help for configen
  -s, --silent                  Silent mode
      --source string           Directory of the source go files (default ".")
      --struct string           Name of the struct to generate config from
      --value-tag string        Tag name for a default value, prepends the default lookup if given
  -v, --version                 version for configen
      --yaml string             Path to YAML config file, set 'true' to enable with default path
      --yaml-tag string         Tag name for a YAML field names (default "yaml")
```

</details>

### Generating multiple versions from one struct

Sometimes you need to generate multiple versions of the config file, for example, for different environments.
To do this, you can use the `--value-tag` flag to specify a custom tag name for default values.

For example, to generate two YAMLs for production and local environments:

```go
package config

//go:generate go run ../../cmd/configen/main.go --source=. --struct=MultiConfig --yaml=production.yaml
//go:generate go run ../../cmd/configen/main.go --source=. --struct=MultiConfig --yaml=local.yaml --yaml-tag=local --value-tag=localDefault

type MultiConfig struct {
	Env                 string `yaml:"env" default:"production" local:"env" localDefault:"development"`
	ProductionOnlyValue string `yaml:"production_only_value" local:"-" default:"very productional value"`
}
```

This will give you two YAML files with different keys presence and values:

```yaml
# production.yaml

env: production
production_only_value: very productional value
```

```yaml
# local.yaml

env: development
```

## Contributing

Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) doc.

## TODO

- Add an empty value check in go generator (for basic types), for example:
  ```golang
  func NewAppConfig(dto appConfig) AppConfig {
    if dto.env == "" {
        dto.env = "development"
    }
   
    return AppConfig{
        env: dto.Env,
    }
  }
  ```
- Make an option to include/exclude origin in generated structs.

## License

[MIT](LICENSE).