package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dhij/ecomm/ecomm-api/server"
)

type handler struct {
	ctx    context.Context
	server *server.Server
}

func NewHandler(server *server.Server) *handler {
	return &handler{
		ctx:    context.Background(),
		server: server,
	}
}

func (h *handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var p ProductReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
	}

	product, err := h.server.CreateProduct(h.ctx, toStoreProduct(p))
	if err != nil {
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) getProduct(w http.ResponseWriter, r *http.Request) {
}
