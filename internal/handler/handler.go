package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/muosilva/url-shortener/internal/domain"
	"github.com/muosilva/url-shortener/internal/service"
	"github.com/muosilva/url-shortener/metrics"
)

const (
	dm = "http://localhost:8080/"
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
		slog.Error("failed to decode request body", "error", err)
		writeJSON(w, domain.JSONResponse{
			Msg:        "Invalid JSON body",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	err = req.Validation(req.Url)
	if err != nil {
		slog.Error("validation error", "error", err)
		writeJSON(w, domain.JSONResponse{
			Msg:        err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	code, err := h.service.Create(r.Context(), req.Url)
	if err != nil {
		slog.Error("failed to create short url", "error", err)
		writeJSON(w, domain.JSONResponse{
			Msg:        "Failed to create short URL",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	metrics.URLsCreatedTotal.Inc()

	resp := domain.JSONResponse{
		Msg:        "Created successfully!",
		StatusCode: http.StatusCreated,
		Url:        dm + code,
	}

	slog.Info("url created")
	writeJSON(w, resp)
}

func (h *Handler) RedirectToOriginalURL(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	originalUrl, err := h.service.Get(r.Context(), code)
	if err != nil {
		statusCode := http.StatusInternalServerError
		msg := "Failed to get the original URL"

		if errors.Is(err, service.ErrURLNotFound) {
			statusCode = http.StatusNotFound
			msg = "URL not found"
		}

		writeJSON(w, domain.JSONResponse{
			Msg:        msg,
			StatusCode: statusCode,
		})
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusFound)
}

func writeJSON(w http.ResponseWriter, response domain.JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	_ = json.NewEncoder(w).Encode(response)
}
