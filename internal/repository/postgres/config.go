package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     int
	MaxOpenConns     int           // Max number of open connections
	MaxIdleConns     int           // Max number of idle connections
	ConnMaxLifetime  time.Duration // Connection maximum lifetime
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
	)
}
