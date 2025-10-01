package cmd

import (
	"strings"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid email with subdomain", "user@mail.example.com", false},
		{"valid email with dots", "first.last@example.com", false},
		{"valid email with plus", "user+tag@example.com", false},
		{"valid email with numbers", "user123@example456.com", false},
		{"empty email", "", true},
		{"no @ symbol", "userexample.com", true},
		{"no domain", "user@", true},
		{"no local part", "@example.com", true},
		{"spaces in email", "user @example.com", true},
		{"multiple @ symbols", "user@@example.com", true},
		{"local part too long", strings.Repeat("a", 65) + "@example.com", true},
		{"email too long", "user@" + strings.Repeat("a", 250) + ".com", true},
		{"invalid characters", "user name@example.com", true},
		{"trailing dot in domain", "user@example.com.", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{"valid US phone", "703-555-5555", false},
		{"valid with parentheses", "(703) 555-5555", false},
		{"valid with spaces", "703 555 5555", false},
		{"valid with extension", "703-555-5555,123", false},
		{"valid with ext keyword", "703-555-5555 ext 123", false},
		{"valid international", "+1-703-555-5555", false},
		{"valid with type prefix", "mobile:703-555-5555", false},
		{"valid with dots", "703.555.5555", false},
		{"valid minimal", "5555555", false},
		{"empty phone", "", true},
		{"too short", "12345", true},
		{"only 6 digits", "123456", true},
		{"invalid characters", "703-ABC-5555", true},
		{"only letters", "abcdefg", true},
		{"special chars", "703-555-5555!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePhoneNumber(tt.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid UUID v4", "550e8400-e29b-41d4-a716-446655440000", false},
		{"valid UUID v1", "c9bf9e58-0000-1000-8000-00805f9b34fb", false},
		{"valid UUID uppercase", "550E8400-E29B-41D4-A716-446655440000", false},
		{"valid UUID no hyphens", "550e8400e29b41d4a716446655440000", false},
		{"empty UUID", "", true},
		{"invalid format", "not-a-uuid", true},
		{"too short", "550e8400-e29b-41d4", true},
		{"too long", "550e8400-e29b-41d4-a716-446655440000-extra", true},
		{"invalid characters", "550e8400-e29b-41d4-a716-44665544000g", true},
		{"missing segments", "550e8400-e29b-a716-446655440000", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUUID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGroupName(t *testing.T) {
	tests := []struct {
		name      string
		groupName string
		wantErr   bool
	}{
		{"valid simple name", "engineering", false},
		{"valid with hyphen", "eng-team", false},
		{"valid with underscore", "eng_team", false},
		{"valid with dot", "eng.team", false},
		{"valid with numbers", "team123", false},
		{"valid email format", "team@example.com", false},
		{"empty name", "", true},
		{"too long", strings.Repeat("a", 61), true},
		{"invalid characters space", "eng team", true},
		{"invalid characters slash", "eng/team", true},
		{"invalid email format", "team@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGroupName(tt.groupName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGroupName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDepartment(t *testing.T) {
	tests := []struct {
		name    string
		dept    string
		wantErr bool
	}{
		{"valid department", "Engineering", false},
		{"valid with spaces", "Human Resources", false},
		{"valid with special chars", "R&D Department", false},
		{"valid long name", "Global Sales and Marketing Division", false},
		{"empty department", "", true},
		{"too long", strings.Repeat("a", 101), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDepartment(tt.dept)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDepartment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no change needed", "Hello World", "Hello World"},
		{"trim whitespace", "  Hello World  ", "Hello World"},
		{"remove null bytes", "Hello\x00World", "HelloWorld"},
		{"keep newlines", "Hello\nWorld", "Hello\nWorld"},
		{"keep tabs", "Hello\tWorld", "Hello\tWorld"},
		{"remove control chars", "Hello\x01\x02World", "HelloWorld"},
		{"empty string", "", ""},
		{"only whitespace", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput() = %q, want %q", result, tt.expected)
			}
		})
	}
}
