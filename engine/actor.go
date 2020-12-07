package engine

type Direction int8

const (
	Up    Direction = 0
	Right           = 1
	Down            = 3
	Left            = 5
)

type Actor struct {
	Ascii    string
	X        int8
	Y        int8
	Passable bool
}

func (a Actor) String() string {
	return a.Ascii
}

func (a *Actor) SetCoords(x int8, y int8) {
	(*a).X = x
	(*a).Y = y
}
