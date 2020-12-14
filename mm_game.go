package main

import (
	"fmt"

	"errors"
	"os"

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
	direction engine.Direction
	engine.Actor
}

var paused = false
var pausedString string
var gameOverString string

var width int
var height int

func main() {
	//mm_render.Client = mm_render.StartMattermostClient("https://localhost:8080")

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

	width = len(gameMap[0])
	height = len(gameMap)

	x := int((width - len("PAUSED") - 6) / 2)
	y := int((height - 5) / 2)

	pausedString, err = createMessage("PAUSED", 10, x, y)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	x = int((width - len("GAME OVER") - 8) / 2)
	gameOverString, err = createMessage("GAME OVER", 13, x, y)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	playerActor := engine.Actor{ASCII: "P", X: 11, Y: 12, Passable: false}
	player1 := player{3, 0, playerActor}
	gameMap.AddActor(player1.Actor)
	enemy1 := engine.Actor{ASCII: "E", X: 5, Y: 6, Passable: true}
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
		if player1.lives == 0 {
			fmt.Printf("\x1b[u%s", gameOverString)
			os.Exit(0)
		}
	}
}

// actorsOnTile is a soon-to-be deprecated function meant only for use while
// debugging. It prints in sequence all actors located on the selected Tile.
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

// togglePaused toggles the game's paused state. If the game is now
// paused, display a "PAUSED" message on the board.
func togglePaused() {
	paused = !paused
	if paused {
		fmt.Printf("\x1b[u%s", pausedString)
	} else {
		fmt.Printf("\x1b[u\x1b[2J")
	}
}

// createMessage creates a message box of desired width, to be drawn
// at a certain x/y position. It returns the built string and any errors
// that may have been found about the rendering environment.
func createMessage(msg string, boxWidth int, x int, y int) (string, error) {
	// TODO: Get rid of ansi escape sequences.
	// Possible ideas: return []string where each elt is a new line in the message box

	if boxWidth < len(msg) {
		return "", errors.New("can't have a box smaller than the text")
	}

	if x < 0 {
		return "", errors.New("message box would be too wide")
	}

	isOdd := (boxWidth - len(msg)) % 2
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("\x1b[%dB\x1b[%dC", y, x))
	for i := 0; i < boxWidth+4; i++ {
		builder.WriteString(" ")
	}
	builder.WriteString(fmt.Sprintf("\n\x1b[%dC +", x))
	for i := 0; i < boxWidth; i++ {
		builder.WriteString("=")
	}
	builder.WriteString(fmt.Sprintf("+ \n\x1b[%dC |", x))
	padding := int((boxWidth - len(msg)) / 2)
	for i := 0; i < padding; i++ {
		builder.WriteString(" ")
	}

	builder.WriteString(msg)

	for i := 0; i < padding+isOdd; i++ {
		builder.WriteString(" ")
	}
	builder.WriteString(fmt.Sprintf("| \n\x1b[%dC +", x))
	for i := 0; i < boxWidth; i++ {
		builder.WriteString("=")
	}
	builder.WriteString(fmt.Sprintf("+ \n\x1b[%dC", x))
	for i := 0; i < boxWidth+4; i++ {
		builder.WriteString(" ")
	}

	return builder.String(), nil
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

// populateBoard is to be run once at the start of the game. It parses the
// background characters of the initial board state and generates actors
// accordingly.
func populateBoard(gm *engine.GameMap) {
	for y := 0; y < len(*gm); y++ {
		for x := 0; x < len((*gm)[0]); x++ {
			bg := (*gm)[y][x].Background
			if bg == "." || bg == "@" {
				(*gm)[y][x].Background = " "
				(*gm).AddActor(engine.Actor{ASCII: bg, X: int8(x), Y: int8(y), Passable: true})
			} else if bg == " " {
				(*gm).AddActor(engine.Actor{ASCII: "", X: int8(x), Y: int8(y), Passable: true})
			}
		}
	}
}

// collision affects the game state if the player collides with
// another actor.
func (p *player) collision(a *engine.Actor) {
	if (*a).ASCII == "." {
		(*p).score += 10
		(*a).ASCII = ""
	} else if (*a).ASCII == "E" {
		if (*p).lives > 0 {
			(*p).lives--
		}
	}
}

// Hud builds a string displaying the game's status, namely the player's
// score and lives. It returns the generated string.
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
	builder.WriteString(fmt.Sprintf("[%s] (%d,%d) ", (*p).ASCII, (*p).Actor.X, (*p).Actor.Y))

	return builder.String()
}
