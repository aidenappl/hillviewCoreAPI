package env

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aidenappl/go-keyring"
)

var (
	Port             = getEnv("PORT", "8000")
	DatabaseDSN      = getEnvOrPanic("DATABASE_DSN")
	JWTSigningKey    = getEnvOrPanic("JWT_SIGNING_KEY")
	AssetDatabaseDSN = getEnvOrPanic("ASSET_DATABASE_DSN")
	SendgridAPIKey   = getEnvOrPanic("SENDGRID_API_KEY")
	HealthCheckURL   = getEnvOrPanic("HEALTH_CHECK_URL")

	// Content storage for image and thumbnail uploads. With CONTENT_S3_ENDPOINT
	// set, uploads go to that S3-compatible store (OpenBucket) with the
	// CONTENT_S3_* credentials (environment, else Keyring); unset, they go to AWS S3 via the default
	// credential chain. Either way the returned URL is ContentPublicBaseURL.
	ContentS3Endpoint        = getEnv("CONTENT_S3_ENDPOINT", "")
	ContentS3Region          = getEnv("CONTENT_S3_REGION", "us-west-1")
	ContentS3Bucket          = getEnv("CONTENT_S3_BUCKET", "content.hillview.tv")
	ContentS3AccessKeyID     = getSecret("CONTENT_S3_ACCESS_KEY_ID")
	ContentS3SecretAccessKey = getSecret("CONTENT_S3_SECRET_ACCESS_KEY")
	ContentPublicBaseURL     = strings.TrimRight(getEnv("CONTENT_PUBLIC_BASE_URL", "https://content.hillview.tv"), "/")

	// Monitor (zone hillview). All optional: unset, the service runs without
	// telemetry rather than failing to boot.
	MonitorIngestURL = getEnv("MONITOR_INGEST_URL", "")
	MonitorAPIKey    = getEnv("MONITOR_API_KEY", "")
	MonitorEnv       = getEnv("MONITOR_ENV", "production")
	MonitorSpoolDir  = getEnv("MONITOR_SPOOL_DIR", "")
	MonitorStdout    = getEnv("MONITOR_STDOUT", "false") == "true"
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

// keyringSecrets holds the secrets granted to this service's Keyring token,
// loaded once at startup. Empty when Keyring is not configured
// (KEYRING_URL / KEYRING_ACCESS_KEY_ID / KEYRING_SECRET_ACCESS_KEY unset),
// disabled (KEYRING_SOURCE=env), or unreachable.
var keyringSecrets = loadKeyring()

func loadKeyring() map[string]string {
	if keyring.EnvOnly() {
		return nil
	}
	client, err := keyring.New()
	if err != nil {
		log.Printf("keyring: not configured, using environment only: %v", err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	secrets, err := client.Load(ctx)
	if err != nil {
		log.Printf("keyring: load failed, using environment only: %v", err)
		return nil
	}
	log.Printf("keyring: loaded %d secrets", len(secrets))
	return secrets
}

// getSecret resolves an optional secret: the environment first (so local runs
// and the stack can override), then Keyring, else "".
func getSecret(key string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return keyringSecrets[key]
}
