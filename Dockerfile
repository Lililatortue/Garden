FROM golang:1.16-alpine AS builder

WORKDIR /app
COPY ./app .

RUN go build -o bloomhub .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/bloomhub .

ENTRYPOINT ["./bloomhub"]
