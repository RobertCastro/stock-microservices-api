# Stock API Service

Microservicio encargado de proporcionar una API para consultar y obtener recomendaciones de stocks.

## Funcionalidades

- Consulta de stocks con filtrado y paginación
- Detalles de stocks específicos
- Generación de recomendaciones de inversión
- Verificaciones de salud del servicio

## Requisitos

- Go 1.23 o superior
- Docker (para desarrollo y despliegue)
- CockroachDB (compartido con el Stock Data Service)

## Configuración

### Variables de entorno

El servicio requiere las siguientes variables de entorno:

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| SERVER_PORT | Puerto en el que se ejecutará el servidor | 8080 |
| DB_HOST | Host de la base de datos | localhost |
| DB_PORT | Puerto de la base de datos | 26257 |
| DB_USER | Usuario de la base de datos | root |
| DB_PASSWORD | Contraseña de la base de datos | - |
| DB_NAME | Nombre de la base de datos | stockdb |
| DB_SSL_MODE | Modo SSL para la conexión a la base de datos | disable |

## Desarrollo local

### Ejecutar con Go

```bash
# Clonar el repositorio
git clone https://github.com/RobertCastro/stock-microservices-api.git
cd stock-microservices-api/stock-api-service

# Instalar dependencias
go mod download

# Ejecutar
go run cmd/api/main.go