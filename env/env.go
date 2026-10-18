package env

import (
	"fmt"
	"os"
	"strings"
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
	// CONTENT_S3_* credentials; unset, they go to AWS S3 via the default
	// credential chain. Either way the returned URL is ContentPublicBaseURL.
	ContentS3Endpoint        = getEnv("CONTENT_S3_ENDPOINT", "")
	ContentS3Region          = getEnv("CONTENT_S3_REGION", "us-west-1")
	ContentS3Bucket          = getEnv("CONTENT_S3_BUCKET", "content.hillview.tv")
	ContentS3AccessKeyID     = getEnv("CONTENT_S3_ACCESS_KEY_ID", "")
	ContentS3SecretAccessKey = getEnv("CONTENT_S3_SECRET_ACCESS_KEY", "")
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
