# Guía: Cómo agregar una nueva Entidad (con GORM)

Sigue estos pasos para extender el proyecto con un nuevo dominio usando la "magia" de GORM.

---

## 1. Definir Entidad de Dominio (`internal/products/product.go`)

Usa `tags` para configurar el comportamiento del ORM.

```go
package products

import "time"

type Product struct {
    ID        int64     `gorm:"primaryKey"`
    Name      string    `gorm:"not null"`
    Price     float64   `gorm:"type:decimal(10,2)"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```

## 2. Definir Puerto (Repositorio) (`internal/products/repository.go`)

```go
package products

import "context"

type ProductRepository interface {
    Save(ctx context.Context, p Product) (Product, error)
}
```

## 3. Implementar Repo con GORM (`infrastructure/postgres/product_repository.go`)

```go
package postgres

import (
    "context"
    "github.com/elosanz/demo/internal/products"
    "gorm.io/gorm"
)

type ProductGORMRepository struct {
    db *gorm.DB
}

func (r *ProductGORMRepository) Save(ctx context.Context, p products.Product) (products.Product, error) {
    err := r.db.WithContext(ctx).Create(&p).Error
    return p, err
}
```

## 4. Registrar en el Entrypoint (`cmd/api/main.go`)

Simplemente agrega tu struct al `AutoMigrate` y GORM creará la tabla al arrancar:

```go
db.AutoMigrate(&user.User{}, &products.Product{})
```

---

## 🧪 Cómo testear (Test de Integración)

Crea un archivo `handler_integration_test.go` en el paquete del handler:
1. Usa `sqlite.Open(":memory:")` para una DB limpia.
2. Llama a `db.AutoMigrate()`.
3. Usa `httptest.NewRecorder()` para validar el flujo HTTP completo.
