package cmd

import (
	"testing"

	admin "google.golang.org/api/admin/directory/v1"
)

func TestRandomPassword(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 8", 8},
		{"length 12", 12},
		{"length 16", 16},
		{"length 32", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password := randomPassword(tt.length)

			if len(password) != tt.length {
				t.Errorf("randomPassword(%d) length = %d, want %d", tt.length, len(password), tt.length)
			}

			// Check that password only contains valid characters
			validChars := "abcdefghijkmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ123456789"
			for _, char := range password {
				found := false
				for _, valid := range validChars {
					if char == valid {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("randomPassword contains invalid character: %c", char)
				}
			}
		})
	}
}

func TestRandomPasswordUniqueness(t *testing.T) {
	// Generate multiple passwords and ensure they're different
	passwords := make(map[string]bool)
	iterations := 100
	length := 12

	for i := 0; i < iterations; i++ {
		password := randomPassword(length)
		if passwords[password] {
			t.Errorf("randomPassword generated duplicate: %s", password)
		}
		passwords[password] = true
	}
}

func TestRandomPasswordZeroLength(t *testing.T) {
	password := randomPassword(0)
	if len(password) != 0 {
		t.Errorf("randomPassword(0) = %q, want empty string", password)
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		fname        string
		lname        string
		primaryEmail string
	}{
		{
			name:         "full user info",
			email:        "john.personal@example.com",
			fname:        "John",
			lname:        "Doe",
			primaryEmail: "john.doe@company.com",
		},
		{
			name:         "empty personal email",
			email:        "",
			fname:        "Jane",
			lname:        "Smith",
			primaryEmail: "jane.smith@company.com",
		},
		{
			name:         "single character names",
			email:        "a@example.com",
			fname:        "A",
			lname:        "B",
			primaryEmail: "ab@company.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &admin.User{
				PrimaryEmail: tt.primaryEmail,
			}

			updateUser(user, tt.email, tt.fname, tt.lname)

			// Check Name is set correctly
			if user.Name == nil {
				t.Fatal("user.Name is nil")
			}
			if user.Name.GivenName != tt.fname {
				t.Errorf("GivenName = %q, want %q", user.Name.GivenName, tt.fname)
			}
			if user.Name.FamilyName != tt.lname {
				t.Errorf("FamilyName = %q, want %q", user.Name.FamilyName, tt.lname)
			}
			expectedFullName := tt.fname + " " + tt.lname
			if user.Name.FullName != expectedFullName {
				t.Errorf("FullName = %q, want %q", user.Name.FullName, expectedFullName)
			}

			// Check Emails are set correctly
			emails, ok := user.Emails.([]admin.UserEmail)
			if !ok {
				t.Fatal("user.Emails is not []admin.UserEmail")
			}
			if len(emails) != 2 {
				t.Errorf("len(Emails) = %d, want 2", len(emails))
				return
			}

			// First email should be personal
			if emails[0].Address != tt.email {
				t.Errorf("Emails[0].Address = %q, want %q", emails[0].Address, tt.email)
			}
			if emails[0].Type != "home" {
				t.Errorf("Emails[0].Type = %q, want %q", emails[0].Type, "home")
			}

			// Second email should be primary
			if emails[1].Address != tt.primaryEmail {
				t.Errorf("Emails[1].Address = %q, want %q", emails[1].Address, tt.primaryEmail)
			}
			if !emails[1].Primary {
				t.Error("Emails[1].Primary = false, want true")
			}
		})
	}
}
