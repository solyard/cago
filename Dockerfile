FROM golang:alpine as build-env
LABEL maintainer="dizstorm@gmail.com"
COPY . /app
WORKDIR /app
RUN go mod download && go build -o /usr/bin/cert-manager-issuer .
ENTRYPOINT ["cert-manager-issuer"]