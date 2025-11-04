package utils

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} 

	return defaultVal
}

func GetEnvAsInt(key string, defaultVal int) int {
	strVal := GetEnv(key, strconv.Itoa(defaultVal))

	if val, err := strconv.Atoi(strVal); err == nil {
		return val
	}

	return defaultVal
}
