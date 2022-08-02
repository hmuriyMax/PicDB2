FROM golang:1.18.4-alpine

RUN mkdir -p /usr/src/app/
WORKDIR /usr/src/app/
COPY . /usr/src/app/

RUN go build main.go
EXPOSE 80
ENTRYPOINT go run main.go 80