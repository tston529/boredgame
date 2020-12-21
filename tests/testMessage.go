package main

import (
	"../engine"
	"fmt"
	"strings"
)

var gameMap string = ":pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pacman::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:\n:pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot::pm_dot:"


func main() {
	gameData := engine.LoadGameData("asciiyaml.yml")
	gameOverString, _ := engine.CreateEmojiMessage("game over", len("game over")+4, gameData.Message)

	y := int((23 - 5) / 2)
	x := int((23-len(strings.Split(strings.Split(gameOverString, "\n")[0], "::")))/2) + 1
	gameOverFrame := engine.OverlayEmojiMessage(gameMap, gameOverString, x, y, "::")
	fmt.Printf("%s\n", gameOverFrame)
}
