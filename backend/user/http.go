package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dnswd/arus/util"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5"
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
		util.Error(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err)
		return
	}

	util.JSON(w, http.StatusCreated, user)
}

func (h *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.ErrorMsg(w, http.StatusNotFound, "user not found")
			return
		}
		util.Error(w, http.StatusInternalServerError, err)
		return
	}

	util.JSON(w, http.StatusOK, user)
}

func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

    // TODO: Delete not deleting
	err := h.service.DeleteUser(r.Context(), id)
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
