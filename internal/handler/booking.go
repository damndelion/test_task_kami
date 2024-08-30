package handler

import (
	"encoding/json"
	"errors"
	"github.com/damndelion/test_task_kami/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (h *Handler) BookingRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateReservation)
	r.Get("/{roomID}", h.GetReservationsByRoomID)
	return r
}

// CreateReservation creates a new reservation.
//
//	@Summary		Create a new reservation
//	@Description	Create a reservation for a room
//	@Tags			Bookings
//	@Accept			json
//	@Produce		json
//	@Param			booking	body		models.BookingCreate	true	"BookingCreate - RoomID is required. StartTime and EndTime must be in the format YYYY-MM-DD-hh:mm:ss and are required."
//	@Success		201		{object}	map[string]int			"id"
//	@Failure		400		{string}	string					"invalid input"
//	@Failure		500		{string}	string					"internal server error"
//	@Router			/api/v1/bookings [post]
func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var input models.BookingCreate

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logs.Errorf("failed to decode request body: %v", err)
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Validate
	if err := h.validate.Struct(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			h.logs.Warnf("validation failed: %v", validationErrors)
			http.Error(w, "validation error: "+validationErrors.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateReservation(r.Context(), input)
	if err != nil {
		h.logs.Errorf("failed to create reservation: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int{"id": id}); err != nil {
		h.logs.Errorf("failed to encode response: %v", err)
	}
}

// GetReservationsByRoomID retrieves reservations by room ID.
//
//	@Summary		Get reservations by room ID
//	@Description	Retrieve all reservations for a specific room
//	@Tags			Bookings
//	@Produce		json
//	@Param			roomID	path		string	true	"Room ID"
//	@Success		200		{array}		models.BookingDTO
//	@Failure		500		{string}	string	"internal server error"
//	@Router			/api/v1/bookings/{roomID} [get]
func (h *Handler) GetReservationsByRoomID(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "roomID")

	reservations, err := h.service.GetReservationByRoomID(r.Context(), roomIDStr)
	if err != nil {
		h.logs.Errorf("failed to get reservations for room ID %s: %v", roomIDStr, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(reservations); err != nil {
		h.logs.Errorf("failed to encode reservations response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
