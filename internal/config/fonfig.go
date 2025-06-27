package config

type Config struct {
	PgDsn string `env:"PG_DSN" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}
