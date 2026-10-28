package utils

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

var (
	envLoaded bool
	envMutex  sync.Mutex
)

// LoadEnv reads environment variables from .env file if it exists and has not been loaded yet.
func LoadEnv() {
	envMutex.Lock()
	defer envMutex.Unlock()

	if envLoaded {
		return
	}

	file, err := os.Open(".env")
	if err != nil {
		envLoaded = true
		return // Ignore if file doesn't exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Remove quotes if present
		if (strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) ||
			(strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) {
			val = val[1 : len(val)-1]
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}

	envLoaded = true
}
