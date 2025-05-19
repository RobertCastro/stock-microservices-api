// Paquete handlers contiene los manejadores de solicitudes HTTP.
package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/client"
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
}

// NewSyncHandler crea una nueva instancia de SyncHandler.
func NewSyncHandler(client *client.ExternalAPIClient) *SyncHandler {
	return &SyncHandler{
		client: client,
	}
}

// SyncStocks maneja la solicitud para sincronizar stocks desde la API externa.
// Operación asíncrona que devuelve inmediatamente y continúa en segundo plano.
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
		Message: "Sincronización iniciada. Se publicará un evento cuando se complete.",
	}

	// Enviar respuesta 202 Accepted
	c.JSON(http.StatusAccepted, response)

	// Ejecutar la sincronización en una goroutine separada
	go func() {
		stocks, err := h.client.FetchAllStocks()
		if err != nil {
			log.Printf("Error al obtener stocks de la API: %v", err)
			return
		}

		if len(stocks) == 0 {
			log.Printf("No se encontraron stocks para sincronizar")
			return
		}

		// TODO publicar los stocks

		log.Printf("Sincronización completada: %d stocks obtenidos", len(stocks))
	}()
}
