package telemetry

import "time"

// Telemetry is the single source of truth shared by ALL services
// Newline-delimited JSON is used over raw TCP
type Telemetry struct {
	TrainID        string    `json:"train_id"`
	Station        string    `json:"station"`
	SpeedKmh       int       `json:"speed_kmh"`
	PassengerCount int       `json:"passenger_count"`
	Timestamp      time.Time `json:"timestamp"`
}
