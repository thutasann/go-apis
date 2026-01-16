package telemetry

import (
	"encoding/json"
	"io"
)

// Decode reads one telemetry message from stream
func Decode(r io.Reader) (Telemetry, error) {
	dec := json.NewDecoder(r)
	var t Telemetry
	err := dec.Decode(&t)
	return t, err
}
