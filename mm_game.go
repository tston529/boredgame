package main

import (
	"./engine"
	"./util"
	"fmt"
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
	enemy1 := engine.Actor{Ascii: "E", X: 6, Y: 5, Passable: true}
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
			}
			if key == keyboard.KeyEsc {
				exit = true
				break
			}
		}
	}()

	for !exit {
		fmt.Printf("\x1b[0E\x1b7%s%s\x1b[2G\x1b[28A", gameMap, player1.Hud())
		//fmt.Println(player1.Hud())
		time.Sleep(100 * time.Millisecond)
	}
}

func validateMove(gm *engine.GameMap, p *player, dir engine.Direction) {
	var startX int8 = (*p).Actor.X
	var startY int8 = (*p).Actor.Y

	newCoord, err := gm.Move(&(*p).Actor, dir)
	// If out of bounds, pac-man loop around
	if err != nil {
		var destX int8
		var destY int8
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
		(*p).collision((*gm)[destY][destX].Actor)
		(*p).Actor.SetCoords(destX, destY)
		(*gm)[destY][destX].Actor = (*p).Actor
		(*gm)[startY][startX].Actor = engine.Actor{}
	} else {
		(*p).collision((*gm)[newCoord.Y][newCoord.X].Actor)
		if (*gm)[newCoord.Y][newCoord.X].Passable() {
			(*p).Actor.SetCoords(newCoord.X, newCoord.Y)
			(*gm)[newCoord.Y][newCoord.X].Actor = (*p).Actor
			(*gm)[startY][startX].Actor = engine.Actor{}
		}
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
				(*gm)[y][x].Actor = engine.Actor{bg, int8(x), int8(y), true}
			} else if bg == " " {
				(*gm)[y][x].Actor = engine.Actor{"", int8(x), int8(y), true}
			}
		}
	}
}

func (p *player) collision(a engine.Actor) {
	if a.Ascii == "." {
		(*p).score += 10
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
	builder.WriteString(fmt.Sprintf("[%s] (%d,%d)\n ", (*p).Ascii, (*p).Actor.X, (*p).Actor.Y))

	return builder.String()
}
