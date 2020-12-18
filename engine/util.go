package engine

import (
	"bytes"
	"io"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"gopkg.in/yaml.v2"
)

type sz struct {
	X int
	Y int
}

type TileData struct {
	ASCII string
	Data map[interface{}]interface{}
}

type ActorData struct {
	ASCII string
	Data map[interface{}]interface{}
}

type MapData struct {
	Filename string
	Data sz
}

type SceneData struct {
	Map MapData
	Tiles map[string]TileData
	Actors map[string]ActorData
}

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

func lineCounter(r io.Reader) (sz, error) {
	buf := make([]byte, 32*1024)
	count := 0
	x := 0
	lineSep := []byte{'\n'}

xy_exit:
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		x = len(strings.Fields(string(buf[:c])))

		switch {
		case err == io.EOF:
			break xy_exit
			//return count, nil

		case err != nil:
			break xy_exit
			//return count, err
		}
	}

	return sz{x, count}, nil
}
