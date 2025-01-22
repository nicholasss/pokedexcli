package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/nicholasss/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	mapPageNum int
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

func commandMap(cfg *config) error {
	// increment page num
	cfg.mapPageNum++

	// fmt.Printf("page #%d\n", cfg.mapPageNum)

	// locally used for findListOffset()
	countPerPage := 20

	baseURL := "https://pokeapi.co/api/v2/location-area"
	offsetQuery := pokeapi.FindListOffset(countPerPage, cfg.mapPageNum)
	fullURL := baseURL + offsetQuery

	body, err := pokeapi.GetBodyFromURL(fullURL)
	if err != nil {
		return fmt.Errorf("unable to get body: %w", err)
	}

	var locationList pokeapi.LocationList
	if err := json.Unmarshal(body, &locationList); err != nil {
		return fmt.Errorf("unable to unmarshal json request: %w", err)
	}

	for _, loc := range locationList.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapB(cfg *config) error {
	// ensure page num doesnt go negative
	if cfg.mapPageNum <= 1 {
		fmt.Println("You are on the first page.")
		return nil
	} else {
		cfg.mapPageNum -= 1
	}
	// fmt.Printf("page #%d\n", cfg.mapPageNum)

	// locally used for findListOffset()
	countPerPage := 20

	baseURL := "https://pokeapi.co/api/v2/location-area"
	offsetQuery := pokeapi.FindListOffset(countPerPage, cfg.mapPageNum)
	fullURL := baseURL + offsetQuery

	body, err := pokeapi.GetBodyFromURL(fullURL)
	if err != nil {
		return fmt.Errorf("unable to get body: %w", err)
	}

	var locationList pokeapi.LocationList
	if err := json.Unmarshal(body, &locationList); err != nil {
		return fmt.Errorf("unable to unmarshal json request: %w", err)
	}

	for _, loc := range locationList.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// local variables struct
	cfg := &config{
		mapPageNum: 0,
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
		if err := validCommand.callback(cfg); err != nil {
			fmt.Println("Error in commands:", err)
		}
	}
}
