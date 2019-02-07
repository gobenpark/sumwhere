package test

import (
	"testing"
)
import "github.com/streadway/amqp"

func TestRabbitMQ(t *testing.T) {
	conn, err := amqp.Dial("amqp://qkrqjadn:1q2w3e4r@192.168.1.11:5672/")
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Error(err)
	}

	err = ch.ExchangeDeclare("Room_1", "fanout", true, false, false, false, nil)
	if err != nil {
		t.Error(err)
	}

	body := "1123123123123123131!"
	err = ch.Publish("Room_1", "", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
}
