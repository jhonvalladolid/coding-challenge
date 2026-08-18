# Decisiones técnicas

## ADR-001 — QR como operación principal

El documento del challenge menciona rotación en la arquitectura y factorización QR en la funcionalidad requerida. Se implementó QR porque está nombrado de forma explícita como requisito funcional. La rotación queda fuera de esta entrega.

## ADR-002 — Arquitectura stateless

No hay persistencia. No se usaron Sequelize, GORM ni PostgreSQL.

## ADR-003 — Go como orquestador

El cliente solo llama a Go. Node calcula estadísticas sobre Q y R. Cada API permanece independiente y testeable.

## ADR-004 — Comunicación HTTP

Go usa `net/http` con timeout configurable. No hay colas ni gRPC. Los errores de red, timeout y contrato se traducen a 503/504/502.

## ADR-005 — Gonum para QR

`gonum.org/v1/gonum/mat.QR` exige `m >= n`. Q es orthonormal m×m y R es m×n. Si `m < n`, la API responde 422 `UNSUPPORTED_MATRIX_DIMENSIONS` sin panic.
