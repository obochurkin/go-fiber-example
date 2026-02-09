package handlers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/obochurkin/go-fiber-example/config"
	"github.com/obochurkin/go-fiber-example/dtos"
	"github.com/obochurkin/go-fiber-example/errors"
	"golang.org/x/crypto/argon2"
)

var secret = []byte(config.GetEnvVariable("JWT_SECRET"))

type AuthController struct{}

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

func VerifyPassword(password, storedPasswordHash, salt string) bool {
	hashedInputPassword, err := HashPassword(password, salt)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(storedPasswordHash), []byte(hashedInputPassword)) == 1
}

func (auth *AuthController) Login(c *fiber.Ctx) error {
	input := dtos.AuthLoginDTO{}
	if err := c.BodyParser(&input); err != nil {
		return errors.NotFoundError()
	}

	user, err := usersRepository.GetOneUserByEmail(input.Email)
	if err != nil {
		return errors.NotFoundError()
	}
	isAuthorized := VerifyPassword(input.Password, user.Password, user.Salt)

	if !isAuthorized {
		return errors.NotFoundError()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "go-fiber-example",
		Subject:   strconv.FormatUint(uint64(user.ID), 10),
		Audience:  []string{"at"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // 72 hours
	})

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Errorf("token.SignedString: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "__Secure-access_token",
		Value:    tokenString,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None", 
		MaxAge:   3600,  // 1 hour
	}

	c.Cookie(&cookie)
	return c.SendStatus(fiber.StatusOK)
}
