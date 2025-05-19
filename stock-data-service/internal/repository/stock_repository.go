// Paquete repository proporciona acceso a la capa de persistencia de datos.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/RobertCastro/stock-microservices-api/stock-data-service/internal/models"
)

// StockRepository maneja las operaciones de base de datos para los stocks.
type StockRepository struct {
	db *sql.DB
}

// NewStockRepository crea una nueva instancia del repositorio de stocks.
func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{
		db: db,
	}
}

// InitDB inicializa la base de datos creando las tablas necesarias.
func (r *StockRepository) InitDB(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS stocks (
        ticker STRING PRIMARY KEY,
        company STRING NOT NULL,
        target_from STRING NOT NULL,
        target_to STRING NOT NULL,
        action STRING NOT NULL,
        brokerage STRING NOT NULL,
        rating_from STRING NOT NULL,
        rating_to STRING NOT NULL,
        time TIMESTAMP NOT NULL,
        created_at TIMESTAMP DEFAULT current_timestamp()
    )
    `

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// SaveStocks guarda múltiples stocks en la base de datos utilizando una transacción.
func (r *StockRepository) SaveStocks(ctx context.Context, stocks []models.Stock) error {
	// Iniciar transacción
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", err)
	}

	// Preparar statement para inserción/actualización
	stmt, err := tx.PrepareContext(ctx, `
        UPSERT INTO stocks (
            ticker, company, target_from, target_to, 
            action, brokerage, rating_from, rating_to, time
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error al preparar el statement: %w", err)
	}
	defer stmt.Close()

	// Insertar cada stock
	for _, stock := range stocks {
		_, err := stmt.ExecContext(
			ctx,
			stock.Ticker,
			stock.Company,
			stock.TargetFrom,
			stock.TargetTo,
			stock.Action,
			stock.Brokerage,
			stock.RatingFrom,
			stock.RatingTo,
			stock.Time,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error al guardar el stock %s: %w", stock.Ticker, err)
		}
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	return nil
}

// Ping verifica la conexión a la base de datos.
func (r *StockRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
