package configs

type Postgres struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Http struct {
	Port string
}

type Logger struct {
	LogLevel string
}

// AppConfigs groups all application configurations.
type AppConfigs struct {
	Postgres Postgres
	Http     Http
	Logger   Logger
}
