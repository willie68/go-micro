package utils

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateID generate an uuid as string without minuses
func GenerateID() string {
	uuidStr := uuid.NewString()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	return uuidStr
}
