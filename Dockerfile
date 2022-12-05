# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine3.15

ENV PORT=9200

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /metrics-filter

CMD [ "/metrics-filter" ]
