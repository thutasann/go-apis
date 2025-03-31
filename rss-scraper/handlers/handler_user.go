package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thutasann/rssagg/converters"
	config "github.com/thutasann/rssagg/internal"
	"github.com/thutasann/rssagg/internal/database"
	"github.com/thutasann/rssagg/utilities"
)

// Handler Struct
type Handler struct {
	// API Config
	API *config.APIConfig
}

// Health Handler
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	utilities.RespondWithJSON(w, 200, struct {
		Message string
		State   int32
	}{Message: "Connected", State: 0})
}

// Create User Handler
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utilities.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := h.API.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		utilities.RespondWithError(w, 500, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	utilities.RespondWithJSON(w, 200, converters.DatabaseUserToUser(user))
}
