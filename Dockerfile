FROM golang:1.22.4-alpine

WORKDIR /src
COPY ./src /src

RUN go mod tidy && go build -o ./main

ENTRYPOINT ./main