FROM golang:latest

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

COPY . .

RUN go mod tidy

CMD [ "go","run","main.go" ]
