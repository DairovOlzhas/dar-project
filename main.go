package main

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	ch *amqp.Channel
	session chan bool
	conn *amqp.Connection
	queue amqp.Queue
	users []User
	user User
	window_width = 800
	window_height  = 500
)

type User struct {
	username string
	ID string
	x,y int
}


func failOnError(err error, msg string){
	if err != nil {
		log.Fatalf("%s: %s", msg, err.Error())
	}
}

func stringFrom(args []string) (s string) {
	if len(args) < 2 || args[1] == "" {
		s = "username"
	}else{
		s = strings.Join(args, " ")
	}
	return
}

func genRandString(l int) string {
	bytes := make([]byte, l)

	for i:=0; i < l; i++ {
		bytes[i] = byte(randInt(65,90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}



func main() {

	err := rabbitMQ()
	defer conn.Close()
	defer ch.Close()
	failOnError(err, "Failed to configure RabbitMQ")

	username := stringFrom(os.Args)

	err = createSession(username)
	failOnError(err, "Failed to create a session")

	users, err = getOnlineUsers()
	failOnError(err, "Failed to get users")

	run()
}

func rabbitMQ() (err error){
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	_, err = ch.QueueDeclare(
		"users",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.ExchangeDeclare(
		"commands",
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	queue, err = ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an User's queue")

	err = ch.QueueBind(
		queue.Name,
		"",
		"commands",
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue")

	return
}

func createSession(username string) (err error) {

	user = User{username,genRandString(32), randInt(50, window_width-50), randInt(50, window_height-50)}
	body, err := json.Marshal(user)

	err = ch.Publish(
		"",
		"users",
		false,
		false,
		amqp.Publishing{
			ContentType:     "application/json",
			ReplyTo: 		 queue.Name,
			CorrelationId:   user.ID,
			Body:            body,
		},
	)
	return
}

func getOnlineUsers() (users []User, err error) {

	msgs, err := ch.Consume(
		"users",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	q, err := ch.QueueInspect("users")
	failOnError(err, "Failed to inspect queue")

	n := q.Messages

	for d := range msgs {
		n--
		user := User{}
		err = json.Unmarshal(d.Body, user)
		failOnError(err, "Failed to receive User")
		users = append(users, user)
		if n == 0 {
			break
		}
	}

	return
}




func run(){

	listenCommands()

	for {
		go checkForOnlineUsers()


	}
}

func listenCommands(){

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			command := strings.Split(string(d.Body), " ")
			if len(command) < 2 {failOnError(errors.New("wrong command"), "Failed to read command")}
			id := command[1]
			switch command[0] {
			case "delete":
				for i, u := range users {
					if u.ID == id {
						users = append(users[:i], users[i+1:]...)
						break
					}
				}
				break
			case "move":
				if len(command) < 3 {failOnError(errors.New("wrong command"), "Failed to read command")}
				switch command[2] {
				case "left":
					for i, u := range users { if u.ID == id { users[i].x--; break } }
					break
				case "right":
					for i, u := range users { if u.ID == id { users[i].x++; break } }
					break
				case "up":
					for i, u := range users { if u.ID == id { users[i].y++; break } }
					break
				case "down":
					for i, u := range users { if u.ID == id { users[i].y--; break } }
					break
				}
				break
			case "online":
				if command[1] == user.ID {
					err = ch.Publish(
						"",
						d.ReplyTo,
						false,
						false,
						amqp.Publishing{
							ContentType:     "text/plain",
							Body:            nil,
						})

				}
			}
			d.Ack(false)
		}
	}()
}

func checkForOnlineUsers(){
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
		)
	failOnError(err, "Failed to declare a queue")

	for _, u := range users {

		err = ch.Publish(
			"commands",
			"",
			false,
			false,
			amqp.Publishing{
				ContentType:     "text/plain",
				ReplyTo:         q.Name,
				Body:            []byte("online "+u.ID),
			})
		failOnError(err, "Failed to pusblish a command")

		msgs, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
			)
		failOnError(err, "Failed to register a consumer")

		cnt := 0

		go func() {
			for d := range msgs {
				cnt++
				d.Ack(false)
			}

		}()

		time.Sleep(10*time.Millisecond)

		if cnt == 0 {
			sendCommand("delete "+u.ID)
		}

	}
}

func sendCommand(command string){
	err := ch.Publish(
		"commands",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(command),
		})

	failOnError(err, "Fail to publish a command")
}
