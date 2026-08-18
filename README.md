## Running with Docker

Requiere Docker Engine y el plugin Docker Compose en un servidor Linux. En esta máquina de desarrollo local **no** está instalado Docker.

Desde la raíz del repositorio:

```bash
cp .env.example .env
docker compose up -d --build
docker compose ps
docker compose logs -f
```

Endpoints en el host:

- Go: `http://localhost:8080`
- Node: `http://localhost:3000` (expuesto para pruebas; en producción podría quedar solo interno)

Dentro de la red de Compose, Go llama a Node en `http://node-api:3000`, no en `localhost`.

Detener:

```bash
docker compose down
```
