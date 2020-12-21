package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"./engine"
	"./mmrender"
	"./util"
	"github.com/eiannone/keyboard"
)

type player struct {
	lives int
	score int
	engine.Actor
}

type enemyState int

const (
	hunt enemyState = iota
	scared
	dead
)

type enemy struct {
	direction        engine.Direction
	state            enemyState
	changeEnemyState chan enemyState
	engine.Actor
}

var paused = false

var width int
var height int

var gameData engine.SceneData

// Whether to play the game locally (cli) or via mattermost
var cli = true

// For sending frame to mattermost as preformatted text
// (wrapped in backticks -> ```\n<frame>\n```)
var preBeginWrap string
var preEndWrap string

var changeEnemyState chan enemyState

func main() {
	mmUser := flag.String("user", "", "The user to receive the DM of the game")
	mmChannel := flag.String("channel", "", "The channel to receive the game message")
	mmPreformatted := flag.Bool("pre", true, "Whether to wrap each frame in backticks to be rendered as preformatted text on Mattermost.")
	flag.Parse()

	if *mmUser != "" && *mmChannel != "" {
		fmt.Println("Can't specify both user and channel, choose one or the other.")
		os.Exit(1)
	}

	if *mmUser != "" || *mmChannel != "" {
		cli = false
	}

	if !cli {
		if *mmPreformatted {
			preBeginWrap = "```\n"
			preEndWrap = "\n```"
		} else {
			preBeginWrap = ""
			preEndWrap = ""
		}
		mmData := mmrender.LoadMattermostData("./tests/mattermost.yml")

		mmrender.StartMattermostClient(mmData.ServerURL, mmData.User, mmData.Pass)
		if *mmUser != "" {
			mmrender.GetDirectMessageChannel(*mmUser)
		} else if *mmChannel != "" {
			mmrender.FindTeam(mmData.TeamName)
			mmrender.GetChannel(*mmChannel)
		}

		mmrender.PostMessage("Gamu Starto desu")
	}

	gameData = engine.LoadGameData("./tests/asciiyaml.yml")

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

	// Create player
	playerX := (gameData.Actors["player"].Data["spawnX"]).(int)
	playerY := (gameData.Actors["player"].Data["spawnY"]).(int)
	playerActor := engine.Actor{ASCII: gameData.Actors["player"].ASCII, X: int8(playerX), Y: int8(playerY)}
	player1 := player{3, 0, playerActor}
	gameMap.AddActor(player1.Actor)

	// Create slice full of enemies, from count in yaml file
	enemyCount := (gameData.Actors["enemy"].Data["enemyAmt"]).(int)
	enemies := make([]enemy, 0, enemyCount)
	enemyX := (gameData.Actors["enemy"].Data["spawnX"]).(int)
	enemyY := (gameData.Actors["enemy"].Data["spawnY"]).(int)

	for i := 0; i < cap(enemies); i++ {
		changeEnemyState := make(chan enemyState)
		enemyActor := engine.Actor{ASCII: gameData.Actors["enemy"].ASCII, X: int8(enemyX), Y: int8(enemyY - 1 + i)}
		enemy := enemy{engine.Up, hunt, changeEnemyState, enemyActor}
		gameMap.AddActor(enemy.Actor)
		enemies = append(enemies, enemy)
	}
	rand.Seed(time.Now().Unix())

	// this channel gets data passed to it when movement occurs.
	// Collision detection only happens when data is pulled off the channel,
	// preventing more than one collision from happening at once.
	movement := make(chan engine.Coord)

	// changeEnemyState gets passed data in cases when collision with something
	// causes a non-specified amount of enemies to have their states altered.
	changeEnemyState = make(chan enemyState)

	var exit = false
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			switch char {
			case 'w', 'W':
				if validateMove(&gameMap, &(player1.Actor), engine.Up) {
					movement <- engine.Coord{X: player1.X, Y: player1.Y}
					time.Sleep(60 * time.Millisecond)
				}
				break
			case 'a', 'A':
				if validateMove(&gameMap, &(player1.Actor), engine.Left) {
					movement <- engine.Coord{X: player1.X, Y: player1.Y}
					time.Sleep(60 * time.Millisecond)
				}
				break
			case 's', 'S':
				if validateMove(&gameMap, &(player1.Actor), engine.Down) {
					movement <- engine.Coord{X: player1.X, Y: player1.Y}
					time.Sleep(60 * time.Millisecond)
				}
				break
			case 'd', 'D':
				if validateMove(&gameMap, &(player1.Actor), engine.Right) {
					movement <- engine.Coord{X: player1.X, Y: player1.Y}
					time.Sleep(60 * time.Millisecond)
				}
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

	// Ensure each tile displays the first actor in its list
	go func() {
		for {
			for y := 0; y < len(gameMap); y++ {
				for x := 0; x < len(gameMap[0]); x++ {
					gameMap[y][x].KeepOnTop()
				}
			}
		}
	}()

	// Rudimentary, random enemy movement. Quickens up their pace
	// a bit when in scared state.
	for i := 0; i < len(enemies); i++ {
		go func(i int, enemies []enemy) {
			for {
				if enemies[i].state == scared {
					go func() {
						time.Sleep(8000 * time.Millisecond)
						enemies[i].changeEnemyState <- hunt
					}()
					enemyMoveRandom(gameMap, &enemies[i])
					movement <- engine.Coord{X: enemies[i].X, Y: enemies[i].Y}
					time.Sleep(time.Duration(180+10*i) * time.Millisecond)
				} else if enemies[i].state == hunt {
					enemyMoveRandom(gameMap, &enemies[i])
					movement <- engine.Coord{X: enemies[i].X, Y: enemies[i].Y}
					time.Sleep(time.Duration(240+10*i) * time.Millisecond)
				}
			}
		}(i, enemies)

	}

	// Handle enemy state changes on a per-enemy basis
	for i := 0; i < len(enemies); i++ {
		go func(i int, enemies []enemy, gameMap *engine.GameMap) {
			for {
				newEnemyState := <-enemies[i].changeEnemyState
				enemies[i].state = newEnemyState

				if newEnemyState == scared {
					enemies[i].ASCII = gameData.Actors["enemy"].Data["scared"].(string)
				} else if newEnemyState == hunt {
					enemies[i].ASCII = gameData.Actors["enemy"].ASCII
				} else if newEnemyState == dead {
					enemies[i].ASCII = gameData.Actors["enemy"].Data["dead"].(string)
					warpActor(&enemies[i].Actor, int8(enemyX), int8(enemyY), gameMap)
				}
			}
		}(i, enemies, &gameMap)
	}

	// a more global changeEnemyState channel will disperse generic state changes to either
	// individual or batches of enemies. Needed for use in collision detection to skirt
	// around the fact that collision() takes in an Actor, not an Enemy.
	go func() {
		for {
			newEnemyState := <-changeEnemyState
			if newEnemyState == dead {
				for i := 0; i < len(enemies); i++ {
					if enemies[i].X == player1.X && enemies[i].Y == player1.Y {
						enemies[i].changeEnemyState <- newEnemyState
					}
				}
			} else {
				for i := 0; i < len(enemies); i++ {
					enemies[i].changeEnemyState <- newEnemyState
				}
			}
		}
	}()

	// Collision handling - each movement (by player or enemy) will trigger
	// data to be pushed onto the movement channel.  When movement is
	// detected this way, collision checking occurs.
	go func() {
		for {
			validMove := <-movement
			if validMove.X == player1.X && validMove.Y == player1.Y && len(gameMap[player1.Y][player1.X].Actors) > 1 {
				for i := range gameMap[player1.Y][player1.X].Actors {
					player1.collision(&(gameMap[player1.Y][player1.X].Actors[i]))
				}
			}
		}
	}()

	// Core loop; continue rendering each frame until all player's lives are lost.
	for !exit {
		if !paused {
			if cli {
				fmt.Printf("\x1b[0E\x1b7%s%s\x1b[K\x1b[2G\x1b[27A", gameMap, player1.Hud())
			} else {
				mmrender.SendNextFrame(fmt.Sprintf("%s%s%s%s", preBeginWrap, gameMap, player1.Hud(), preEndWrap))
			}
			time.Sleep(100 * time.Millisecond)
		} else {
			y := int((height - 5) / 2)
			x := int((width-len(strings.Split(pausedString, "\n")[0]))/2) + 1
			pausedFrame := engine.OverlayMessage(gameMap.String(), pausedString, x, y)

			if cli {
				fmt.Printf("\x1b[0E\x1b7%s%s\x1b[K\n\x1b[2G\x1b[27A", pausedFrame, player1.Hud())
			} else {
				mmrender.SendNextFrame(fmt.Sprintf("%s%s%s%s", preBeginWrap, pausedFrame, player1.Hud(), preEndWrap))
			}
			for paused {
				time.Sleep(100 * time.Millisecond)
			}
		}
		if player1.lives == 0 {
			y := int((height - 5) / 2)
			x := int((width-len(strings.Split(gameOverString, "\n")[0]))/2) + 1
			gameOverFrame := engine.OverlayMessage(gameMap.String(), gameOverString, x, y)

			if cli {
				fmt.Printf("\x1b[0E\x1b7%s%s\x1b[K\n", gameOverFrame, player1.Hud())
			} else {
				mmrender.SendNextFrame(fmt.Sprintf("%s%s%s%s", preBeginWrap, gameOverFrame, player1.Hud(), preEndWrap))
			}

			os.Exit(0)
		}
	}
}

// togglePaused toggles the game's paused state.
func togglePaused() {
	paused = !paused
}

func warpActor(a *engine.Actor, destX int8, destY int8, gm *engine.GameMap) bool {
	var startX int8 = a.X
	var startY int8 = a.Y
	if passable, ok := (*gm)[destY][destX].Data["passable"]; ok && passable.(bool) {
		a.SetCoords(destX, destY)
		gm.AddActor(*a)
		(*gm)[startY][startX].Actors = (*gm)[startY][startX].Actors[1:]
		return true
	}
	return false
}

// validateMove moves an actor  within the board, handling obstacles and
// the classic "pac-man warp."
func validateMove(gm *engine.GameMap, a *engine.Actor, dir engine.Direction) bool {

	if paused {
		return false
	}

	var destX int8 = a.X
	var destY int8 = a.Y

	// If out of bounds, pac-man loop around
	if newCoord, err := gm.InBounds(a, dir); err != nil {
		switch dir {
		case engine.Left:
			destX = int8(len((*gm)[0]) - 1)
			destY = a.Y
			break
		case engine.Right:
			destX = 0
			destY = a.Y
			break
		case engine.Up:
			destX = a.X
			destY = int8(len(*gm) - 1)
			break
		case engine.Down:
			destX = a.X
			destY = 0
			break
		}
	} else {
		destX = newCoord.X
		destY = newCoord.Y
	}

	return warpActor(a, destX, destY, gm)
}

// enemyMoveRandom is a helper function to validate a random movement
// for an enemy, changing directions if it hits a wall.
func enemyMoveRandom(gm engine.GameMap, e *enemy) {
	if e.X == 11 && (e.Y == 8 || e.Y == 9) {
		e.direction = engine.Up
		validateMove(&gm, &(e.Actor), e.direction)
		return
	}
	for !validateMove(&gm, &(e.Actor), e.direction) {
		x := engine.Direction(rand.Intn(4))
		if x != e.direction {
			e.direction = x
		}
	}
}

// collision affects the game state if the player collides with
// another actor.
func (p *player) collision(a *engine.Actor) {
	switch a.ASCII {
	case gameData.Actors["dot"].ASCII:
		p.score += a.Data["score"].(int)
		a.ASCII = "" // This is grounds for actor deletion in the KeepOnTop() routine
		break
	case gameData.Actors["puck"].ASCII:
		p.score += a.Data["score"].(int)
		a.ASCII = "" // This is grounds for actor deletion in the KeepOnTop() routine
		changeEnemyState <- scared
		break
	case gameData.Actors["enemy"].ASCII:
		if p.lives > 0 {
			p.lives--
		}
		break
	case gameData.Actors["enemy"].Data["scared"].(string):
		p.score += 250
		changeEnemyState <- dead
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
