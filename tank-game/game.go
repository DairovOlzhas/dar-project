package tank_game

import (
	tl "github.com/JoelOtter/termloop"
)

var (
	fps = 60.0
	Level *tl.BaseLevel
	Game *tl.Game
)

type GameClass struct {
	*tl.Game
	fps int
}

func CreateGame(fps int) *GameClass {
	return &GameClass{
		Game: tl.NewGame(),
		fps:  fps,
	}
}

func (g *GameClass) StartGame() {


	g.Screen().SetFps(fps)
	//Game.Screen().EnablePixelMode()

	Level = tl.NewBaseLevel(tl.Cell{Bg: tl.ColorWhite})
	g.Screen().SetLevel(Level)

	//	Load map
	//loadmap()

	//  Load players
	getOnlinePlayers()
	//  Start listening
	listenCommands()

	//CreatePlayer()
	//
	//CheckOnlinePlayers()
	//  Menu

	StartMenu()
	g.Start()
	//  Run

}

func loadmap(){
	for i := 0; i < 4; i++ {
		Level.AddEntity(tl.NewRectangle(randInt(0,50), randInt(0, 50),
			randInt(5,10), randInt(5, 10), tl.ColorBlack))
	}
}

