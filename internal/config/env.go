package config

import (
	"github.com/Tsisar/extended-log-go/log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Warnf("%s not found in environment variables, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

func getStringSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		log.Warnf("%s not found in environment variables, using default: %v", key, defaultValue)
		return defaultValue
	}
	return strings.Split(value, ",")
}

func getUint64(key string, defaultValue uint64) uint64 {
	value := os.Getenv(key)
	if value == "" {
		log.Warnf("%s not found in environment variables, using default: %d", key, defaultValue)
		return defaultValue
	}
	// First try to parse as uint64
	if parsedValue, err := strconv.ParseUint(value, 10, 64); err == nil {
		return parsedValue
	}
	// If that fails, try to parse as float64 and then cast to uint64
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return uint64(floatValue)
	} else {
		log.Warnf("Failed to parse %s, using default: %d: %v", key, defaultValue, err)
		return defaultValue
	}
}

func getInt64(key string, defaultValue int) int64 {
	value := os.Getenv(key)
	if value == "" {
		log.Warnf("%s not found in environment variables, using default: %d", key, defaultValue)
		return int64(defaultValue)
	}
	// First try to parse as int64
	if parsedValue, err := strconv.ParseInt(value, 10, 64); err == nil {
		return parsedValue
	}
	// If that fails, try to parse as float64 and then cast to int64
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return int64(floatValue)
	} else {
		log.Warnf("Failed to parse %s, using default: %d: %v", key, defaultValue, err)
		return int64(defaultValue)
	}
}

func getInt(key string, defaultValue int) int {
	return int(getInt64(key, defaultValue))
}

func getFloat64(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		log.Warnf("%s not found in environment variables, using default: %f", key, defaultValue)
		return defaultValue
	}
	// First try to parse as float64
	if parsedValue, err := strconv.ParseFloat(value, 64); err == nil {
		return parsedValue
	}
	// If that fails, try to parse as int64 and then cast to float64
	if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
		return float64(intValue)
	} else {
		log.Warnf("Failed to parse %s, using default: %f: %v", key, defaultValue, err)
		return defaultValue
	}
}

func getBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		log.Warnf("%s not found in environment variables, using default: %t", key, defaultValue)
		return defaultValue
	}
	// First try to parse as bool
	if parsedValue, err := strconv.ParseBool(value); err == nil {
		return parsedValue
	}
	// If that fails, try to parse as int and then cast to bool
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue != 0
	} else {
		log.Warnf("Failed to parse %s, using default: %t: %v", key, defaultValue, err)
		return defaultValue
	}
}

func getDuration(key string, defaultValue int) time.Duration {
	interval := getInt64(key, defaultValue)
	return time.Duration(interval) * time.Minute
}
