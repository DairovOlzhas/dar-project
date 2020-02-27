package tank_game

import tl "github.com/JoelOtter/termloop"

type Player struct{
	username string
	ID string
	tank Tank
	score int
}
type Tank struct {
	text tl.Text
}