package main

import (
	"net/http/httptest"
	"testing"
	
	"github.com/gofiber/fiber/v2"
)

func TestEndpoints(t *testing.T) {
	app := Init()

	// Test GET api/v1/users

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}