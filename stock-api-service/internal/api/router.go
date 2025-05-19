// Paquete api proporciona la configuración y las rutas del API.
package api

import (
	"net/http"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/api/handlers"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/api/middlewares"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/client"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/health"
	"github.com/gin-gonic/gin"
)

// Router maneja la configuración de rutas de la API.
type Router struct {
	syncHandler   *handlers.SyncHandler
	healthHandler *health.HealthHandler
}

// NewRouter crea una nueva instancia del router.
func NewRouter(client *client.ExternalAPIClient) *Router {
	return &Router{
		syncHandler:   handlers.NewSyncHandler(client),
		healthHandler: health.NewHealthHandler(),
	}
}

// SetupRoutes configura las rutas del router.
func (r *Router) SetupRoutes(router *gin.Engine) {
	// Middleware para todas las rutas
	router.Use(middlewares.Logger())
	router.Use(middlewares.CORS())

	// Rutas para la API
	api := router.Group("/api/v1")
	{
		// Ruta para sincronización
		api.POST("/sync", r.syncHandler.SyncStocks)
	}

	// Rutas para health checks
	router.GET("/health", r.healthHandler.BasicHealth)
	router.GET("/health/detailed", r.healthHandler.DetailedHealth)
}

// SetupServer configura y devuelve un servidor HTTP listo para usar.
func (r *Router) SetupServer(port string) *http.Server {
	if port == "" {
		port = "8080"
	}

	// Crear router Gin
	router := gin.Default()

	// Configurar rutas
	r.SetupRoutes(router)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server
}
