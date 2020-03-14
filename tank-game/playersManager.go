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
	CurrentPlayerID = ""
)

type playersManager struct {
	players         map[string] *Player
	playerToDelete  map[string] bool
	CurrentPlayerID string
}

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
const (
	MOVE   = 0
	BULLET = 1
	DELETE = 2
	CHECK  = 3
) // command action
func SendCommand(c Command){

	body, _ := json.Marshal(c)
	_ = Channel().Publish(
		CommandsExchangeName(),
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


func (pm *playersManager) getOnlinePlayers() {


	msgs, err := Channel().Consume(
		OnlinePlayersQueueName(),
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

			if _, prs := pm.playerToDelete[playerOnline.ID]; prs{
				d.Ack(false)
				log.Println("info deleted from queue "+d.CorrelationId)
			}else{
				if _, prs := pm.players[playerOnline.ID]; !prs {
					pm.players[playerOnline.ID] = &Player{
						ID: 	  playerOnline.ID,
						Username: playerOnline.Username,
						Tank:     NewTankXY(playerOnline.X, playerOnline.Y, tl.Cell{Bg: playerOnline.Color}, playerOnline.Direction),
						Score:    playerOnline.Score,

					}
					Level().AddEntity(pm.players[playerOnline.ID])
				}
				d.Nack(false, true)
			}
		}
	}()
}
func (pm *playersManager) checkOnlinePlayers(){

	err := Channel().Publish(
		"",
		CheckOnlinePlayersQueue(),
		false,
		false,
		amqp.Publishing{
			ContentType:     "plain/text",
			CorrelationId:   pm.CurrentPlayerID,
			ReplyTo:         ReceiverQueue().Name,
		})
	failOnError(err, "Failed to publish a massage", "Massage published")


	msgs, err := Channel().Consume(
		CheckOnlinePlayersQueue(),
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
			if d.CorrelationId != pm.CurrentPlayerID {
				q, err := Channel().QueueDeclare("", false, true, false, false, nil)
				failOnError(err, "Failed to declare queue", "")

				body, _ := json.Marshal(Command{ID: d.CorrelationId, Action: CHECK,})
				for i:=0; i < 3; i++ {
					err = Channel().Publish(
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

				msgs, err := Channel().Consume(q.Name,"",true,false,false,false,nil,)
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
func (pm *playersManager) listenCommands() {
	msgs, err := Channel().Consume(
		ReceiverQueue().Name,
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
				if _, ok := pm.players[a.ID]; ok {
					p := pm.players[a.ID]
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
				Level().AddEntity(b)
			case DELETE:
				log.Println("info delete received "+d.CorrelationId)
				pm.playerToDelete[a.ID] = true
				Level().RemoveEntity(players[a.ID])
				delete(players, a.ID)
			case CHECK:
				if a.ID == pm.CurrentPlayerID {
					err = Channel().Publish("", d.ReplyTo, false,false,amqp.Publishing{})
					failOnError(err, "Failed to publish CHECK" + pm.CurrentPlayerID, "CHECK published" + pm.CurrentPlayerID)
				}
			}

			d.Ack(false)
		}
	}()
}

