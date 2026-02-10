package handlers

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/obochurkin/go-fiber-example/dtos"
	"github.com/obochurkin/go-fiber-example/errors"
	"github.com/obochurkin/go-fiber-example/repositories"
)

type UserController struct{}

var usersRepository = &repositories.UsersRepository{}

func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	users, err := usersRepository.FindAll()
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (uc *UserController) GetUserById(c *fiber.Ctx) error {
	userID := c.Locals("id").(uint)
	user, err := usersRepository.FindByID(userID)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	userInput := dtos.CreateUserDTO{}
	if err := c.BodyParser(&userInput); err != nil {
		return errors.BadRequestError()
	}

	// Step 1: Check email exists
	count, err := usersRepository.IsEmailExists(userInput.Email)
	if err != nil {
		log.Errorf("Error checking existing email: %v", err)
		return errors.InternalError()
	}

	if count > 0 {
		return c.SendStatus(fiber.StatusCreated)
	}

	// Step 2: Check if password has been pwned
	if err := uc.checkPwnedPassword(c.Context(), userInput.Password); err != nil {
		return err
	}

	salt, err := GenerateSalt(userInput.Email)
	if err != nil {
		log.Errorf("Error generating salt: %v", err)
		return errors.InternalError()
	}

	stringifiedSalt := base64.StdEncoding.EncodeToString(salt)

	hashedPassword, err := HashPassword(userInput.Password, stringifiedSalt)
	if err != nil {
		return errors.InternalError()
	}

	userInput.Password = hashedPassword

	if err := usersRepository.Create(userInput, stringifiedSalt); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (uc *UserController) checkPwnedPassword(ctx context.Context, password string) error {
	const pwnedAPIURL = "https://api.pwnedpasswords.com/range/"
	hash := sha1.Sum([]byte(password))
	hashString := strings.ToUpper(hex.EncodeToString(hash[:]))
	prefix := hashString[:5]
	suffix := hashString[5:]

	// Create request with context. inherits:
	// timeout, 
	// cancellation, 
	// request ID
	req, err := http.NewRequestWithContext(ctx, "GET", pwnedAPIURL+prefix, nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return nil
	}

	client := &http.Client{
		Timeout: 5 * time.Second, // Fallback timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error calling Pwned Passwords API: %v", err)
		return nil
	}
	defer resp.Body.Close()

	//check response status
	if resp.StatusCode != http.StatusOK {
		log.Errorf("Unexpected response from Pwned Passwords API: %d", resp.StatusCode)
		return nil
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// Check if context was cancelled (user disconnected/timeout)
		select {
		case <-ctx.Done():
			log.Warnf("Password check cancelled: %v", ctx.Err())
			return nil
		default:
		}

		text := scanner.Text()
		expectedLines := strings.Split(text, ":")
		if len(expectedLines) != 2 {
			log.Warnf("Unexpected line format: %s", text)
			continue
		}
		externalHashSuffix := expectedLines[0]

		if externalHashSuffix == suffix {
			log.Infof("Password found in breach database")
			return errors.BadRequestError()
		}
	}

	// Check for scanner errors or context cancellation
	if err := scanner.Err(); err != nil {
		log.Errorf("Error scanning response: %v", err)
		return nil
	}

	log.Infof("Password not found in breach database")
	return nil
}

// func (uc *UserController) UpdateUser(c *fiber.Ctx) error {

// }
