package handler

import (
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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, domain.JSONResponse{
			Msg:        "Invalid JSON body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err = req.Validation(req.Url)
	if err != nil {
		writeJSON(w, domain.JSONResponse{
			Msg:        err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	code, err := h.service.Create(r.Context(), req.Url)
	if err != nil {
		writeJSON(w, domain.JSONResponse{
			Msg:        "Failed to create short URL",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	fmt.Println("Handler Code: ", code)
	fmt.Println("Handler Code Size: ", len(code))

	code = "http://localhost:8080/" + code

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
