// Paquete handlers contiene los manejadores de solicitudes HTTP.
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/algorithm"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/models"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// RecommendationHandler maneja las solicitudes relacionadas con recomendaciones de stocks.
type RecommendationHandler struct {
	repo        *repository.StockRepository
	recommender *algorithm.StockRecommender
}

// NewRecommendationHandler crea una nueva instancia de RecommendationHandler.
func NewRecommendationHandler(repo *repository.StockRepository) *RecommendationHandler {
	return &RecommendationHandler{
		repo:        repo,
		recommender: algorithm.NewStockRecommender(),
	}
}

// GetRecommendations genera recomendaciones de stocks.
func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	// Obtener stocks recientes para análisis (últimos 30 días)
	endDate := time.Now()
	startDate := endDate.AddDate(0, -1, 0)

	stocks, err := h.repo.GetStocksByDateRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener stocks para recomendaciones: " + err.Error(),
		})
		return
	}

	// Generar recomendaciones
	recommendationResults := h.recommender.GenerateRecommendations(stocks, 10)

	// Crear respuesta
	response := models.RecommendationResponse{
		Recommendations: recommendationResults,
		GeneratedAt:     time.Now(),
		Count:           len(recommendationResults),
		Message:         h.generateResponseMessage(len(recommendationResults)),
	}

	c.JSON(http.StatusOK, response)
}

// generateResponseMessage genera un mensaje para la respuesta.
func (h *RecommendationHandler) generateResponseMessage(count int) string {
	if count == 0 {
		return "No se encontraron recomendaciones para hoy. Intente más tarde cuando haya nuevas actualizaciones."
	} else if count == 1 {
		return "Se encontró 1 recomendación de inversión para hoy."
	} else {
		return fmt.Sprintf("Se encontraron %d recomendaciones de inversión para hoy.", count)
	}
}
