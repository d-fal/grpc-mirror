FROM golang:1.12-alpine

RUN mkdir -p /app
ADD . /app

WORKDIR /app


COPY . .

RUN apk update && \
    apk add  git 

RUN go get -d ./...

RUN go build -o main .




ENTRYPOINT ["/app/main"]