# Ejemplos de Uso de la API (Go Port)

Esta guía muestra cómo interactuar con la API usando `curl` o `httpie`.

---

## 👤 Dominio: Usuarios

### 1. Crear Usuario
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Elo Sanz",
    "email": "elo@example.com",
    "phone": "555-1234",
    "website": "elosanz.com"
  }'
```

### 2. Obtener Todos los Usuarios (paginado)
```bash
curl "http://localhost:8080/api/users?page=1&size=5"
```

### 3. Sincronizar desde API Externa (JSONPlaceholder)
```bash
curl -X POST http://localhost:8080/api/users/sync/1
```

---

## 🥒 Dominio: Rick and Morty

### 4. Buscar Personajes
```bash
curl "http://localhost:8080/api/rickandmorty/characters?name=rick&status=alive"
```

### 5. Obtener Personaje por ID
```bash
curl http://localhost:8080/api/rickandmorty/characters/1
```

---

## 📦 Dominio: Storage (S3)

### 6. Listar Archivos
```bash
curl http://localhost:8080/api/storage/files
```

### 7. Subir Archivo
```bash
curl -X POST http://localhost:8080/api/storage/upload \
  -F "file=@/ruta/a/tu/archivo.txt"
```

---

## 🔍 Inspección de Base de Datos (SQLite)

A diferencia de Java/H2, no hay una consola web activa por defecto, pero puedes inspeccionar el archivo de base de datos directamente si usas `DATABASE_PATH`.

Si usas `:memory:` (default), los datos se pierden al reiniciar. Para persistir localmente y debuguear:

1. Levanta la app con un path:
   ```bash
   DATABASE_PATH=./demo.db air
   ```
2. Usa el CLI de sqlite3:
   ```bash
   sqlite3 ./demo.db "SELECT * FROM users;"
   ```

## ⚡ Logging Diferencial

La aplicación incluye un middleware de logging que muestra cada request en la terminal donde corre `air`:

```text
2026/04/15 21:10:00 INFO request method=GET path=/api/users status=200 duration=1.2ms
```
