package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thutasann/rssagg/converters"
	"github.com/thutasann/rssagg/internal/database"
	"github.com/thutasann/rssagg/utilities"
)

// Handle Create Feed
func (h *Handler) HandleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decorder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decorder.Decode(&params)
	if err != nil {
		utilities.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON : %v", err))
		return
	}

	feed, err := h.API.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		utilities.RespondWithError(w, 400, fmt.Sprintf("Couldn't create Feed: %v", err))
		return
	}

	utilities.RespondWithJSON(w, 201, converters.DatabaseFeedToFeed(feed))
}

// Handle Get Feeds
func (h *Handler) HandlerGetFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := h.API.DB.GetFeeds(r.Context())
	if err != nil {
		utilities.RespondWithError(w, 500, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}
	utilities.RespondWithJSON(w, 200, converters.DatabaseFeedsToFeeds(feeds))
}
