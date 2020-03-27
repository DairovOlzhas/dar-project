package game

import (
	"github.com/streadway/amqp"
	"log"
)

var (
	conn               	*amqp.Connection
	ch                 	*amqp.Channel
	rabbitMQURL        		= 	"amqp://tanks:sQcp3CHep58G@35.184.207.230:5672/"
	onlinePlayersQueue 		= 	"onlinePlayers"
	commandsExchange   		= 	"commands"
	receiverQueue      	amqp.Queue
	checkOnlinePlayersQueue = 	"checkOnline"
)

func failOnError(err error, msg string, ok string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
	if ok != ""{
		log.Printf(" [*] " + ok)
	}
}

func RabbitMQ(){
	conn, err := amqp.Dial(rabbitMQURL)
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
	failOnError(err, "Failed to open a Channel",  commandsExchange + " exchange declared")

	QueueDeclare(onlinePlayersQueue, true)
	receiverQueue = QueueDeclare("", true)
	QueueBind(receiverQueue.Name, commandsExchange)
	QueueDeclare(checkOnlinePlayersQueue, true)
}

func CloseConnectionAndChannel() (err error){
	err = ch.Close()
	failOnError(err, "Failed to close channel", "Channel closed")
	err = conn.Close()
	failOnError(err, "Failed to close connection", "Connection closed")
	return
}


func QueueDeclare(qname string, autodelete bool) (queue amqp.Queue) {
	queue, err := ch.QueueDeclare(
		qname,
		false,
		autodelete,
		false,
		false,
		nil,
		)
	failOnError(err, "Failed to declare Queue " + queue.Name, queue.Name + " queue declared")
	return
}

func QueueBind(qname string, exchangeName string){
	err := ch.QueueBind(
		qname,
		"",
		exchangeName,
		false,
		nil)
	failOnError(err, "Failed to bind " + qname + " queue to " +exchangeName+ " exchange.",
		qname + " queue binded to " +exchangeName+ " exchange")
}


func Consumer(qname string) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		qname,
		"",
		false,
		false,
		false,
		false,
		nil,)
	failOnError(err, "Failed to register consumer to " + qname + " queue", "")
	return msgs
}

func Publisher(exchange string, key string, publishing amqp.Publishing){
	err := ch.Publish(
		exchange,
		key,
		false,
		false,
		publishing)
	failOnError(err, "Failed to publish message", "")
}


func OnlinePlayersQueue() string{ return onlinePlayersQueue}

func CommandsExchange() string{ return commandsExchange}

func ReceiverQueue() string{ return receiverQueue.Name}

func CheckOnlinePlayersQueue() string{ return checkOnlinePlayersQueue}
