package engine

// Direction is a simple type to help abstract cardinal directions as an enum.
type Direction int8

// List of cardinal Directions as a quasi-enum
const (
	Up    Direction = 0
	Right           = 1
	Down            = 3
	Left            = 5
)

// Actor is any game object that will have some kind of logic associated with it.
type Actor struct {
	ASCII    string
	X        int8
	Y        int8
	Passable bool
}

func (a Actor) String() string {
	return a.ASCII
}

// SetCoords is a helper function to change the x, y coordinates of the Actor.
// You do not generally want to call this by itself, as inherently it does nothing
// to the board state.
func (a *Actor) SetCoords(x int8, y int8) {
	(*a).X = x
	(*a).Y = y
}
