package cmd

import (
	"reflect"
	"testing"

	admin "google.golang.org/api/admin/directory/v1"
)

func TestParsePhone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []admin.UserPhone
	}{
		{
			name:  "single phone",
			input: "mobile:555-1234",
			expected: []admin.UserPhone{
				{Type: "mobile", Value: "555-1234"},
			},
		},
		{
			name:  "multiple phones",
			input: "mobile:555-1234;work:555-5678",
			expected: []admin.UserPhone{
				{Type: "mobile", Value: "555-1234"},
				{Type: "work", Value: "555-5678"},
			},
		},
		{
			name:  "phone with extension",
			input: "work:555-1234,123",
			expected: []admin.UserPhone{
				{Type: "work", Value: "555-1234,123"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePhone(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parsePhone(%q) returned %d phones, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i].Type != tt.expected[i].Type || result[i].Value != tt.expected[i].Value {
					t.Errorf("parsePhone(%q)[%d] = {Type:%q, Value:%q}, want {Type:%q, Value:%q}",
						tt.input, i, result[i].Type, result[i].Value, tt.expected[i].Type, tt.expected[i].Value)
				}
			}
		})
	}
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []admin.UserAddress
	}{
		{
			name:  "simple address",
			input: "New York, NY",
			expected: []admin.UserAddress{
				{Formatted: "New York, NY"},
			},
		},
		{
			name:  "full address",
			input: "123 Main St, New York, NY 10001",
			expected: []admin.UserAddress{
				{Formatted: "123 Main St, New York, NY 10001"},
			},
		},
		{
			name:     "empty address",
			input:    "",
			expected: []admin.UserAddress{{Formatted: ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAddress(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseAddress(%q) = %+v, want %+v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseManager(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"valid email", "manager@example.com"},
		{"simple name", "manager"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseManager(tt.input)
			// Just verify it doesn't crash and returns a slice
			if result == nil {
				t.Errorf("parseManager(%q) returned nil", tt.input)
			}
		})
	}
}

func TestParseOrg(t *testing.T) {
	tests := []struct {
		name  string
		dept  string
		title string
	}{
		{"both dept and title", "Engineering", "Software Engineer"},
		{"only dept", "Sales", ""},
		{"only title", "", "Manager"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseOrg(&orgArgs{Dept: tt.dept, Title: tt.title})
			if len(result) == 0 {
				t.Errorf("parseOrg() returned empty slice")
			}
			if !result[0].Primary {
				t.Errorf("parseOrg() first org should be primary")
			}
		})
	}
}

func TestParseType(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"staff", "staff"},
		{"contractor", "contractor"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseType(tt.input)
			if result == nil {
				t.Errorf("parseType(%q) returned nil", tt.input)
			}
			if _, ok := result["Employee_Type"]; !ok {
				t.Errorf("parseType(%q) missing 'Employee_Type' key", tt.input)
			}
		})
	}
}

func TestParseID(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000"},
		{"simple id", "12345"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseID(tt.input)
			if result == nil {
				t.Errorf("parseID(%q) returned nil", tt.input)
			}
		})
	}
}
