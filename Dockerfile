FROM hashicorp/terraform:1.14.1 AS terraform-bin
FROM golang:1.21-bullseye as test

COPY --from=terraform-bin /bin/terraform /usr/local/bin/terraform

WORKDIR /build

COPY . .
RUN go install github.com/jstemmer/go-junit-report/v2@latest && go mod download
