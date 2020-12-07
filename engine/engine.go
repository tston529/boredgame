package engine

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type GameMap [][]Tile
type Coord struct {
	X int8
	Y int8
}

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
			tileAscii := string(line[i])
			_, ok := wallData[tileAscii]
			gameMap[currLine][i] = Tile{Background: tileAscii, passable: !ok}
		}
		currLine += 1
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

func (gm GameMap) AddActor(a Actor) {
	if gm[a.X][a.Y].String() == " " {
		gm[a.X][a.Y].Actor = a
	}
}

func (gm GameMap) Move(a *Actor, dir Direction) (Coord, error) {
	var startX int8 = a.X
	var startY int8 = a.Y

	var destX int8 = a.X
	var destY int8 = a.Y
	switch dir {
	case Up:
		destY -= 1
		break
	case Right:
		destX += 1
		break
	case Down:
		destY += 1
		break
	case Left:
		destX -= 1
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
