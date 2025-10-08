package cmd

import (
	"fmt"
	"testing"

	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
)

// mockAdminClient implements adminClientInterface for testing
type mockAdminClient struct {
	insertUserFunc   func(*admin.User) (*admin.User, error)
	getUserFunc      func(string) (*admin.User, error)
	listUsersFunc    func() (*admin.Users, error)
	insertMemberFunc func(string, *admin.Member) (*admin.Member, error)
	listMembersFunc  func(string) (*admin.Members, error)
}

func (m *mockAdminClient) InsertUser(user *admin.User) (*admin.User, error) {
	if m.insertUserFunc != nil {
		return m.insertUserFunc(user)
	}
	return user, nil
}

func (m *mockAdminClient) GetUser(email string) (*admin.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(email)
	}
	return nil, &googleapi.Error{Code: 404, Message: "User not found"}
}

func (m *mockAdminClient) ListUsers() (*admin.Users, error) {
	if m.listUsersFunc != nil {
		return m.listUsersFunc()
	}
	return &admin.Users{Users: []*admin.User{}}, nil
}

func (m *mockAdminClient) InsertMember(groupEmail string, member *admin.Member) (*admin.Member, error) {
	if m.insertMemberFunc != nil {
		return m.insertMemberFunc(groupEmail, member)
	}
	return member, nil
}

func (m *mockAdminClient) ListMembers(groupEmail string) (*admin.Members, error) {
	if m.listMembersFunc != nil {
		return m.listMembersFunc(groupEmail)
	}
	return &admin.Members{Members: []*admin.Member{}}, nil
}

// TestCreateUserWithClient tests the createUserWithClient function
func TestCreateUserWithClient(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		flags         createUserFlags
		wantErr       bool
		errContains   string
		setupMock     func(*mockAdminClient)
		verifyInsert  func(*testing.T, *admin.User)
		verifyMembers func(*testing.T, string, *admin.Member)
	}{
		{
			name: "successful user creation without groups",
			args: []string{"newuser@example.com"},
			flags: createUserFlags{
				personalEmail: "personal@example.com",
				firstName:     "John",
				lastName:      "Doe",
			},
			wantErr: false,
			setupMock: func(m *mockAdminClient) {
				m.insertUserFunc = func(u *admin.User) (*admin.User, error) {
					return u, nil
				}
			},
			verifyInsert: func(t *testing.T, u *admin.User) {
				if u.PrimaryEmail != "newuser@example.com" {
					t.Errorf("Expected PrimaryEmail newuser@example.com, got %s", u.PrimaryEmail)
				}
				if !u.ChangePasswordAtNextLogin {
					t.Error("Expected ChangePasswordAtNextLogin to be true")
				}
				if u.Password == "" {
					t.Error("Expected password to be set")
				}
				if u.Name == nil || u.Name.GivenName != "John" {
					t.Error("Expected GivenName to be John")
				}
			},
		},
		{
			name: "successful user creation with group",
			args: []string{"newuser@example.com"},
			flags: createUserFlags{
				groups:        []string{"admins@example.com"},
				personalEmail: "personal@example.com",
				firstName:     "Jane",
				lastName:      "Smith",
			},
			wantErr: false,
			setupMock: func(m *mockAdminClient) {
				m.insertUserFunc = func(u *admin.User) (*admin.User, error) {
					return u, nil
				}
				m.insertMemberFunc = func(groupEmail string, member *admin.Member) (*admin.Member, error) {
					return member, nil
				}
			},
			verifyMembers: func(t *testing.T, groupEmail string, member *admin.Member) {
				if groupEmail != "admins@example.com" {
					t.Errorf("Expected group admins@example.com, got %s", groupEmail)
				}
				if member.Email != "newuser@example.com" {
					t.Errorf("Expected member email newuser@example.com, got %s", member.Email)
				}
			},
		},
		{
			name:        "missing email argument",
			args:        []string{},
			flags:       createUserFlags{},
			wantErr:     true,
			errContains: "email is a required argument",
		},
		{
			name: "invalid email format",
			args: []string{"invalid-email"},
			flags: createUserFlags{
				personalEmail: "personal@example.com",
				firstName:     "John",
				lastName:      "Doe",
			},
			wantErr:     true,
			errContains: "invalid email address",
		},
		{
			name: "missing required flags",
			args: []string{"user@example.com"},
			flags: createUserFlags{
				firstName: "John",
				// Missing lastName and personalEmail
			},
			wantErr:     true,
			errContains: "all user details must be provided via flags",
		},
		{
			name: "user creation API error",
			args: []string{"newuser@example.com"},
			flags: createUserFlags{
				personalEmail: "personal@example.com",
				firstName:     "John",
				lastName:      "Doe",
			},
			wantErr:     true,
			errContains: "unable to create user",
			setupMock: func(m *mockAdminClient) {
				m.insertUserFunc = func(u *admin.User) (*admin.User, error) {
					return nil, &googleapi.Error{Code: 409, Message: "User already exists"}
				}
			},
		},
		{
			name: "invalid group name",
			args: []string{"newuser@example.com"},
			flags: createUserFlags{
				groups:        []string{"invalid group!"}, // Contains invalid characters
				personalEmail: "personal@example.com",
				firstName:     "John",
				lastName:      "Doe",
			},
			wantErr:     true,
			errContains: "invalid group name",
			setupMock: func(m *mockAdminClient) {
				m.insertUserFunc = func(u *admin.User) (*admin.User, error) {
					return u, nil
				}
			},
		},
		{
			name: "group member insertion error",
			args: []string{"newuser@example.com"},
			flags: createUserFlags{
				groups:        []string{"admins@example.com"},
				personalEmail: "personal@example.com",
				firstName:     "John",
				lastName:      "Doe",
			},
			wantErr:     true,
			errContains: "unable to add",
			setupMock: func(m *mockAdminClient) {
				m.insertUserFunc = func(u *admin.User) (*admin.User, error) {
					return u, nil
				}
				m.insertMemberFunc = func(groupEmail string, member *admin.Member) (*admin.Member, error) {
					return nil, &googleapi.Error{Code: 404, Message: "Group not found"}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &mockAdminClient{}
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			// Override insertUserFunc to verify the user if needed
			if tt.verifyInsert != nil {
				originalInsert := mockClient.insertUserFunc
				mockClient.insertUserFunc = func(u *admin.User) (*admin.User, error) {
					tt.verifyInsert(t, u)
					if originalInsert != nil {
						return originalInsert(u)
					}
					return u, nil
				}
			}

			// Override insertMemberFunc to verify the member if needed
			if tt.verifyMembers != nil {
				originalMember := mockClient.insertMemberFunc
				mockClient.insertMemberFunc = func(groupEmail string, member *admin.Member) (*admin.Member, error) {
					tt.verifyMembers(t, groupEmail, member)
					if originalMember != nil {
						return originalMember(groupEmail, member)
					}
					return member, nil
				}
			}

			// Execute the function
			err := createUserWithClient(mockClient, tt.args, tt.flags)

			// Verify results
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errContains)
				} else if tt.errContains != "" {
					errMsg := fmt.Sprintf("%v", err)
					found := false
					// Simple contains check
					for i := 0; i <= len(errMsg)-len(tt.errContains); i++ {
						if errMsg[i:i+len(tt.errContains)] == tt.errContains {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error containing '%s', got '%s'", tt.errContains, errMsg)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
