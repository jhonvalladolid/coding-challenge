# Arquitectura

## Componentes

- **go-api** (Fiber): borde público. Valida, factoriza QR y llama a Node.
- **node-api** (Express): servicio interno. Estadísticas y diagonalidad.

No hay base de datos, colas ni caché.

## Comunicación

```text
Cliente
  POST /api/v1/matrices/qr
    → go-api
         → Validate
         → Gonum QR
         → POST {STATISTICS_API_URL}/api/v1/statistics
              X-Request-ID, Content-Type: application/json
         → node-api
    ← { originalMatrix, factorization, statistics, diagonal, meta }
```

`STATISTICS_API_URL` es solo la base (`http://localhost:3000` o `http://node-api:3000`). El path `/api/v1/statistics` lo añade el client de Go. El cliente HTTP no elige la URL de Node (evita SSRF).

## Contenedores

Misma red de Compose. Healthchecks independientes sobre `GET /health`. Go usa `depends_on` con `condition: service_healthy`, pero también maneja Node caído con 503.

## Trazas

`X-Request-ID` se genera o reutiliza en Go y se propaga a Node.
