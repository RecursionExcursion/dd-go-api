package lib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var envLoaded bool = false

func EnvGet(key string) (string, error) {
	if !envLoaded {
		if err := EnvLoader(); err != nil {
			return "", err
		}
	}

	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("env key %v not set", key)
	}
	return val, nil
}

/* EnvGetOrPanic loads EnvVars if not already loaded, then retrieves the var or panics if not set */
func EnvGetOrPanic(key string) string {
	if !envLoaded {
		if err := EnvLoader(); err != nil {
			panic(err)
		}
	}

	val := os.Getenv(key)
	if val == "" {
		log.Panicf("Env key %v not set", key)
	}
	return val
}

func EnvGetOrFallback(key string, fallback string) string {
	if !envLoaded {
		if err := EnvLoader(); err != nil {
			return fallback
		}
	}

	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func EnvGetOrDefault(key string) string {
	if !envLoaded {
		if err := EnvLoader(); err != nil {
			return ""
		}
	}
	return os.Getenv(key)
}

/* EnvLoader loads files into env
 * Variables are expected to be in <key> = <value> format
 * If no files are provided by the called the loader will
 * attempt to load a .env file.
 */
func EnvLoader(files ...string) error {
	defer func() {
		envLoaded = true
	}()

	if len(files) == 0 {
		return loadOsEnv(".env")
	}

	for _, f := range files {
		if err := loadOsEnv(f); err != nil {
			return err
		}
	}
	return nil
}

func loadOsEnv(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		/* TODO Need to impl logic to cut out trailing comments
		 * '#' should only be permitted inside quotes
		 * if not then everything after needs to be treated as a comment
		 */
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		val = strings.Trim(val, `"'`)

		os.Setenv(key, val)
	}

	return scanner.Err()
}
