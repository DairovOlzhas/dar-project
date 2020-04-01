package game

import (
	tl "github.com/JoelOtter/termloop"
)

var (
	borderColor = tl.ColorBlack
)

func NewWalls() []*tl.Rectangle {
	walls := make([]*tl.Rectangle, 0)

	gX, gY := Game().Size()

	walls = append(
		walls,
		tl.NewRectangle(-1,-1,gX+2,1, borderColor),
		tl.NewRectangle(-1,-1,1,gY+2, borderColor),
		tl.NewRectangle(gX,-1,1,gY+2, borderColor),
		tl.NewRectangle(-1,gY,gX+2,1, borderColor),
		)

	return walls
}


