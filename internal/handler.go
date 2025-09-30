package internal

import (
	"io"
	"log"
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
	c.JSON(http.StatusOK, gin.H{"status": "okk"})
}

// MutateHandler logs the incoming request body and returns 200 OK or the mutation response.
func MutateHandler(svc *MutateService) gin.HandlerFunc {
	log.Println("MUTATE HANDLER")
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		defer c.Request.Body.Close()

		respBody, err := svc.Mutate(body)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if len(respBody) == 0 {
			c.Status(http.StatusOK)
			return
		}

		c.Data(http.StatusOK, "application/json", respBody)
	}
}
