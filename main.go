package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tobiaspartzsch/pokedex/internal/pokeapi"
	"github.com/tobiaspartzsch/pokedex/internal/pokecache"
)

type cliCommand struct {
	description string
	callback    func(*Config, []string) error
}

type Config struct {
	PokeAPIConfig pokeapi.Config
	Commands      map[string]cliCommand
	PokeCache     *pokecache.Cache
}

func main() {
	cfg := Config{
		PokeAPIConfig: pokeapi.Config{
			Next:     "https://pokeapi.co/api/v2/location-area", // Initialize with the base URL
			Previous: "",                                        // No previous page initially
		},
		Commands:  map[string]cliCommand{},
		PokeCache: pokecache.NewCache(10),
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
		"explore": {
			description: "Explore within a region (shows all Pokemon there)",
			callback:    commandExplore,
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
		command.callback(&cfg, cleanInput[1:])
	}
}

func commandExit(cfg *Config, args []string) error {
	msg := "Closing the Pokedex... Goodbye!"
	fmt.Println(msg)
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config, args []string) error {
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

func commandMapb(cfg *Config, args []string) error {
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

func commandExplore(cfg *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore command requires a location name")
	}
	locationName := args[0]
	fmt.Printf("Exploring %s...\n", locationName)

	url := "https://pokeapi.co/api/v2/location-area/" + locationName

	fetchFn := func(u string) (interface{}, error) {
		return pokeapi.GetLocationAreaDetails(u)
	}

	var locationArea pokeapi.LocationArea
	err := fetchAndCacheData(cfg, url, fetchFn, &locationArea)
	if err != nil {
		return fmt.Errorf("failed to explore location area: %w", err)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandPrintHelp(cfg *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println("")

	for command, definition := range cfg.Commands {
		fmt.Println(command + ": " + definition.description)
	}
	return nil
}

// helper functions

func fetchAndPrintLocationAreas(cfg *Config, url string) error {
	fetchFn := func(u string) (any, error) {
		return pokeapi.GetLocationAreas(u)
	}

	var locationAreas pokeapi.LocationAreas
	err := fetchAndCacheData(cfg, url, fetchFn, &locationAreas)
	if err != nil {
		return fmt.Errorf("failed to fetch and print location areas: %w", err)
	}

	cfg.PokeAPIConfig.Next = locationAreas.Next
	cfg.PokeAPIConfig.Previous = locationAreas.Previous

	for _, result := range locationAreas.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func fetchAndCacheData(
	cfg *Config,
	url string,
	fetchFunc func(string) (any, error),
	target any) error {

	var rawData []byte
	var exists bool

	if rawData, exists = cfg.PokeCache.Get(url); !exists {
		// Data not in cache, fetch it
		data, err := fetchFunc(url) // Call the specific API fetch function
		if err != nil {
			return err // Error from the fetchFunc is already descriptive
		}

		// Marshal the fetched data to store in cache
		rawData, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshalling data from %s: %w", url, err)
		}

		// Add to cache
		if !cfg.PokeCache.Add(url, rawData) {
			return fmt.Errorf("couldn't add url %v to cache", url)
		}
	}

	// Unmarshal the (potentially cached) raw data into the target struct
	err := json.Unmarshal(rawData, target)
	if err != nil {
		return fmt.Errorf("error unmarshalling raw data from %s: %w", url, err)
	}

	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
