package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dnswd/arus/util"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type HTTPHandler struct {
	service Service // Depends on interface
}

type UserResponse struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
	Balance   string
}

func toResponse(u *User) UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		Balance:   u.Balance.Text('f'),
	}
}

func NewHandler(service Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email   string
		Name    string
		Balance string
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.CreateUser(r.Context(), req.Email, req.Name, req.Balance)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err)
		return
	}

	util.JSON(w, http.StatusCreated, toResponse(user))
}

func (h *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, err := uuid.Parse(idString)
	if err != nil {
		util.Error(w, http.StatusBadRequest, err)
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.ErrorMsg(w, http.StatusNotFound, "user not found")
			return
		}
		util.Error(w, http.StatusInternalServerError, err)
		return
	}

	util.JSON(w, http.StatusOK, toResponse(user))
}

func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")

	id, err := uuid.Parse(idString)
	if err != nil {
		util.Error(w, http.StatusBadRequest, err)
	}

	err = h.service.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.ErrorMsg(w, http.StatusNotFound, "user not found")
			return
		}
		util.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
