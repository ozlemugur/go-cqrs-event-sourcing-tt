// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	"github.com/gin-contrib/cors"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
	_ "github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/docs"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/internal/usecase"
)

// NewRouter -.
// Swagger spec:
// @title       Wallet Management Service
// @description Wallet Management
// @version     1.0
// @host        localhost:8081
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.WalletHandler) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe --Checks if the container is still alive or stuck. If a container fails the liveness probe, Kubernetes will restart it.
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	l.Info("Wallet Management Service http://localhost:8081/swagger/index.html")

	// Routers
	h := handler.Group("/v1")
	{
		newWalletRoutes(h, t, l)
	}

}
