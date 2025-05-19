// Paquete health proporciona verificaciones de salud del servicio.
package health

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// HealthStatus representa el estado de salud del servicio.
type HealthStatus struct {
	Status         string            `json:"status"`
	Components     map[string]string `json:"components,omitempty"`
	APICredentials bool              `json:"api_credentials_configured"`
	Timestamp      time.Time         `json:"timestamp"`
	Version        string            `json:"version"`
}

// HealthHandler maneja las verificaciones de salud del servicio.
type HealthHandler struct {
	repo *repository.StockRepository
}

// NewHealthHandler crea una nueva instancia de HealthHandler.
func NewHealthHandler(repo *repository.StockRepository) *HealthHandler {
	return &HealthHandler{
		repo: repo,
	}
}

// BasicHealth verifica el estado básico del servicio.
func (h *HealthHandler) BasicHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DetailedHealth verifica el estado detallado del servicio.
func (h *HealthHandler) DetailedHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	status := HealthStatus{
		Status:     "ok",
		Components: make(map[string]string),
		Timestamp:  time.Now(),
		Version:    "1.0.0",
	}

	// Verificar conexión a la base de datos
	if err := h.repo.Ping(ctx); err != nil {
		status.Status = "degradado"
		status.Components["database"] = "error: " + err.Error()
	} else {
		status.Components["database"] = "ok"
	}

	// Verificar configuración de API
	apiToken := os.Getenv("STOCK_API_AUTH_TOKEN")
	apiURL := os.Getenv("STOCK_API_BASE_URL")

	if apiToken == "" || apiURL == "" {
		status.APICredentials = false
		status.Components["api_config"] = "faltan credenciales"
		status.Status = "degradado"
	} else {
		status.APICredentials = true
		status.Components["api_config"] = "configurado"
	}

	c.JSON(http.StatusOK, status)
}
