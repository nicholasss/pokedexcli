package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	pokeapi "github.com/nicholasss/pokedexcli/internal/pokeapi"
	pokecache "github.com/nicholasss/pokedexcli/internal/pokecache"
)

// =====
// Types
// =====
type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	mapNURL string
	mapPURL string
	cache   *pokecache.Cache
}

// =====================
// Initializing Commands
// =====================
var validCommands map[string]cliCommand

func init() {

	validCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "Lists Pokemon that live in the area",
			callback:    commandExplore,
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

// =================
// Utility Functions
// =================
func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func commandHelp(cfg *config, optional string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, ci := range validCommands {
		fmt.Printf("%s: %s\n", ci.name, ci.description)
	}

	return nil
}

func requestThroughCache(URL string, cfg *config) ([]byte, error) {
	reqData, inCache := cfg.cache.Get(URL)
	// fmt.Println(" %%% Looking at:", URL)

	if inCache {
		// fmt.Println(" %%% USED CACHED DATA")
		return reqData, nil
	}

	// fmt.Println(" %%% REQUESTED NEW DATA")
	return pokeapi.RequestGETBody(URL)
}

// =================
// Command Functions
// =================
func commandExit(cfg *config, optional string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandExplore(cfg *config, name string) error {

	return nil
}

func unmarshalLocationList(data []byte) (pokeapi.LocationList, error) {
	var locationList pokeapi.LocationList
	if err := json.Unmarshal(data, &locationList); err != nil {
		return pokeapi.LocationList{}, fmt.Errorf("unable to unmarshal json request: %w", err)
	}

	return locationList, nil
}

func commandMap(cfg *config, optional string) error {
	// checking URL
	var URL string
	if cfg.mapNURL == "null" {
		URL = pokeapi.LocationAreaListURL
	} else {
		URL = cfg.mapNURL
	}

	// check the cache
	body, err := requestThroughCache(URL, cfg)
	if err != nil {
		return fmt.Errorf("unable to request data: %w", err)
	}

	// add to cache
	cfg.cache.Add(URL, body)

	locationList, err := unmarshalLocationList(body)
	if err != nil {
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

func commandMapB(cfg *config, optional string) error {
	// checking for URL
	var URL string
	if cfg.mapPURL == "null" {
		URL = pokeapi.LocationAreaListURL
	} else {
		URL = cfg.mapPURL
	}

	// check the cache
	body, err := requestThroughCache(URL, cfg)
	if err != nil {
		return fmt.Errorf("unable to request body: %w", err)
	}

	// add to cache
	cfg.cache.Add(URL, body)

	locationList, err := unmarshalLocationList(body)
	if err != nil {
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

	const interval = (10 * time.Minute)

	// local variables struct
	cfg := &config{
		mapNURL: "null",
		mapPURL: "null",
		cache:   pokecache.NewCache(interval),
	}

	scanner := bufio.NewScanner(os.Stdin)

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

		var optional string
		if len(args) > 1 {
			optional = args[1]
		} else {
			optional = ""
		}

		validCommand, exists := validCommands[command]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		// pass in local variables struct, and optional argument
		if err := validCommand.callback(cfg, optional); err != nil {
			fmt.Println("Error in commands:", err)
		}
	}
}
