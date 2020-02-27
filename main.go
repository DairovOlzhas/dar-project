package main

import (
	tank_game "github.com/dairovolzhas/dar-project/tank-game"
	"log"
)

func failOnError(err error, msg string, ok string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
	log.Printf(" [*] " + ok)
}


func main() {

	err := tank_game.RabbitMQ()
	failOnError(err, "Failed to configure RabbitMQ", "RabbitMQ configured")
	defer tank_game.CloseConnectionAndChannel()
	
}
