package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Next     string
	Previous string
}

func GetLocationAreas(url string) (LocationAreas, error) {
	las := LocationAreas{}
	err := fetchAndUnmarshall(url, &las)
	if err != nil {
		return LocationAreas{}, err
	}
	return las, nil
}

func GetLocationAreaDetails(url string) (LocationArea, error) {
	la := LocationArea{}
	err := fetchAndUnmarshall(url, &la)
	if err != nil {
		return LocationArea{}, err
	}
	return la, nil
}

func GetPokemon(url string) (Pokemon, error) {
	p := Pokemon{}
	err := fetchAndUnmarshall(url, &p)
	if err != nil {
		return Pokemon{}, err
	}
	return p, nil
}

func fetchAndUnmarshall(url string, target any) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while trying to get url %s: %w", url, err)
	}
	// Always close the response body when you're done with it!
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error while trying to read response body %v: %w", res.Body, err)
	}

	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("error while trying to unmarshal body into target: %w", err)
	}

	return nil
}
