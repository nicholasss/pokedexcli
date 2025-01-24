package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// =========
// Constants
// =========
const LocationAreaURL = "https://pokeapi.co/api/v2/location-area"
const LocationAreaListURL = LocationAreaURL + "?offset=0&limit=20"
const LocationAreaInfoURL = LocationAreaURL + "/" // is this needed? who knows

const PokemonInfoURL = "https://pokeapi.co/api/v2/pokemon/"

// =====
// Types
// =====
type LocationList struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"location.url"`
	PokemonList []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonInfo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	TypeList []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

// ===================
// Unmarshal Functions
// ===================

// Unmarshals data to a PokemonInfo struct.
func UnmarshalPokemonInfo(data []byte) (PokemonInfo, error) {
	var pokemonInfo PokemonInfo
	if err := json.Unmarshal(data, &pokemonInfo); err != nil {
		return PokemonInfo{}, fmt.Errorf("unable to unmarshal json request: %w", err)
	}

	return pokemonInfo, nil
}

// Unmarshals data to a LocationInfo struct.
func UnmarshalLocationInfo(data []byte) (LocationInfo, error) {
	var locationInfo LocationInfo
	if err := json.Unmarshal(data, &locationInfo); err != nil {
		return LocationInfo{}, fmt.Errorf("unable to unmarshal json request: %w", err)
	}

	return locationInfo, nil
}

// Unmarshals data to a LocationList struct.
func UnmarshalLocationList(data []byte) (LocationList, error) {
	var locationList LocationList
	if err := json.Unmarshal(data, &locationList); err != nil {
		return LocationList{}, fmt.Errorf("unable to unmarshal json request: %w", err)
	}

	return locationList, nil
}

// =================
// Network Functions
// =================

// provides the body of from a get request
func RequestGETBody(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to perform GET with address '%s': %w", URL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to ReadAll from response body: %w", err)
	}

	return body, nil
}
