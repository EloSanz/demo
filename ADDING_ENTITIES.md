# Guía: Cómo agregar una nueva Entidad (Go Clean Architecture)

Sigue estos pasos para extender el proyecto con un nuevo dominio (ej. `Products`).

---

## 1. Definir Entidad de Dominio (`internal/products/product.go`)

Define la estructura básica y errores de negocio.

```go
package products

import "errors"

type Product struct {
    ID          int64
    Name        string
    Price       float64
}

var ErrProductNotFound = errors.New("product not found")
```

## 2. Definir Puertos (Interfaces) (`internal/products/repository.go`)

```go
package products

import "context"

type ProductRepository interface {
    Save(ctx context.Context, p Product) (Product, error)
    FindByID(ctx context.Context, id int64) (*Product, error)
}
```

## 3. Implementar Adaptador (Infraestructura) (`infrastructure/postgres/product_repository.go`)

```go
package postgres

import (
    "context"
    "database/sql"
    "github.com/elosanz/demo/internal/products"
)

type ProductSQLiteRepository struct {
    db *sql.DB
}

func NewProductSQLiteRepository(db *sql.DB) products.ProductRepository {
    return &ProductSQLiteRepository{db: db}
}

func (r *ProductSQLiteRepository) Save(ctx context.Context, p products.Product) (products.Product, error) {
    // Implementar ejecución SQL...
}

func (r *ProductSQLiteRepository) FindByID(ctx context.Context, id int64) (*products.Product, error) {
    // Implementar ejecución SQL...
}
```

## 4. Definir e Implementar Service (`internal/products/service.go`)

```go
package products

import "context"

type ProductService interface {
    Create(ctx context.Context, name string, price float64) (Product, error)
}

type productService struct {
    repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
    return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, name string, price float64) (Product, error) {
    return s.repo.Save(ctx, Product{Name: name, Price: price})
}
```

## 5. Implementar Handler y DTOs (`internal/products/handler/`)

Crea `dto.go` para los requests/responses y `handler.go` para la lógica HTTP usando `pkg/web`.

## 6. Registrar en el Entrypoint (`cmd/api/main.go`)

1. Agrega la migración SQL en `runMigrations`.
2. Instancia Repository → Service → Handler.
3. Registra la ruta: `mux.HandleFunc("POST /api/products", web.Adapt(productH.Create))`.

---

## 💡 Tips Go vs Java

- **No hay Annotations**: Todo es explícito (SQL, JSON tags).
- **Interfaces en el Consumidor**: Las interfaces (`ProductRepository`) viven en la carpeta de dominio (`internal/products`), no con la implementación.
- **Error Handling**: Siempre retorna `error` y úsalo con `errors.Is` en los handlers.
