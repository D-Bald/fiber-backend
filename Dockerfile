# syntax=docker/dockerfile:1
FROM golang:latest

RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE ${FIBER_PORT}

CMD ["/app/main"]