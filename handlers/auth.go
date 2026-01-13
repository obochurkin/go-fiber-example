package handlers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

type PasswordHashingConfig struct {
	Memory     uint32
	Time       uint32
	Threads    uint8
	SaltLength uint32
	KeyLength  uint32
}

var defaultPasswordHashingConfig = PasswordHashingConfig{
	Memory:     64 * 1024, // 64 MB
	Time:       3,         // iterations
	Threads:    2,         // parallelism
	SaltLength: 32,        // length of the salt in bytes
	KeyLength:  32,        // length of the generated key in bytes
}

// GenerateSalt creates a unique salt for password hashing
func GenerateSalt(email string) ([]byte, error) {
	salt := make([]byte, defaultPasswordHashingConfig.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

func HashPassword(password string, salt string) (string, error) {
	hashedPassword := argon2.IDKey(
		[]byte(password),
		[]byte(salt), 
		defaultPasswordHashingConfig.Time, 
		defaultPasswordHashingConfig.Memory, 
		defaultPasswordHashingConfig.Threads, 
		defaultPasswordHashingConfig.KeyLength,
	)

	return base64.RawStdEncoding.EncodeToString(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password, salt string) bool {
	hashedInputPassword, err := HashPassword(password, salt)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(hashedInputPassword)) == 1
}