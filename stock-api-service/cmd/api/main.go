package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/api"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/client"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Nota: No se pudo cargar el archivo .env: %v", err)
	}

	// Inicializar configuración
	cfg := config.NewConfig()

	// Crear cliente de API externa
	externalClient := client.NewExternalAPIClient(cfg.StockAPIBaseURL, cfg.StockAPIToken)

	// Configurar servidor HTTP con Gin
	router := api.NewRouter(externalClient)
	server := router.SetupServer(cfg.ServerPort)

	// Arrancar servidor en una goroutine
	go func() {
		log.Printf("Servidor iniciado en el puerto %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Configurar graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Cerrar con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error al cerrar el servidor: %v", err)
	}
}
