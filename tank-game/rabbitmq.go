package tank_game

import (
	"github.com/streadway/amqp"
	"log"
)

var (
	conn *amqp.Connection
	ch *amqp.Channel
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	onlinePlayersQueue = "onlinePlayers"
	commandsExchange = "commands"
)

func failOnError(err error, msg string, ok string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
	log.Printf(" [*] " + ok)
}

func RabbitMQ() (err error){
	conn, err = amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect RbbitMQ", "Connected to RabbitMQ")

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
		false,
		false,
		false,
		nil,
		)

	return
}

func CloseConnectionAndChannel() (err error){
	err = ch.Close()
	failOnError(err, "Failed to close channel", "Channel closed")
	err = conn.Close()
	failOnError(err, "Failed to close connection", "Connection closed")
	return
}

func GetOnlinePlayersQueueName() string {
	return onlinePlayersQueue
}
func GetCommandsExchangeName() string {
	return commandsExchange
}