[![CircleCI](https://circleci.com/gh/RoboCup-SSL/ssl-game-controller/tree/master.svg?style=svg)](https://circleci.com/gh/RoboCup-SSL/ssl-game-controller/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/RoboCup-SSL/ssl-game-controller?style=flat-square)](https://goreportcard.com/report/github.com/RoboCup-SSL/ssl-game-controller)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/RoboCup-SSL/ssl-game-controller/internal/app/controller)
[![Release](https://img.shields.io/github/release/RoboCup-SSL/ssl-game-controller.svg?style=flat-square)](https://github.com/RoboCup-SSL/ssl-game-controller/releases/latest)

# ssl-game-controller

The [ssl-refbox](https://github.com/RoboCup-SSL/ssl-refbox) replacement that will be introduced at RoboCup 2019

## Usage
If you just want to use this app, simply download the latest [release binary](https://github.com/RoboCup-SSL/ssl-game-controller/releases/latest). The binary is self-contained. No dependencies are required.

### Reference Clients
There are some reference clients:
 * [ssl-ref-client](./cmd/ssl-ref-client): A client that receives referee messages
 * [ssl-auto-ref-client](./cmd/ssl-auto-ref-client/README.md): A client that connects to the controller as an autoRef

## Development

### Requirements
You need to install following dependencies first: 
 * Go >= 1.9
 * NPM

### Prepare
Download and install to [GOPATH](https://github.com/golang/go/wiki/GOPATH):
```bash
go get -u github.com/RoboCup-SSL/ssl-game-controller/...
```
Switch to project root directory
```bash
cd $GOPATH/src/github.com/RoboCup-SSL/ssl-game-controller/
```
Download dependencies for frontend
```bash
npm install
```

### Run
Run the backend:
```bash
go run cmd/ssl-game-controller/main.go
```

Run the UI:
```bash
# compile and hot-reload
npm run serve
```
Or use the provided IntelliJ run configurations.

### Build self-contained release binary
First, build the UI resources
```bash
# compile and minify UI
npm run build
```
Then build the backend with `packr`
```bash
# get packr
go get github.com/gobuffalo/packr/packr
# install the binary
cd cmd/ssl-game-controller
packr install
```