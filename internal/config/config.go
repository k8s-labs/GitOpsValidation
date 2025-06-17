package config

import (
	"os"
	"strconv"
)

type Config struct {
	Daemon        bool
	Sleep         int
	Namespace     string
	Source        string
	Kustomization string
}

func LoadConfig() *Config {
	return &Config{
		Daemon:        getEnvBool("GOV_DAEMON", false),
		Sleep:         getEnvInt("GOV_SLEEP", 60),
		Namespace:     getEnvString("GOV_NAMESPACE", "flux-system"),
		Source:        getEnvString("GOV_SOURCE", "gitops"),
		Kustomization: getEnvString("GOV_KUSTOMIZATION", "gitops"),
	}
}

func getEnvBool(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return b
}

func getEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}

func getEnvString(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
