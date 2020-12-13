package main

import (
	"fmt"

	"./engine"
	"./util"
	"github.com/eiannone/keyboard"

	//"strconv"
	"strings"
	"time"
	//"gopkg.in/yaml.v2"
)

type player struct {
	lives int
	score int
	engine.Actor
}

type enemy struct {
	engine.Actor
}

var paused = false

func main() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	var wallData = make(map[string]bool)
	wallData["["] = true
	wallData["="] = true
	wallData["]"] = true
	wallData["#"] = true
	wallData["+"] = true
	wallData["|"] = true

	gameMap, err := engine.LoadMap("./maps/ascii_map.txt", wallData)
	if err != nil || gameMap == nil {
		fmt.Println("Failed to load map. Exiting...")
	}
	populateBoard(&gameMap)

	playerActor := engine.Actor{Ascii: "P", X: 11, Y: 12, Passable: false}
	player1 := player{3, 0, playerActor}
	gameMap.AddActor(player1.Actor)
	enemy1 := engine.Actor{Ascii: "E", X: 5, Y: 6, Passable: true}
	gameMap.AddActor(enemy1)

	var exit = false
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			switch char {
			case 'w':
				validateMove(&gameMap, &player1, engine.Up)
				break
			case 'a':
				validateMove(&gameMap, &player1, engine.Left)
				break
			case 's':
				validateMove(&gameMap, &player1, engine.Down)
				break
			case 'd':
				validateMove(&gameMap, &player1, engine.Right)
				break
			case 'p':
				togglePaused()
				break
			}
			if key == keyboard.KeyEsc {
				exit = true
				break
			}
		}
	}()

	go func() {
		for {
			for y := 0; y < len(gameMap); y++ {
				for x := 0; x < len((gameMap)[0]); x++ {
					gameMap[y][x].KeepOnTop()
				}
			}
		}
	}()

	for !exit {
		if !paused {
			fmt.Printf("\x1b[0E\x1b7%s%s\x1b[K%s\n\x1b[2G\x1b[28A", gameMap, player1.Hud(), actorsOnTile(player1.X, player1.Y, &gameMap))
			//fmt.Println(player1.Hud())
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func actorsOnTile(x int8, y int8, gm *engine.GameMap) string {
	builder := strings.Builder{}
	l := len((*gm)[y][x].Actors)
	builder.WriteString("Actors on tile: [")
	if l > 0 {
		for x, a := range (*gm)[y][x].Actors {
			builder.WriteString(fmt.Sprintf("'%s'", a.String()))
			if x != l-1 {
				builder.WriteString(", ")
			}
		}
	}
	builder.WriteString("]")
	return builder.String()
}

func togglePaused() {
	paused = !paused
	if paused {
		fmt.Printf("\x1b[s\x1b[11B\x1b[5C              \n\x1b[5C +==========+ \n\x1b[5C |  PAUSED  | \n\x1b[5C +==========+ \n\x1b[5C              ")
	} else {
		fmt.Printf("\x1b[u\x1b[2J")
	}
}

func validateMove(gm *engine.GameMap, p *player, dir engine.Direction) {

	if paused {
		return
	}

	var startX int8 = (*p).Actor.X
	var startY int8 = (*p).Actor.Y

	var destX int8 = (*p).Actor.X
	var destY int8 = (*p).Actor.Y

	newCoord, err := gm.Move(&(*p).Actor, dir)
	// If out of bounds, pac-man loop around
	if err != nil {
		switch dir {
		case engine.Left:
			destX = int8(len((*gm)[0]) - 1)
			destY = (*p).Actor.Y
			break
		case engine.Right:
			destX = 0
			destY = (*p).Actor.Y
			break
		case engine.Up:
			destX = (*p).Actor.X
			destY = int8(len(*gm) - 1)
			break
		case engine.Down:
			destX = (*p).Actor.X
			destY = 0
			break
		}
	} else {
		destX = newCoord.X
		destY = newCoord.Y
	}

	if len((*gm)[destY][destX].Actors) > 0 {
		(*p).collision(&(*gm)[destY][destX].Actors[0])
	}
	if (*gm)[destY][destX].Passable() {
		(*p).Actor.SetCoords(destX, destY)
		(*gm).AddActor((*p).Actor)
		(*gm)[startY][startX].Actors = (*gm)[startY][startX].Actors[1:]
	}
}

func enemyMove(gm engine.GameMap, e *enemy) {

}

func populateBoard(gm *engine.GameMap) {
	for y := 0; y < len(*gm); y++ {
		for x := 0; x < len((*gm)[0]); x++ {
			bg := (*gm)[y][x].Background
			if bg == "." || bg == "@" {
				(*gm)[y][x].Background = " "
				(*gm).AddActor(engine.Actor{Ascii: bg, X: int8(x), Y: int8(y), Passable: true})
			} else if bg == " " {
				(*gm).AddActor(engine.Actor{Ascii: "", X: int8(x), Y: int8(y), Passable: true})
			}
		}
	}
}

func (p *player) collision(a *engine.Actor) {
	if (*a).Ascii == "." {
		(*p).score += 10
		(*a).Ascii = ""
	} else if (*a).Ascii == "E" {
		if (*p).lives > 0 {
			(*p).lives -= 1
		}
	}
}

func (p *player) Hud() string {
	builder := strings.Builder{}
	builder.WriteString("+--------------------------------+\n| LIVES: ")
	for i := 0; i < (*p).lives; i++ {
		builder.WriteString("<3 ")
	}
	for i := (*p).lives; i < 3; i++ {
		builder.WriteString("   ")
	}
	builder.WriteString(util.FixedLengthString(13, fmt.Sprintf("| score: %d", (*p).score)))
	builder.WriteString("  |\n+--------------------------------+\n")
	builder.WriteString(fmt.Sprintf("[%s] (%d,%d) ", (*p).Ascii, (*p).Actor.X, (*p).Actor.Y))

	return builder.String()
}
