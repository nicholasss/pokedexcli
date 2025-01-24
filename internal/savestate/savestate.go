package savestate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	pokedex "github.com/nicholasss/pokedexcli/internal/pokedex"
)

type SaveFile struct {
	SaveTime    time.Time        `json:"save_time"`
	PokedexData *pokedex.Pokedex `json:"pokedex_data"`
}

var mux *sync.Mutex

func SavePokedex(path string, pokedex *pokedex.Pokedex) error {
	_, namesInList := pokedex.GetAll()
	if !namesInList {
		fmt.Println("There are no Pokemon to save!")
		return nil
	}

	// only lock if there is anything to actually save
	mux.Lock()
	defer mux.Unlock()

	newSave := SaveFile{
		SaveTime:    time.Now(),
		PokedexData: pokedex,
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(newSave, "", "\t")
	if err != nil {
		return err
	}
	dataReader := bytes.NewReader(data)

	_, err = io.Copy(file, dataReader)
	return err
}

func LoadPokedex(path string) (*pokedex.Pokedex, error) {
	mux.Lock()
	defer mux.Unlock()

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		fmt.Println("There is no save file to load.")
		return &pokedex.Pokedex{}, err
	} else if err != nil {
		return &pokedex.Pokedex{}, err
	}
	defer file.Close()

	loadedPokedex := &pokedex.Pokedex{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loadedPokedex); err != nil {
		return &pokedex.Pokedex{}, nil
	}

	return loadedPokedex, nil
}
