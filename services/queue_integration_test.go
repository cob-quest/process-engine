// +build integration

package services

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/streadway/amqp"
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
	suite.QueueName = "test_queue"
	suite.ExchangeName = "test_exchange"
	suite.RoutingKey = "test.key"

	conn, err := amqp.Dial("amqp://test:test@rabbitmq:5672/")
	suite.NoError(err)
	suite.Connection = conn

	ch, err := conn.Channel()
	suite.NoError(err)
	suite.Channel = ch
}

func (suite *RabbitMQSuite) SetupTest() {
	// Setup code before each test here
	err := suite.Channel.ExchangeDeclare(
		suite.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	suite.NoError(err)

	_, err = suite.Channel.QueueDeclare(
		suite.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	suite.NoError(err)

	err = suite.Channel.QueueBind(
		suite.QueueName,
		suite.RoutingKey,
		suite.ExchangeName,
		false,
		nil,
	)
	suite.NoError(err)
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

// Example of a test function
func (suite *RabbitMQSuite) TestExample() {
	// Your test code here
	// Use suite.Channel to interact with RabbitMQ
	suite.True(true) // Dummy assertion
}
