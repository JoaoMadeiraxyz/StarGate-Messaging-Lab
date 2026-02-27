package env

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// LoadEnv carrega as variáveis de ambiente do arquivo .env
func LoadEnv() error {
	// Tenta carregar o arquivo .env. Se não existir, continua sem erro
	_ = godotenv.Load()

	return nil
}

// GetString retorna o valor de uma variável de ambiente ou o valor padrão
func GetString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetBool retorna o valor booleano de uma variável de ambiente ou o valor padrão
func GetBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolVal, err := strconv.ParseBool(value)
		if err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// GetInt retorna o valor inteiro de uma variável de ambiente ou o valor padrão
func GetInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intVal, err := strconv.Atoi(value)
		if err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetUint64 retorna o valor uint64 de uma variável de ambiente ou o valor padrão
func GetUint64(key string, defaultValue uint64) uint64 {
	if value, exists := os.LookupEnv(key); exists {
		uint64Val, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			return uint64Val
		}
	}
	return defaultValue
}

// GetDuration retorna o valor de duração de uma variável de ambiente ou o valor padrão
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		durationVal, err := time.ParseDuration(value)
		if err == nil {
			return durationVal
		}
	}
	return defaultValue
}
