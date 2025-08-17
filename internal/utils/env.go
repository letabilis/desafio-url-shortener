package utils

import (
	"fmt"
	"os"
)

func LoadEnv(envars ...string) (map[string]string, error) {
	values := make(map[string]string)
	for _, envar := range envars {
		value, ok := os.LookupEnv(envar)
		if !ok {
			return values, fmt.Errorf("missing environment variable %s", envar)
		}
		values[envar] = value
	}
	return values, nil
}
