// Package interfaces provides all interface definitions for the application.
// This central package contains interfaces for all layers: domain, repository, service, platform, and config.
//
// Organization:
// - config.go: Configuration-related interfaces
// - domain.go: Domain entity interfaces
// - platform.go: Platform service interfaces (editor, etc.)
// - repository.go: Data access layer interfaces
// - service.go: Business logic layer interfaces (service layer)
//
// Note: Interface implementations are located in their respective packages
// (e.g., ConfigProvider implementation is in internal/config/toml).
package interfaces

// All interfaces are defined in their respective files and automatically available
// when importing this package. No re-export needed in Go.
