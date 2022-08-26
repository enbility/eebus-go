# eebus-go

![Build Status](https://github.com/DerAndereAndi/eebus-go/actions/workflows/default.yml/badge.svg?branch=dev)

The goal is to provide a basic EEBUS implementation

## Introduction

This repository contains:

- adoptions of the SPINE and SHIP EEBUS model definitions, there are likely issues and some models are not 100% correct
- (De-)serialization for EEBUS specific JSON format requirements
- Certificate support incl. creating a compatible cert and key
- mDNS Support (announcement and connecting to an announced SKI)
- ... work in progress

## Usage

### HEMS

```sh
Usage: go run cmd/hems/main.go <serverport> <remoteski> <certfile> <keyfile>
```

### EVSE

```sh
Usage: go run cmd/evse/main.go <serverport> <remoteski> <certfile> <keyfile>
```

### Explanation

The remoteski is from the eebus service to connect to.
If no certfile or keyfile are provided, they are generated and printed in the console so they can be saved in a file and later used again. The local SKI is also printed.

## Roadmap - Spine specification implementation

### General request processing

- [x] Request / Handle acknowledgement
- [x] Use maximum response delay to timeout requests
- [X] Send error result when processing failed
- [ ] Sending heartbeats

### Node Management

- [ ] Detailed Discovery
  - [X] Request and process full data
  - [X] Response full data 
  - [ ] Request and process partial data 
  - [ ] Response partial data
  - [ ] Notify subscribers 
- [ ] Destination List
  - [ ] Request and process full data 
  - [X] Response full data 
  - [ ] Request and process partial data 
  - [ ] Response partial data 
  - [ ] Notify subscribers 
- [ ] Binding
- [ ] Subscription
  - [X] Add subscription
  - [x] Delete subscription
  - [ ] Request and process full data 
  - [X] Response full data 
  - [ ] Request and process partial data 
  - [ ] Response partial data 
- [ ] Use Case Discovery
  - [X] Request and process full data 
  - [X] Response full data 
  - [ ] Request and process partial data 
  - [ ] Response partial data 


### General feature implementation

- [X] Request and process full data
- [X] Response full data
- [ ] Request partial data 
- [ ] Process partial data 
  - [X] Delete Selectors
  - [X] Update Selectors
  - [ ] Elements
- [ ] Response partial data
- [ ] Process write call
- [ ] Request subscription 
- [X] Notify subscribers 
- [ ] Handle incoming error results
- [X] Handle incoming success results

### Feature with partial data support

- ElectricalConnection


