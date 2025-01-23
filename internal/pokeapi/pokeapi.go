package pokeapi

import (
	"fmt"
	"io"
	"net/http"
)

const LocationAreaURL = "https://pokeapi.co/api/v2/location-area"

// location area
// field names need to be public with upper case for json package
//
// https://pokeapi.co/api/v2/location-area/
type LocationList struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

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
