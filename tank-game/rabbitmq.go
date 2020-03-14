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
	checkOnlinePlayersQueue = "checkOnline"
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
		true,
		false,
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

	_, err = ch.QueueDeclare(
		checkOnlinePlayersQueue,
		false,
		true,
		false,
		false,
		nil,
		)
	failOnError(err, "Failed to declare checkOnlinePlayersQueue", "checkOnlinePlayersQueue declared")
	return
}

func CloseConnectionAndChannel() (err error){
	err = ch.Close()
	failOnError(err, "Failed to close channel", "Channel closed")
	err = conn.Close()
	failOnError(err, "Failed to close connection", "Connection closed")
	return
}

func Channel() *amqp.Channel{
	return ch
}
func OnlinePlayersQueueName() string {
	return onlinePlayersQueue
}
func CommandsExchangeName() string {
	return commandsExchange
}
func ReceiverQueue() amqp.Queue {
	return receiverQueue
}
func CheckOnlinePlayersQueue() string {
	return checkOnlinePlayersQueue
}
