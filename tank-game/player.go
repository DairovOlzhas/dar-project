package tank_game

import (
	"encoding/json"
	tl "github.com/JoelOtter/termloop"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
)

type Player struct {
	*Tank
	ID string
	Username string
	PreX  int
	PreY  int
	Score int
}

func failOnError(err error, msg string, ok string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
	log.Printf(" [*] " + ok)
}
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
func genRandString(l int) string {
	bytes := make([]byte, l)

	for i:=0; i < l; i++ {
		bytes[i] = byte(randInt(65,90))
	}
	return string(bytes)
}

func CreatePlayer() error {
	playerOnline := PlayerOnline{
		ID:        genRandString(32),
		Username:  genRandString(10),
		Score:     0,
		X:         randInt(0,25),
		Y:         randInt(0,25),
		Direction: UP,
		Color:     tl.RgbTo256Color(randInt(0,256), randInt(0,256),randInt(0,256)),
	}

	CurrentPlayerID = playerOnline.ID


	body, err := json.Marshal(playerOnline)
	failOnError(err, "Failed to marshal onlineplayer", "onlineplayer marsheled")

	err = Channel().Publish(
		"",
		OnlinePlayersQueueName(),
		false,
		false,
		amqp.Publishing{
			ContentType:     "",
			ReplyTo:         ReceiverQueueName(),
			Body:            body,
		})
	log.Printf("Online player sended %v", playerOnline)
	failOnError(err, "Failed to pusblish onlineplayer", "onlineplayer published")

	return err
}

func (p *Player) Tick(event tl.Event) {
	if p.ID == CurrentPlayerID {
		if event.Type == tl.EventKey {

			var bulletX, bulletY, bulletDirection int
			bulletDirection = p.Tank.GetDirection()
			p.PreX, p.PreY = p.Position()

			switch event.Key {

			case tl.KeyArrowUp:
				log.Println("Up pressed")
				SendCommand(Command{ID: p.ID, Action: TANK, X: p.PreX, Y: p.PreY, Direction: UP})
			case tl.KeyArrowDown:
				log.Println("Down pressed")
				SendCommand(Command{ID: p.ID, Action: TANK, X: p.PreX, Y: p.PreY, Direction: DOWN})
			case tl.KeyArrowRight:
				log.Println("Right pressed")
				SendCommand(Command{ID: p.ID, Action: TANK, X: p.PreX, Y: p.PreY, Direction: RIGHT})
			case tl.KeyArrowLeft:
				log.Println("Left pressed")
				SendCommand(Command{ID: p.ID, Action: TANK, X: p.PreX, Y: p.PreY, Direction: LEFT})

			case tl.KeySpace:
				switch bulletDirection {

				case UP:
					bulletX = p.PreX + 4
					bulletY = p.PreY
				case DOWN:
					bulletX = p.PreX + 4
					bulletY = p.PreY + 9
				case LEFT:
					bulletX = p.PreX
					bulletY = p.PreY + 4
				case RIGHT:
					bulletX = p.PreX + 9
					bulletY = p.PreY + 4
				}
				SendCommand(Command{
					ID:        p.ID,
					Action:    BULLET,
					X:         bulletX,
					Y:         bulletY,
					Direction: p.direction,
				})
			}

		}
	}

}

func (p *Player) Draw(screen *tl.Screen) {

	tX, tY := p.Position()
	sX, sY := screen.Size()

	if tX < 0 {
		p.SetPosition(tX+1, tY)
	}
	if tX > sX-9 {
		p.SetPosition(tX-1, tY)
	}
	if tY < 0 {
		p.SetPosition(tX, tY+1)
	}
	if tY > sY-9 {
		p.SetPosition(tX, tY-1)
	}
	p.Entity.Draw(screen)
}


//TODO
//func (p *Player) Collide(collision tl.Physical) {
//
//	if _, ok := collision.(tank.Bullet); ok {
//
//		// remove from screen
//		Level.RemoveEntity(p)
//
//	} else if _, ok := collision.(tank.Tank); ok {
//		p.SetPosition(p.preX, p.preY)
//	}
//
//}

