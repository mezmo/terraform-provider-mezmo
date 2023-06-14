FROM golang:1.20-bullseye as test

WORKDIR /build

COPY . .
RUN go mod download
