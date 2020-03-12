package tank_game

import (
	"fmt"
	tl "github.com/JoelOtter/termloop"
	"log"
)

var (
	Menuhidden = false
)

type Menu struct {
	items []*tl.Text
	index int
}

const (
	START_OR_RESUME_GAME = 0
	CHANGENAME 	= 	1
	EXIT 		=	2
)

func (m *Menu) Tick(ev tl.Event) {
	if _, prs := players[CurrentPlayerID]; !prs {
		m.items[0].SetText("Start game")
		Menuhidden = false
	}else{
		m.items[0].SetText("Resume")
	}
	if ev.Type == tl.EventKey {
		switch ev.Key {
		case tl.Key(127):
			log.Println(ev.Ch)
			Menuhidden = false
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
			switch m.index {
			case START_OR_RESUME_GAME:
				if _, prs := players[CurrentPlayerID]; !prs {
					CreatePlayer()
					CheckOnlinePlayers()
				}
				Menuhidden = true
			case EXIT:
				Game.Stop()
			}
		default:
			log.Println(ev)
		}
	}
}
func (m *Menu) Draw(s *tl.Screen) {
	if !Menuhidden {
		Level.SetOffset(0,0)
		sx, sy := s.Size()
		//dx, dy := Level.Offset()
		for i,_ := range m.items {
			ix, iy := m.items[i].Size()
			m.items[i].SetPosition(sx/2-ix/2, sy/2-iy/2+i-len(m.items))
			m.items[i].Draw(s)
		}
	}
}

func CreateMenu(arg_items []string) int {

	menu := Menu{
		items:    make([]*tl.Text, len(arg_items)),
		index: 	  0,
	}
	for i, item := range arg_items {
		menu.items[i] = tl.NewText(0,i, item, tl.ColorBlack, tl.ColorWhite)
		if i == menu.index {
			menu.items[i] = tl.NewText(0,i, item, tl.ColorWhite, tl.ColorBlack)
		}
	}
	Level.AddEntity(&menu)
	return menu.index
}

func StartMenu(){

	items := make([]string,3)
	items[0] = "Start Game"
	items[1] = "Set Name"
	items[2] = "Exit"

	fmt.Println("selected item: "+items[CreateMenu(items)])

}