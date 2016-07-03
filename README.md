# Messenger example in GO

[![Apache 2.0 License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://github.com/szpakas/example-messenger-go/blob/master/LICENSE)

`example-messenger-go` is a proof of concept implementation of messenger app in GO.

## Running tests

    $ go test -v ./

## Building
[govend](https://github.com/govend/govend) is used for vendoring.

```bash
govend -v 

go build
```

## Configuration

[12-factor](http://12factor.net/config) principles are followed for app configuration. All config values are stored in environmental variables prefixed with "APP_".
[envconfig](https://github.com/vrischmann/envconfig) package is used for configuration processing.

Configuration options
```go
// HTTPHost is address on which HTTP server endpoint is listening.
HTTPHost string `envconfig:"default=0.0.0.0"`

// HTTPPort is a port number on which HTTP server endpoint is listening.
HTTPPort int `envconfig:"default=8080"`

// LogLevel is a minimal log severity required for the message to be logged.
// Valid levels: [debug, info, warn, error, fatal, panic].
LogLevel string `envconfig:"default=info"`
```

## Using docker

build locally

    $ docker build -t fakepushprovider .

## License

Apache 2.0, see [LICENSE](./LICENSE).
