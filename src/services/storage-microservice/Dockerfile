FROM golang:1.23-alpine AS builder

WORKDIR /service/app


COPY /services/storage-microservice .

COPY /pkg/logger ../../pkg/logger


RUN ls -l

RUN go mod download

WORKDIR /service/app/cmd/main

RUN go build -o ./storage-service


FROM alpine:latest

WORKDIR /app/configs

COPY --from=builder /service/app/configs .

WORKDIR /app

COPY --from=builder /service/app/cmd/main/storage-service .

EXPOSE 8085
EXPOSE 9091

CMD ["./storage-service"]