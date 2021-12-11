package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateID() string {
	uuidStr := uuid.NewString()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	return uuidStr
}
