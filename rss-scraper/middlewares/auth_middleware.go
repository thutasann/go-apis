package middlewares

import (
	"fmt"
	"net/http"

	config "github.com/thutasann/rssagg/internal"
	"github.com/thutasann/rssagg/internal/auth"
	"github.com/thutasann/rssagg/internal/database"
	"github.com/thutasann/rssagg/utilities"
)

// Authed Handler Function Signature
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

type Handler struct {
	Cfg *config.APIConfig
}

// Authentication Middleware
//
// Parameters:
//
// - Authed Handler that needs three parameters
func (h *Handler) AuthMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			utilities.RespondWithError(w, 403, fmt.Sprintf("auth error: %v", err))
			return
		}

		user, err := h.Cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			utilities.RespondWithError(w, 400, fmt.Sprintf("couldn't get user %v", err))
			return
		}

		handler(w, r, user)
	}
}
