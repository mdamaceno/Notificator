package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
}

var (
	uri          = "amqp://guest:guest@localhost:5672/"
	exchangeName = "notificator"
	exchangeType = "direct"
	queueName    = "notificator-queue"
	bindingKey   = "notificator-key"
	continuous   = false
	body         = flag.String("m", "Message", "body of message")
	WarnLog      = log.New(os.Stderr, "[WARNING] ", log.LstdFlags|log.Lmsgprefix)
	ErrLog       = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)
	Log          = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix)
)

func init() {
	flag.Parse()
}

func main() {
	var err error

	producer, err := NewProducer(uri, exchangeName, exchangeType, queueName, bindingKey, "notificator-producer")
	defer producer.channel.Close()

	err = producer.channel.PublishWithContext(
		context.Background(),
		exchangeName, // publish to an exchange
		bindingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(*body),
		},
	)

	if err != nil {
		ErrLog.Fatalf("%s", err)
	}

	Log.Printf("Message sent: %s", *body)
}

func openConnection(amqpURI string) (*amqp.Connection, error) {
	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	Log.Printf("dialing %q", amqpURI)
	conn, err := amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	return conn, nil
}

func NewProducer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Producer, error) {
	var err error

	p := &Producer{
		conn:    nil,
		channel: nil,
		tag:     "",
	}

	p.conn, err = openConnection(amqpURI)

	if err != nil {
		ErrLog.Fatalf("%s", err)
	}

	Log.Printf("got Connection, getting Channel")

	p.channel, err = p.conn.Channel()

	if err != nil {
		ErrLog.Fatalf("%s", err)
	}

	Log.Printf("got Channel, declaring Exchange (%q)", exchange)

	if err = p.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	Log.Printf("declared Exchange, declaring Queue %q", queueName)

	queue, err := p.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)

	if err != nil {
		ErrLog.Fatalf("%s", err)
	}

	err = p.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	)

	if err != nil {
		ErrLog.Fatalf("%s", err)
	}

	Log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", ctag)

	if err := p.channel.Confirm(false); err != nil {
		ErrLog.Fatalf("producer: channel could not be put into confirm mode: %s", err)
	}

	return p, nil
}
