// Paquete repository proporciona acceso a la capa de persistencia de datos.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/models"
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

// GetStocks recupera stocks con paginación y ordenamiento.
func (r *StockRepository) GetStocks(ctx context.Context, orderBy string, sortOrder string, offset, limit int) ([]models.Stock, error) {
	// Establecer valores predeterminados si no se proporcionan
	if orderBy == "" {
		orderBy = "time"
	}
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	// Consulta con ordenamiento y paginación
	query := fmt.Sprintf(`
		SELECT 
			ticker, company, target_from, target_to, 
			action, brokerage, rating_from, rating_to, time
		FROM stocks
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, orderBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al consultar stocks: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker,
			&stock.Company,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		); err != nil {
			return nil, fmt.Errorf("error al escanear stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar stocks: %w", err)
	}

	return stocks, nil
}

// CountStocks cuenta el total de stocks en la base de datos.
func (r *StockRepository) CountStocks(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stocks").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar stocks: %w", err)
	}
	return count, nil
}

// GetStocksByBrokerage recupera stocks filtrados por casa de bolsa con paginación.
func (r *StockRepository) GetStocksByBrokerage(ctx context.Context, brokerage string, offset, limit int) ([]models.Stock, error) {
	query := `
		SELECT 
			ticker, company, target_from, target_to, 
			action, brokerage, rating_from, rating_to, time
		FROM stocks
		WHERE brokerage = $1
		ORDER BY time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, brokerage, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al consultar stocks por casa de bolsa: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker,
			&stock.Company,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		); err != nil {
			return nil, fmt.Errorf("error al escanear stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar stocks: %w", err)
	}

	return stocks, nil
}

// CountStocksByBrokerage cuenta el total de stocks para una casa de bolsa específica.
func (r *StockRepository) CountStocksByBrokerage(ctx context.Context, brokerage string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stocks WHERE brokerage = $1", brokerage).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar stocks por casa de bolsa: %w", err)
	}
	return count, nil
}

// GetStocksByTickerPattern recupera stocks cuyo ticker coincida con un patrón.
func (r *StockRepository) GetStocksByTickerPattern(ctx context.Context, tickerPattern string, offset, limit int) ([]models.Stock, error) {
	query := `
		SELECT 
			ticker, company, target_from, target_to, 
			action, brokerage, rating_from, rating_to, time
		FROM stocks
		WHERE ticker ILIKE $1
		ORDER BY time DESC
		LIMIT $2 OFFSET $3
	`

	// Realizamos búsqueda parcial
	pattern := "%" + tickerPattern + "%"

	rows, err := r.db.QueryContext(ctx, query, pattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al consultar stocks por patrón de ticker: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker,
			&stock.Company,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		); err != nil {
			return nil, fmt.Errorf("error al escanear stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar stocks: %w", err)
	}

	return stocks, nil
}

// CountStocksByTickerPattern cuenta el total de stocks que coinciden con un patrón de ticker.
func (r *StockRepository) CountStocksByTickerPattern(ctx context.Context, tickerPattern string) (int, error) {
	var count int

	pattern := "%" + tickerPattern + "%"

	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stocks WHERE ticker ILIKE $1", pattern).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar stocks por patrón de ticker: %w", err)
	}
	return count, nil
}

// GetStocksByRating recupera stocks por su rating (from o to).
func (r *StockRepository) GetStocksByRating(ctx context.Context, rating string, offset, limit int) ([]models.Stock, error) {
	query := `
		SELECT 
			ticker, company, target_from, target_to, 
			action, brokerage, rating_from, rating_to, time
		FROM stocks
		WHERE rating_from = $1 OR rating_to = $1
		ORDER BY time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, rating, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error al consultar stocks por rating: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker,
			&stock.Company,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		); err != nil {
			return nil, fmt.Errorf("error al escanear stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar stocks: %w", err)
	}

	return stocks, nil
}

// CountStocksByRating cuenta el total de stocks con un rating específico.
func (r *StockRepository) CountStocksByRating(ctx context.Context, rating string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM stocks WHERE rating_from = $1 OR rating_to = $1", rating).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error al contar stocks por rating: %w", err)
	}
	return count, nil
}

// GetStockByTicker obtiene un stock por su ticker.
func (r *StockRepository) GetStockByTicker(ctx context.Context, ticker string) (models.Stock, error) {
	var stock models.Stock

	query := `
	SELECT 
		ticker, company, target_from, target_to, 
		action, brokerage, rating_from, rating_to, time
	FROM stocks 
	WHERE ticker = $1
	`

	err := r.db.QueryRowContext(ctx, query, ticker).Scan(
		&stock.Ticker,
		&stock.Company,
		&stock.TargetFrom,
		&stock.TargetTo,
		&stock.Action,
		&stock.Brokerage,
		&stock.RatingFrom,
		&stock.RatingTo,
		&stock.Time,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return stock, fmt.Errorf("stock no encontrado: %s", ticker)
		}
		return stock, fmt.Errorf("error al obtener stock: %w", err)
	}

	return stock, nil
}

// GetStocksByDateRange recupera stocks en un rango de fechas específico.
func (r *StockRepository) GetStocksByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.Stock, error) {
	query := `
		SELECT 
			ticker, company, target_from, target_to, 
			action, brokerage, rating_from, rating_to, time
		FROM stocks
		WHERE time BETWEEN $1 AND $2
		ORDER BY time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error al consultar stocks por rango de fechas: %w", err)
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker,
			&stock.Company,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
		); err != nil {
			return nil, fmt.Errorf("error al escanear stock: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar stocks: %w", err)
	}

	return stocks, nil
}

// Ping verifica la conexión a la base de datos.
func (r *StockRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
