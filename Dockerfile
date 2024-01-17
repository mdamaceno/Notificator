FROM golang:1.21-alpine

WORKDIR /opt/app/api

RUN go install github.com/cosmtrek/air@latest

ENTRYPOINT ["air", "--build.cmd", "go mod tidy && go build -o ./.bin/notificator-api ./cmd/api/main.go", "--build.bin", "./.bin/notificator-api"]

EXPOSE 8080
