package utils

import (
	"strings"
)

func GetBearerToken(authString string) (string, error) {
	if authString == "" {
		return "", nil
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authString, prefix) {
		return "", nil
	}
	token := strings.TrimPrefix(authString, prefix)
	return token, nil
}
