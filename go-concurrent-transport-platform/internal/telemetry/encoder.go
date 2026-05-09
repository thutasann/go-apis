package telemetry

import (
	"encoding/json"
	"io"
)

// Encode writes telemetry
func Encode(w io.Writer, t Telemetry) error {
	enc := json.NewEncoder(w)
	return enc.Encode(t)
}
