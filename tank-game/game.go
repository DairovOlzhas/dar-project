package tank_game

import (
	tl "github.com/JoelOtter/termloop"
)

var (
	fps = 120.0
	Level *tl.BaseLevel
)


func StartGame() {

	game := tl.NewGame()
	game.Screen().SetFps(fps)
	game.Screen().EnablePixelMode()
	game.SetDebugOn(true)
	Level = tl.NewBaseLevel(tl.Cell{})
	game.Screen().SetLevel(Level)

	//	Load map

	//  Load players
	getOnlinePlayers()
	//  Start listening
	listenCommands()


	CreatePlayer()

	game.Start()
	//  Menu
	//  Run

}



