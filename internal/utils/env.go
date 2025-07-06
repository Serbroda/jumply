package utils

import (
	"flag"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"sync"
)

var (
	once        sync.Once
	environment string
)

func LoadEnv() {
	environment = os.Getenv("ENV")
	if environment == "" {
		environment = "development"
	}
	_ = godotenv.Load(".env." + environment + ".local")
	if environment != "test" {
		_ = godotenv.Load(".env.local")
	}
	_ = godotenv.Load(".env." + environment)
	_ = godotenv.Load()
}

func getEnv(key string) (string, bool) {
	once.Do(LoadEnv)
	return os.LookupEnv(key)
}

func getEnvWithFallback(key, fallback string) string {
	if value, ok := getEnv(key); ok {
		return value
	}
	return fallback
}

func parseEnvInt(value string, bitSize int) (int64, error) {
	return strconv.ParseInt(value, 10, bitSize)
}

func mustParseInt(value string, bitSize int) int64 {
	intValue, err := parseEnvInt(value, bitSize)
	if err != nil {
		panic("invalid integer value: " + value)
	}
	return intValue
}

func mustParseBool(value string) bool {
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	panic("invalid boolean value: " + value)
}

func GetStringFallback(key, fallback string) string {
	return getEnvWithFallback(key, fallback)
}

func GetBoolFallback(key string, fallback bool) bool {
	if value, ok := getEnv(key); ok {
		return value == "true"
	}
	return fallback
}

func GetInt32Fallback(key string, fallback int) int {
	if value, ok := getEnv(key); ok {
		if intValue, err := parseEnvInt(value, 32); err == nil {
			return int(intValue)
		}
	}
	return fallback
}

func GetInt64Fallback(key string, fallback int64) int64 {
	if value, ok := getEnv(key); ok {
		if intValue, err := parseEnvInt(value, 64); err == nil {
			return intValue
		}
	}
	return fallback
}

func MustGetString(key string) string {
	return mustGetEnv(key)
}

func MustGetBool(key string) bool {
	return mustParseBool(mustGetEnv(key))
}

func MustGetInt32(key string) int32 {
	return int32(mustParseInt(mustGetEnv(key), 32))
}

func MustGetInt64(key string) int64 {
	return mustParseInt(mustGetEnv(key), 64)
}

func mustGetEnv(key string) string {
	value, ok := getEnv(key)
	if !ok {
		panic("mandatory env " + key + " not found")
	}
	return value
}

// StringFlag bindet Long- und optional Short-Variante auf dieselbe Variable.
func StringFlag(ptr *string, defaultVal, name, short, usage string, mandatory bool) {
	flag.StringVar(ptr, name, defaultVal, usage)
	if short != "" {
		flag.StringVar(ptr, short, defaultVal, usage+" (shorthand)")
	}
	if mandatory && ptr == nil {
		panic("mandatory flag " + name + " not found")
	}
}

// IntFlag bindet Long- und optional Short-Variante für Integer.
func IntFlag(ptr *int, defaultVal int, name, short, usage string, mandatory bool) {
	flag.IntVar(ptr, name, defaultVal, usage)
	if short != "" {
		flag.IntVar(ptr, short, defaultVal, usage+" (shorthand)")
	}
	if mandatory && ptr == nil {
		panic("mandatory flag " + name + " not found")
	}
}

// BoolFlag bindet Long- und optional Short-Variante für Boolean.
func BoolFlag(ptr *bool, defaultVal bool, name, short, usage string, mandatory bool) {
	flag.BoolVar(ptr, name, defaultVal, usage)
	if short != "" {
		flag.BoolVar(ptr, short, defaultVal, usage+" (shorthand)")
	}
	if mandatory && ptr == nil {
		panic("mandatory flag " + name + " not found")
	}
}
