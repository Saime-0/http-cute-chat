package utils

import (
	"crypto/sha1"
	"fmt"
)

func HashPassword(passwd, globalSalt string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(globalSalt + passwd)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(""))), nil
}
