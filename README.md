# eebus-go

The goal is to provide a basic EEBUS implementation

## Introduction

This repository contains:

- adoptions of the SPINE and SHIP EEBUS model definitions, there are likely issues and some models are not 100% correct
- (De-)serialization for EEBUS specific JSON format requirements
- ... work in progress

## Usage

```sh
Usage: go run cmd/prototype/main.go <command>
Commands:
  browse
  connect <host> <port>
```

- `browse` will search mdns for EEBUS services and try to connect to them
- `connect` will connect to a specific EEBUS service by provide a host/ip adress and port address
