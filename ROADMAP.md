# Gapi Platform Roadmap

## Phase 1: Core Features

### 1.2 Resource Enhancement

Status: ⚪ Not Started

- Add JSON serialization groups

```go
type ResourceMetadata struct {
    Groups      map[string][]string
    Validation  map[string]string
    Middleware  []fiber.Handler
}
```

### 1.3 Request/Response Handling

Status: ⚪ Not Started

- Standardize error responses (RFC 7807)
- Input validation using validator/v10
- Middleware interface and registry
  - Create middleware abstraction layer
  - Provide examples
  - Documentation

### 1.4 Pagination & Filtering

Status: ⚪ Not Started

- Offset/limit pagination
- Filter query parser
- Sort query parser

### 1.5 Basic Caching

Status: ⚪ Not Started

- In-memory cache using go-cache
- HTTP cache headers support
- Cache invalidation on write operations

## Phase 2: Developer Experience

### 2.1 OpenAPI Documentation

Status: ⚪ Not Started

- Generate OpenAPI specs from resources
- Swagger UI integration
- Resource metadata for documentation

### 2.2 CLI Tools

Status: ⚪ Not Started

- Resource generator
- Project scaffolding
- Configuration generator

## Phase 3: Events & Real-time

### 3.1 Event System

Status: ⚪ Not Started

```go
type Event struct {
    Type       string
    Resource   string
    ID         string
    Payload    interface{}
    Timestamp  time.Time
}
```

### 3.2 ⚪ NATS Integration


- NATS connection manager
- Resource event publishers
- Subscription handlers
