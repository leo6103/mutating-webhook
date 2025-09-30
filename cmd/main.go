package main

import (
	"log"

	"github.com/gin-gonic/gin"

	internal "mutating-webhook/internal"
)

const (
	logFileName = "webhook.log"
	tlsPort     = "9443"
	tlsCertPath = "/tls/tls.crt"
	tlsKeyPath  = "/tls/tls.key"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	svc := internal.NewMutateService(logFileName)
	internal.RegisterRoutes(r, svc)

	addr := ":" + tlsPort
	if err := r.RunTLS(addr, tlsCertPath, tlsKeyPath); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
