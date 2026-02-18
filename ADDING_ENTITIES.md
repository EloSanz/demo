# 🚀 Guía Rápida: Agregar Nueva Entidad

Esta guía te muestra cómo agregar rápidamente una nueva entidad (ej: `Product`) a tu template.

## 📋 Checklist

- [ ] 1. Crear entidad en `domain/`
- [ ] 2. Crear DTOs en `dto/`
- [ ] 3. Crear repository en `repository/`
- [ ] 4. Crear service interface en `service/`
- [ ] 5. Crear service implementation en `service/impl/`
- [ ] 6. Crear controller en `controller/`
- [ ] 7. (Opcional) Agregar cliente HTTP si consume API externa

---

## 📝 Ejemplo: Entidad Product

### 1. Domain Entity (`domain/Product.java`)

```java
package com.example.demo.domain;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.math.BigDecimal;

@Entity
@Table(name = "products")
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Product {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(nullable = false)
    private String name;
    
    private String description;
    
    @Column(nullable = false)
    private BigDecimal price;
    
    private Integer stock;
}
```

### 2. DTOs (`dto/ProductRequest.java` y `dto/ProductResponse.java`)

**ProductRequest.java:**
```java
package com.example.demo.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.math.BigDecimal;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ProductRequest {
    
    @NotBlank(message = "Name is required")
    private String name;
    
    private String description;
    
    @NotNull(message = "Price is required")
    @Positive(message = "Price must be positive")
    private BigDecimal price;
    
    @Positive(message = "Stock must be positive")
    private Integer stock;
}
```

**ProductResponse.java:**
```java
package com.example.demo.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.math.BigDecimal;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ProductResponse {
    private Long id;
    private String name;
    private String description;
    private BigDecimal price;
    private Integer stock;
}
```

### 3. Repository (`repository/ProductRepository.java`)

```java
package com.example.demo.repository;

import com.example.demo.domain.Product;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import java.util.List;

@Repository
public interface ProductRepository extends JpaRepository<Product, Long> {
    
    List<Product> findByNameContainingIgnoreCase(String name);
    
    List<Product> findByPriceLessThan(BigDecimal maxPrice);
}
```

### 4. Service Interface (`service/ProductService.java`)

```java
package com.example.demo.service;

import com.example.demo.dto.ProductRequest;
import com.example.demo.dto.ProductResponse;
import java.util.List;

public interface ProductService {
    
    List<ProductResponse> getAllProducts();
    
    ProductResponse getProductById(Long id);
    
    ProductResponse createProduct(ProductRequest request);
    
    ProductResponse updateProduct(Long id, ProductRequest request);
    
    void deleteProduct(Long id);
    
    List<ProductResponse> searchByName(String name);
}
```

### 5. Service Implementation (`service/impl/ProductServiceImpl.java`)

```java
package com.example.demo.service.impl;

import com.example.demo.domain.Product;
import com.example.demo.dto.ProductRequest;
import com.example.demo.dto.ProductResponse;
import com.example.demo.repository.ProductRepository;
import com.example.demo.service.ProductService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
@Transactional
public class ProductServiceImpl implements ProductService {
    
    private final ProductRepository productRepository;
    
    @Override
    @Transactional(readOnly = true)
    public List<ProductResponse> getAllProducts() {
        log.info("Fetching all products");
        return productRepository.findAll().stream()
                .map(this::mapToResponse)
                .collect(Collectors.toList());
    }
    
    @Override
    @Transactional(readOnly = true)
    public ProductResponse getProductById(Long id) {
        log.info("Fetching product with id: {}", id);
        Product product = productRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Product not found with id: " + id));
        return mapToResponse(product);
    }
    
    @Override
    public ProductResponse createProduct(ProductRequest request) {
        log.info("Creating new product: {}", request.getName());
        
        Product product = Product.builder()
                .name(request.getName())
                .description(request.getDescription())
                .price(request.getPrice())
                .stock(request.getStock())
                .build();
        
        Product savedProduct = productRepository.save(product);
        log.info("Product created with id: {}", savedProduct.getId());
        
        return mapToResponse(savedProduct);
    }
    
    @Override
    public ProductResponse updateProduct(Long id, ProductRequest request) {
        log.info("Updating product with id: {}", id);
        
        Product product = productRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Product not found with id: " + id));
        
        product.setName(request.getName());
        product.setDescription(request.getDescription());
        product.setPrice(request.getPrice());
        product.setStock(request.getStock());
        
        Product updatedProduct = productRepository.save(product);
        log.info("Product updated with id: {}", updatedProduct.getId());
        
        return mapToResponse(updatedProduct);
    }
    
    @Override
    public void deleteProduct(Long id) {
        log.info("Deleting product with id: {}", id);
        
        if (!productRepository.existsById(id)) {
            throw new RuntimeException("Product not found with id: " + id);
        }
        
        productRepository.deleteById(id);
        log.info("Product deleted with id: {}", id);
    }
    
    @Override
    @Transactional(readOnly = true)
    public List<ProductResponse> searchByName(String name) {
        log.info("Searching products by name: {}", name);
        return productRepository.findByNameContainingIgnoreCase(name).stream()
                .map(this::mapToResponse)
                .collect(Collectors.toList());
    }
    
    private ProductResponse mapToResponse(Product product) {
        return ProductResponse.builder()
                .id(product.getId())
                .name(product.getName())
                .description(product.getDescription())
                .price(product.getPrice())
                .stock(product.getStock())
                .build();
    }
}
```

### 6. Controller (`controller/ProductController.java`)

