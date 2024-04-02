package env

import (
	"fmt"
	"os"
)

var (
	Port             = getEnv("PORT", "8000")
	DatabaseDSN      = getEnvOrPanic("DATABASE_DSN")
	JWTSigningKey    = getEnvOrPanic("JWT_SIGNING_KEY")
	AssetDatabaseDSN = getEnvOrPanic("ASSET_DATABASE_DSN")
	SendgridAPIKey   = getEnvOrPanic("SENDGRID_API_KEY")
)

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Missing required environment variable '%v'\n", key))
	}
	return value
}
