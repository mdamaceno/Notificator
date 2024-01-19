# Notificator

A simple API service to deliver messages by email, sms and whatsapp.

## How to use

Create an `.env` file in the root of the project. The variables accepted by the app are:

```
# App info
APP_ID=notificator
APP_PORT=8080

# Database config
DB_HOST=db
DB_USER=postgres
DB_PASS=postgres
DB_NAME=notificator_dev
DB_PORT=5432
POSTGRES_DRIVER=postgres
POSTGRES_SOURCE=postgresql://postgres:postgres@db:5432/notificator_dev?sslmode=disable

# Testing database config
DB_TEST_HOST=dbtest
DB_TEST_USER=postgres
DB_TEST_PASS=postgres
DB_TEST_NAME=notificator_test
DB_TEST_PORT=5432

# Redis config
REDIS_PORT=63790

# Email service config
EMAIL_HOST=mailcatcher
EMAIL_PORT="1025"
EMAIL_FROM=from@email.com

# SMS service config
SMS_FROM_NUMBER=
SMS_USERNAME=
SMS_PASSWORD=

# Whatsapp service config
WHATSAPP_FROM_NUMBER=
WHATSAPP_USERNAME=${SMS_USERNAME} # this copies the value from SMS_USERNAME
WHATSAPP_PASSWORD=${SMS_PASSWORD} # this copies the value from SMS_PASSWORD
```

You must have Docker and Docker Compose installed. Just type `docker compose up` in the terminal to run it.

Following this instructions, the API will be available to accept requests on http://localhost:8080.

For example:

```
Content-Type: application/json

POST /v1/message

{
    "service": ["email", "sms", "whatsapp"],
    "title": "Hey!",
    "body": "What's up?",
    "receivers": [
        "johndoe@email.com",
        "maryjane@email.com",
        "+5511988880000"
    ]
}
```

You can use whatever tool to test the request: curl, Insomnia, Postman, Bruno...

**Mailcatcher** is in charge to intercept messages that use SMTP protocol on port 1025. You can change it if you want to use any other email service.
The emails sent will be available on http://localhost:1080.

**SMS** and **Whatsapp** messages are delivered by Twilio. You must have a Twilio account configured.

Running all the containers, you can use the **producer** to send messages to the consumer using AMQP protocol.

Open another tab in you terminal and run:

```
go run cmd/producer/main.go -s "email,sms,whatsapp" -t "Title" -m "Body" -r "email1@email.com,email2@email.com"
```

## Test

To test the app, just run `docker compose exec app01 go test ./...`.

## TODO

- Slack Integration
- Integration tests
