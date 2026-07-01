package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"messaging/internal/configs"
	"messaging/internal/mailer"
	"messaging/internal/models"

	"github.com/danmaina/infra/v2/rabbitmq"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestRabbitMQConsumerAndMailer(t *testing.T) {
	// Load .env from the parent directory
	_ = godotenv.Load("../.env")

	// 1. Read configs
	config, err := configs.ReadConfigs()
	if err != nil {
		t.Fatalf("Failed to read configs: %v", err)
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(config.RabbitMQ.URL)
	if err != nil {
		t.Fatalf("RabbitMQ not available at %s, cannot run integration test: %v", config.RabbitMQ.URL, err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	queueName := "NOTIFICATIONS_TEST"

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to declare queue: %v", err)
	}

	// 2. Setup the consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processed := make(chan bool)

	consumer := &rabbitmq.RabbitMQConsumer{
		URL:       config.RabbitMQ.URL,
		QueueName: queueName,
		Handler: func(body []byte) error {
			// This calls the mailer logic to dispatch the email
			err := mailer.ProcessEmailRequest(body, config)
			if err == nil {
				processed <- true
			} else {
				t.Errorf("ProcessEmailRequest failed: %v", err)
				processed <- false
			}
			return err
		},
	}

	go func() {
		err := consumer.Start(ctx)
		if err != nil {
			t.Logf("Consumer stopped: %v", err)
		}
	}()

	// Give consumer a moment to start
	time.Sleep(1 * time.Second)

	// 3. Publish a test message
	testMsg := models.EmailMessage{
		To:      []string{"test@example.com"},
		From:    "sender@example.com",
		Subject: "Integration Test Email",
		Message: "This is a test email dispatched from the integration test.",
	}

	body, _ := json.Marshal(testMsg)

	err = ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// 4. Wait for processing
	select {
	case success := <-processed:
		if !success {
			t.Fatal("Failed to process and dispatch email.")
		}
		t.Log("Successfully received RabbitMQ message and dispatched email!")
	case <-time.After(30 * time.Second):
		t.Fatal("Test timed out waiting for message processing")
	}
}
