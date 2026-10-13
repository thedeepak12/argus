package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

var DefaultConfig = &PasswordConfig{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		DefaultConfig.Time,
		DefaultConfig.Memory,
		DefaultConfig.Threads,
		DefaultConfig.KeyLen,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	fullHash := fmt.Sprintf(format, argon2.Version, DefaultConfig.Memory, DefaultConfig.Time, DefaultConfig.Threads, b64Salt, b64Hash)

	return fullHash, nil
}

func VerifyPassword(hash, password string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("Invalid hash format")
	}

	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		uint32(len(decodedHash)),
	)

	return subtle.ConstantTimeCompare(decodedHash, otherHash) == 1, nil
}
