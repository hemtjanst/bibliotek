# 🏛️ Bibliotek 📚 [![Build Status](https://travis-ci.org/hemtjanst/bibliotek.svg?branch=master)](https://travis-ci.org/hemtjanst/bibliotek) ![GitHub release](https://img.shields.io/github/release/hemtjanst/bibliotek.svg) [![hemtjanst godoc](https://godoc.org/github.com/hemtjanst?status.svg)](https://godoc.org/github.com/hemtjanst/bibliotek)

Bibliotek ('library' in Swedish) provides common primitives and utilities for
integrating with and extending the Hemtjänst platform.

## Usage

### `client` [![client godoc](https://godoc.org/github.com/hemtjanst?status.svg)](https://godoc.org/github.com/hemtjanst/bibliotek/client)

The [client package](https://godoc.org/github.com/hemtjanst/bibliotek/client)
should be used if what you want to do is publish/control
your own devices but do not care about other devices in the system.

For example, lets say you want to control devices connected to a Z-Wave
network, or republish external datasources like an HTTP API as a device (e.g
a temperature sensor).

### `server` [![server godoc](https://godoc.org/github.com/hemtjanst?status.svg)](https://godoc.org/github.com/hemtjanst/bibliotek/server)

The [server package](https://godoc.org/github.com/hemtjanst/bibliotek/server)
allows you to fetch all devices, their updates and send them commands. It does
not allow for the creation of devices.

It can be used to implement things like a HomeKit bridge or to watch for and
republish device data, like sensor readings, to another platform (e.g Prometheus).

Take a look at [`cmd/explorer`](../master/cmd/explorer/main.go) on how to
use it.

### `transport/mqtt` [![client godoc](https://godoc.org/github.com/hemtjanst?status.svg)](https://godoc.org/github.com/hemtjanst/bibliotek/transport/mqtt)

The [transport/mqtt package](https://godoc.org/github.com/hemtjanst/bibliotek/transport/mqtt)
contains everything needed to transport device data over MQTT. It
also contains the [`Flags()`](https://godoc.org/github.com/hemtjanst/bibliotek/transport/mqtt#Flags)
and [`MustFlags()`](https://godoc.org/github.com/hemtjanst/bibliotek/transport/mqtt#MustFlags)
helpers that provide a set of flags to configure common MQTT options.

## Contributing

Contributions are very welcome. Do ensure you run the tests:

```sh
$ go test -race -v -coverprofile=profile.cov $(go list ./...)
...
```

Note that in CI the `_integration_test.go` will run too. You can
run them locally as well by setting `BIBLIOTEK_TEST_INTEGRATION=1`.
These tests require Docker to run as they spin up an actual MQTT
broker.

```sh
$ env BIBLIOTEK_TEST_INTEGRATION=1 go test -race -v -coverprofile=profile.cov $(go list ./...)
...
```
