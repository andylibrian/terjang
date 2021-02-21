# Terjang
Scalable HTTP load testing tool built on Vegeta


## Build for production

```bash
make
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

```
./bin/terjang server
```

Run a worker on another terminal:
```
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


