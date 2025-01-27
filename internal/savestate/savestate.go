package savestate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	pokedex "github.com/nicholasss/pokedexcli/internal/pokedex"
)

type SaveFile struct {
	SaveTime    time.Time `json:"save_time"`
	PokedexList []string  `json:"pokedex_list"`
}

var mux sync.Mutex

func init() {
	mux = sync.Mutex{}
}

func SavePokedex(path string, pokedex *pokedex.Pokedex) error {
	pokedexList, ok := pokedex.GetAll()
	if !ok {
		return errors.New("unable to get list from pokedex")
	}

	// only lock if there is anything to actually save
	mux.Lock()
	defer mux.Unlock()

	newSave := SaveFile{
		SaveTime:    time.Now(),
		PokedexList: pokedexList,
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(newSave)
	if err != nil {
		return err
	}

	dataReader := bytes.NewReader(data)

	_, err = io.Copy(file, dataReader)
	if err != nil {
		return err
	}

	fmt.Println("Saved Pokedex successfully.")
	return nil
}

func LoadPokedex(path string, pokedex *pokedex.Pokedex) error {
	mux.Lock()
	defer mux.Unlock()

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		fmt.Println("There is no save file to load.")
		return err
	} else if err != nil {
		return err
	}
	defer file.Close()

	var oldSave SaveFile
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&oldSave); err != nil {
		return err
	}

	ok := pokedex.AddList(oldSave.PokedexList)
	if !ok {
		return errors.New("unable to add list to pokedex:")
	}

	time := oldSave.SaveTime
	fmt.Println("Loaded save file from:", time.Local().String())
	return nil
}
