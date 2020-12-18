package engine

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"fmt"
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
func LoadMap(filename string, sd *SceneData) (GameMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	file.Seek(0, 0)

	var xy sz

	if (*sd).Map.Data.X > 0 && (*sd).Map.Data.Y > 0 {
		xy.X = (*sd).Map.Data.X
		xy.Y = (*sd).Map.Data.Y
	} else {
		xy, err = lineCounter(file)
		if err != nil {
			fmt.Printf("error counting lines in map file")
			return nil, err
		}
	}

	gameMap := make(GameMap, xy.Y)

	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	currLine := 0
	for scanner.Scan() {
		if currLine >= int(xy.X) {
			break
		}
		line := scanner.Text()
		lineSlice := strings.Fields(line)
		gameMap[currLine] = make([]Tile, xy.Y)

		for i := 0; i < len(lineSlice); i++ {
			tileASCII := string(lineSlice[i])
			tileActors := []Actor{}
			if tileVal, ok := (*sd).Tiles[tileASCII]; ok {
				gameMap[currLine][i] = Tile{Background: tileVal.ASCII, Actors: tileActors, Data: tileVal.Data}
				continue
			} 
			if actorVal, ok := (*sd).Actors[tileASCII]; ok {
				tileActors = append(tileActors, Actor{ASCII: actorVal.ASCII, Data: actorVal.Data})
				gameMap[currLine][i] = Tile{Background: (*sd).Tiles["blank"].ASCII, Actors: tileActors, Data: (*sd).Tiles["blank"].Data}
			} 
		}
		currLine++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	//os.Exit(1)
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

// InBounds is a function that determines whether a movement option for an actor
// in a chosen direction is valid and returns a Coord object containing the
// new x, y values.
func (gm GameMap) InBounds(a *Actor, dir Direction) (Coord, error) {
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
		return Coord{startX, startY}, errors.New("out of bounds")
	}

	return Coord{destX, destY}, nil
}
