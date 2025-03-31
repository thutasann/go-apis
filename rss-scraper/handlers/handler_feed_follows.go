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

// Handle Create Feed Follow
func (h *Handler) HandleCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utilities.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed_follow, err := h.API.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		utilities.RespondWithError(w, 500, fmt.Sprintf("Couldn't create feed_follows: %v", err))
		return
	}

	utilities.RespondWithJSON(w, 200, converters.DatabaseFeedFollowToFeedFollow(feed_follow))
}

// Handle Get Feed Follow by User ID
func (h *Handler) HandleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follow, err := h.API.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		utilities.RespondWithError(w, 500, fmt.Sprintf("cannot get feed follows: %v", err))
		return
	}
	utilities.RespondWithJSON(w, 200, converters.DatabaseFeedFollowsToFeedFollows(feed_follow))
}
