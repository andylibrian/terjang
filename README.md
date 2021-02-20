# Terjang
Scalable HTTP load testing tool built on Vegeta


## Build for production

```bash
make
```

## Development

Build backend:

```bash
go build -o ./bin/terjang ./cmd/terjang/
```

Run backend:

```
./bin/terjang server
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


