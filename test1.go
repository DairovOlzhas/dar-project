package main

import (
	"github.com/streadway/amqp"
	"os"
)

var (
	mp map[string]*class
)


type class struct {
	A int    `json:"a"`
	B string `json:"b"`
	//c tl.Attr
}

type response1 struct {
	Page   int
	Fruits []string
}

func main() {

	_ = ch.Publish(
		"",
		"asdf",
		false,
		false,
		amqp.Publishing{
			CorrelationId:   "asdf",
			ReplyTo:         "asdfasdf",
			Body:            []byte(os.Args[1]),
		},
	)
}