package handler

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/damndelion/test_task_kami/internal/middleware"
	"github.com/damndelion/test_task_kami/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	service *service.Service
	logs    *zap.SugaredLogger
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

	return r
}
