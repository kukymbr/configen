# configen

The `configen` is a code generator for the [Golang](https://go.dev).

Golang code generator

## Why?

## Installation

The go 1.24 is a minimal requirement for the `configen`. 

The `go tool` is a preferred way to install:

```shell
go get -tool github.com/kukymbr/configen/cmd/sqlamble@latest
```

## Usage

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

1. ...
2. Add the go file with a `//go:generate` directive:
   ```go
    package sql  

   //go:generate go tool configen --package=mypkg
   ```
3. Run the `go generate` command:
   ```shell
   go generate ./...
   ```
4. ...

See the [example](example) directory for usage and generated code example.

## Contributing

Please refer the [CONTRIBUTING.md](CONTRIBUTING.md) doc.

## License

[MIT](LICENSE).