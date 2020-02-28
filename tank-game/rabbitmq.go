package tank_game

import (
	"github.com/streadway/amqp"
)

var (
	conn               *amqp.Connection
	ch                 *amqp.Channel
	rabbitMQURL        = "amqp://guest:guest@localhost:5672/"
	onlinePlayersQueue = "onlinePlayers"
	commandsExchange   = "commands"
	receiverQueue      amqp.Queue
)


func RabbitMQ() (err error){
	conn, err = amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect RabbitMQ", "Connected to RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a Channel", "Channel opened")

	err = ch.ExchangeDeclare(
		commandsExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
		)

	_, err = ch.QueueDeclare(
		onlinePlayersQueue,
		false,
		true,
		false,
		false,
		nil,
		)
	failOnError(err, "Failed to declare onlinePlayersQueue", "onlinePlayersQueue declared")

	receiverQueue, err = ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	failOnError(err, "Failed to declare receiverQueue", "receiverQueue declared")

	err = ch.QueueBind(
		receiverQueue.Name,
		"",
		commandsExchange,
		false,
		nil)
	failOnError(err, "Failed to bind receiverQueue", "receiverQueue bound")

	return
}

func CloseConnectionAndChannel() (err error){
	err = ch.Close()
	failOnError(err, "Failed to close channel", "Channel closed")
	err = conn.Close()
	failOnError(err, "Failed to close connection", "Connection closed")
	return
}

func OnlinePlayersQueueName() string {
	return onlinePlayersQueue
}
func CommandsExchangeName() string {
	return commandsExchange
}
func Channel() *amqp.Channel{
	return ch
}
func ReceiverQueueName() string{
	return receiverQueue.Name
}