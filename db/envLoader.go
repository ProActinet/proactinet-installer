package db

import (
	"log"

	"github.com/joho/godotenv"
)

func loadEnv() map[string]string {
	envVars, err := godotenv.Read()
	if err != nil {
		log.Println("Warning: No .env file found or unable to load it")
	}
	return envVars
}

func EnvLoader() string {
	log.Println("Loading environment variables...")
	envVars := loadEnv()
	return envVars["DSN"]
	// You can now use envVars which contains all environment variables.
	// For example: log.Println(envVars)
}