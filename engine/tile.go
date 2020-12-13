package engine

import (
	"fmt"
)

type Tile struct {
	Background string
	Actors     []Actor
	passable   bool
}

func (t Tile) Passable() bool {
	if len(t.Actors) == 0 {
		return t.passable
	} else if t.Actors[0].String() != "" {
		return t.Actors[0].Passable && t.passable
	}
	return t.passable
}

func (t Tile) String() string {
	if len(t.Actors) != 0 && t.Actors[0].String() != "" {
		return fmt.Sprintf("%s", t.Actors[0])
	}
	return fmt.Sprint(t.Background)
}

func (t *Tile) KeepOnTop() {
	if len((*t).Actors) > 1 && ((*t).Actors[0].String() == "" || (*t).Actors[0].String() == " ") {
		(*t).Actors = append((*t).Actors[:0], (*t).Actors[1:]...)
		if len((*t).Actors) > 1 {
			(*t).Actors = (*t).Actors[:len((*t).Actors)-1]
		}
	}
}
