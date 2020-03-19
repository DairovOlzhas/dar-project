package game

import (
	tl "github.com/JoelOtter/termloop"
)

type Bullet struct {
	*tl.Entity
	direction int
	owner string
}

func NewBullet(x, y, d int, owner string) Bullet {
	b := Bullet{
		Entity:    tl.NewEntity(x, y, 1, 1),
		direction: d,
		owner: 		owner,
	}
	b.SetCell(0, 0, &tl.Cell{Fg: tl.ColorBlack, Bg: tl.ColorBlack})

	return b
}

func (b Bullet) Draw(screen *tl.Screen) {

	bX, bY := b.Position()
	screenX, screenY := Game().Size()

	if bX > screenX || bX < 0 || bY > screenY || bY < 0 {
		screen.RemoveEntity(b)
		screen.Level().RemoveEntity(b)

	}

	switch b.direction {

	case UP:
		b.SetPosition(bX, bY-1)
	case DOWN:
		b.SetPosition(bX, bY+1)
	case LEFT:
		b.SetPosition(bX-1, bY)
	case RIGHT:
		b.SetPosition(bX+1, bY)
	}
	b.Entity.Draw(screen)

}

func (b Bullet) Tick(event tl.Event) {
	if b.owner == Game().currentPlayerID {
		for id, p := range Game().onlinePlayers {
			if id != Game().currentPlayerID {

				px, py := p.Position()
				pw, ph := p.Size()

				bx, by := b.Position()
				bw, bh := b.Size()

				if px < bx+bw && px+pw > bx &&
					py < by+bh && py+ph > by {
					if p.HP > 5 {
						Command{ID: id, Action: ATTACKED}.Send()
						Command{ID: b.owner, Action:HITTED}.Send()
					}else{
						Command{ID: p.ID, Action: DELETE}.Send()
						Command{ID: b.owner, Action:KILL}.Send()
					}
					Game().Level().RemoveEntity(b)
				}
			}
		}
	}
}

