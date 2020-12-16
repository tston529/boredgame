package main
import (
	"os"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Tile struct {
	ASCII string
	Data map[interface{}]interface{}
}

type Actor struct {
	ASCII string
	Data map[interface{}]interface{}
}

type T struct {
	Tiles map[string]Tile
	Actors map[string]Actor
}

func main() {
	file, err := ioutil.ReadFile("testyaml.yml")
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}

	t := T{}
	err = yaml.Unmarshal(file, &t)
    if err != nil {
        fmt.Printf("error: %v", err)
    }

	for key, _ := range t.Tiles {
		fmt.Printf("%v : %v\n", key, t.Tiles[key].ASCII)
	}
	for key, _ := range t.Actors {
		fmt.Printf("%v : %v\n", key, t.Actors[key].ASCII)
	}

}