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
			randInt(1,4),
			),
		ID:       genRandString(32),
		Username: Username(),
		Score:    0,
		HP: 	  100,
	}
	for player.CollideWorker(player.direction) {
		player.direction = randInt(1,4)
		player.SetPosition(	randInt(0,gX-9), randInt(0, gY-9))
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
				if p.direction == UP {
					command.Y -= 1
					if p.CollideWorker(UP) {
						return
					}
				}
				command.Direction = UP
			case tl.KeyArrowDown:
				if p.direction == DOWN {
					command.Y += 1
					if p.CollideWorker(DOWN) {
						return
					}
				}
				command.Direction = DOWN
			case tl.KeyArrowRight:
				if p.direction == RIGHT {
					command.X += 2
					if p.CollideWorker(RIGHT) {
						return
					}
				}
				command.Direction = RIGHT
			case tl.KeyArrowLeft:
				if p.direction == LEFT {
					command.X -= 2
					if p.CollideWorker(LEFT) {
						return
					}
				}
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
					bulletY = p.preY + 6
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
			command.Send()
		}
	}
}

func (p *player) Draw(screen *tl.Screen) {

	if Menuhidden {
		x,y := p.Position()
		teX := len([]rune(p.Username))
		tl.NewText(x+3,y-2,strconv.Itoa(p.HP), tl.ColorBlue, tl.ColorWhite).Draw(screen)
		tl.NewText(x+4-teX/2,y-1, p.Username, tl.ColorBlack, tl.ColorWhite).Draw(screen)
		p.Entity.Draw(screen)
	}
}

func (p *player) CollideWorker(direction int) bool {
	if p.ID == Game().currentPlayerID {
		for _, c := range Game().onlinePlayers {
			if c.ID != p.ID && collided(p, c, direction) {
				return true
			}
		}
		for _, c := range Game().walls {
			if collided(p, c, direction) {
				return true
			}
		}
	}
	return false
}

func collided(p tl.Physical, c tl.Physical, direction int) bool {
	px, py := p.Position()
	cx, cy := c.Position()
	pw, ph := p.Size()
	cw, ch := c.Size()

	switch direction {
	case UP:
		py -= 1
	case DOWN:
		py += 1
	case LEFT:
		px -= 2
	case RIGHT:
		px += 2
	}

	if px < cx+cw && px+pw > cx &&
		py < cy+ch && py+ph > cy {
		return true
	}
	return false
}







