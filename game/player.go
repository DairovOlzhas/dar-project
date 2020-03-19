package game

import (
	tl "github.com/JoelOtter/termloop"
	"github.com/streadway/amqp"
	"math/rand"
	"strconv"
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
	preX, preY	int
	ID 			string
	Username 	string
	Score		int
	HP 			int
}

func Player(id string) *player {
	return &player{
		Tank:     NewTankXY(0,0, tl.Cell{Bg: tl.ColorBlack}, UP),
		ID:       id,
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
		Username: genRandString(10),
		Score:    0,
		HP: 	  100,
	}

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
		p.preX, p.preY = p.Position()
		command := Command{ID: p.ID, Action: TANK, X: p.preX, Y: p.preY, Direction: p.direction}

		if event.Type == tl.EventKey {

			switch event.Key {
			case tl.KeyArrowUp:
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
				var bulletX, bulletY, bulletDirection  int
				bulletDirection = p.Tank.GetDirection()

				switch bulletDirection {
				case UP:
					bulletX = p.preX + 4
					bulletY = p.preY - 1
				case DOWN:
					bulletX = p.preX + 4
					bulletY = p.preY + 10
				case LEFT:
					bulletX = p.preX - 1
					bulletY = p.preY + 2
				case RIGHT:
					bulletX = p.preX + 10
					bulletY = p.preY + 2
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
		sX, sY := screen.Size()
		Game().Level().SetOffset(sX/2-tX-5, sY/2-tY-5)
	}

	if Menuhidden {
		x,y := p.Position()
		teX := len(p.Username)
		tl.NewText(x+3,y-2,strconv.Itoa(p.HP), tl.ColorBlue, tl.ColorWhite).Draw(screen)
		tl.NewText(x+5-teX/2,y-1, p.Username, tl.ColorBlack, tl.ColorWhite).Draw(screen)
		p.Entity.Draw(screen)
	}
}

func (p *player) Collide(collision tl.Physical) {
	if p.ID == Game().currentPlayerID {
		if _, ok := collision.(Bullet); ok {
			//Command{ID: p.ID, Action: DELETE,}.Send()
		}else if _, pl := collision.(*player); pl {
			//preX, preY := p.Position()
			switch p.direction {
			case UP:
				Command{ID: p.ID, Action: TANK, X: p.preX, Y:p.preY-1, Direction: p.direction}.Send()
			case DOWN:
				Command{ID: p.ID, Action: TANK, X: p.preX, Y:p.preY+1, Direction: p.direction}.Send()
			case RIGHT:
				Command{ID: p.ID, Action: TANK, X: p.preX-1, Y:p.preY, Direction: p.direction}.Send()
			case LEFT:
				Command{ID: p.ID, Action: TANK, X: p.preX+1, Y:p.preY, Direction: p.direction}.Send()
			}
		}
	}
}







