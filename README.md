# 🏛️ Bibliotek 📚 [![CI](https://github.com/hemtjanst/bibliotek/workflows/CI/badge.svg?branch=master)](https://github.com/hemtjanst/bibliotek/actions?query=workflow%3ACI) ![GitHub release](https://img.shields.io/github/release/hemtjanst/bibliotek.svg) [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/lib.hemtjan.st/)

Bibliotek ('library' in Swedish) provides common primitives and utilities for
integrating with and extending the Hemtjänst platform.

## Usage

### `client` [![client godoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/lib.hemtjan.st/client)

The [client package](https://pkg.go.dev/lib.hemtjan.st/client)
should be used if what you want to do is publish/control
your own devices but do not care about other devices in the system.

For example, lets say you want to control devices connected to a Z-Wave
network, or republish external datasources like an HTTP API as a device (e.g
a temperature sensor).

### `server` [![server godoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/lib.hemtjan.st/server)

The [server package](https://pkg.go.dev/lib.hemtjan.st/server)
allows you to fetch all devices, their updates and send them commands. It does
not allow for the creation of devices.

It can be used to implement things like a HomeKit bridge or to watch for and
republish device data, like sensor readings, to another platform (e.g Prometheus).

Take a look at [`cmd/explorer`](../master/cmd/explorer/main.go) on how to
use it.

### `transport/mqtt` [![mqtt transport godoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/lib.hemtjan.st/transport/mqtt)

The [transport/mqtt package](https://pkg.go.dev/lib.hemtjan.st/transport/mqtt)
contains everything needed to transport device data over MQTT. It
also contains the [`Flags()`](https://pkg.go.dev/lib.hemtjan.st/transport/mqtt#Flags)
and [`MustFlags()`](https://pkg.go.dev/lib.hemtjan.st/transport/mqtt#MustFlags)
helpers that provide a set of flags to configure common MQTT options.

## Contributing

Contributions are very welcome. Do ensure you run the tests:

```sh
$ go test -race -v -coverprofile=profile.cov $(go list ./...)
...
```

Note that in CI the `_integration_test.go` will run too. You can
run them locally as well by setting `BIBLIOTEK_TEST_INTEGRATION=1`.

In order for these test to complete successfully you'll have to spin
up an MQTT broker and if the MQTT broker is not found on `localhost:1883`
specify a `host:port` string in the `MQTT_ADDRESS` environment variable.

```sh
$ env BIBLIOTEK_TEST_INTEGRATION=1 go test -race -v -coverprofile=profile.cov $(go list ./...)
...
```
