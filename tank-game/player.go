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
	ID       string
	Username string
	preY 	 int
	preX 	 int
	Score    int
}

func failOnError(err error, msg string, ok string) {
	if err != nil {
		log.Fatalf("%s: error %s", msg, err)
	}
	if len(ok) != 0{
		log.Printf(" [*] " + ok)
	}
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
		X:         randInt(15,30),
		Y:         randInt(15,30),
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
			Body:            body,
		})
	failOnError(err, "Failed to pusblish onlineplayer", "onlineplayer published")

	return err
}

func (p *Player) Tick(event tl.Event) {
	if p.ID == CurrentPlayerID && Menuhidden{
		if event.Type == tl.EventKey {

			var bulletX, bulletY, bulletDirection  int
			bulletDirection = p.Tank.GetDirection()

			p.preX, p.preY = p.Position()

			switch event.Key {
			case tl.KeyArrowUp:
				SendCommand(Command{ID: p.ID, Action: TANK, X: p.preX, Y: p.preY - 1, Direction: UP})
			case tl.KeyArrowDown:
				SendCommand(Command{ID: p.ID, Action: TANK, X: p.preX, Y: p.preY + 1, Direction: DOWN})
			case tl.KeyArrowRight:
				SendCommand(Command{ID: p.ID, Action: TANK, X: p.preX + 1, Y: p.preY, Direction: RIGHT})
			case tl.KeyArrowLeft:
				SendCommand(Command{ID: p.ID, Action: TANK, X: p.preX - 1, Y: p.preY, Direction: LEFT})
			case tl.KeySpace:
				switch bulletDirection {
				case UP:
					bulletX = p.preX + 4
					bulletY = p.preY - 1
				case DOWN:
					bulletX = p.preX + 4
					bulletY = p.preY + 10
				case LEFT:
					bulletX = p.preX - 1
					bulletY = p.preY + 4
				case RIGHT:
					bulletX = p.preX + 10
					bulletY = p.preY + 4
				}
				SendCommand(Command{
					ID:        p.ID,
					Action:    BULLET,
					X:         bulletX,
					Y:         bulletY,
					Direction: bulletDirection,
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
		//p.SetPosition(sX, tY)
	}
	if tX > sX-9 {
		p.SetPosition(tX-1, tY)
		//p.SetPosition(0, tY)
	}
	if tY < 0 {
		p.SetPosition(tX, tY+1)
		//p.SetPosition(tX, sY)
	}
	if tY > sY-9 {
		p.SetPosition(tX, tY-1)
		//p.SetPosition(tX, 0)
	}
	if p.ID == CurrentPlayerID && Menuhidden{
		Level.SetOffset(sX/2-tX-5, sY/2-tY-5)
	}
	if Menuhidden {
		p.Entity.Draw(screen)
	}
}

func (p *Player) Collide(collision tl.Physical) {
	if p.ID == CurrentPlayerID {
		if _, ok := collision.(Bullet); ok {
			SendCommand(Command{ID:p.ID, Action:DELETE,})
		}else if _, pl := collision.(*Player); pl {
			switch p.direction {
			case UP:
				SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY+1, Direction: p.direction})
			case DOWN:
				SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY-1, Direction: p.direction})
			case RIGHT:
				SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX-1, Y:p.preY, Direction: p.direction})
			case LEFT:
				SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX+1, Y:p.preY, Direction: p.direction})
			}
		}else if _, pl := collision.(*tl.Rectangle); pl {
			switch p.direction {
			case UP:
				SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY, Direction: p.direction})
			case DOWN:
				SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY, Direction: p.direction})
			case RIGHT:
				SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY, Direction: p.direction})
			case LEFT:
				SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY, Direction: p.direction})
			}
		}
	}
}

