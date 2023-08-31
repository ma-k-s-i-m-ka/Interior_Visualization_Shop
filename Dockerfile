FROM golang:1.20-alpine AS builder

WORKDIR /usr/local/src

COPY . ./

RUN go build -o ./bin/app app/cmd/main.go

FROM alpine

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/.env /
COPY --from=builder /usr/local/src/config.yml /

CMD ["/app"]