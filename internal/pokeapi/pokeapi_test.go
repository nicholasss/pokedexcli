package pokeapi_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	pokeapi "github.com/nicholasss/pokedexcli/internal/pokeapi"
)

func TestGetBody(t *testing.T) {
	// test cases

	googleRobotsFileURI := "./googlerobots.txt"
	grbytes, err := os.ReadFile(googleRobotsFileURI)
	if err != nil {
		fmt.Println("unable to open and read file file:", googleRobotsFileURI)
		return
	}

	cases := []struct {
		input    string
		expected []byte
	}{
		{
			input:    "https://www.google.com/robots.txt",
			expected: grbytes,
		},
	}

	for _, c := range cases {
		actual, err := pokeapi.RequestGETBody(c.input)
		if err != nil {
			fmt.Println("failure with request:", err)
		}

		fmt.Printf("actual: %s", actual)
		fmt.Printf("expected: %s", c.expected)

		if !bytes.Equal(actual, c.expected) {
			t.Errorf("actual request body different from expected request body")
			return
		}
	}
}
