package main

import (
	"crypto/sha256"
	"encoding/json"
	"net/http"

	"github.com/ivanrad/base62"
)

// GetShortCode generates a short, base62-encoded, string from a given URL.
func GetShortCode(longURL string) string {
	hash := sha256.Sum256([]byte(longURL))

	code := base62.EncodeToString(hash[0:8])

	return code

}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
