FROM golang:1.24-alpine AS app_builder

WORKDIR /app
COPY ./app .

RUN go build -o bloomhub .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=app_builder /app/bloomhub .
COPY ./web ./web

ENTRYPOINT ["./bloomhub"]
