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

func init() {
	mux = &sync.Mutex{}
}

func SavePokedex(path string, pokedex *pokedex.Pokedex) error {
	_, namesInList := pokedex.GetAll()
	if !namesInList {
		// TODO: ensure that calling function can check whether it saved or not
		// fmt.Println("There are no Pokemon to save!")
		return nil
	}

	// only lock if there is anything to actually save
	mux.Lock()
	defer mux.Unlock()

	newSave := SaveFile{
		SaveTime:    time.Now(),
		PokedexData: pokedex,
	}

	fmt.Printf("%+v\n", pokedex)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(newSave)
	if err != nil {
		return err
	}

	fmt.Println(data)

	dataReader := bytes.NewReader(data)

	_, err = io.Copy(file, dataReader)
	if err != nil {
		return err
	}

	fmt.Println("Saved Pokedex successfully.")
	return nil
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

	loadedPokedex := pokedex.NewPokedex()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loadedPokedex); err != nil {
		return &pokedex.Pokedex{}, nil
	}

	fmt.Println("Loaded Pokedex succsesfully.")
	return loadedPokedex, nil
}
