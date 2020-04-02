package main

import (
	"github.com/dairovolzhas/dar-project/game"
	"math/rand"
	"time"
)

func main() {
	//file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()
	//
	//log.SetOutput(file)
	//
	rand.Seed(time.Now().UnixNano())

	game.RabbitMQ()
	defer game.CloseConnectionAndChannel()

	game.Game().Start()
}