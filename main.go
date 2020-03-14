package main

import (
	tg "github.com/dairovolzhas/dar-project/tank-game"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	fps = 60
)

func failOnError(err error, msg string, ok string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
	log.Printf(" [*] " + ok)
}


func main() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	rand.Seed(time.Now().UnixNano())

	err = tg.RabbitMQ()
	failOnError(err, "Failed to configure RabbitMQ", "RabbitMQ configured")
	defer tg.CloseConnectionAndChannel()

	tg.CreateGame(fps)
	tg.StartGame()

}
