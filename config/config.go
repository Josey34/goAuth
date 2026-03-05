package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           int
	DBPath         string
	JWTSecret      string
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	BcryptCost     int
	AllowedOrigins []string
	LogLevel       string
	RateLimitRPS   float64
	RateLimitBurst int
}

func Load() (*Config, error) {
	godotenv.Load()
	dbPath := os.Getenv("DB_PATH")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	accessTTL, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	refreshTTL, err := parseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}

	bcryptCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil {
		return nil, err
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	rateLimitRPS, err := strconv.ParseFloat(os.Getenv("RATE_LIMIT_RPS"), 64)
	if err != nil || rateLimitRPS == 0 {
		rateLimitRPS = 10
	}

	rateLimitBurst, err := strconv.Atoi(os.Getenv("RATE_LIMIT_BURST"))
	if err != nil || rateLimitBurst == 0 {
		rateLimitBurst = 20
	}

	return &Config{
		Port:           port,
		DBPath:         dbPath,
		JWTSecret:      jwtSecret,
		AccessTTL:      accessTTL,
		RefreshTTL:     refreshTTL,
		BcryptCost:     bcryptCost,
		AllowedOrigins: strings.Split(allowedOrigins, ","),
		LogLevel:       logLevel,
		RateLimitRPS:   rateLimitRPS,
		RateLimitBurst: rateLimitBurst,
	}, nil
}

func parseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}
