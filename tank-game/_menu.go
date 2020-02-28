package tank_game

import (
	"fmt"
	tl "github.com/JoelOtter/termloop"
)

type Menu struct {
	items []*tl.Text
	index int
	selected bool

}

func (m *Menu) Tick(ev tl.Event) {
	if ev.Type == tl.EventKey {
		switch ev.Key {
		case tl.KeyArrowUp:
			if m.index > 0 {
				m.index -= 1
				m.items[m.index].SetColor(tl.ColorWhite, tl.ColorBlack)
				m.items[m.index+1].SetColor(tl.ColorBlack, tl.ColorWhite)
			}
		case tl.KeyArrowDown:
			if m.index < len(m.items)-1 {
				m.index += 1
				m.items[m.index].SetColor(tl.ColorWhite, tl.ColorBlack)
				m.items[m.index-1].SetColor(tl.ColorBlack, tl.ColorWhite)
			}
		case tl.KeyEnter:
			m.selected = true
			Game().Screen().RemoveEntity(m)
		}
	}
}
func (m *Menu) Draw(s *tl.Screen) {
	for i,_ := range m.items {
		m.items[i].Draw(s)
	}
}

func CreateMenu(arg_items []string) int {

	menu := Menu{
		items:    make([]*tl.Text, len(arg_items)),
		selected: false,
		index: 	  0,
	}
	for i, item := range arg_items {
		menu.items[i] = tl.NewText(0,i, item, tl.ColorBlack, tl.ColorWhite)
		if i == menu.index {
			menu.items[i] = tl.NewText(0,i, item, tl.ColorWhite, tl.ColorBlack)
		}
	}
	return menu.index
}

func StartMenu(){

	items := make([]string,3)
	items[0] = "Start Game"
	items[1] = "Set Name"
	items[2] = "Exit"

	fmt.Println("selected item: "+items[CreateMenu(items)])

}