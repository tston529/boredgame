package mattermost-game-engine

import (
	"fmt"
)

// Tile is a single cell in the GameMap.
// Each Tile has a background which will be rendered
// provided no Actors exist in the Tile.
type Tile struct {
	Background string
	Actors     []Actor
	Data       map[interface{}]interface{}
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
	if len(t.Actors) == 1 && (t.Actors[0].String() == "" || t.Actors[0].String() == " ") {
		t.Actors = append(t.Actors[:0])
		return
	}
	if len(t.Actors) > 1 && (t.Actors[0].String() == "" || t.Actors[0].String() == " ") {
		t.Actors = append(t.Actors[:0], t.Actors[1:]...)
		if len(t.Actors) > 1 {
			t.Actors = t.Actors[:len(t.Actors)-1]
		}
	}
}
