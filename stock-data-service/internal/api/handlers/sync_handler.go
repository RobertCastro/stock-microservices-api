// Paquete handlers contiene los manejadores de solicitudes HTTP.
package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-data-service/internal/client"
	"github.com/RobertCastro/stock-microservices-api/stock-data-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// SyncResponse representa la respuesta a una operación de sincronización.
type SyncResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SyncHandler maneja las solicitudes de sincronización con la API externa.
type SyncHandler struct {
	client *client.ExternalAPIClient
	repo   *repository.StockRepository
}

// NewSyncHandler crea una nueva instancia de SyncHandler.
func NewSyncHandler(client *client.ExternalAPIClient, repo *repository.StockRepository) *SyncHandler {
	return &SyncHandler{
		client: client,
		repo:   repo,
	}
}

// SyncStocks maneja la solicitud para sincronizar stocks desde la API externa.
// Esta es una operación asíncrona que devuelve inmediatamente y continúa en segundo plano.
func (h *SyncHandler) SyncStocks(c *gin.Context) {
	// Verificar que el token de autenticación esté configurado
	apiToken := os.Getenv("STOCK_API_AUTH_TOKEN")
	if apiToken == "" {
		response := SyncResponse{
			Status:  "error",
			Message: "Error de configuración: No se encontró el token de autenticación para la API",
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Responder inmediatamente indicando que la sincronización ha comenzado
	response := SyncResponse{
		Status:  "accepted",
		Message: "Sincronización iniciada, esto puede tomar varios minutos",
	}

	// Enviar respuesta 202 Accepted
	c.JSON(http.StatusAccepted, response)

	// Ejecutar la sincronización en una goroutine separada
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		// Obtener todos los stocks de la API externa
		stocks, err := h.client.FetchAllStocks()
		if err != nil {
			log.Printf("Error al obtener stocks de la API: %v", err)
			return
		}

		if len(stocks) == 0 {
			log.Printf("No se encontraron stocks para sincronizar")
			return
		}

		// Guardar los stocks en la base de datos
		if err := h.repo.SaveStocks(ctx, stocks); err != nil {
			log.Printf("Error al guardar stocks en la base de datos: %v", err)
			return
		}

		log.Printf("Sincronización completada: %d stocks guardados", len(stocks))
	}()
}
