# Ejemplos de Uso de la API (Go Port)

Esta guía muestra cómo interactuar con los nuevos endpoints de infraestructura y observación.

---

## 📊 Observabilidad y Monitoreo

### 1. Métricas de Prometheus
La aplicación expone métricas nativas para ser consumidas por un servidor de Prometheus.

```bash
curl http://localhost:8080/metrics
```
*Tip: Busca `http_requests_total` para ver cuántos hits recibió cada endpoint.*

### 2. Health Check
Verifica que la app esté viva y tenga conexión a la base de datos (Postgres o SQLite).

```bash
curl http://localhost:8080/health
```

---

## 🛡️ Pruebas de Resiliencia

### 1. Test de Graceful Shutdown
Para verificar que el apagado no corta conexiones activas:

1.  Llama al endpoint lento:
    ```bash
    curl http://localhost:8080/api/test/slow
    ```
2.  Apaga el servicio inmediatamente (ej: `docker-compose stop api`).
3.  El `curl` debería terminar con éxito antes de que la app se cierre.

---

## 🧪 Comandos de Testing

### Suite completa
```bash
go test ./... -v
```

### Solo un paquete (sin cache)
```bash
go test -count=1 ./internal/user/handler/... -v
```
