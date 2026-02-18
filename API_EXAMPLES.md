# 🧪 API Testing Examples

Ejemplos de requests para probar todos los endpoints del template.

## 📋 User Endpoints

### 1. Crear Usuario
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "555-1234",
    "website": "johndoe.com"
  }'
```

**Respuesta esperada:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "555-1234",
  "website": "johndoe.com"
}
```

---

### 2. Obtener Todos los Usuarios
```bash
curl http://localhost:8080/api/users
```

**Respuesta esperada:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "555-1234",
    "website": "johndoe.com"
  }
]
```

---

### 3. Obtener Usuario por ID
```bash
curl http://localhost:8080/api/users/1
```

---

### 4. Actualizar Usuario
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "email": "john@example.com",
    "phone": "555-9999",
    "website": "johnupdated.com"
  }'
```

---

### 5. Eliminar Usuario
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

**Respuesta esperada:** `204 No Content`

---

## 🌐 External API Endpoints

### 6. Obtener Usuarios de API Externa (JSONPlaceholder)
```bash
curl http://localhost:8080/api/users/external
```

**Respuesta esperada:**
```json
[
  {
    "id": 1,
    "name": "Leanne Graham",
    "email": "Sincere@april.biz",
    "phone": "1-770-736-8031 x56442",
    "website": "hildegard.org"
  },
  ...
]
```

---

### 7. Sincronizar Usuario de API Externa a BD Local
```bash
curl -X POST http://localhost:8080/api/users/sync/1
```

**Respuesta esperada:**
```json
{
  "id": 1,
  "name": "Leanne Graham",
  "email": "Sincere@april.biz",
  "phone": "1-770-736-8031 x56442",
  "website": "hildegard.org"
}
```

---

## ❌ Error Handling Examples

### Validación - Email Inválido
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "invalid-email",
    "phone": "555-1234"
  }'
```

**Respuesta esperada:**
```json
{
  "timestamp": "2026-02-17T12:56:37.123",
  "status": 400,
  "errors": {
    "email": "Email must be valid"
  }
}
```

---

### Campo Requerido Faltante
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "555-1234"
  }'
```

**Respuesta esperada:**
```json
{
  "timestamp": "2026-02-17T12:56:37.123",
  "status": 400,
  "errors": {
    "name": "Name is required",
    "email": "Email is required"
  }
}
```

---

### Usuario No Encontrado
```bash
curl http://localhost:8080/api/users/999
```

**Respuesta esperada:**
```json
{
  "timestamp": "2026-02-17T12:56:37.123",
  "status": 500,
  "message": "User not found with id: 999"
}
```

---

## 🧪 Testing con HTTPie (alternativa más legible)

Si tienes HTTPie instalado (`brew install httpie`):

### Crear Usuario
```bash
http POST localhost:8080/api/users \
  name="Jane Doe" \
  email="jane@example.com" \
  phone="555-5678" \
  website="janedoe.com"
```

### Obtener Todos
```bash
http GET localhost:8080/api/users
```

### Actualizar
```bash
http PUT localhost:8080/api/users/1 \
  name="Jane Updated" \
  email="jane@example.com" \
  phone="555-9999"
```

### Eliminar
```bash
http DELETE localhost:8080/api/users/1
```

---

## 🔍 Acceso a H2 Console

1. Abre tu navegador en: `http://localhost:8080/h2-console`
2. Configura:
   - **JDBC URL:** `jdbc:h2:mem:testdb`
   - **Username:** `sa`
   - **Password:** *(dejar vacío)*
3. Haz clic en "Connect"

### Queries SQL de Ejemplo

```sql
-- Ver todos los usuarios
SELECT * FROM users;

-- Buscar por email
SELECT * FROM users WHERE email LIKE '%example%';

-- Contar usuarios
SELECT COUNT(*) FROM users;

-- Eliminar todos los usuarios
DELETE FROM users;
```

---

## 📊 Postman Collection (JSON)

Puedes importar esta colección en Postman:

```json
{
  "info": {
    "name": "Spring Boot Template API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Create User",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"John Doe\",\n  \"email\": \"john@example.com\",\n  \"phone\": \"555-1234\",\n  \"website\": \"johndoe.com\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/users",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "users"]
        }
      }
    },
    {
      "name": "Get All Users",
      "request": {
        "method": "GET",
        "url": {
          "raw": "http://localhost:8080/api/users",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "users"]
        }
      }
    },
    {
      "name": "Get User by ID",
      "request": {
        "method": "GET",
        "url": {
          "raw": "http://localhost:8080/api/users/1",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "users", "1"]
        }
      }
    },
    {
      "name": "Update User",
      "request": {
        "method": "PUT",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"John Updated\",\n  \"email\": \"john@example.com\",\n  \"phone\": \"555-9999\",\n  \"website\": \"johnupdated.com\"\n}"
        },
        "url": {
          "raw": "http://localhost:8080/api/users/1",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "users", "1"]
        }
      }
    },
    {
      "name": "Delete User",
      "request": {
        "method": "DELETE",
        "url": {
          "raw": "http://localhost:8080/api/users/1",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "users", "1"]
        }
      }
    },
    {
      "name": "Get External Users",
      "request": {
        "method": "GET",
        "url": {
          "raw": "http://localhost:8080/api/users/external",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "users", "external"]
        }
      }
    },
    {
      "name": "Sync User from External API",
      "request": {
        "method": "POST",
        "url": {
          "raw": "http://localhost:8080/api/users/sync/1",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8080",
          "path": ["api", "users", "sync", "1"]
        }
      }
    }
  ]
}
```

---

## 🎯 Quick Test Script

Guarda esto como `test-api.sh` y ejecútalo:

```bash
#!/bin/bash

echo "🧪 Testing Spring Boot Template API"
echo "===================================="
echo ""

echo "1️⃣ Creating user..."
curl -s -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","phone":"555-0000","website":"test.com"}' \
  | jq .

echo ""
echo "2️⃣ Getting all users..."
curl -s http://localhost:8080/api/users | jq .

echo ""
echo "3️⃣ Getting external users..."
curl -s http://localhost:8080/api/users/external | jq '. | length'

echo ""
echo "4️⃣ Syncing user from external API..."
curl -s -X POST http://localhost:8080/api/users/sync/2 | jq .

echo ""
echo "✅ Tests completed!"
```

Ejecutar:
```bash
chmod +x test-api.sh
./test-api.sh
```

---

## 📝 Notas

- Todos los ejemplos asumen que la aplicación está corriendo en `http://localhost:8080`
- Los ejemplos usan `curl`, pero puedes usar Postman, Insomnia, o cualquier cliente HTTP
- La base de datos H2 es in-memory, por lo que los datos se pierden al reiniciar la aplicación
- Para producción, cambia a PostgreSQL/MySQL y configura persistencia
