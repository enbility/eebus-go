# eebus-go

[![Build Status](https://github.com/enbility/eebus-go/actions/workflows/default.yml/badge.svg?branch=dev)](https://github.com/enbility/eebus-go/actions/workflows/default.yml/badge.svg?branch=dev)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://godoc.org/github.com/enbility/eebus-go)
[![Coverage Status](https://coveralls.io/repos/github/enbility/eebus-go/badge.svg?branch=dev)](https://coveralls.io/github/enbility/eebus-go?branch=dev)
[![Go report](https://goreportcard.com/badge/github.com/enbility/eebus-go)](https://goreportcard.com/report/github.com/enbility/eebus-go)
[![CodeFactor](https://www.codefactor.io/repository/github/enbility/eebus-go/badge)](https://www.codefactor.io/repository/github/enbility/eebus-go)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/enbility/eebus-go)

This library provides a foundation for implementing [EEBUS](https://eebus.org) use cases in [go](https://golang.org). It uses the SHIP implementation [ship-go](https://github.com/enbility/ship-go) and the SPINE implementation [spine-go](https://github.com/enbility/spine-go). Both repositories started as part of this repository, before they were moved into their own separate repositories and go packages.

Basic understanding of the EEBUS concepts SHIP and SPINE to use this library is required. Please check the corresponding specifications on the [EEBUS downloads website](https://www.eebus.org/media-downloads/).

## Introduction

The supported functionality contains:

- Support for SHIP 1.0.1 via [ship-go](https://github.com/enbility/ship-go)
- Support for SPINE 1.3.0 via [spine-go](https://github.com/enbility/spine-go)
- Certificate handling
- mDNS Support, incl. avahi support
- Connection (websocket) handling, including reconnection and double connections
- Support for handling pairing of devices

## Packages

- `api`: global API interface definitions and eebus service configuration
- `features/client`: provides feature helpers with the local SPINE feature having the client role and the remote SPINE feature being the server for easy access to commonly used functions
- `features/server`: provides feature helpers with the local SPINE feature having the server role for easy access to commonly used functions
- `service`: central package which provides access to SHIP and SPINE. Use this to create the EEBUS service, its configuration and connect to remote EEBUS services
- `usecases`: containing actor and use case based implementations with use case scenario based APIs and events

## Examples

The examples folder contains a few demo applications using this stack. These do **not** provide complete implementations of any use case, but are intended as usage guidelines for the eebus-go stack in general and to have a quick demo.

Therefore, please do not expect any of these examples to provide any meaningful functionality on their own, but instead view them as rough guidelines on the functionality you could implement using the eebus-go stack.

### Controlbox

This includes example code for sending an LPC limit 5 seconds after connecting to a compatible device that can receive LPC limits.

#### First Run

```sh
go run examples/controlbox/main.go -port 4713
```

`4713` is the example server port that this process should listen on

The certificate and key and the local SKI will be generated and printed. You should then save the certificate and the key to a file.

#### General Usage

```sh
Usage: go run examples/controlbox/main.go -port <serverport> [-certpath <certfile>] [-keypath <keyfile>] [-remoteski <remoteski>] [-target <pairing target>] [-target <pairing target>...]
```

- `remoteski` is the SKI of the remote device or service you want to connect to
- `certfile` is a local file containing the generated certificate in the first usage run
- `keyfile` is a local file containing the generated key in the first usage run
- `pairing target` can be a SHIP QR Code or in this format: `SKI=...,Fingerprint=...,ShipID=...,Secret=hex`

### HEMS

This includes example code for accepting LPC and LPP limits from a control box, receiving and printing data to the console from battery (VABD) and pv inverters (VAPD) and grid connection point data (MGCP).

#### First Run

```sh
go run examples/hems/main.go -port 4714
```

`4714` is the example server port that this process should listen on

The certificate and key and the local SKI will be generated and printed. You should then save the certificate and the key to a file.

#### General Usage

```sh
Usage: go run examples/hems/main.go -port <serverport> [-certpath <certfile>] [-keypath <keyfile>] [-remoteski <remoteski>] [-secret <secret>]
```

- `remoteski` is the SKI of the remote device or service you want to connect to
- `certfile` is a local file containing the generated certificate in the first usage run
- `keyfile` is a local file containing the generated key in the first usage run
- `secret` is a hexadecimal secret key as specified by SHIP Pairing Service Specificiation

### EVSE

This includes example code for accepting LPC from a control box.

#### First Run

```sh
go run examples/hems/main.go -port 4715
```

`4715` is the example server port that this process should listen on

The certificate and key and the local SKI will be generated and printed. You should then save the certificate and the key to a file.

#### General Usage

```sh
Usage: go run examples/evse/main.go -port <serverport> [-certpath <certfile>] [-keypath <keyfile>] [-remoteski <remoteski>] [-secret <secret>]
```

- `remoteski` is the SKI of the remote device or service you want to connect to
- `certfile` is a local file containing the generated certificate in the first usage run
- `keyfile` is a local file containing the generated key in the first usage run
- `secret` is a hexadecimal secret key as specified by SHIP Pairing Service Specificiation

### Explanation

The remoteski is from the eebus service to connect to.
If no certfile or keyfile are provided, they are generated and printed in the console so they can be saved in a file and later used again. The local SKI is also printed.

## SHIP implementation notes

- Double connection handling is not implemented according to SHIP 12.2.2. Instead the connection initiated by the higher SKI will be kept. Much simpler and always works
- PIN Verification is _NOT_ supported other than SHIP 13.4.4.3.5.1 _"none"_ PIN state is supported!
- Access Methods SHIP 13.4.6 only supports the most basic scenario and only works after PIN verification state is completed.
- Supported registration mechanisms (SHIP 5):
  - auto accept (without any interaction mechanism!)
  - user verification

This approach has been tested with:

- Elli Charger Connect
- Porsche Mobile Charger Connect
- SMA Home Energy Manager 2.0

## Interfaces

### Verbose logging

Use `SetLogger` on `Service` to set the logger which needs to conform to the `logging.Logging` interface of [ship-go](https://github.com/enbility/ship-go).

Example:

```go
configuration = service.NewConfiguration(...)
h.myService = service.NewEEBUSService(configuration, h)
h.myService.SetLogging(h)
```
