package utils

import (
	"github.com/rs/xid"
)

// GenerateID generate an uuid as string without minuses
func GenerateID() string {
	return xid.New().String()
}
