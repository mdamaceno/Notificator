FROM golang:1.21-alpine

WORKDIR /opt/app/api

RUN go install github.com/cosmtrek/air@latest

ENTRYPOINT ["air", "--build.cmd", "go mod tidy && go build -o ./.bin/notificator ./main.go", "--build.bin", "./.bin/notificator"]

EXPOSE 8080
