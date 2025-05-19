// Paquete client proporciona funcionalidades de la api externa de stocks
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/RobertCastro/stock-microservices-api/stock-api-service/internal/models"
)

// APIError representa un error devuelto por la API externa.
type APIError struct {
	StatusCode int
	Body       string
	URL        string
}

// Error implementa la interfaz error.
func (e *APIError) Error() string {
	return fmt.Sprintf("API retornó estado %d para URL %s: %s", e.StatusCode, e.URL, e.Body)
}

// ExternalAPIClient maneja la comunicación con la API externa de stocks.
type ExternalAPIClient struct {
	httpClient *http.Client
	baseURL    string
	authToken  string
}

// NewExternalAPIClient crea un nuevo cliente para la API externa.
func NewExternalAPIClient(baseURL, authToken string) *ExternalAPIClient {
	// Usar valores predeterminados si no se proporcionan
	if baseURL == "" {
		baseURL = "https://api.stockapi.com/v1/stocks"
	}

	return &ExternalAPIClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   baseURL,
		authToken: authToken,
	}
}

// FetchStocks obtiene una página de stocks desde la API externa.
func (c *ExternalAPIClient) FetchStocks(nextPage string) ([]models.Stock, string, error) {
	if c.authToken == "" {
		return nil, "", fmt.Errorf("no se ha configurado el token de autenticación (STOCK_API_AUTH_TOKEN)")
	}

	// Construir URL con parámetros de paginación si es necesario
	reqURL := c.baseURL
	if nextPage != "" {
		params := url.Values{}
		params.Add("next_page", nextPage)
		reqURL = fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
	}

	// Crear la solicitud con encabezados de autenticación
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error al crear la solicitud: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.authToken)
	req.Header.Add("Content-Type", "application/json")

	// Realizar la solicitud
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error al realizar la solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error al leer el cuerpo de la respuesta: %w", err)
	}

	// Verificar si la respuesta tiene un código de estado exitoso
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusGone {
			return nil, "", fmt.Errorf("el recurso de la API ya no está disponible (410 Gone). El endpoint de la API podría estar obsoleto o haber sido movido")
		}

		return nil, "", &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(bodyBytes),
			URL:        reqURL,
		}
	}

	// Decodificar la respuesta JSON
	var apiResp models.APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, "", fmt.Errorf("error al decodificar la respuesta: %w", err)
	}

	return apiResp.Items, apiResp.NextPage, nil
}

// FetchAllStocks recupera todos los stocks paginando automáticamente.
func (c *ExternalAPIClient) FetchAllStocks() ([]models.Stock, error) {
	var allStocks []models.Stock
	nextPage := ""
	maxRetries := 3
	retryCount := 0

	for {
		stocks, newNextPage, err := c.FetchStocks(nextPage)
		if err != nil {
			// Manejar el caso especial de recurso no disponible
			if err.Error() == "el recurso de la API ya no está disponible (410 Gone). El endpoint de la API podría estar obsoleto o haber sido movido" {
				return nil, err
			}

			// Reintentar en caso de error
			retryCount++
			if retryCount <= maxRetries {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, err
		}

		// Restablecer el contador de reintentos en caso de éxito
		retryCount = 0

		// Agregar los stocks recuperados al resultado
		if len(stocks) > 0 {
			allStocks = append(allStocks, stocks...)
		}

		// Si no hay más páginas, terminar
		if newNextPage == "" {
			break
		}

		// Continuar con la siguiente página
		nextPage = newNextPage
	}

	return allStocks, nil
}
