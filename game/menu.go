package game

import (
	"fmt"
	tl "github.com/JoelOtter/termloop"
)

var (
	Menuhidden   = false
	nameChanging = false
	new_username = ""
)

type menu struct {
	items []*tl.Text
	index int
}

const (
	START_OR_RESUME_GAME = 0
	CHANGENAME 	= 	1
	EXIT 		=	2
)

func (m *menu) Tick(ev tl.Event) {
	if _, prs := Game().onlinePlayers[Game().currentPlayerID]; !prs {
		m.items[0].SetText("Start game")
		Menuhidden = false
	}else{
		m.items[0].SetText("Resume")
	}
	if ev.Type == tl.EventKey {
		if !nameChanging {
			switch ev.Key {
			case tl.Key(127):
				Menuhidden = false
				m.index = 0
			case tl.KeyArrowUp:
				if m.index > 0 {
					m.index -= 1
				}
			case tl.KeyArrowDown:
				if m.index < len(m.items)-1 {
					m.index += 1
				}
			case tl.KeyEnter:
				switch m.index {
				case START_OR_RESUME_GAME:
					if _, prs := Game().onlinePlayers[Game().currentPlayerID]; !prs {
						NewPlayer()
					}
					Menuhidden = true
				case CHANGENAME:
					nameChanging = true
					new_username = ""
				case EXIT:
					Game().Stop()
				}
			}
		} else {
			switch ev.Key {
			case tl.Key(127):
				nameChanging = false
			case tl.Key(65522):
				if len(new_username) > 0 {
					new_username = new_username[0:len(new_username)-1]
				}
			case tl.KeyEnter:
				Game().onlinePlayers[Game().currentPlayerID].Username = new_username
				nameChanging = false
			default:
				if ((ev.Ch >= 'a' && ev.Ch <= 'z' )|| (ev.Ch >= 'A' && ev.Ch <= 'Z')) && len(new_username) < 10{
					new_username = new_username + string(ev.Ch)
				}
			}
		}
	}
}
func (m *menu) Draw(s *tl.Screen) {
	if !Menuhidden {
		Game().Level().SetOffset(0,0)
		sx, sy := s.Size()
		if nameChanging {

			text := tl.NewText(0,0, "Enter username(at most 10 chars):", tl.ColorBlack, tl.ColorWhite)
			ix, iy := text.Size()
			text.SetPosition(sx/2-ix/2, sy/2-iy/2)
			text.Draw(s)

			username := tl.NewText(0,0, new_username, tl.ColorBlack, tl.ColorWhite)
			ix, iy = username.Size()
			username.SetPosition(sx/2-ix/2, sy/2-iy/2+1)
			username.Draw(s)

		} else {
			for i,_ := range m.items {
				if i == m.index {
					m.items[i].SetColor(tl.ColorWhite, tl.ColorBlack)
				} else {
					m.items[i].SetColor(tl.ColorBlack, tl.ColorWhite)
				}

				ix, iy := m.items[i].Size()
				m.items[i].SetPosition(sx/2-ix/2, sy/2-iy/2+i-len(m.items))
				m.items[i].Draw(s)
			}
		}
	}
}

func CreateMenu(arg_items []string) int {

	menu := menu{
		items:    make([]*tl.Text, len(arg_items)),
		index: 	  0,
	}
	for i, item := range arg_items {
		menu.items[i] = tl.NewText(0,i, item, tl.ColorBlack, tl.ColorWhite)
		if i == menu.index {
			menu.items[i] = tl.NewText(0,i, item, tl.ColorWhite, tl.ColorBlack)
		}
	}
	Game().Level().AddEntity(&menu)
	return menu.index
}
func StartMenu(){

	items := make([]string,3)
	items[0] = "Start Game"
	items[1] = "Set Name"
	items[2] = "Exit"

	fmt.Println("selected item: "+items[CreateMenu(items)])

}

