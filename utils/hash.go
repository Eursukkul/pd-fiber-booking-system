package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func HashID(id int) string {
	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(id)))
	return hex.EncodeToString(hash.Sum(nil))
}