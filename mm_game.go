package main

import (
	"fmt"
	"os"

	"./engine"
	"./util"
	//"./mm_render"
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

var width int
var height int

var gameData engine.SceneData

func main() {
	/*mm_render.Client = mm_render.StartMattermostClient("<chat server>")
	mm_render.UserLogin("<username>", "<password>")
	mm_render.SetupGracefulShutdown()
	channel, _ := mm_render.Client.CreateDirectChannel(mm_render.MyUser.Id, os.Args[1])
	mm_render.PostMessage(mm_render.MyUser.Id, channel.Id, "Gamu Starto desu")*/

	//os.Exit(0)
	gameData = engine.LoadGameData("./tests/testyaml.yml")

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	// Generate game map from loaded game data
	gameMap, err := engine.LoadMap(gameData.Map.Filename, &gameData)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	width = len(gameMap[0])
	height = len(gameMap)

	// Initialize pop-up message strings
	pausedString, err := engine.CreateMessage("PAUSED", len("PAUSED")+4)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	gameOverString, err := engine.CreateMessage("GAME OVER", len("GAME OVER")+4)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	playerActor := engine.Actor{ASCII: gameData.Actors["player"].ASCII, X: 11, Y: 12}
	player1 := player{3, 0, playerActor}
	gameMap.AddActor(player1.Actor)
	enemy1 := engine.Actor{ASCII: gameData.Actors["enemy"].ASCII, X: 9, Y: 12}
	gameMap.AddActor(enemy1)

	var exit = false
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			switch char {
			case 'w', 'W':
				validateMove(&gameMap, &player1, engine.Up)
				break
			case 'a', 'A':
				validateMove(&gameMap, &player1, engine.Left)
				break
			case 's', 'S':
				validateMove(&gameMap, &player1, engine.Down)
				break
			case 'd', 'D':
				validateMove(&gameMap, &player1, engine.Right)
				break
			case 'p', 'P':
				togglePaused()
				break
			}
			if key == keyboard.KeyEsc {
				paused = false
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
			//mm_render.SendNextFrame(fmt.Sprintf("```\n%s%s\n```", gameMap, player1.Hud()))
			time.Sleep(100 * time.Millisecond)
		} else {
			y := int((height - 5) / 2)
			x := int((width-len(strings.Split(pausedString, "\n")[0]))/2) + 1
			pausedFrame := engine.OverlayMessage(gameMap.String(), pausedString, x, y)

			fmt.Printf("\x1b[0E\x1b7%s%s\x1b[K%s\n\x1b[2G\x1b[28A", pausedFrame, player1.Hud(), actorsOnTile(player1.X, player1.Y, &gameMap))
			// mm_render.SendNextFrame(fmt.Sprintf("```\n%s%s\n```", pausedFrame, player1.Hud()))
			for paused {
				time.Sleep(100 * time.Millisecond)
			}
		}
		if player1.lives == 0 {
			y := int((height - 5) / 2)
			x := int((width-len(strings.Split(gameOverString, "\n")[0]))/2) + 1
			gameOverFrame := engine.OverlayMessage(gameMap.String(), gameOverString, x, y)
			fmt.Printf("\x1b[0E\x1b7%s%s\x1b[K%s\n\x1b[2G\x1b[28A", gameOverFrame, player1.Hud(), actorsOnTile(player1.X, player1.Y, &gameMap))
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

// togglePaused toggles the game's paused state.
func togglePaused() {
	paused = !paused
}

// validateMove moves the player within the board, handling obstacles and
// the classic "pac-man warp."
func validateMove(gm *engine.GameMap, p *player, dir engine.Direction) {

	if paused {
		return
	}

	var startX int8 = (*p).Actor.X
	var startY int8 = (*p).Actor.Y

	var destX int8 = (*p).Actor.X
	var destY int8 = (*p).Actor.Y

	// If out of bounds, pac-man loop around
	if newCoord, err := gm.InBounds(&(*p).Actor, dir); err != nil {
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
	if passable, ok := (*gm)[destY][destX].Data["passable"]; ok && passable.(bool) {
		(*p).Actor.SetCoords(destX, destY)
		gm.AddActor((*p).Actor)
		(*gm)[startY][startX].Actors = (*gm)[startY][startX].Actors[1:]
	}
}

func enemyMove(gm engine.GameMap, e *enemy) {

}

// collision affects the game state if the player collides with
// another actor.
func (p *player) collision(a *engine.Actor) {
	switch a.ASCII {
	case gameData.Actors["dot"].ASCII, gameData.Actors["puck"].ASCII:
		p.score += a.Data["score"].(int)
		a.ASCII = gameData.Actors["blank"].ASCII
		break
	case gameData.Actors["enemy"].ASCII:
		if (*p).lives > 0 {
			(*p).lives--
		}
		break
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
