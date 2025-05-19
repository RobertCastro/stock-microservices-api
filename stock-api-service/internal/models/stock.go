// Paquete models define las estructuras de datos utilizadas en la aplicación.
package models

import (
	"time"
)

// Stock representa la información de una acción en bolsa.
type Stock struct {
	// Símbolo o ticker de la acción
	Ticker string `json:"ticker"`
	// Nombre de la compañía
	Company string `json:"company"`
	// Precio objetivo anterior
	TargetFrom string `json:"target_from"`
	// Precio objetivo actual
	TargetTo string `json:"target_to"`
	// Tipo de acción realizada sobre la recomendación (upgraded, downgraded, etc.)
	Action string `json:"action"`
	// Casa de bolsa que emitió la recomendación
	Brokerage string `json:"brokerage"`
	// Calificación anterior
	RatingFrom string `json:"rating_from"`
	// Calificación actual
	RatingTo string `json:"rating_to"`
	// Fecha y hora de la actualización
	Time time.Time `json:"time"`
}

// Pagination contiene la información de paginación para las consultas.
type Pagination struct {
	// Página actual
	Page int
	// Cantidad de elementos por página
	Limit int
	// Desplazamiento para la consulta SQL
	Offset int
}

// StockListResponse representa la respuesta para el listado de stocks.
type StockListResponse struct {
	// Lista de stocks
	Stocks []Stock `json:"stocks"`
	// Total de stocks que coinciden con los criterios de búsqueda
	TotalStocks int `json:"total_stocks"`
	// Total de páginas disponibles
	TotalPages int `json:"total_pages"`
	// Página actual
	CurrentPage int `json:"current_page"`
	// Cantidad de elementos por página
	ItemsPerPage int `json:"items_per_page"`
}

// RecommendationResult representa el resultado de una recomendación.
type RecommendationResult struct {
	// Stock recomendado
	Stock Stock `json:"stock"`
	// Puntuación de la recomendación
	Score float64 `json:"score"`
	// Explicación de la recomendación
	Rationale string `json:"rationale"`
	// Retorno potencial estimado
	PotentialReturn string `json:"potential_return"`
}

// RecommendationResponse representa la respuesta del servicio de recomendaciones.
type RecommendationResponse struct {
	// Lista de recomendaciones
	Recommendations []RecommendationResult `json:"recommendations"`
	// Fecha y hora de generación
	GeneratedAt time.Time `json:"generated_at"`
	// Cantidad de recomendaciones
	Count int `json:"count"`
	// Mensaje informativo
	Message string `json:"message"`
}
