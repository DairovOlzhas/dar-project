package tank_game

import tl "github.com/JoelOtter/termloop"

type Gamer struct{
	username string
	ID string
	tank Tank
}
type Tank struct {
	text tl.Text
}