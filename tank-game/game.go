package tank_game

import (
	tl "github.com/JoelOtter/termloop"
)

var (
	fps = 	60.0
	//Level() 	*tl.BaseLevel
	g 		*GameClass
)


type GameClass struct {
	*tl.Game
	fps            int
	playersManager playersManager
}

func CreateGame(fps int) *GameClass {
	g = &GameClass{
		Game: 			tl.NewGame(),
		fps:  			fps,
		playersManager: playersManager{
			players:         make(map[string]*Player),
			playerToDelete:  make(map[string]bool),
			CurrentPlayerID: "",
		},
	}
	return g
}

func StartGame() {


	g.Screen().SetFps(fps)

	g.Screen().SetLevel(tl.NewBaseLevel(tl.Cell{Bg: tl.ColorWhite}))

	//loadmap()

	g.playersManager.getOnlinePlayers()
	g.playersManager.listenCommands()

	startMenu()

	g.Start()
}

func Game() *GameClass{
	return g
}

func Level() tl.Level {
	return g.Screen().Level()
}

func loadmap(){
	for i := 0; i < 4; i++ {
		Level().AddEntity(tl.NewRectangle(randInt(0,50), randInt(0, 50),
			randInt(5,10), randInt(5, 10), tl.ColorBlack))
	}
}

