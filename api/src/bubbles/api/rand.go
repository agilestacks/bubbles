package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func RandomString(randomBytesLen int) (string, error) {
	buf := make([]byte, randomBytesLen)
	read, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("Unable to generate random string: random read error (read %d bytes): %v", read, err)
	}
	return hex.EncodeToString(buf), nil
}
