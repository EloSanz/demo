# Ejemplos de Uso de la API (Go Port)

## ✉️ Notificaciones Asincrónicas

Este endpoint pone un mensaje en cola (SQS o Memoria) para ser procesado por un worker en segundo plano.

```bash
curl -X POST http://localhost:8080/api/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "welcome_email",
    "content": "Hola Elo! Bienvenido al sistema asincrónico"
  }'
```
**Respuesta (202 Accepted):** `{"status":"queued"}`

---

## 📊 Observabilidad y Monitoreo

### 1. Métricas de Prometheus
```bash
curl http://localhost:8080/metrics
```

### 2. Health Check
```bash
curl http://localhost:8080/health
```

---

## 🛡️ Pruebas de Resiliencia

### 1. Test de Graceful Shutdown (Lento)
```bash
curl http://localhost:8080/api/test/slow
```
*Si apagas la app mientras corre, el worker de notificaciones y este request terminarán antes del cierre.*

---

## 🧪 Testing Unitario e Integración
```bash
go test ./... -v
```
