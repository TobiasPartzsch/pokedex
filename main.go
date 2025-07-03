package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	Pokedex       map[string]pokeapi.Pokemon
}

func main() {
	cfg := Config{
		PokeAPIConfig: pokeapi.Config{
			Next:     "https://pokeapi.co/api/v2/location-area", // Initialize with the base URL
			Previous: "",                                        // No previous page initially
		},
		Commands:  map[string]cliCommand{},
		PokeCache: pokecache.NewCache(10),
		Pokedex:   make(map[string]pokeapi.Pokemon),
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
		"catch": {
			description: "Throw a Pokeball at a Pokemon to catch it",
			callback:    commandCatch,
		},
		"inspect": {
			description: "Inspect a Pokemon in your Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			description: "Lists all the caught Pokemon in your Pokedex",
			callback:    commandPokedex,
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
		err := command.callback(&cfg, cleanInput[1:])
		if err != nil {
			fmt.Printf("Error executing %s, %v\n", cleanInput[0], err.Error())
		}
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

	fetchFn := func(u string) (any, error) {
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

func commandCatch(cfg *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("catch command requires a pokemon name")
	}
	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	if _, caught := cfg.Pokedex[pokemonName]; caught {
		fmt.Printf("%s is already in your Pokedex!\n", pokemonName)
		return nil
	}

	urlPokemon := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	fetchFnPokemon := func(u string) (any, error) {
		return pokeapi.GetPokemon(u)
	}

	var pokemon pokeapi.Pokemon
	err := fetchAndCacheData(cfg, urlPokemon, fetchFnPokemon, &pokemon)
	if err != nil {
		fmt.Printf("commandCatch: err %v\n", err)
		return fmt.Errorf("failed to get pokemon information: %w", err)
	}
	fmt.Printf("%s has %d base experience\n", pokemonName, pokemon.BaseExperience)

	urlSpecies := pokemon.Species.URL
	fetchFnSpecies := func(u string) (any, error) {
		return pokeapi.GetPokemonSpecies(u)
	}

	var species pokeapi.PokemonSpecies
	err = fetchAndCacheData(cfg, urlSpecies, fetchFnSpecies, &species)
	if err != nil {
		return fmt.Errorf("failed to get Pokemon species information: %w", err)
	}

	captureRate := species.CaptureRate
	fmt.Printf("%s has a capture rate of %d (out of 255).\n", pokemonName, captureRate)

	if rand.Intn(256) < captureRate {
		cfg.Pokedex[pokemonName] = pokemon
		fmt.Printf("%s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandPokedex(cfg *Config, args []string) error {
	fmt.Println("Your Pokedex:")
	for pokemonName := range cfg.Pokedex {
		fmt.Println(" - " + pokemonName)
	}
	return nil
}

func commandInspect(cfg *Config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("explore command requires a Pokemon name")
	}
	pokemonName := args[0]
	pokemon, caught := cfg.Pokedex[pokemonName]
	if !caught {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Inspecting %s...\n", pokemonName)
	pokemon.PrintDetails()

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
