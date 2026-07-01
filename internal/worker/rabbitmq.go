package worker

import (
	"context"
	"github.com/danmaina/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"messaging/internal/configs"
	"messaging/internal/mailer"
)

// StartConsumer connects to RabbitMQ and starts consuming messages from the queue
func StartConsumer(ctx context.Context, config *configs.Config) error {
	conn, err := amqp.Dial(config.RabbitMQ.URL)
	if err != nil {
		logger.ERR("Failed to connect to RabbitMQ: ", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.ERR("Failed to open a channel: ", err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"NOTIFICATIONS", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		logger.ERR("Failed to declare a queue: ", err)
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.ERR("Failed to register a consumer: ", err)
		return err
	}

	logger.INFO("RabbitMQ Consumer started. Waiting for messages on NOTIFICATIONS queue...")

	for {
		select {
		case <-ctx.Done():
			logger.INFO("Context done, stopping consumer")
			return nil
		case d, ok := <-msgs:
			if !ok {
				logger.ERR("RabbitMQ channel closed")
				return nil
			}
			err := mailer.ProcessEmailRequest(d.Body, config)
			if err != nil {
				// We can implement retry logic or DLQ here, but for now we just log and reject/requeue
				logger.ERR("Error processing message: ", err)
				d.Nack(false, false) // reject, do not requeue for now to prevent infinite loops on bad payloads
			} else {
				d.Ack(false)
			}
		}
	}
}
