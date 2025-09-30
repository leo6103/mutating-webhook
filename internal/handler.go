package internal

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes attaches the webhook endpoints to the Gin engine.
func RegisterRoutes(r *gin.Engine, svc *MutateService) {
	r.GET("/healthz", HealthzHandler)
	r.POST("/mutate", MutateHandler(svc))
}

// HealthzHandler returns a simple JSON payload for readiness checks.
func HealthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// MutateHandler logs the incoming request body and returns 200 OK.
func MutateHandler(svc *MutateService) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		defer c.Request.Body.Close()

		if err := svc.Mutate(body); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}
