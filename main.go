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
	fmt.Println("tentacool")
	fmt.Println("tentacruel")
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

func commandPrintHelp(cfg *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println("")

	for command, definition := range cfg.Commands {
		fmt.Println(command + ": " + definition.description)
	}
	return nil
}

func fetchAndPrintLocationAreas(cfg *Config, url string) error {
	var locationsDataRaw []byte
	var exists bool

	if locationsDataRaw, exists = cfg.PokeCache.Get(url); !exists {
		locationsDataResp, err := pokeapi.GetLocationAreas(url)
		if err != nil {
			return fmt.Errorf("failed to get location areas: %w", err)
		}
		locationsDataRaw, err = json.Marshal(locationsDataResp)
		if err != nil {
			return fmt.Errorf(
				"error trying to marshal response %v!: %v",
				locationsDataResp, err,
			)
		}
		success := cfg.PokeCache.Add(url, locationsDataRaw)
		if !success {
			return fmt.Errorf("couldn't add url %v to chache", url)
		}
	}
	var finalLocationAreasResp pokeapi.LocationAreas
	err := json.Unmarshal(locationsDataRaw, &finalLocationAreasResp)
	if err != nil {
		return fmt.Errorf(
			"error trying to unmarshal raw data %v!: %v",
			locationsDataRaw, err,
		)
	}

	cfg.PokeAPIConfig.Next = finalLocationAreasResp.Next
	cfg.PokeAPIConfig.Previous = finalLocationAreasResp.Previous

	for _, result := range finalLocationAreasResp.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
