package engine

import (
	"fmt"
)

// Tile is a single cell in the GameMap.
// Each Tile has a background which will be rendered
// provided no Actors exist in the Tile.
type Tile struct {
	Background string
	Actors     []Actor
	passable   bool
}

// Passable checks whether the Tile is solid or otherwise.
// Returns passable state as a boolean.
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

// KeepOnTop is run to ensure the first position in the array of Actors
// is always populated, if the array contains at least one actor.
// Since the engine renders only the first Actor in the array per-Tile,
// This helps keep something on top if the first Actor in the Tile
// is removed from the array.
func (t *Tile) KeepOnTop() {
	if len((*t).Actors) > 1 && ((*t).Actors[0].String() == "" || (*t).Actors[0].String() == " ") {
		(*t).Actors = append((*t).Actors[:0], (*t).Actors[1:]...)
		if len((*t).Actors) > 1 {
			(*t).Actors = (*t).Actors[:len((*t).Actors)-1]
		}
	}
}
