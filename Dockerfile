FROM golang:1.15-alpine as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 go build -a -o terjang ./cmd/terjang/



FROM alpine:3.13

RUN apk --no-cache add ca-certificates
USER nobody
COPY --from=builder --chown=nobody:nobody /workspace/terjang .

ENTRYPOINT ["./terjang"]
