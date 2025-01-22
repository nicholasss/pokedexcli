package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	mapPageNum int
}

// location area
// field names need to be public with upper case for json package
//
// https://pokeapi.co/api/v2/location-area/
type locationList struct {
	Count int `json:"count"`
	// Next     string `json:"next"`
	// Previous string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// initialized in func init()
var validCommands map[string]cliCommand

func init() {

	validCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Lists map areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists previous map areas",
			callback:    commandMapB,
		},
	}

}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(cfg *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, ci := range validCommands {
		fmt.Printf("%s: %s\n", ci.name, ci.description)
	}

	return nil
}

func findListOffset(countPerPage int, pageNum int) string {
	baseQuery := "?offset="

	var offset int
	if pageNum == 1 {
		offset = 0
	} else {
		offset = (pageNum - 1) * countPerPage
	}

	// fmt.Printf("pageCount:%d, pageNum:%d --> offset:%d\n", countPerPage, pageNum, offset)
	return fmt.Sprintf(baseQuery+"%d", offset)
}

func commandMap(cfg *config) error {
	countPerPage := 20

	baseURL := "https://pokeapi.co/api/v2/location-area"
	offsetQuery := findListOffset(countPerPage, cfg.mapPageNum)
	fullURL := baseURL + offsetQuery

	fmt.Println(fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("unable to perform GET with address '%s': %w", fullURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to ReadAll from response body: %w", err)
	}

	var locationList locationList
	if err := json.Unmarshal(body, &locationList); err != nil {
		return fmt.Errorf("unable to unmarshal json request: %w", err)
	}

	fmt.Printf("%+v\n", locationList)

	// end of func cleanup
	cfg.mapPageNum++

	return nil
}

func commandMapB(cfg *config) error {
	// goes back a page using the commandMap func
	// ensure it doesnt go negative
	if cfg.mapPageNum <= 0 {
		cfg.mapPageNum = 0
	} else {
		cfg.mapPageNum--
	}

	commandMapB(cfg)

	// undo the increment within commandMapB
	// ensure it doesnt go negative
	if cfg.mapPageNum <= 0 {
		cfg.mapPageNum = 0
	} else {
		cfg.mapPageNum--
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// local variables struct
	cfg := &config{
		mapPageNum: 1,
	}

	for {
		fmt.Print("Pokedex > ") // prompt

		if ok := scanner.Scan(); !ok {
			continue // no text provided
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error occured: %s\n", err)
		}

		args := cleanInput(scanner.Text())
		command := args[0]

		validCommand, exists := validCommands[command]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		// pass in local variables struct
		validCommand.callback(cfg)
	}
}
