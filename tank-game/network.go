package tank_game

import (
	"encoding/json"
	tl "github.com/JoelOtter/termloop"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var (
	players         = make(map[string] *Player)
	playerToDelete  = make(map[string] bool)
	CurrentPlayerID = ""
)

type PlayerOnline struct{
	ID        string
	Username  string
	Score     int
	X         int
	Y         int
	Direction int
	Color     tl.Attr
}

type Command struct {
	ID string
	Action int
	X int
	Y int
	Direction int
}
// command action
const (
	TANK = 0
	BULLET = 1
	DELETE = 2
	POSITION = 3
)


func getOnlinePlayers() {

	msgs, err := ch.Consume(
		onlinePlayersQueue,
		"",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to register a consumer(getting Online players)", "Consumer registered(getting Online players)")

	go func() {

		for d := range msgs {

			playerOnline := PlayerOnline{}

			_ = json.Unmarshal(d.Body, &playerOnline)
			//failOnError(err, "Failed to convert getter player", "Player getted")
			if _, prs := players[playerOnline.ID]; !prs {
				log.Printf("Online player received %v", playerOnline)
				players[playerOnline.ID] = &Player{
					ID: 	  playerOnline.ID,
					Username: playerOnline.Username,
					//Tank:     NewTankXY(0, 0, tl.Cell{Fg: 0, Color: tl.ColorBlack, Ch: 0,}),
					Tank:     NewTankXY(playerOnline.X, playerOnline.Y, tl.Cell{Bg: playerOnline.Color}, playerOnline.Direction),
					PreX:     0,
					PreY:     0,
					Score:    playerOnline.Score,

				}

				log.Printf("%v", players[playerOnline.ID])
				Level.AddEntity(players[playerOnline.ID])
			}

			if _, prs := playerToDelete[playerOnline.ID]; prs{
				d.Ack(false)
			}else{
				d.Nack(false, true)
			}
		}
	}()
}

func listenCommands() {
	msgs, err := ch.Consume(
		receiverQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to register a Consumer(listencommands)", "Consumer registered(listencommands)")

	// sleep to get Online players
	time.Sleep(5*time.Second)
	go func() {

		for d := range msgs {
			a := Command{}
			err = json.Unmarshal(d.Body, &a)
			failOnError(err, "Failed to unmarshal command", "Command getted")

			switch a.Action {
			case TANK:
				p := players[a.ID]
				p.PreX, p.PreY = p.Position()
				cell := tl.Cell{Bg: p.color}
				switch a.Direction {
				case UP:
					TankUp(p.Tank, cell)
					p.SetPosition(p.PreX, p.PreY-1)
					p.SetPosition(a.X, a.Y-1)
				case DOWN:
					TankDown(p.Tank, cell)
					p.SetPosition(p.PreX, p.PreY+1)
					p.SetPosition(a.X, a.Y+1)
				case LEFT:
					TankLeft(p.Tank, cell)
					p.SetPosition(p.PreX-1, p.PreY)
					p.SetPosition(a.X-1, a.Y)
				case RIGHT:
					TankRight(p.Tank, cell)
					p.SetPosition(p.PreX+1, p.PreY)
					p.SetPosition(a.X+1, a.Y)
				}
			case BULLET:
				b := NewBullet(a.X, a.Y, a.Direction)
				Level.AddEntity(b)
			case DELETE:
				Level.RemoveEntity(players[a.ID])
				delete(players, a.ID)
			}

			d.Ack(false)
		}
	}()
}

func SendCommand(c Command){

	body, _ := json.Marshal(c)
	//failOnError(err, "Failed to marshal a command", "Command marsheled")

	_ = ch.Publish(
		commandsExchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:     "application/json",
			Body:            body,
		})
}

