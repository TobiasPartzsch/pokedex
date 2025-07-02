package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tobiaspartzsch/pokedex/internal/pokeapi"
)

type cliCommand struct {
	description string
	callback    func(*Config) error
}

type Config struct {
	PokeAPIConfig pokeapi.Config
	Commands      map[string]cliCommand
}

func main() {
	cfg := Config{
		PokeAPIConfig: pokeapi.Config{
			Next:     "https://pokeapi.co/api/v2/location-area", // Initialize with the base URL
			Previous: "",                                        // No previous page initially
		},
		Commands: map[string]cliCommand{},
	}
	cfg.Commands = map[string]cliCommand{
		"help": {
			description: "Displays a help message",
			callback:    commandPrintHelp,
		},
		"exit": {
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			description: "Display the next 20 area locations",
			callback:    commandMap,
		},
		"mapb": {
			description: "Display the previous 20 area locations",
			callback:    commandMapb,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scan := scanner.Scan()
		if !scan {
			log.Fatal("scanner finished")
		}
		text := scanner.Text()
		cleanInput := cleanInput(text)
		if len(cleanInput) == 0 {
			continue
		}
		command, exists := cfg.Commands[cleanInput[0]]
		if !exists {
			fmt.Printf("Unknown command: %s\n", cleanInput[0])
			continue
		}
		command.callback(&cfg)
	}
}

func commandExit(cfg *Config) error {
	msg := "Closing the Pokedex... Goodbye!"
	fmt.Println(msg)
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config) error {
	next := cfg.PokeAPIConfig.Next
	if next == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	return fetchAndPrintLocationAreas(
		cfg,
		next,
	)
}

func commandMapb(cfg *Config) error {
	previous := cfg.PokeAPIConfig.Previous
	if previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	return fetchAndPrintLocationAreas(
		cfg,
		previous,
	)
}

func commandPrintHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println("")

	for command, definition := range cfg.Commands {
		fmt.Println(command + ": " + definition.description)
	}
	return nil
}

func fetchAndPrintLocationAreas(cfg *Config, url string) error {
	locationsData, err := pokeapi.GetLocationAreas(url)
	if err != nil {
		return fmt.Errorf("failed to get location areas: %w", err)
	}

	cfg.PokeAPIConfig.Next = locationsData.Next
	cfg.PokeAPIConfig.Previous = locationsData.Previous

	for _, result := range locationsData.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
