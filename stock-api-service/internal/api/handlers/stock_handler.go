// Paquete handlers contiene los manejadores de solicitudes HTTP.
package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/models"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// StockHandler maneja las solicitudes relacionadas con stocks.
type StockHandler struct {
	repo *repository.StockRepository
}

// NewStockHandler crea una nueva instancia de StockHandler.
func NewStockHandler(repo *repository.StockRepository) *StockHandler {
	return &StockHandler{
		repo: repo,
	}
}

// ListStocks maneja la solicitud para listar stocks con filtros y paginación.
func (h *StockHandler) ListStocks(c *gin.Context) {
	// Extraer parámetros de filtrado
	brokerage := c.Query("brokerage")
	ticker := c.Query("ticker")
	rating := c.Query("rating")

	// Parsear parámetros de paginación
	pagination := h.parsePagination(c)

	// Extraer parámetros de ordenamiento
	orderBy := c.Query("order_by")
	sortOrder := c.Query("sort")

	// Validar y establecer valores predeterminados para ordenamiento
	if orderBy == "" {
		orderBy = "time"
	} else {
		// Validar que el campo de ordenamiento sea válido
		validFields := map[string]bool{
			"ticker": true, "company": true, "brokerage": true,
			"rating_from": true, "rating_to": true, "time": true,
		}
		if !validFields[orderBy] {
			orderBy = "time"
		}
	}

	if sortOrder == "" {
		sortOrder = "DESC"
	} else {
		sortOrder = strings.ToUpper(sortOrder)
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "DESC"
		}
	}

	var stocks []models.Stock
	var totalStocks int
	var err error

	// Obtener stocks según los filtros
	if ticker != "" {
		stocks, err = h.repo.GetStocksByTickerPattern(c.Request.Context(), ticker, pagination.Offset, pagination.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al obtener stocks por ticker: " + err.Error(),
			})
			return
		}

		totalStocks, err = h.repo.CountStocksByTickerPattern(c.Request.Context(), ticker)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al contar stocks por ticker: " + err.Error(),
			})
			return
		}
	} else if brokerage != "" {
		stocks, err = h.repo.GetStocksByBrokerage(c.Request.Context(), brokerage, pagination.Offset, pagination.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al obtener stocks por brokerage: " + err.Error(),
			})
			return
		}

		totalStocks, err = h.repo.CountStocksByBrokerage(c.Request.Context(), brokerage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al contar stocks por brokerage: " + err.Error(),
			})
			return
		}
	} else if rating != "" {
		stocks, err = h.repo.GetStocksByRating(c.Request.Context(), rating, pagination.Offset, pagination.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al obtener stocks por rating: " + err.Error(),
			})
			return
		}

		totalStocks, err = h.repo.CountStocksByRating(c.Request.Context(), rating)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al contar stocks por rating: " + err.Error(),
			})
			return
		}
	} else {
		// Sin filtros, obtener todos los stocks
		stocks, err = h.repo.GetStocks(c.Request.Context(), orderBy, sortOrder, pagination.Offset, pagination.Limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al obtener stocks: " + err.Error(),
			})
			return
		}

		totalStocks, err = h.repo.CountStocks(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al contar stocks: " + err.Error(),
			})
			return
		}
	}

	// Calcular el total de páginas
	totalPages := (totalStocks + pagination.Limit - 1) / pagination.Limit

	// Crear respuesta
	response := models.StockListResponse{
		Stocks:       stocks,
		TotalStocks:  totalStocks,
		TotalPages:   totalPages,
		CurrentPage:  pagination.Page,
		ItemsPerPage: pagination.Limit,
	}

	c.JSON(http.StatusOK, response)
}

// GetStockDetails maneja la solicitud para obtener los detalles de un stock específico por ticker.
func (h *StockHandler) GetStockDetails(c *gin.Context) {
	ticker := c.Param("ticker")

	if ticker == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Se requiere especificar un ticker",
		})
		return
	}

	// Obtener stock por ticker exacto
	stock, err := h.repo.GetStockByTicker(c.Request.Context(), ticker)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Stock no encontrado: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// parsePagination extrae y valida los parámetros de paginación de la solicitud.
func (h *StockHandler) parsePagination(c *gin.Context) models.Pagination {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := c.Query("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			if s > 100 {
				s = 100
			}
			pageSize = s
		}
	}

	// Calcular offset para la consulta a la BD
	offset := (page - 1) * pageSize

	return models.Pagination{
		Page:   page,
		Limit:  pageSize,
		Offset: offset,
	}
}
