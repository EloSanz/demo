# Ejemplos de Uso de la API (Go Port)

Esta guía muestra cómo interactuar con la API. Gracias a **GORM**, la base de datos se gestiona sola.

---

## 👤 Dominio: Usuarios

### 1. Crear Usuario (Auditoría Automática)
Al crear un usuario, GORM setea `created_at` y `updated_at` automáticamente.

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "GORM User", "email": "gorm@example.com"}'
```

---

## 🔍 Inspección y Debugging

### Logs de SQL
La aplicación está configurada para imprimir todas las queries SQL que GORM ejecuta en la terminal. Verás algo como:
```text
[0.452ms] [rows:1] INSERT INTO "users" ("name","email",...) VALUES (...)
```

### Inspeccionar SQLite
Si usas un archivo (default `demo.db`):
```bash
sqlite3 demo.db .tables
sqlite3 demo.db "SELECT * FROM users;"
```

---

## 🧪 Tests de Integración
Para validar que todo el stack funciona sin levantar la app manualmente:
```bash
# Corre todos los tests levantando DBs en memoria y mocks de APIs externas
go test ./... -v
```
