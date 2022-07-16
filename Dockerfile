# syntax=docker/dockerfile:1

FROM golang:1.18

WORKDIR /app

COPY . ./

ENV GOFLAGS=-mod=vendor
RUN go build -o mend-home-assignment cmd/main.go

EXPOSE 443

CMD ./mend-home-assignment
