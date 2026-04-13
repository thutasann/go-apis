package types

// Login Payload
type Login struct {
	ClientID int    `json:"clientID"`
	Username string `json:"username"`
}

// Websocket Message Payload
type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

// Player's Position
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Player State
type PlayerState struct {
	Health   int      `json:"health"`
	Position Position `json:"position"`
}
