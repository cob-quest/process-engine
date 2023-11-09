// +build integration

package services

import (
	"context"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/suite"
)

type RabbitMQSuite struct {
	suite.Suite
	Connection *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
	ExchangeName string
	RoutingKey string
}

func (suite *RabbitMQSuite) SetupSuite() {
	// Setup code for the entire suite here
	suite.QueueName = "test.queue.publish"
	suite.ExchangeName = "topic.router"
	suite.RoutingKey = "test.routing.key"

	conn, err := amqp.Dial("amqp://test:test@rabbitmq:5672/")
	suite.NoError(err)
	suite.Connection = conn

	ch, err := conn.Channel()
	suite.NoError(err)
	suite.Channel = ch
}

func (suite *RabbitMQSuite) TearDownTest() {
	// Teardown code after each test here
	err := suite.Channel.QueueUnbind(
		suite.QueueName,
		suite.RoutingKey,
		suite.ExchangeName,
		nil,
	)
	suite.NoError(err)

	_, err = suite.Channel.QueueDelete(suite.QueueName, false, false, false)
	suite.NoError(err)

	err = suite.Channel.ExchangeDelete(suite.ExchangeName, false, false)
	suite.NoError(err)
}

func (suite *RabbitMQSuite) TearDownSuite() {
	// Teardown code for the entire suite here
	err := suite.Channel.Close()
	suite.NoError(err)

	err = suite.Connection.Close()
	suite.NoError(err)
}

func TestRabbitMQSuite(t *testing.T) {
	suite.Run(t, new(RabbitMQSuite))
}

func (suite *RabbitMQSuite) TestPublish() {
    // Set up test message
    testMessage := []byte("Test Publish message")

    // Declare a test exchange
    err := suite.Channel.ExchangeDeclare(
        "topic.router", // exchange name
        "topic",                // exchange type
        true,                   // durable
        false,                  // auto-deleted
        false,                  // internal
        false,                  // no-wait
        nil,                    // arguments
    )
    suite.NoError(err)

    // Declare a test queue
    q, err := suite.Channel.QueueDeclare(
        "test.queue.publish", // queue name
        false,                // durable
        false,                // delete when unused
        true,                 // exclusive
        false,                // no-wait
        nil,                  // arguments
    )
    suite.NoError(err)

    // Bind the test queue to the test exchange
    err = suite.Channel.QueueBind(
        q.Name,                  // queue name
        "test.routing.key",      // routing key
        "topic.router", // exchange
        false,
        nil,
    )
    suite.NoError(err)

    // Run the Publish function
    ctx := context.Background()
    Publish(suite.Channel, ctx, testMessage, "test.routing.key")

    // Try to consume the message
    msgs, err := suite.Channel.Consume(
        q.Name,  // queue
        "",      // consumer
        true,    // auto-ack
        false,   // exclusive
        false,   // no-local
        false,   // no-wait
        nil,     // args
    )
    suite.NoError(err)

    // Use a select statement to wait for a message or timeout
    select {
    case d := <-msgs:
        suite.Equal(testMessage, d.Body, "The message body should be equal to the published message.")
    case <-time.After(5 * time.Second):
        suite.FailNow("Failed to receive message in time")
    }

    // Cleanup
    suite.Channel.QueueDelete(q.Name, false, false, true)
    suite.Channel.ExchangeDelete("test.exchange.publish", false, false)
}
