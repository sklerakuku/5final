package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabasePath         string
	ServerPort           string
	JWTSecret            string
	GRPCAddress          string
	ComputingPower       int
	TimeAdditionMS       int
	TimeSubtractionMS    int
	TimeMultiplicationMS int
	TimeDivisionMS       int
}

func Load() *Config {
	return &Config{
		DatabasePath:         getEnv("DATABASE_PATH", "./calculator.db"),
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		JWTSecret:            getEnv("JWT_SECRET", "secret"),
		GRPCAddress:          getEnv("GRPC_ADDRESS", "localhost:50051"),
		ComputingPower:       getEnvAsInt("COMPUTING_POWER", 2),
		TimeAdditionMS:       getEnvAsInt("TIME_ADDITION_MS", 1000),
		TimeSubtractionMS:    getEnvAsInt("TIME_SUBTRACTION_MS", 1000),
		TimeMultiplicationMS: getEnvAsInt("TIME_MULTIPLICATION_MS", 1000),
		TimeDivisionMS:       getEnvAsInt("TIME_DIVISION_MS", 1000),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
