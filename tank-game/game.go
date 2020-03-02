package tank_game

import (
	tl "github.com/JoelOtter/termloop"
)

var (
	fps = 60.0
	Level *tl.BaseLevel
)


func StartGame() {

	game := tl.NewGame()
	game.Screen().SetFps(fps)
	game.Screen().EnablePixelMode()

	Level = tl.NewBaseLevel(tl.Cell{})
	game.Screen().SetLevel(Level)

	//	Load map

	//  Load players
	getOnlinePlayers()
	//  Start listening
	listenCommands()


	CreatePlayer()

	CheckOnlinePlayers()

	game.Start()
	//  Menu
	//  Run

}



