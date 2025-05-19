// Paquete api proporciona la configuración y las rutas del API.
package api

import (
	"net/http"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/api/handlers"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/api/middlewares"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/health"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// Router maneja la configuración de rutas de la API.
type Router struct {
	stockHandler          *handlers.StockHandler
	recommendationHandler *handlers.RecommendationHandler
	healthHandler         *health.HealthHandler
}

// NewRouter crea una nueva instancia del router.
func NewRouter(repo *repository.StockRepository) *Router {
	return &Router{
		stockHandler:          handlers.NewStockHandler(repo),
		recommendationHandler: handlers.NewRecommendationHandler(repo),
		healthHandler:         health.NewHealthHandler(repo),
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
		// Rutas para stocks
		api.GET("/stocks", r.stockHandler.ListStocks)
		api.GET("/stocks/:ticker", r.stockHandler.GetStockDetails)

		// Ruta para recomendaciones
		api.GET("/recommendations", r.recommendationHandler.GetRecommendations)
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
