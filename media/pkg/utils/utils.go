package utils

import (
	"encoding/hex"
	"github.com/google/uuid"
)

// RandString implements encoding.BinaryMarshaler to New UUID
func RandString() string {
	pid := uuid.New()
	b, _ := pid.MarshalBinary()
	return hex.EncodeToString(b)
}
