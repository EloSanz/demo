# 🚀 FastAPI Hexagonal Observability Boilerplate

Un boilerplate de grado producción para microservicios de alto rendimiento basados en **FastAPI**, **Arquitectura Hexagonal (Puertos y Adaptadores)** y un stack de **Observabilidad Nativo de Nube** utilizando **Prometheus**, **Loki**, **Promtail** y **Grafana**.

---

## 🏗️ Arquitectura del Proyecto

El proyecto sigue una estructura limpia basada en **Arquitectura Hexagonal**, separando estrictamente las reglas de negocio de la infraestructura:

```text
├── src/
│   ├── core/                        # Funcionalidades transversales (Core)
│   │   ├── logging/                 # Configuración de logs estructurados JSON
│   │   └── middleware/              # Middleware de logs e inyección de request_id
│   │
│   ├── postback/                    # Dominio de Postbacks (Negocio)
│   │   ├── ports/                   # Interfaces abstractas (Puertos)
│   │   │   ├── queue_port.py
│   │   │   └── transaction_repository_port.py
│   │   ├── exceptions.py            # Excepciones de negocio dedicadas
│   │   ├── handlers.py              # Manejadores HTTP para mapeo de excepciones
│   │   ├── schemas.py               # Modelos Pydantic de entrada/salida
│   │   ├── service.py               # Casos de uso de negocio (Aplicación)
│   │   ├── router.py                # Controladores / Endpoints FastAPI
│   │   └── dependencies.py          # Inyección de dependencias de puertos
│   │
│   └── infrastructure/              # Implementaciones tecnológicas (Adaptadores)
│       ├── queue/                   # Adaptador de colas (e.g., SQS)
│       └── repository/              # Adaptador de almacenamiento (e.g., In-Memory)
│
├── tests/                           # Suite de pruebas unitarias e integración
└── pyproject.toml                   # Gestión de dependencias y scripts con uv
```

---

## ⚡ Requisitos Previos

Asegúrate de tener instalados los siguientes componentes antes de comenzar:
- **Python 3.13 o 3.14**
- **[uv](https://github.com/astral-sh/uv)** (Administrador de dependencias ultra-rápido en Rust)
- **Docker & Docker Compose**

---

## 🛠️ Configuración Rápida en Local

### 1. Sincronizar el entorno de desarrollo
Clona el repositorio e inicializa el entorno virtual con todas las dependencias:
```bash
uv sync
```
*Esto creará la carpeta `.venv` y descargará los paquetes necesarios en milisegundos.*

### 2. Correr las pruebas unitarias
```bash
uv run task test
```

### 3. Verificar la cobertura de código (Excluyendo boilerplate de infraestructura)
```bash
uv run task coverage
```

---

## 🐳 Despliegue con Docker Compose (Observabilidad Completa)

El proyecto incluye un entorno integrado de observabilidad local. Para levantarlo:

```bash
docker-compose up --build -d
```

Esto desplegará **5 contenedores** interconectados:
1. **`postback_api` (FastAPI)**: Disponible en `http://localhost:8000`
2. **`prometheus`**: Disponible en `http://localhost:9090` (Scrapea métricas de rendimiento).
3. **`loki`**: Disponible en `http://localhost:3100` (Almacenamiento de logs centralizados).
4. **`promtail`**: Colector que lee los logs estructurados y los envía a Loki en tiempo real.
5. **`grafana`**: Disponible en `http://localhost:3000` (Visualización de tableros).

### Probar Endpoints de la API
* **Verificación de Salud:**
  ```bash
  curl -i http://localhost:8000/test
  ```
* **Registrar conversión exitosa (iOS Determinístico):**
  ```bash
  curl -i -X POST http://localhost:8000/conversions \
    -H "Content-Type: application/json" \
    -d '{"transaction_id": "tx_xyz", "campaign_id": "camp_abc", "os": "ios", "device_id": "idfa-999", "revenue": 15.50}'
  ```

---

## 🔍 Logs JSON Estructurados

La API emite logs estructurados nativos de producción. Al ejecutar `docker logs postback_api`, verás los logs en perfecto formato JSON auto-descriptivo con trazabilidad distribuida por `request_id`:

```json
{"method": "GET", "path": "/test", "event": "request_started", "request_id": "f2ca62f0-47f2-47a4-968f-b22d752708f2", "level": "info", "timestamp": "2026-05-17T20:37:09.831176Z"}
{"duration": 0.003206, "status": 200, "event": "request_finished", "request_id": "f2ca62f0-47f2-47a4-968f-b22d752708f2", "level": "info", "timestamp": "2026-05-17T20:37:09.834151Z"}
```

---

## 📈 Consultas en Grafana / Loki (LogQL)

Para explorar tus logs estructurados en el explorador de Grafana (`http://localhost:3000`), puedes utilizar los siguientes filtros **LogQL**:

* **Filtrar por contenedor:**
  ```logql
  {container="postback_api"}
  ```
* **Parsear logs JSON estructurados:**
  ```logql
  {container="postback_api"} | json
  ```
* **Filtrar peticiones con errores (Status >= 400):**
  ```logql
  {container="postback_api"} | json | status >= 400
  ```
* **Buscar por transacción específica:**
  ```logql
  {container="postback_api"} | json | transaction_id = "tx_xyz"
  ```
