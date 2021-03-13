# Terjang
Scalable HTTP load testing tool built on [Vegeta](https://github.com/tsenart/vegeta).

[![Build status](https://img.shields.io/github/workflow/status/andylibrian/terjang/CI?style=flat)](https://github.com/andylibrian/terjang/actions)


## Features

- Scalable: support multiple node of workers
- Web UI with detailed report
- Extensible:
  - Start and stop load test via HTTP API
  - Receive real time results via websocket

![Demo](docs/demo.gif?raw=true "Demo")

## Install

### Pre-compiled

Download the pre-compiled executables from the [releases page](https://github.com/andylibrian/terjang/releases) and copy to the desired location.

### From source

```bash
git clone git@github.com:andylibrian/terjang.git
cd terjang
make
```

## Usage

### Quick start on your local machine

Open a terminal, and run:

```bash
terjang server
```

Open another terminal, and run:

```bash
terjang worker
```

Then open [http://localhost:9009](http://localhost:9009)

### See more options

```bash
terjang -h
```

## Development

Prerequisites:
- Go (tested version: 1.15)
- Node (tested version: 14)

Build backend:

```bash
go build -o ./bin/terjang ./cmd/terjang/
```

Run server:

```bash
./bin/terjang server
```

Run a worker on another terminal:

```bash
./bin/terjang worker
```

Install frontend:

```bash
cd web
npm install
```

Run frontend:

```bash
./node_modules/.bin/vue-cli-service serve
```

Then open http://localhost:8080


