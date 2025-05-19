// Paquete main es el punto de entrada principal de la aplicación.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/api"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/config"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/database"
	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Nota: No se pudo cargar el archivo .env: %v", err)
	}

	// Inicializar configuración
	cfg := config.NewConfig()

	// Conectar a la base de datos
	db, err := database.Connect(cfg.GetDBConnectionString())
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Crear repositorio de stocks
	repo := repository.NewStockRepository(db)

	// Configurar servidor HTTP con Gin
	router := api.NewRouter(repo)
	server := router.SetupServer(cfg.ServerPort)

	// Arrancar servidor en una goroutine
	go func() {
		log.Printf("Servidor iniciado en el puerto %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Configurar graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Apagando servidor...")

	// Cerrar con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error al cerrar el servidor: %v", err)
	}

	log.Println("Servidor apagado correctamente")
}
