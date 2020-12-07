package engine

import (
	"fmt"
)

type Tile struct {
	Background string
	Actor      Actor
	passable   bool
}

func (t Tile) Passable() bool {
	if t.Actor.String() != "" {
		return t.Actor.Passable && t.passable
	}
	return t.passable
}

func (t Tile) String() string {
	if t.Actor.String() != "" {
		return fmt.Sprintf("%s", t.Actor)
	}
	return fmt.Sprint(t.Background)
}
