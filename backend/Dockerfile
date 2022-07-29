FROM golang:1.19 AS builder
ARG SERVICE_NAME
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY ./libraries ./libraries
COPY ./services ./services
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/service ./services/$SERVICE_NAME

FROM alpine:3.16
COPY --from=builder /src/bin /bin
CMD ["service"]