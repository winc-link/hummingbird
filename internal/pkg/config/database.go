package config

// DatabaseInfo interfaces provides an abstraction for obtaining the database configuration information.
type DatabaseInfo interface {
	// GetDatabaseInfo returns a database information map.
	GetDatabaseInfo() map[string]Database
}

type DatabaseConfig interface {
	Configuration
	DatabaseInfo
}
