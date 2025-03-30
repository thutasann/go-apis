package handlers

import (
	"net/http"

	"github.com/thutasann/rssagg/utilities"
)

// Handler Readiness
func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	utilities.RespondWithJSON(w, 200, struct{}{})
}
