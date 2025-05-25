// Package interfaces provides all interface definitions for the application.
// This central package contains interfaces for all layers: domain, repository, usecase, platform, and config.
//
// Organization:
// - config.go: Configuration-related interfaces
// - domain.go: Domain entity interfaces
// - platform.go: Platform service interfaces (editor, etc.)
// - repository.go: Data access layer interfaces
// - usecase.go: Business logic layer interfaces
// - provider.go: Interface implementations
package interfaces

// Re-export all interfaces for convenience
// This allows consumers to import just "github.com/hirotoni/memov2/internal/interfaces"
// instead of importing specific files

// All interfaces are defined in their respective files and automatically available
// when importing this package. No re-export needed in Go.
