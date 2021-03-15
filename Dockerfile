FROM alpine:3.13

RUN apk --no-cache add ca-certificates
USER nobody
COPY terjang .

ENTRYPOINT ["./terjang"]
