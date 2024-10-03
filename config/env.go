package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Config struct {
	PublicHost string
	Port string
	
	DBUser string
	DBPassword string
	DBAddress string
	DBName string

	JWTExpirationInSeconds int64
	JWTSecret string
}

// create singleton variable to store the config
var Envs = initConfig()

func initConfig() Config {


	godotenv.Load()

	return Config {
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port: getEnv("PORT", "8080"),
		DBUser: getEnv("DB_USER", "remote-admin"),
		DBPassword: getEnv("DB_PASSWORD", "droltbd2"),
		DBAddress: fmt.Sprintf("%s:%s", getEnv("DB_HOST", "34.71.193.30"), getEnv("DB_PORT", "3306")),

		DBName: getEnv("DB_NAME", "ecom"),
		JWTSecret: getEnv("JWT_SECRET", "not-secret-secret-anymore?"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}