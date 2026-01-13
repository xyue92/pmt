package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// GenerateID creates a unique 7-character ID
func GenerateID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	// Combine timestamp and random bytes for uniqueness
	combined := fmt.Sprintf("%d%s", timestamp, hex.EncodeToString(randomBytes))
	hash := hashString(combined)

	// Return first 7 characters
	if len(hash) >= 7 {
		return hash[:7]
	}
	return hash
}

// MatchIDPrefix checks if a full ID starts with the given prefix
func MatchIDPrefix(fullID, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(fullID), strings.ToLower(prefix))
}

// Simple hash function to generate consistent IDs
func hashString(s string) string {
	hash := uint64(5381)
	for _, c := range s {
		hash = ((hash << 5) + hash) + uint64(c)
	}
	return fmt.Sprintf("%x", hash)
}
