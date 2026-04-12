package types

// Login Payload
type Login struct {
	ClientID int    `json:"clientID"`
	Username string `json:"username"`
}
