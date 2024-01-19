package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	uri          = "amqp://guest:guest@localhost:5672/"
	exchangeName = "notificator"
	exchangeType = "direct"
	queueName    = "notificator-queue"
	bindingKey   = "notificator-key"
	ctag         = "notificator-producer"
	continuous   = false
	service      = flag.String("s", "", "service names separated by comma")
	title        = flag.String("t", "", "title of the message")
	body         = flag.String("m", "", "body of the message")
	receivers    = flag.String("r", "", "receivers separated by comma")
	WarnLog      = log.New(os.Stderr, "[WARNING] ", log.LstdFlags|log.Lmsgprefix)
	ErrLog       = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)
	Log          = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix)
)

func init() {
	flag.Parse()
}

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
}

type Message struct {
	Service   []string
	Title     string
	Body      string
	Receivers []string
}

func main() {
	var err error

	message := Message{
		Service:   strings.Split(*service, ","),
		Title:     *title,
		Body:      *body,
		Receivers: strings.Split(*receivers, ","),
	}

	err = validateMessageOpts(&message)

	if err != nil {
		ErrLog.Fatalf("%s", err)
	}

	producer, err := NewProducer(uri, exchangeName, exchangeType, queueName, bindingKey, ctag)
	defer producer.channel.Close()

	requestBody, err := json.Marshal(message)

	err = producer.channel.PublishWithContext(
		context.Background(),
		exchangeName, // publish to an exchange
		bindingKey,   // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(requestBody),
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

func validateMessageOpts(message *Message) error {
	if message == nil {
		return fmt.Errorf("message is nil")
	}

	if len(message.Service) == 0 {
		return fmt.Errorf("service is empty")
	}

	if len(message.Title) == 0 {
		return fmt.Errorf("title is empty")
	}

	if len(message.Body) == 0 {
		return fmt.Errorf("body is empty")
	}

	if len(message.Receivers) == 0 {
		return fmt.Errorf("receivers is empty")
	}

	return nil
}

func NewProducer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Producer, error) {
	var err error

	p := &Producer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
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

	Log.Printf("Queue bound to Exchange, starting Produce (producer tag %q)", p.tag)

	if err := p.channel.Confirm(false); err != nil {
		ErrLog.Fatalf("producer: channel could not be put into confirm mode: %s", err)
	}

	return p, nil
}
