---
description: Proceso estándar para aplicar cambios y verificar la integridad del proyecto.
---

// turbo-all
1. Ejecutar `./gradlew spotlessApply` para asegurar el formato correcto.
2. Ejecutar `./gradlew check` para correr todos los tests (unitarios e integración) y validaciones estáticas.
3. Verificar que el resultado sea BUILD SUCCESSFUL.
