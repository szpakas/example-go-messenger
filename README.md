# Messenger example in GO [![Build Status](https://travis-ci.org/szpakas/example-go-messenger.svg?branch=master)](https://travis-ci.org/szpakas/example-go-messenger)

[![Apache 2.0 License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://github.com/szpakas/example-messenger-go/blob/master/LICENSE) [![codecov](https://codecov.io/gh/szpakas/example-go-messenger/branch/master/graph/badge.svg)](https://codecov.io/gh/szpakas/example-go-messenger) [![Go Report Card](https://goreportcard.com/badge/github.com/szpakas/example-go-messenger)](https://goreportcard.com/report/github.com/szpakas/example-go-messenger)

`example-messenger-go` is a proof of concept implementation of messenger app in GO.

## Running tests

    $ go test -v ./

test coverage:

    go test -covermode=count -coverprofile=count.out ./ && go tool cover -html=count.out

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

## Endpoints

Swagger 2.0 is used for REST endpoint documentation. It's available at /v1/swagger.json.

[swagger-go](https://github.com/go-swagger/go-swagger) project is used for documentation through annotations.

To regenerate:
```bash

swagger generate spec -o ./swagger.json
```

## Using docker

build locally

    $ docker build -t szpakas/example-go-messenger .

or pull from docker hub

    $ docker pull szpakas/example-go-messenger

```bash
docker run --rm -i \
    -p 8080:8080 \
    -e "APP_LOG_LEVEL=debug" \
    --name example-go-messenger szpakas/example-go-messenger
```
    
## License

Apache 2.0, see [LICENSE](./LICENSE).
