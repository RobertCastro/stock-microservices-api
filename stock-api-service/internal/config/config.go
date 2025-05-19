// Paquete config proporciona funcionalidades para la configuración de la aplicación.
package config

import (
	"fmt"
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
	// Configuración de la base de datos
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// NewConfig crea una nueva instancia de configuración con valores predeterminados
// y los sobrescribe con variables de entorno si están disponibles.
func NewConfig() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		StockAPIBaseURL: getEnv("STOCK_API_BASE_URL", "https://api.stockapi.com/v1/stocks"),
		StockAPIToken:   getEnv("STOCK_API_AUTH_TOKEN", ""),

		// Configuración de base de datos
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "26257"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "stockdb"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

// GetDBConnectionString construye la cadena de conexión para la base de datos.
func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

// getEnv obtiene el valor de una variable de entorno o devuelve un valor predeterminado.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
