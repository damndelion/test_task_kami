package handler

import (
	"net/http"

	"github.com/damndelion/test_task_kami/internal/middleware"
	"github.com/damndelion/test_task_kami/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Handler struct {
	service  *service.Service
	logs     *zap.SugaredLogger
	validate *validator.Validate
}

func NewHandler(logs *zap.SugaredLogger, service *service.Service) *Handler {
	return &Handler{
		logs:    logs,
		service: service,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORSMiddleware())

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Health check
	r.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/reservation", h.BookingRoutes())
	})

	return r
}
