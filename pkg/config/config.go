package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type EnvTypeSet interface {
	string | bool
}

func MustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("config.MustGet: required var \"%s\" is not exists", key))
	}

	return val
}

func Get(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	return val
}

func GetBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	res, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}

	return res
}

func GetInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	res, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		return defaultVal
	}

	return int(res)
}

func LoadDotEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("config.LoadDotEnv: file %s not found", path)
			return nil
		}

		return fmt.Errorf("config.LoadDotEnv: cant read file: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	splitted := strings.Split(string(data), "\n")
	for _, line := range splitted {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		key, val, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("config.LoadDotEnv: invalid format: %s", line)
		}

		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)

		if os.Getenv(key) == "" {
			err = os.Setenv(key, val)
			if err != nil {
				return fmt.Errorf("config.LoadDotEnv: cant set env var: %w", err)
			}
		}
	}

	return nil
}
