package game

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
)

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

type player struct {
	*Tank

	ID 			string
	Username 	*tl.Text
	Score		int
	HP 			int
}

func Player(id string) *player {
	return &player{
		Tank:     NewTankXY(0,0, tl.Cell{Bg: tl.ColorBlack}, UP),
		ID:       id,
		Username: tl.NewText(0,0,"", tl.ColorBlack, tl.ColorWhite),
		Score:    0,
	}
}

func NewPlayer() *player{
	gX, gY := Game().Size()

	player := &player{
		Tank:     NewTankXY(
			randInt(0,gX-9),
			randInt(0, gY-9),
			tl.Cell{Bg:tl.RgbTo256Color(randInt(0,255), randInt(0,255),randInt(0,255))},
			UP,
			),
		ID:       genRandString(32),
		Score:    0,
	}
	x,y := player.Position()
	player.Username = tl.NewText(x,y,genRandString(10), tl.ColorBlack, tl.ColorWhite)

	Game().onlinePlayers[player.ID] = player
	game.Screen().Level().AddEntity(player)

	Publisher("", OnlinePlayersQueue(), amqp.Publishing{Body: []byte(player.ID)})
	Game().currentPlayerID = player.ID


	Publisher("", CheckOnlinePlayersQueue(),
		amqp.Publishing{
			ContentType:   "plain/text",
			CorrelationId: Game().currentPlayerID,
			ReplyTo:       ReceiverQueue(),
		})


	return player
}

func (p *player) Tick(event tl.Event) {
	if p.ID == Game().currentPlayerID && Menuhidden {
		preX, preY := p.Position()
		command := Command{ID: p.ID, Action: TANK, X: preX, Y: preY, Direction: p.direction}
		Command{ID: p.ID, Action: NAME}.Send()

		if event.Type == tl.EventKey {

			var bulletX, bulletY, bulletDirection  int
			bulletDirection = p.Tank.GetDirection()

			switch event.Key {
			case tl.KeyArrowUp:
				log.Println("UP")
				command.Y -= 1
				command.Direction = UP
			case tl.KeyArrowDown:
				command.Y += 1
				command.Direction = DOWN
			case tl.KeyArrowRight:
				command.X += 1
				command.Direction = RIGHT
			case tl.KeyArrowLeft:
				command.X -= 1
				command.Direction = LEFT
			case tl.KeySpace:
				switch bulletDirection {
				case UP:
					bulletX = preX + 4
					bulletY = preY - 1
				case DOWN:
					bulletX = preX + 4
					bulletY = preY + 10
				case LEFT:
					bulletX = preX - 1
					bulletY = preY + 4
				case RIGHT:
					bulletX = preX + 10
					bulletY = preY + 4
				}
				command.Action = BULLET
				command.X = bulletX
				command.Y = bulletY
				command.Direction = bulletDirection
			}
		}
		command.Send()
	}
}

func (p *player) Draw(screen *tl.Screen) {

	tX, tY := p.Position()
	sX, sY := screen.Size()
	gX, gY := Game().Size()
	if tX < 0 {
		p.SetPosition(tX+1, tY)
		//p.SetPosition(sX, tY)
	}
	if tX > gX-9 {
		p.SetPosition(tX-1, tY)
		//p.SetPosition(0, tY)
	}
	if tY < 0 {
		p.SetPosition(tX, tY+1)
		//p.SetPosition(tX, sY)
	}
	if tY > gY-9 {
		p.SetPosition(tX, tY-1)
		//p.SetPosition(tX, 0)
	}
	if p.ID == Game().currentPlayerID && Menuhidden {
		Game().Level().SetOffset(sX/2-tX-5, sY/2-tY-5)
	}
	x,y := p.Position()
	teX,_ := p.Username.Size()
	p.Username.SetPosition(x+5-teX/2,y)
	if Menuhidden {
		p.Username.Draw(screen)
		p.Entity.Draw(screen)
	}
}

func (p *player) Collide(collision tl.Physical) {
	if p.ID == Game().currentPlayerID {
		if _, ok := collision.(Bullet); ok {
			Command{ID: p.ID, Action: DELETE,}.Send()
		}else if _, pl := collision.(*player); pl {
			preX, preY := p.Position()
			switch p.direction {
			case UP:
				Command{ID: p.ID, Action: TANK, X: preX, Y:preY+1, Direction: p.direction}.Send()
			case DOWN:
				Command{ID: p.ID, Action: TANK, X: preX, Y:preY-1, Direction: p.direction}.Send()
			case RIGHT:
				Command{ID: p.ID, Action: TANK, X: preX-1, Y:preY, Direction: p.direction}.Send()
			case LEFT:
				Command{ID: p.ID, Action: TANK, X: preX+1, Y:preY, Direction: p.direction}.Send()
			}
		}
		//else
		//if _, pl := collision.(*tl.Rectangle); pl {
		//	switch p.direction {
		//	case UP:
		//		SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY, Direction: p.direction})
		//	case DOWN:
		//		SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY, Direction: p.direction})
		//	case RIGHT:
		//		SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY, Direction: p.direction})
		//	case LEFT:
		//		SendCommand(Command{ID:p.ID, Action:TANK, X: p.preX, Y:p.preY, Direction: p.direction})
		//	}
		//}
	}
}







