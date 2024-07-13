FROM golang:1.21-bullseye as test

WORKDIR /build

COPY . .
RUN go install github.com/jstemmer/go-junit-report/v2@latest && go mod download
