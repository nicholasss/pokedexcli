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
	mapNURL string
	mapPURL string
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
	// checking url
	var url string
	if cfg.mapNURL == "null" {
		url = pokeapi.LocationAreaURL
	} else {
		url = cfg.mapNURL
	}

	body, err := pokeapi.RequestGETBody(url)
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

	if locationList.Next != nil {
		cfg.mapNURL = *locationList.Next
	}
	if locationList.Previous != nil {
		cfg.mapPURL = *locationList.Previous
	}

	return nil
}

func commandMapB(cfg *config) error {
	// checking for url
	var url string
	if cfg.mapPURL == "null" {
		url = pokeapi.LocationAreaURL
	} else {
		url = cfg.mapPURL
	}

	body, err := pokeapi.RequestGETBody(url)
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

	if locationList.Next != nil {
		cfg.mapNURL = *locationList.Next
	}
	if locationList.Previous != nil {
		cfg.mapPURL = *locationList.Previous
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// local variables struct
	cfg := &config{
		mapNURL: "null",
		mapPURL: "null",
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
