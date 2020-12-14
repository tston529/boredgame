package engine

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// GameMap is a type representing the board - it is a 2D array of Tiles.
type GameMap [][]Tile

// Coord is a type which is an abstraction for an X and Y coordinate.
type Coord struct {
	X int8
	Y int8
}

// LoadMap reads a map from a file and returns a 2d array of Tiles which contain
// the characters as Backgrounds.
func LoadMap(filename string, wallData map[string]bool) (GameMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	file.Seek(0, 0)

	xy, cErr := lineCounter(file)
	if cErr != nil {
		return nil, cErr
	}

	gameMap := make(GameMap, xy.x, xy.y)

	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	currLine := 0
	for scanner.Scan() {
		if currLine >= int(xy.x) {
			break
		}
		line := scanner.Text()
		gameMap[currLine] = make([]Tile, xy.y)
		for i := 0; i < len(line); i++ {
			tileASCII := string(line[i])
			_, ok := wallData[tileASCII]
			tileActors := []Actor{}
			gameMap[currLine][i] = Tile{Background: tileASCII, passable: !ok, Actors: tileActors}
		}
		currLine++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return gameMap, nil
}

func (gm GameMap) String() (output string) {
	builder := strings.Builder{}
	for y := 0; y < len(gm); y++ {
		for x := 0; x < len(gm[y]); x++ {
			builder.WriteString(gm[y][x].String())
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// AddActor adds an Actor to the Tile with the same
// x, y coordinates as the passed Actor.
func (gm GameMap) AddActor(a Actor) {
	if len(gm[a.Y][a.X].Actors) != 0 {
		gm[a.Y][a.X].Actors = append(gm[a.Y][a.X].Actors, Actor{})
		copy(gm[a.Y][a.X].Actors[1:], gm[a.Y][a.X].Actors[0:])
		gm[a.Y][a.X].Actors[0] = a
	} else {
		gm[a.Y][a.X].Actors = []Actor{a}
	}
}

// Move is a function that determines whether a movement option for an actor
// in a chosen direction is valid and returns a Coord object containing the
// new x, y values.
func (gm GameMap) Move(a *Actor, dir Direction) (Coord, error) {
	var startX int8 = a.X
	var startY int8 = a.Y

	var destX int8 = a.X
	var destY int8 = a.Y
	switch dir {
	case Up:
		destY--
		break
	case Right:
		destX++
		break
	case Down:
		destY++
		break
	case Left:
		destX--
		break
	}

	if destY < 0 || destX < 0 || int(destY) > len(gm)-1 || int(destX) > len(gm[0])-1 {
		return Coord{}, errors.New("out of bounds")
	}

	if gm[destY][destX].Passable() {
		return Coord{destX, destY}, nil
	}

	return Coord{startX, startY}, nil
}
