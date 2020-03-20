package game

import (
	"encoding/json"
	tl "github.com/JoelOtter/termloop"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var (
	game 	*gameClass
	fps 			= 	90.0 // should be float
	gameWidth 		=	100
	gameHeight 		=	100
	backgroundColor = 	tl.ColorWhite
	playersToDelete = 	make(map[string]bool)
)

type gameClass struct {
	*tl.Game

	width, height 	int
	level 			*tl.BaseLevel

	currentPlayerID string

	walls 			[]*tl.Rectangle
	onlinePlayers 	map[string]*player
}

func Game() *gameClass {
	if game == nil {
		game = &gameClass{
			Game:   tl.NewGame(),
			width:  gameWidth,
			height: gameHeight,
			level: tl.NewBaseLevel(tl.Cell{Bg: backgroundColor}),
			onlinePlayers: make(map[string]*player),
		}
		game.walls = NewWalls()
		var s string
		log.Println(s)
	}
	return game
}

func (g *gameClass) Size() (int,int){
	return g.width, g.height
}

func (g *gameClass) Level() *tl.BaseLevel {
	return g.level
}

func (g *gameClass) Start() {

	g.Screen().SetFps(fps)

	for _, w := range g.walls {
		g.level.AddEntity(w)
	}
	g.Screen().SetLevel(g.level)
	//TODO map loading


	//TODO load online players
	g.getOnlinePlayers()

	//TODO check online users
	g.checkOnlinePlayers()

	//TODO listen commands
	g.listenCommands()

	//TODO start menu
	StartMenu()


	g.Game.Start()
}


func (g *gameClass) getOnlinePlayers() {

	msgs := Consumer(onlinePlayersQueue)

	go func() {
		for d := range msgs {
			id := string(d.Body)

			if _, prs := playersToDelete[id]; prs{
				d.Ack(false)
			}else{
				if _, prs := g.onlinePlayers[id]; !prs {
					g.onlinePlayers[id] = Player(id)
					game.Screen().Level().AddEntity(g.onlinePlayers[id])
				}

				d.Nack(false, true)
			}
		}
	}()
}

func (g *gameClass) checkOnlinePlayers() {

	msgs := Consumer(CheckOnlinePlayersQueue())

	go func(){
		for d := range msgs {
			if d.CorrelationId != g.currentPlayerID {
				q := QueueDeclare("", true)

				for i:=0; i < 3; i++ {
					Command{ID: d.CorrelationId, Action: CHECK, ReplyTo: q.Name}.Send()
				}

				msgs := Consumer(q.Name)

				cnt := 0
				go func() {
					for range msgs {
						cnt++
						log.Println("warning "+d.CorrelationId)
					}
				}()

				time.Sleep(100*time.Millisecond)

				if cnt == 0 {
					Command{ID:d.CorrelationId, Action:DELETE,}.Send()
					d.Ack(false)
				} else {
					d.Nack(false, true)
				}

			}
		}
	}()

}

func (g *gameClass) listenCommands() {

	msgs := Consumer(ReceiverQueue())

	go func() {

		for d := range msgs {
			a := Command{}
			err := json.Unmarshal(d.Body, &a)
			failOnError(err, "Failed to unmarshal command", "")

			switch a.Action {
			case KILL:
				if _, ok := g.onlinePlayers[a.ID]; ok {
					g.onlinePlayers[a.ID].HP += 5
					g.onlinePlayers[a.ID].Score += 1
				}
			case HITTED:
				if _, ok := g.onlinePlayers[a.ID]; ok {
					g.onlinePlayers[a.ID].HP += 5
				}
			case ATTACKED:
				if _, ok := g.onlinePlayers[a.ID]; ok {
					g.onlinePlayers[a.ID].HP -= 5
				}
			case TANK:
				if _, ok := g.onlinePlayers[a.ID]; ok {
					p := g.onlinePlayers[a.ID]

					p.Username = a.Username
					p.Score = a.Score
					p.HP = a.HP
					p.color = a.Color
					p.preX, p.preY = p.Position()
					p.SetPosition(a.X, a.Y)
					switch a.Direction {
					case UP:
						TankUp(p.Tank)
					case DOWN:
						TankDown(p.Tank)
					case LEFT:
						TankLeft(p.Tank)
					case RIGHT:
						TankRight(p.Tank)
					}
					if p.ID == Game().currentPlayerID && Menuhidden {
						sX, sY := Game().Screen().Size()
						tX, tY := p.Position()
						Game().Level().SetOffset(sX/2-tX-5, sY/2-tY-5)
					}

				}
			case BULLET:
				b := NewBullet(a.X, a.Y, a.Direction, a.ID)
				g.level.AddEntity(b)
			case DELETE:
				if _, ok := g.onlinePlayers[a.ID]; ok {
					log.Println("info delete received "+d.CorrelationId)
					playersToDelete[a.ID] = true
					g.level.RemoveEntity(g.onlinePlayers[a.ID])
					delete(g.onlinePlayers, a.ID)
				}
			case CHECK:
				if a.ID == g.currentPlayerID {
					Publisher("", d.ReplyTo, amqp.Publishing{})
				}
			}

			d.Ack(false)
		}
	}()

}


const (
	TANK = 0
	BULLET = 1
	DELETE = 2
	CHECK = 3
	HITTED = 4
	ATTACKED = 5
	KILL = 6
) // command action

type Command struct {
	ID string
	Action int
	ReplyTo string
	X, Y, Direction, Score, HP int
	Username string
	Color tl.Attr
}

func (c Command) Send(){


	if c.Action == TANK {
		p := Game().onlinePlayers[Game().currentPlayerID]
		c.HP = p.HP
		c.Score = p.Score
		c.Username = Username()
		c.Color = p.color
	}

	body, _ := json.Marshal(c)

	msg := amqp.Publishing{
		ContentType:     "application/json",
		Body:            body,
	}

	if c.Action == CHECK {
		msg.ReplyTo = c.ReplyTo
	}

	Publisher(CommandsExchange(),"", msg)

}
