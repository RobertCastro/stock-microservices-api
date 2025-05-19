# Stock API Service

Microservicio responsable de interactuar con la API externa de stocks y sincronizar datos con el sistema.

## Funcionalidades

- Sincronización de datos de stocks desde una API externa
- Verificaciones de salud del servicio

## Requisitos

- Go 1.23 o superior
- Docker (para desarrollo y despliegue)

## Configuración

### Variables de entorno

El servicio requiere las siguientes variables de entorno:

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| SERVER_PORT | Puerto en el que se ejecutará el servidor | 8080 |
| STOCK_API_BASE_URL | URL base de la API externa de stocks | https://api.stockapi.com/v1/stocks |
| STOCK_API_AUTH_TOKEN | Token de autenticación para la API externa | - |

## Desarrollo local

### Ejecutar con Go

```bash
# Clonar el repositorio
git clone https://github.com/RobertCastro/stock-microservices-api.git
cd stock-microservices-api/stock-data-service

# Instalar dependencias
go mod download

# Ejecutar
go run cmd/api/main.go