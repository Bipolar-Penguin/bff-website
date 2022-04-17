package amqp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/streadway/amqp"

	"github.com/Bipolar-Penguin/bff-website/pkg/domain"
)

type RabbitBroker struct {
	connString string
	logger     log.Logger
}

func NewRabbitBroker(connString string, logger log.Logger) *RabbitBroker {
	return &RabbitBroker{connString, logger}
}

func (b *RabbitBroker) PublishEvent(event domain.Event) {
	b.logger.Log("event", fmt.Sprintf("%v", event))

	jsonBody, err := json.Marshal(event)
	if err != nil {
		b.logger.Log("error", err)
	}

	conn, err := amqp.Dial(b.connString)
	if err != nil {
		b.logger.Log("error", err)
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		b.logger.Log("error", err)
		time.Sleep(5 * time.Second)
	}

	_, err = ch.QueueDeclare(
		"trading_session.update", // name
		false,                    // durable
		true,                     // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		map[string]interface{}{"x-message-ttl": 10000}, // arguments
	)

	err = ch.Publish(
		"",                       // exchange
		"trading_session.update", // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonBody),
		})
	if err != nil {
		b.logger.Log("error", err)
	}
}
