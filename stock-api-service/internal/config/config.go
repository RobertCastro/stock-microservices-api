// Paquete config proporciona funcionalidades para la configuración de la aplicación.
package config

import (
	"os"
)

// Config contiene la configuración de la aplicación.
type Config struct {
	// Puerto del servidor
	ServerPort string
	// URL base de la API externa de stocks
	StockAPIBaseURL string
	// Token de autenticación para la API externa
	StockAPIToken string
}

// NewConfig crea una nueva instancia de configuración con valores predeterminados
// y los sobrescribe con variables de entorno si están disponibles.
func NewConfig() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		StockAPIBaseURL: getEnv("STOCK_API_BASE_URL", "https://api.stockapi.com/v1/stocks"),
		StockAPIToken:   getEnv("STOCK_API_AUTH_TOKEN", ""),
	}
}

// getEnv obtiene el valor de una variable de entorno o devuelve un valor predeterminado.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
