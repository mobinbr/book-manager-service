package config

type Config struct {
	Database struct {
		Host     string `env:"DATABASE_HOST" env-default:"localhost" env-description:"Database host for service"`
		Port     int    `env:"DATABASE_PORT" env-default:"5432" env-description:"Database port for service"`
		Name     string `env:"DATABASE_NAME" env-default:"book_manager_db" env-description:"Database name for service"`
		Username string `env:"DATABASE_USERNAME" env-default:"postgres" env-description:"Database username for service"`
		Password string `env:"DATABASE_PASSWORD" env-default:"postgres" env-description:"Database password for service"`
	}
}
