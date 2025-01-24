package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	pokeapi "github.com/nicholasss/pokedexcli/internal/pokeapi"
	pokecache "github.com/nicholasss/pokedexcli/internal/pokecache"
	"github.com/nicholasss/pokedexcli/internal/pokedex"
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
	cache   *pokecache.Cache
	mapNURL string
	mapPURL string
	pokedex *pokedex.Pokedex
}

// =====================
// Initializing Commands
// =====================
var validCommands map[string]cliCommand

func init() {

	validCommands = map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "Attempts to catch a given Pokemon",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "Lists Pokemon that live in a given area",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"inspect": {
			name:        "inspect",
			description: "Provides info for a given Pokemon",
			callback:    commandInspect,
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
		"pokedex": {
			name:        "pokedex",
			description: "Lists all caught Pokemon",
			callback:    commandPokedex,
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
func commandCatch(cfg *config, name string) error {
	if name == "" {
		fmt.Println("Please provide the name of a Pokemon to catch.")
		return nil
	}

	URL := pokeapi.PokemonInfoURL + name + "/"

	data, err := requestThroughCache(URL, cfg)
	if err != nil {
		return fmt.Errorf("unable to request through cache: %w", err)
	}

	cfg.cache.Add(URL, data)

	pokemon, err := pokeapi.UnmarshalPokemonInfo(data)
	if err != nil {
		return fmt.Errorf("unable to unmarshal pokemon info: %w", err)
	}

	caught := pokedex.AttemptCatch(pokemon)
	if !caught {
		return nil
	}

	// Add to pokedex if caught
	cfg.pokedex.Add(name, pokemon)

	return nil
}

func commandExit(cfg *config, optional string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandExplore(cfg *config, name string) error {
	if name == "" {
		fmt.Println("Please provide the name of an area to explore.")
		return nil
	}

	// creating URL
	URL := pokeapi.LocationAreaInfoURL + name + "/"

	// requesting through cache
	data, err := requestThroughCache(URL, cfg)
	if err != nil {
		return fmt.Errorf("unable to request through cache: %w", err)
	}

	// TODO: check to ensure that cache doesnt do more work if replacing with same data
	// adding back to cache
	cfg.cache.Add(URL, data)

	locationInfo, err := pokeapi.UnmarshalLocationInfo(data)
	if err != nil {
		return err
	}

	// print out list of pokemon here
	for _, pokemon := range locationInfo.PokemonList {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func commandHelp(cfg *config, optional string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, ci := range validCommands {
		fmt.Printf("%s: %s\n", ci.name, ci.description)
	}

	return nil
}

func commandInspect(cfg *config, name string) error {
	if name == "" {
		fmt.Println("Please provide the name of a Pokemon you have caught.")
		return nil
	}

	pokemon, wasCaught := cfg.pokedex.Get(name)
	if !wasCaught {
		return nil
	}

	statList := pokemon.StatList
	typeList := pokemon.TypeList

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Printf("Stats:\n")
	for _, stat := range statList {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Printf("Types:\n")
	for _, pType := range typeList {
		fmt.Printf("  -%s\n", pType.PType.Name)
	}

	return nil
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

	locationList, err := pokeapi.UnmarshalLocationList(body)
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

	locationList, err := pokeapi.UnmarshalLocationList(body)
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

func commandPokedex(cfg *config, optional string) error {
	pokemonList, namesInList := cfg.pokedex.GetAll()
	if !namesInList {
		fmt.Println("You have not caught any Pokemon yet!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, name := range pokemonList {
		fmt.Printf("  -%s\n", name)
	}

	return nil
}

// =============
// Main Function
// =============
func main() {

	const interval = (10 * time.Minute)

	// local variables struct
	cfg := &config{
		cache:   pokecache.NewCache(interval),
		mapNURL: "null",
		mapPURL: "null",
		pokedex: pokedex.NewPokedex(),
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

		// adds a line between last command and the next prompt
		fmt.Printf("\n")
	}
}
