// Paquete models define las estructuras de datos que se utilizan en la aplicación.
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

// APIResponse representa la respuesta de la API externa.
type APIResponse struct {
	// Lista de stocks en la respuesta
	Items []Stock `json:"items"`
	// Token para la siguiente página de resultados
	NextPage string `json:"next_page"`
}
