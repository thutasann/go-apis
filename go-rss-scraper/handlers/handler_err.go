package handlers

import (
	"net/http"

	"github.com/thutasann/rssagg/utilities"
)

// Error Handler
func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utilities.RespondWithError(w, 500, "Something went Wrong")
}
