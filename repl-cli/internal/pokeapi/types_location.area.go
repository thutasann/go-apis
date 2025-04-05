package pokeapi

// Location Area Response Type
type LocationAreaResp struct {
	Count    int     `json:"count"`    // Count
	Next     *string `json:"next"`     // Nest response URL
	Previous *string `json:"previous"` // Previous response URL
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"` // Results Arrray
}