```java
package com.example.demo.controller;

import com.example.demo.dto.ProductRequest;
import com.example.demo.dto.ProductResponse;
import com.example.demo.service.ProductService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import java.util.List;

@RestController
@RequestMapping("/api/products")
@RequiredArgsConstructor
@Slf4j
public class ProductController {
    
    private final ProductService productService;
    
    @GetMapping
    public ResponseEntity<List<ProductResponse>> getAllProducts() {
        log.info("GET /api/products - Fetching all products");
        List<ProductResponse> products = productService.getAllProducts();
        return ResponseEntity.ok(products);
    }
    
    @GetMapping("/{id}")
    public ResponseEntity<ProductResponse> getProductById(@PathVariable Long id) {
        log.info("GET /api/products/{} - Fetching product", id);
        ProductResponse product = productService.getProductById(id);
        return ResponseEntity.ok(product);
    }
    
    @PostMapping
    public ResponseEntity<ProductResponse> createProduct(@Valid @RequestBody ProductRequest request) {
        log.info("POST /api/products - Creating product: {}", request.getName());
        ProductResponse createdProduct = productService.createProduct(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(createdProduct);
    }
    
    @PutMapping("/{id}")
    public ResponseEntity<ProductResponse> updateProduct(
            @PathVariable Long id,
            @Valid @RequestBody ProductRequest request) {
        log.info("PUT /api/products/{} - Updating product", id);
        ProductResponse updatedProduct = productService.updateProduct(id, request);
        return ResponseEntity.ok(updatedProduct);
    }
    
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteProduct(@PathVariable Long id) {
        log.info("DELETE /api/products/{} - Deleting product", id);
        productService.deleteProduct(id);
        return ResponseEntity.noContent().build();
    }
    
    @GetMapping("/search")
    public ResponseEntity<List<ProductResponse>> searchProducts(@RequestParam String name) {
        log.info("GET /api/products/search?name={}", name);
        List<ProductResponse> products = productService.searchByName(name);
        return ResponseEntity.ok(products);
    }
}
```

---

## 🌐 Agregar Cliente HTTP para API Externa

### 1. Agregar configuración en `application.yml`

```yaml
external:
  api:
    store-api:
      base-url: https://fakestoreapi.com
```

### 2. Crear HTTP Interface (`client/StoreApiClient.java`)

```java
package com.example.demo.client;

import com.example.demo.dto.ProductResponse;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.service.annotation.GetExchange;
import org.springframework.web.service.annotation.HttpExchange;
import java.util.List;

@HttpExchange
public interface StoreApiClient {
    
    @GetExchange("/products")
    List<ProductResponse> getAllProducts();
    
    @GetExchange("/products/{id}")
    ProductResponse getProductById(@PathVariable Long id);
}
```

### 3. Configurar Bean en `HttpClientConfig.java`

```java
@Value("${external.api.store-api.base-url}")
private String storeApiBaseUrl;

@Bean
public StoreApiClient storeApiClient() {
    WebClient webClient = WebClient.builder()
            .baseUrl(storeApiBaseUrl)
            .build();
    
    HttpServiceProxyFactory factory = HttpServiceProxyFactory
            .builderFor(WebClientAdapter.create(webClient))
            .build();
    
    return factory.createClient(StoreApiClient.class);
}
```

### 4. Usar en el Service

```java
@Service
@RequiredArgsConstructor
public class ProductServiceImpl implements ProductService {
    
    private final ProductRepository productRepository;
    private final StoreApiClient storeApiClient; // Inyectar
    
    public List<ProductResponse> fetchFromExternalApi() {
        log.info("Fetching products from external API");
        return storeApiClient.getAllProducts();
    }
}
```

---

## 🎯 Tips para Take-Home Challenges

### 1. **Cambiar Base de Datos**

**Para PostgreSQL:**
```yaml
# application.yml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: postgres
    password: yourpassword
```

```gradle
// build.gradle
runtimeOnly 'org.postgresql:postgresql'
```

### 2. **Agregar Paginación**

```java
// Repository
Page<Product> findAll(Pageable pageable);

// Service
Page<ProductResponse> getAllProducts(int page, int size);

// Controller
@GetMapping
public ResponseEntity<Page<ProductResponse>> getAllProducts(
    @RequestParam(defaultValue = "0") int page,
    @RequestParam(defaultValue = "10") int size) {
    return ResponseEntity.ok(productService.getAllProducts(page, size));
}
```

### 3. **Agregar Swagger/OpenAPI**

```gradle
// build.gradle
implementation 'org.springdoc:springdoc-openapi-starter-webmvc-ui:2.3.0'
```

Accede a: `http://localhost:8080/swagger-ui.html`

### 4. **Agregar Tests**

```java
@SpringBootTest
@AutoConfigureMockMvc
class ProductControllerTest {
    
    @Autowired
    private MockMvc mockMvc;
    
    @Test
    void shouldCreateProduct() throws Exception {
        mockMvc.perform(post("/api/products")
                .contentType(MediaType.APPLICATION_JSON)
                .content("{\"name\":\"Test\",\"price\":10.0}"))
                .andExpect(status().isCreated());
    }
}
```

---

## ⚡ Comandos Rápidos

```bash
# Compilar
./gradlew build

# Ejecutar
./gradlew bootRun

# Tests
./gradlew test

# Limpiar y compilar
./gradlew clean build
```

---

## 📚 Recursos Adicionales

- [Spring Boot 4 Documentation](https://docs.spring.io/spring-boot/docs/current/reference/html/)
- [HTTP Interface Guide](https://docs.spring.io/spring-framework/reference/integration/rest-clients.html#rest-http-interface)
- [Spring Data JPA](https://docs.spring.io/spring-data/jpa/docs/current/reference/html/)
