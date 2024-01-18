FROM golang:1.21-alpine

WORKDIR /opt/app/consumer

RUN go install github.com/cosmtrek/air@latest

ENTRYPOINT ["air", "--build.exclude_dir", "cmd/api,cmd/producer,docker,sql", "--build.cmd", "go mod tidy && go build -o ./.bin/notificator-consumer ./cmd/consumer/main.go", "--build.bin", "./.bin/notificator-consumer"]
