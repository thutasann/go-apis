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
