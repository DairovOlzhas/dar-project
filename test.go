package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

var (
	chh *amqp.Channel
)


func main() {

	conn, _  := amqp.Dial("amqp://guest:guest@localhost:5672/")
	chh, _ = conn.Channel()

	_, _ = chh.QueueDeclare(
		"dfg",
		false,
		false,
		true,
		false,
		nil,
	)

	go func() {
		for d := range chh.NotifyReturn(make(chan amqp.Return)) {
			fmt.Println(string(d.Body))
		}
	}()

	for i:=0; i < 10; i++ {
		_, err := chh.QueueDeclarePassive(
			"dfg",
			false,
			false,
			true,
			false,
			nil,
			)
		if err != nil {
			fmt.Println("asdfg")
		}else{
			fmt.Println("kkkk")
		}
	}

	forever := make(chan bool)
	<- forever
}