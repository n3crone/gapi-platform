# GAPI Platform

A lightweight, extensible REST API framework for Go built on top of Fiber and GORM.

[![Go Report Card](https://goreportcard.com/badge/github.com/n3crone/gapi-platform)](https://goreportcard.com/report/github.com/n3crone/gapi-platform)
[![GoDoc](https://godoc.org/github.com/n3crone/gapi-platform?status.svg)](https://godoc.org/github.com/n3crone/gapi-platform)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue)](https://go.dev/)


> ⚠️ **Alpha Status**: This framework is currently in alpha. While functional, it's not yet recommended for production use.

## Overview

GAPI Platform helps you build REST APIs quickly with a clean, resource-oriented architecture. It provides automatic CRUD endpoints while maintaining flexibility for customization.

## Features

- 🚀 Quick API setup with minimal boilerplate
- 📦 Automatic CRUD operations
- 🛠 Resource-based architecture
- 🔌 Built-in MySQL database integration via GORM
- 🎯 Type-safe request/response handling
- 📝 Structured logging with zerolog
- ⚡ High-performance web server using Fiber

## Requirements

- Go 1.21 or higher
- MySQL 5.7 or higher

## Installation

```bash
go get -u github.com/n3crone/gapi-platform
```

## Quick Start

1. Create a new project:

```bash
mkdir my-api && cd my-api
go mod init my-api
```

2. Create your first API:

```go
package main

import (
    "github.com/n3crone/gapi-platform/pkg/core"
    "github.com/n3crone/gapi-platform/pkg/resource"
)

type User struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name"`
}

func (a *User) CreateResource(rm *resource.ResourceManager) *resource.Resource {
    return rm.CreateResource(a)
}

func main() {
    app, err := core.New(core.Config{
        DatabaseUri: "user:pass@tcp(localhost:3306)/dbname",
    })
    if err != nil {
        panic(err)
    }
    
    if err := app.Migrate(&User{}); err != nil {
        panic(err)
    }

    app.RegisterResource(&User{})
    app.Fiber.Listen(":3000")
}
```

That's it! Your API now has the following endpoints:

- `GET /users` - List all users
- `POST /users` - Create a new user
- `GET /users/:id` - Get a specific user
- `PUT /users/:id` - Update a user
- `DELETE /users/:id` - Delete a user

## Project Structure

```bash
gapi-platform/
└── pkg/
    ├── core/        # Main application core
    ├── database/    # Database connectivity
    ├── resource/    # Resource management
    └── state/       # State providers and processors
```

## Configuration Options

```go
type Config struct {
    Fiber       *fiber.Config // Custom Fiber settings
    DatabaseUri string        // Database connection string
    LogLevel    zerolog.Level // Logging level
    LogFormat   string        // Log format (json/console)
}
```

## 🚧 Examples

For complete examples, check our demo repository:

- 🚧 [gapi-platform-examples](https://github.com/n3crone/gapi-platform-examples)

## 🚧 Documentation

- 🚧 [API Reference](https://godoc.org/github.com/n3crone/gapi-platform)
- 🚧 [Wiki](https://github.com/n3crone/gapi-platform/wiki)
- 🚧 Contributing Guide

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

Built with:

- [Fiber](https://github.com/gofiber/fiber)
- [GORM](https://gorm.io)
- [zerolog](https://github.com/rs/zerolog)
