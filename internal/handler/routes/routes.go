package routes

import (
	"time"

	"go-api-mongodb-sample/internal/core/db"
	"go-api-mongodb-sample/internal/core/config"
	"go-api-mongodb-sample/internal/handler/middlewares"
	"go-api-mongodb-sample/internal/pkg/healthcheck"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func NewRouter(dbConn *db.DBConn) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middlewares.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(db.DialMongoMiddleware(dbConn, config.CF.MongoDB.Schema.DBName))

	r.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			// healthcheck
			healthCheckEndpoint := healthcheck.NewHealthCheckEndpoint()
			r.Get("/healthcheck", healthCheckEndpoint.HealthCheck())
		})
	})
	return r
}
