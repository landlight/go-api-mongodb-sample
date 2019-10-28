package healthcheck 

import (
	"net/http"
	"go-api-mongodb-sample/internal/core/config"

	"github.com/go-chi/render"
)

type HealthCheckEndpoint struct{}

func NewHealthCheckEndpoint() *HealthCheckEndpoint {
	return &HealthCheckEndpoint{}
}

func (h *HealthCheckEndpoint) HealthCheck() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
	render.Status(r, config.RR.Internal.Success.HTTPStatusCode())
	render.JSON(w, r, config.RR.Internal.Success)
	return
}
	return http.HandlerFunc(fn)
}