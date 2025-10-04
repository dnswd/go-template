package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type HTTPHandler struct {
	service Service // Depends on interface
}

func NewHandler(service Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email string `json:"email"`
        Name  string `json:"name"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    user, err := h.service.CreateUser(r.Context(), req.Email, req.Name)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    user, err := h.service.GetUser(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    // TODO: standardized response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}