FROM golang:1.24-alpine

WORKDIR /test

COPY ./app .

ENTRYPOINT ["go", "test"]