package main

import (
	"fmt"

	"github.com/mdamaceno/notificator/app/controllers"
	"github.com/mdamaceno/notificator/config"
	"github.com/mdamaceno/notificator/internal/db"
	"github.com/mdamaceno/notificator/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	tag        string
	deliveries <-chan amqp.Delivery
	done       chan error
}

var (
	uri                = "amqp://guest:guest@broker:5672/"
	exchangeName       = "notificator"
	exchangeType       = "direct"
	queueName          = "notificator-queue"
	bindingKey         = "notificator-key"
	consumerTag        = "notificator-consumer"
	connectionName     = "notificator-consumer"
	deliveryCount  int = 0
)

func main() {
	dbconn, err := config.InitDB()

	if err != nil {
		helpers.ErrLog.Printf("Database: %s", err)
	}

	defer dbconn.Close()

	q := db.New(dbconn)

	c, err := NewConsumer(uri, exchangeName, exchangeType, queueName, bindingKey, consumerTag)
	if err != nil {
		helpers.ErrLog.Fatalf("%s", err)
	}

	c.done = make(chan error)
	helpers.Log.Println("running consumer...")

	messageController := controllers.MessageController{DB: dbconn, Queries: q}
	messageController.Consume(c.deliveries, c.done)

	<-c.done
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	var err error

	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName(connectionName)
	helpers.Log.Printf("dialing %q", amqpURI)

	c.conn, err = amqp.DialConfig(amqpURI, config)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		helpers.Log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	helpers.Log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}
	helpers.Log.Printf("got Channel, declaring Exchange (%q)", exchange)

	if err = c.channel.ExchangeDeclare(
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

	helpers.Log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	helpers.Log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	helpers.Log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	c.deliveries = deliveries

	return c, nil
}
