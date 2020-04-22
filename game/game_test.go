package game

import (
	tl "github.com/JoelOtter/termloop"
	"testing"
)

func TestGameClass_Size(t *testing.T) {
	gotw, goth := Game().Size()
	if goth != gameHeight || gotw != gameWidth {
		t.Errorf("Game().Size() = %d, %d; want %d, %d", gotw, goth, gameWidth, gameHeight)
	}
}


func TestTankUpCanvas(t *testing.T) {
	canvas := TankUpCanvas(tl.Cell{Bg: tl.ColorBlack})
	for _, lines := range canvas{
		for _, cell := range lines{
			if cell.Bg != tl.ColorBlack{
				t.Errorf("Tank Up Canvas not right")
			}
		}
	}
}