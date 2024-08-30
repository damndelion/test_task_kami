package handler

import (
	"bytes"
	"encoding/json"
	"github.com/damndelion/test_task_kami/configs"
	"github.com/damndelion/test_task_kami/internal/infrastructure/database"
	"github.com/damndelion/test_task_kami/internal/models"
	"github.com/damndelion/test_task_kami/internal/repository"
	"github.com/damndelion/test_task_kami/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateReservation(t *testing.T) {
	newReservation := models.BookingCreate{
		RoomID:    "B-101",
		StartTime: "2024-08-31-15:00:00",
		EndTime:   "2024-08-31-16:00:00",
	}
	writer := makeRequest("POST", "/api/v1/reservation", newReservation)
	assert.Equal(t, http.StatusCreated, writer.Code)

	var response map[string]int
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Contains(t, response, "id")
}

func TestCreateReservationValidationError(t *testing.T) {
	invalidReservation := models.BookingCreate{
		RoomID:    "room1",
		StartTime: "invalid-format",
		EndTime:   "2024-08-30-16:00:00",
	}
	writer := makeRequest("POST", "/api/v1/reservation", invalidReservation)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestGetReservationsByRoomID(t *testing.T) {
	roomID := "A-101"
	writer := makeRequest("GET", "/api/v1/reservation/"+roomID, nil)
	assert.Equal(t, http.StatusOK, writer.Code)

	var response []models.BookingDTO
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.NotEmpty(t, response)
}

func TestGetReservationsByRoomIDNotFound(t *testing.T) {
	roomID := "z-101"
	writer := makeRequest("GET", "/api/v1/reservation/"+roomID, nil)
	assert.Equal(t, http.StatusNotFound, writer.Code)
}

func TestCreateReservationEndTimeBeforeStartTime(t *testing.T) {
	invalidReservation := models.BookingCreate{
		RoomID:    "B-102",
		StartTime: "2024-08-31-16:00:00",
		EndTime:   "2024-08-31-15:00:00",
	}
	writer := makeRequest("POST", "/api/v1/reservation", invalidReservation)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestCreateReservationOverlap(t *testing.T) {
	initialReservation := models.BookingCreate{
		RoomID:    "D-101",
		StartTime: "2024-08-31-12:00:00",
		EndTime:   "2024-08-31-13:00:00",
	}
	writer := makeRequest("POST", "/api/v1/reservation", initialReservation)
	assert.Equal(t, http.StatusCreated, writer.Code)

	var response map[string]int
	json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Contains(t, response, "id")

	overlappingReservation := models.BookingCreate{
		RoomID:    "D-101",
		StartTime: "2024-08-31-12:30:00",
		EndTime:   "2024-08-31-13:30:00",
	}
	writer = makeRequest("POST", "/api/v1/reservation", overlappingReservation)
	assert.Equal(t, http.StatusConflict, writer.Code)

}

func makeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

// Initialize the router with Chi
func router() *chi.Mux {
	// Read configs
	config, err := configs.InitConfigs()
	if err != nil {
		log.Fatalf("error with loading config: %s", err.Error())
	}

	// Connect to DB
	db, err := database.NewPostgresDB(config.Postgres)
	if err != nil {
		log.Fatalf("error connection to db: %s", err.Error())

	}
	r := chi.NewRouter()
	logs := zap.NewNop().Sugar()
	repo := repository.NewRepository(db)
	mockService := service.NewService(logs, repo)

	h := NewHandler(logs, mockService)

	// Middleware setup
	r.Use(middleware.Logger) // Logs each request to the console

	// Define routes
	r.Post("/api/v1/reservation", h.CreateReservation)
	r.Get("/api/v1/reservation/{roomID}", h.GetReservationsByRoomID)

	return r
}
