package engine

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type sz struct {
	X int
	Y int
}

// TileData holds unmarshalled tile  metadata
// from the game's yaml file
type TileData struct {
	ASCII string
	Data  map[interface{}]interface{}
}

// ActorData holds unmarshalled actor metadata
// from the game's yaml file
type ActorData struct {
	ASCII string
	Data  map[interface{}]interface{}
}

// MapData holds unmarshalled game map metadata
// from the game's yaml file
type MapData struct {
	Filename string
	Data     sz
}

// SceneData holds unmarshalled data from the game's yaml file,
// including map, tile and actor metadata
type SceneData struct {
	Map    MapData
	Tiles  map[string]TileData
	Actors map[string]ActorData
}

// LoadGameData reads the file as designated by 'filename'
// and unmarshalls it into a returned SceneData object.
func LoadGameData(filename string) SceneData {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}

	data := SceneData{}
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	if _, ok := data.Tiles["blank"]; !ok {
		fmt.Printf("error: tiles list in game data file '%s' needs a 'blank' tile\n", filename)
		os.Exit(1)
	}

	return data
}

// CreateMessage creates a message box of desired width, to be drawn
// at a certain x/y position. It returns the built string and any errors
// that may have been found about the rendering environment.
func CreateMessage(msg string, boxWidth int) (string, error) {
	if boxWidth < len(msg) {
		return "", errors.New("can't have a box smaller than the text")
	}

	isOdd := (boxWidth - len(msg)) % 2
	builder := strings.Builder{}

	for i := 0; i < boxWidth+4; i++ {
		builder.WriteString(" ")
	}

	builder.WriteString("\n +")
	for i := 0; i < boxWidth; i++ {
		builder.WriteString("=")
	}
	builder.WriteString("+ \n |")
	padding := int((boxWidth - len(msg)) / 2)
	for i := 0; i < padding; i++ {
		builder.WriteString(" ")
	}

	builder.WriteString(msg)

	for i := 0; i < padding+isOdd; i++ {
		builder.WriteString(" ")
	}
	builder.WriteString("| \n +")
	for i := 0; i < boxWidth; i++ {
		builder.WriteString("=")
	}
	builder.WriteString("+ \n")
	for i := 0; i < boxWidth+4; i++ {
		builder.WriteString(" ")
	}

	return builder.String(), nil
}

// OverlayMessage inserts a message box into the current frame.
// It returns the newly-built frame as a string.
func OverlayMessage(base string, msg string, x int, y int) string {
	baseStrSlice := strings.Split(base, "\n")
	msgSlice := strings.Split(msg, "\n")
	for i := 0; i < len(msgSlice); i++ {
		b := string(baseStrSlice[i+y])
		l := int(len(msgSlice[i]))
		baseStrSlice[i+y] = b[:x-1] + string(msgSlice[i]) + b[l+x-1:]
	}
	return strings.Join(baseStrSlice, "\n")
}
