package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// List Location Areas Service
func (c *Client) GetPokemon(pokemonName *string) (Pokemon, error) {
	endPoint := "/pokemon/" + *pokemonName
	fullURL := baseURL + endPoint
	fmt.Println(fullURL)

	// check the cache
	data, ok := c.cache.Get(fullURL)
	if ok {
		fmt.Println("cache hit")
		pokemon := Pokemon{}
		err := json.Unmarshal(data, &pokemon)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}
	fmt.Println("cache miss!")

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do((req))
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		return Pokemon{}, fmt.Errorf("bad status code :%v", resp.StatusCode)
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	pokemon := Pokemon{}
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(fullURL, data)

	return pokemon, nil
}
