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
	MOVE   = 0
	BULLET = 1
	DELETE = 2
	CHECK  = 3
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

			if _, prs := playerToDelete[playerOnline.ID]; prs{
				d.Ack(false)
				log.Println("info deleted from queue "+d.CorrelationId)
			}else{
				if _, prs := players[playerOnline.ID]; !prs {
					log.Printf("Online player received %v", playerOnline)
					players[playerOnline.ID] = &Player{
						ID: 	  playerOnline.ID,
						Username: playerOnline.Username,
						Tank:     NewTankXY(playerOnline.X, playerOnline.Y, tl.Cell{Bg: playerOnline.Color}, playerOnline.Direction),
						Score:    playerOnline.Score,

					}
					Level.AddEntity(players[playerOnline.ID])
				}
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
			failOnError(err, "Failed to unmarshal command", "")

			switch a.Action {
			case MOVE:
				if _, ok := players[a.ID]; ok {
					p := players[a.ID]
					cell := tl.Cell{Bg: p.color}
					p.SetPosition(a.X, a.Y)
					switch a.Direction {
					case UP:
						TankUp(p.Tank, cell)
					case DOWN:
						TankDown(p.Tank, cell)
					case LEFT:
						TankLeft(p.Tank, cell)
					case RIGHT:
						TankRight(p.Tank, cell)
					}
				}
			case BULLET:
				b := NewBullet(a.X, a.Y, a.Direction)
				Level.AddEntity(b)
			case DELETE:
				log.Println("info delete received "+d.CorrelationId)
				playerToDelete[a.ID] = true
				Level.RemoveEntity(players[a.ID])
				delete(players, a.ID)
			case CHECK:
				if a.ID == CurrentPlayerID {
					err = ch.Publish("", d.ReplyTo, false,false,amqp.Publishing{})
					failOnError(err, "Failed to publish CHECK" + CurrentPlayerID, "CHECK published" + CurrentPlayerID)
				}
			}

			d.Ack(false)
		}
	}()
}

func SendCommand(c Command){

	body, _ := json.Marshal(c)
	_ = ch.Publish(
		commandsExchange,
		"",
		false,
		false,
		amqp.Publishing{
			//failOnError(err, "Failed to marshal a command", "Command marsheled")
			ContentType:     "application/json",
			Body:            body,
		})
	if c.Action == DELETE {
		log.Println("info delete sended "+c.ID)
	}
}

func CheckOnlinePlayers(){

	err := ch.Publish(
		"",
		checkOnlinePlayersQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:     "plain/text",
			CorrelationId:   CurrentPlayerID,
			ReplyTo:         receiverQueue.Name,
		})
	failOnError(err, "Failed to publish a massage", "Massage published")


	msgs, err := ch.Consume(
		checkOnlinePlayersQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer", "Consumer registered")

	go func(){
		for d := range msgs {
			if d.CorrelationId != CurrentPlayerID {
				q, err := ch.QueueDeclare("", false, true, false, false, nil)
				failOnError(err, "Failed to declare queue", "")

				body, _ := json.Marshal(Command{ID: d.CorrelationId, Action: CHECK,})
				for i:=0; i < 3; i++ {
					err = ch.Publish(
						"commands",
						"",
						false,
						false,
						amqp.Publishing{
							ContentType:     "application/json",
							ReplyTo:         q.Name,
							Body:            body,
						})
					failOnError(err, "Failed to publish check command", "")
				}

				msgs, err := ch.Consume(q.Name,"",true,false,false,false,nil,)
				failOnError(err, "Failed to register a consumer", "")

				cnt := 0
				go func() {
					for range msgs {
						cnt++
						log.Println("warning "+d.CorrelationId)
					}
				}()

				time.Sleep(100*time.Millisecond)
				if cnt == 0 {
					SendCommand(Command{ID:d.CorrelationId, Action:DELETE,})
					d.Ack(false)
				} else {
					d.Nack(false, true)
				}

			}
		}
	}()
}

