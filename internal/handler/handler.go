package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muosilva/url-shortener/internal/domain"
	"github.com/muosilva/url-shortener/internal/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(svc service.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

func (h *Handler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateURLRequest
	var ctx context.Context

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}

	err = req.Validation(req.Url)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	code, err := h.service.Create(ctx, req.Url)
	if err != nil {
		return
	}

	resp := domain.JSONResponse{
		Msg:        "Created successfully!",
		StatusCode: 201,
		Code:       code,
	}

	writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, response domain.JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	_ = json.NewEncoder(w).Encode(response)
}
