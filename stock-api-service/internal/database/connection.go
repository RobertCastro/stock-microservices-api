// Paquete database proporciona funcionalidades para la conexión y gestión de la base de datos.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Connect establece una conexión con la base de datos CockroachDB.
func Connect(connectionString string) (*sql.DB, error) {
	// Conectar a la base de datos
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	// Configurar la conexión
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verificar la conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al verificar la conexión a la base de datos: %w", err)
	}

	log.Println("Conexión exitosa a CockroachDB")
	return db, nil
}
