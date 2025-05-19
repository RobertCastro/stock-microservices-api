// Paquete database proporciona funcionalidades para la conexión y gestión de la base de datos.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Connect establece una conexión con la base de datos CockroachDB.
func Connect(connectionString string) (*sql.DB, error) {
	// Extraer la parte base de la cadena de conexión para crear la base de datos si no existe
	baseConnectionString := extractBaseConnectionString(connectionString)

	// Conectar a la base de datos principal para crear la base de datos específica si no existe
	tempDB, err := sql.Open("postgres", baseConnectionString)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	defer tempDB.Close()

	// Extraer el nombre de la base de datos de la cadena de conexión
	dbName := extractDBName(connectionString)
	if dbName == "" {
		return nil, fmt.Errorf("no se pudo extraer el nombre de la base de datos de la cadena de conexión")
	}

	// Crear la base de datos si no existe
	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return nil, fmt.Errorf("error al crear la base de datos: %w", err)
	}

	// Conectar a la base de datos específica
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos específica: %w", err)
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

// extractBaseConnectionString extrae la parte base de la cadena de conexión.
func extractBaseConnectionString(connectionString string) string {
	// Buscar el último '/' que precede al nombre de la base de datos
	dbNameIndex := strings.LastIndex(connectionString, "/")
	if dbNameIndex == -1 {
		return connectionString
	}

	// Verificar si hay parámetros adicionales después del nombre de la base de datos
	questionMarkIndex := strings.Index(connectionString[dbNameIndex:], "?")
	if questionMarkIndex != -1 {
		return connectionString[:dbNameIndex+1] + connectionString[dbNameIndex+questionMarkIndex:]
	}

	return connectionString[:dbNameIndex+1]
}

// extractDBName extrae el nombre de la base de datos de la cadena de conexión.
func extractDBName(connectionString string) string {
	// Buscar el último '/' que precede al nombre de la base de datos
	dbNameIndex := strings.LastIndex(connectionString, "/")
	if dbNameIndex == -1 {
		return ""
	}

	// Extraer la parte después del último '/'
	dbNamePart := connectionString[dbNameIndex+1:]

	// Si hay parámetros adicionales, eliminarlos
	questionMarkIndex := strings.Index(dbNamePart, "?")
	if questionMarkIndex != -1 {
		return dbNamePart[:questionMarkIndex]
	}

	return dbNamePart
}
