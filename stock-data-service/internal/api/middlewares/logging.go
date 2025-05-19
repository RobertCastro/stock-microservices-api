// Paquete middlewares proporciona middlewares para el router Gin.
package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger es un middleware que registra información sobre las solicitudes HTTP.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tiempo de inicio
		start := time.Now()

		// Procesar la solicitud
		c.Next()

		// Tiempo de finalización
		duration := time.Since(start)

		// Datos de la solicitud
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		method := c.Request.Method
		statusCode := c.Writer.Status()

		// Registrar la información
		log.Printf("%s %s %d %s", method, path, statusCode, duration)
	}
}
