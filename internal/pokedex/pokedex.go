package pokedex

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	pokeapi "github.com/nicholasss/pokedexcli/internal/pokeapi"
)

type Pokedex struct {
	entries map[string]pokeapi.PokemonInfo
	mux     sync.Mutex
}

func chanceCatch(baseXP int) bool {
	// artificial limits on the potential min and max baseXP value
	var max float32 = 800.0
	var min float32 = 20.0

	// min-max normalization formula
	var probability float32 = 1.0 - ((float32(baseXP) - min) / (max - min))

	// the higher the xp the more likely to return true
	return rand.Float32() <= probability
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		entries: make(map[string]pokeapi.PokemonInfo),
		mux:     sync.Mutex{},
	}
}

// adds pokemon to the pokedex for finding later
func (p *Pokedex) Add(name string, pokemonStruct pokeapi.PokemonInfo) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.entries[name] = pokemonStruct
	return
}

// adds list of pokemon to the pokedex, for loading from save
func (p *Pokedex) AddList(nameList []string) bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	for _, name := range nameList {
		nameURL := pokeapi.PokemonInfoURL + name
		pokemonData, err := pokeapi.RequestGETBody(nameURL)
		if err != nil {
			fmt.Println("unable to request data when loading:", err)
			return false
		}

		pokemonStruct, err := pokeapi.UnmarshalPokemonInfo(pokemonData)
		if err != nil {
			fmt.Println("unable to unmarshal data when loading:", err)
			return false
		}
		p.entries[name] = pokemonStruct
	}

	return true
}

// finds pokemon in the Pokedex struct and if found returns the info struct
func (p *Pokedex) Get(name string) (pokeapi.PokemonInfo, bool) {
	p.mux.Lock()
	defer p.mux.Unlock()

	pokemon, ok := p.entries[name]
	if !ok {
		fmt.Printf("%s is not in the Pokedex.\nYou need to catch them first!\n", name)
		return pokeapi.PokemonInfo{}, false
	}

	// fmt.Printf("%s was found in the Pokedex.\n", name)
	return pokemon, true
}

// returns a list of all Pokemon in the Pokemon entries
func (p *Pokedex) GetAll() ([]string, bool) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.entries) <= 0 {
		return []string{}, false
	}

	var names []string
	for pokemon := range p.entries {
		names = append(names, pokemon)
	}

	return names, true
}

// will perform an attempt to 'catch' a pokemon and
// add to the pokedex
func AttemptCatch(pokemon pokeapi.PokemonInfo) bool {
	baseXP := pokemon.BaseExperience
	wasCaught := chanceCatch(baseXP)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	time.Sleep(time.Second)

	if wasCaught {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		return true
	}

	fmt.Printf("%s escaped!\n", pokemon.Name)
	return false
}
