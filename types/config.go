package types

// PostgreSQL
type DatabaseConfig struct {
	SQLHost     string
	SQLPort     uint16
	SQLUsername string
	SQLPassword string
	SQLDatabase string
	MinConns    int32
	MaxConns    int32
}

type APIConfig struct {
	ListenAddress string
	ListenPort    uint16
}
