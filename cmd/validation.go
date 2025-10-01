package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// Email validation regex - RFC 5322 simplified
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

// Phone number validation - supports common formats
var phoneRegex = regexp.MustCompile(`^[\d\s\-\(\)\+\.ext,]+$`)

// ValidateEmail validates an email address format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email address cannot be empty")
	}

	email = strings.TrimSpace(email)

	if len(email) > 254 {
		return fmt.Errorf("email address exceeds maximum length of 254 characters")
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email address format: %s", email)
	}

	// Check local part length (before @)
	parts := strings.Split(email, "@")
	if len(parts[0]) > 64 {
		return fmt.Errorf("email local part exceeds maximum length of 64 characters")
	}

	return nil
}

// ValidatePhoneNumber validates a phone number format
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number cannot be empty")
	}

	phone = strings.TrimSpace(phone)

	// Allow parsing format like "mobile:703-555-5555" or "work:301-684-8080,555"
	parts := strings.Split(phone, ":")
	if len(parts) == 2 {
		phone = parts[1]
	}

	if !phoneRegex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format: %s", phone)
	}

	// Check minimum length (at least 7 digits for a valid phone number)
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	if len(digits) < 7 {
		return fmt.Errorf("phone number must contain at least 7 digits")
	}

	return nil
}

// ValidateUUID validates a UUID format
func ValidateUUID(id string) error {
	if id == "" {
		return fmt.Errorf("UUID cannot be empty")
	}

	id = strings.TrimSpace(id)

	_, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %s", id)
	}

	return nil
}

// ValidateGroupName validates a group name
func ValidateGroupName(name string) error {
	if name == "" {
		return fmt.Errorf("group name cannot be empty")
	}

	name = strings.TrimSpace(name)

	// If it contains @, validate as email
	if strings.Contains(name, "@") {
		return ValidateEmail(name)
	}

	// Otherwise validate as a simple group name
	if len(name) > 60 {
		return fmt.Errorf("group name exceeds maximum length of 60 characters")
	}

	// Allow alphanumeric, hyphens, underscores, and dots
	validName := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("group name contains invalid characters: %s (only alphanumeric, dots, hyphens, and underscores allowed)", name)
	}

	return nil
}

// ValidateDepartment validates a department name
func ValidateDepartment(dept string) error {
	if dept == "" {
		return fmt.Errorf("department name cannot be empty")
	}

	dept = strings.TrimSpace(dept)

	if len(dept) > 100 {
		return fmt.Errorf("department name exceeds maximum length of 100 characters")
	}

	return nil
}

// SanitizeInput removes leading/trailing whitespace and potentially dangerous characters
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove other control characters
	input = strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, input)

	return input
}
