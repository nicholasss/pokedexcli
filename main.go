package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

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
	}

}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, ci := range validCommands {
		fmt.Printf("%s: %s\n", ci.name, ci.description)
	}

	return nil
}

func main() {
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

		validCommand, exists := validCommands[command]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		validCommand.callback()
	}
}
