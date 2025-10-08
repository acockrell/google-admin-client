package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	admin "google.golang.org/api/admin/directory/v1"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// mockServer creates a test HTTP server that mocks Google API responses
type mockServer struct {
	*httptest.Server
	mu       sync.Mutex
	requests []*http.Request
	handler  http.HandlerFunc
}

// newMockServer creates a new mock server for testing Google API calls
func newMockServer(handler http.HandlerFunc) *mockServer {
	ms := &mockServer{
		requests: make([]*http.Request, 0),
		handler:  handler,
	}

	ms.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Store request for verification (protected by mutex)
		ms.mu.Lock()
		ms.requests = append(ms.requests, r)
		ms.mu.Unlock()
		// Call the handler
		if ms.handler != nil {
			ms.handler(w, r)
		}
	}))

	return ms
}

// getLastRequest returns the last request made to the mock server
func (ms *mockServer) getLastRequest() *http.Request {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if len(ms.requests) == 0 {
		return nil
	}
	return ms.requests[len(ms.requests)-1]
}

// createMockAdminClient creates an admin client that uses the mock server
func createMockAdminClient(t *testing.T, server *httptest.Server) *admin.Service {
	t.Helper()

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(server.URL)
			},
		},
	}

	srv, err := admin.NewService(context.Background(),
		option.WithHTTPClient(client),
		option.WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("Failed to create mock admin client: %v", err)
	}

	return srv
}

// createMockCalendarClient creates a calendar client that uses the mock server
func createMockCalendarClient(t *testing.T, server *httptest.Server) *calendar.Service {
	t.Helper()

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(server.URL)
			},
		},
	}

	srv, err := calendar.NewService(context.Background(),
		option.WithHTTPClient(client),
		option.WithEndpoint(server.URL))
	if err != nil {
		t.Fatalf("Failed to create mock calendar client: %v", err)
	}

	return srv
}

// TestUserOperations_Integration tests user creation, listing, and updates with mocked API
func TestUserOperations_Integration(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		// Create mock server that responds to user creation
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method and path
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}

			// Read and parse the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			var requestUser admin.User
			if err := json.Unmarshal(body, &requestUser); err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Verify user data
			if requestUser.PrimaryEmail == "" {
				t.Error("PrimaryEmail should not be empty")
			}
			if !requestUser.ChangePasswordAtNextLogin {
				t.Error("ChangePasswordAtNextLogin should be true")
			}
			if requestUser.Password == "" {
				t.Error("Password should not be empty")
			}

			// Return successful response
			responseUser := admin.User{
				PrimaryEmail:              requestUser.PrimaryEmail,
				ChangePasswordAtNextLogin: requestUser.ChangePasswordAtNextLogin,
				Password:                  requestUser.Password,
				Name:                      requestUser.Name,
				Emails:                    requestUser.Emails,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(responseUser); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test user creation
		testUser := &admin.User{
			PrimaryEmail:              "testuser@example.com",
			ChangePasswordAtNextLogin: true,
			Password:                  randomPassword(12),
		}
		updateUser(testUser, "personal@example.com", "Test", "User")

		result, err := client.Users.Insert(testUser).Do()
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Verify response
		if result.PrimaryEmail != testUser.PrimaryEmail {
			t.Errorf("Expected email %s, got %s", testUser.PrimaryEmail, result.PrimaryEmail)
		}
		if !result.ChangePasswordAtNextLogin {
			t.Error("Expected ChangePasswordAtNextLogin to be true")
		}

		// Verify request was made
		lastReq := mockServer.getLastRequest()
		if lastReq == nil {
			t.Fatal("No request was made to the mock server")
		}
	})

	t.Run("ListUsers", func(t *testing.T) {
		// Create mock server that responds to user listing
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "GET" {
				t.Errorf("Expected GET request, got %s", r.Method)
			}

			// Return list of users
			users := &admin.Users{
				Users: []*admin.User{
					{
						PrimaryEmail: "user1@example.com",
						Name: &admin.UserName{
							GivenName:  "User",
							FamilyName: "One",
							FullName:   "User One",
						},
					},
					{
						PrimaryEmail: "user2@example.com",
						Name: &admin.UserName{
							GivenName:  "User",
							FamilyName: "Two",
							FullName:   "User Two",
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(users); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test user listing
		result, err := client.Users.List().Customer("my_customer").Do()
		if err != nil {
			t.Fatalf("Failed to list users: %v", err)
		}

		// Verify response
		if len(result.Users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(result.Users))
		}
		if result.Users[0].PrimaryEmail != "user1@example.com" {
			t.Errorf("Expected first user email user1@example.com, got %s", result.Users[0].PrimaryEmail)
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		// Create mock server that responds to user get request
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "GET" {
				t.Errorf("Expected GET request, got %s", r.Method)
			}

			// Verify URL contains user email
			if !strings.Contains(r.URL.Path, "testuser@example.com") {
				t.Errorf("Expected URL to contain user email, got %s", r.URL.Path)
			}

			// Return user data
			user := &admin.User{
				PrimaryEmail: "testuser@example.com",
				Name: &admin.UserName{
					GivenName:  "Test",
					FamilyName: "User",
					FullName:   "Test User",
				},
				OrgUnitPath: "/",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(user); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test getting user
		result, err := client.Users.Get("testuser@example.com").Do()
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		// Verify response
		if result.PrimaryEmail != "testuser@example.com" {
			t.Errorf("Expected email testuser@example.com, got %s", result.PrimaryEmail)
		}
		if result.Name.FullName != "Test User" {
			t.Errorf("Expected name 'Test User', got %s", result.Name.FullName)
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Create mock server that responds to user update
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "PUT" && r.Method != "PATCH" {
				t.Errorf("Expected PUT or PATCH request, got %s", r.Method)
			}

			// Read and parse the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			var requestUser admin.User
			if err := json.Unmarshal(body, &requestUser); err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Return updated user
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(requestUser); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test updating user
		updateData := &admin.User{
			Name: &admin.UserName{
				GivenName:  "Updated",
				FamilyName: "Name",
				FullName:   "Updated Name",
			},
		}

		result, err := client.Users.Update("testuser@example.com", updateData).Do()
		if err != nil {
			t.Fatalf("Failed to update user: %v", err)
		}

		// Verify response
		if result.Name.FullName != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %s", result.Name.FullName)
		}
	})
}

// TestGroupOperations_Integration tests group operations with mocked API
func TestGroupOperations_Integration(t *testing.T) {
	t.Run("ListGroups", func(t *testing.T) {
		// Create mock server that responds to group listing
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "GET" {
				t.Errorf("Expected GET request, got %s", r.Method)
			}

			// Return list of groups
			groups := &admin.Groups{
				Groups: []*admin.Group{
					{
						Email:       "group1@example.com",
						Name:        "Group One",
						Description: "First test group",
					},
					{
						Email:       "group2@example.com",
						Name:        "Group Two",
						Description: "Second test group",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(groups); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test group listing
		result, err := client.Groups.List().Customer("my_customer").Do()
		if err != nil {
			t.Fatalf("Failed to list groups: %v", err)
		}

		// Verify response
		if len(result.Groups) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(result.Groups))
		}
		if result.Groups[0].Email != "group1@example.com" {
			t.Errorf("Expected first group email group1@example.com, got %s", result.Groups[0].Email)
		}
	})

	t.Run("GetGroup", func(t *testing.T) {
		// Create mock server that responds to group get request
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "GET" {
				t.Errorf("Expected GET request, got %s", r.Method)
			}

			// Return group data
			group := &admin.Group{
				Email:       "testgroup@example.com",
				Name:        "Test Group",
				Description: "A test group",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(group); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test getting group
		result, err := client.Groups.Get("testgroup@example.com").Do()
		if err != nil {
			t.Fatalf("Failed to get group: %v", err)
		}

		// Verify response
		if result.Email != "testgroup@example.com" {
			t.Errorf("Expected email testgroup@example.com, got %s", result.Email)
		}
		if result.Name != "Test Group" {
			t.Errorf("Expected name 'Test Group', got %s", result.Name)
		}
	})

	t.Run("ListGroupMembers", func(t *testing.T) {
		// Create mock server that responds to member listing
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "GET" {
				t.Errorf("Expected GET request, got %s", r.Method)
			}

			// Return list of members
			members := &admin.Members{
				Members: []*admin.Member{
					{
						Email: "member1@example.com",
						Role:  "MEMBER",
						Type:  "USER",
					},
					{
						Email: "member2@example.com",
						Role:  "OWNER",
						Type:  "USER",
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(members); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test listing group members
		result, err := client.Members.List("testgroup@example.com").Do()
		if err != nil {
			t.Fatalf("Failed to list group members: %v", err)
		}

		// Verify response
		if len(result.Members) != 2 {
			t.Errorf("Expected 2 members, got %d", len(result.Members))
		}
		if result.Members[0].Email != "member1@example.com" {
			t.Errorf("Expected first member email member1@example.com, got %s", result.Members[0].Email)
		}
		if result.Members[1].Role != "OWNER" {
			t.Errorf("Expected second member role OWNER, got %s", result.Members[1].Role)
		}
	})

	t.Run("InsertGroupMember", func(t *testing.T) {
		// Create mock server that responds to member insertion
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}

			// Read and parse the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			var requestMember admin.Member
			if err := json.Unmarshal(body, &requestMember); err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Verify member data
			if requestMember.Email == "" {
				t.Error("Member email should not be empty")
			}

			// Return successful response
			responseMember := admin.Member{
				Email: requestMember.Email,
				Role:  "MEMBER",
				Type:  "USER",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(responseMember); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test inserting group member
		member := &admin.Member{
			Email: "newmember@example.com",
		}

		result, err := client.Members.Insert("testgroup@example.com", member).Do()
		if err != nil {
			t.Fatalf("Failed to insert group member: %v", err)
		}

		// Verify response
		if result.Email != "newmember@example.com" {
			t.Errorf("Expected email newmember@example.com, got %s", result.Email)
		}
	})
}

// TestCalendarOperations_Integration tests calendar operations with mocked API
func TestCalendarOperations_Integration(t *testing.T) {
	t.Run("CreateCalendar", func(t *testing.T) {
		// Create mock server that responds to calendar creation
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}

			// Read and parse the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			var requestCalendar calendar.Calendar
			if err := json.Unmarshal(body, &requestCalendar); err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Verify calendar data
			if requestCalendar.Summary == "" {
				t.Error("Calendar summary should not be empty")
			}

			// Return successful response
			responseCalendar := calendar.Calendar{
				Id:          "test-calendar-id",
				Summary:     requestCalendar.Summary,
				Description: requestCalendar.Description,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(responseCalendar); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockCalendarClient(t, mockServer.Server)

		// Test calendar creation
		testCalendar := &calendar.Calendar{
			Summary:     "Test Calendar",
			Description: "A test calendar",
		}

		result, err := client.Calendars.Insert(testCalendar).Do()
		if err != nil {
			t.Fatalf("Failed to create calendar: %v", err)
		}

		// Verify response
		if result.Summary != "Test Calendar" {
			t.Errorf("Expected summary 'Test Calendar', got %s", result.Summary)
		}
		if result.Id == "" {
			t.Error("Calendar ID should not be empty")
		}
	})

	t.Run("CreateEvent", func(t *testing.T) {
		// Create mock server that responds to event creation
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}

			// Read and parse the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			var requestEvent calendar.Event
			if err := json.Unmarshal(body, &requestEvent); err != nil {
				t.Fatalf("Failed to unmarshal request: %v", err)
			}

			// Verify event data
			if requestEvent.Summary == "" {
				t.Error("Event summary should not be empty")
			}

			// Return successful response
			responseEvent := calendar.Event{
				Id:          "test-event-id",
				Summary:     requestEvent.Summary,
				Description: requestEvent.Description,
				Start:       requestEvent.Start,
				End:         requestEvent.End,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(responseEvent); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockCalendarClient(t, mockServer.Server)

		// Test event creation
		testEvent := &calendar.Event{
			Summary:     "Test Event",
			Description: "A test event",
			Start: &calendar.EventDateTime{
				DateTime: "2024-10-08T10:00:00Z",
			},
			End: &calendar.EventDateTime{
				DateTime: "2024-10-08T11:00:00Z",
			},
		}

		result, err := client.Events.Insert("primary", testEvent).Do()
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}

		// Verify response
		if result.Summary != "Test Event" {
			t.Errorf("Expected summary 'Test Event', got %s", result.Summary)
		}
		if result.Id == "" {
			t.Error("Event ID should not be empty")
		}
	})

	t.Run("ListEvents", func(t *testing.T) {
		// Create mock server that responds to event listing
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method
			if r.Method != "GET" {
				t.Errorf("Expected GET request, got %s", r.Method)
			}

			// Return list of events
			events := &calendar.Events{
				Items: []*calendar.Event{
					{
						Id:      "event1",
						Summary: "Event One",
						Start: &calendar.EventDateTime{
							DateTime: "2024-10-08T10:00:00Z",
						},
					},
					{
						Id:      "event2",
						Summary: "Event Two",
						Start: &calendar.EventDateTime{
							DateTime: "2024-10-08T14:00:00Z",
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(events); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockCalendarClient(t, mockServer.Server)

		// Test event listing
		result, err := client.Events.List("primary").Do()
		if err != nil {
			t.Fatalf("Failed to list events: %v", err)
		}

		// Verify response
		if len(result.Items) != 2 {
			t.Errorf("Expected 2 events, got %d", len(result.Items))
		}
		if result.Items[0].Summary != "Event One" {
			t.Errorf("Expected first event summary 'Event One', got %s", result.Items[0].Summary)
		}
	})
}

// TestErrorHandling_Integration tests error handling with mocked API
func TestErrorHandling_Integration(t *testing.T) {
	t.Run("UserNotFound", func(t *testing.T) {
		// Create mock server that returns 404
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			errorResponse := map[string]interface{}{
				"error": map[string]interface{}{
					"code":    404,
					"message": "User not found",
				},
			}
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				t.Fatalf("Failed to encode error response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test getting non-existent user
		_, err := client.Users.Get("nonexistent@example.com").Do()
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		// Create mock server that returns 400
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := map[string]interface{}{
				"error": map[string]interface{}{
					"code":    400,
					"message": "Invalid request",
				},
			}
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				t.Fatalf("Failed to encode error response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Test creating user with invalid data
		invalidUser := &admin.User{} // Empty user
		_, err := client.Users.Insert(invalidUser).Do()
		if err == nil {
			t.Error("Expected error for invalid user data, got nil")
		}
	})
}

// TestConcurrentOperations_Integration tests concurrent API operations
func TestConcurrentOperations_Integration(t *testing.T) {
	t.Run("ConcurrentGroupListing", func(t *testing.T) {
		// Track concurrent requests (using atomic for thread safety)
		var requestCount int32

		// Create mock server that responds to group listing
		mockServer := newMockServer(func(w http.ResponseWriter, r *http.Request) {
			count := atomic.AddInt32(&requestCount, 1)

			// Return list of groups
			groups := &admin.Groups{
				Groups: []*admin.Group{
					{
						Email: fmt.Sprintf("group%d@example.com", count),
						Name:  fmt.Sprintf("Group %d", count),
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(groups); err != nil {
				t.Fatalf("Failed to encode response: %v", err)
			}
		})
		defer mockServer.Close()

		// Create mock client
		client := createMockAdminClient(t, mockServer.Server)

		// Make multiple concurrent requests
		done := make(chan bool)
		for i := 0; i < 3; i++ {
			go func() {
				_, err := client.Groups.List().Customer("my_customer").Do()
				if err != nil {
					t.Errorf("Failed to list groups: %v", err)
				}
				done <- true
			}()
		}

		// Wait for all requests to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		// Verify all requests were made
		finalCount := atomic.LoadInt32(&requestCount)
		if finalCount != 3 {
			t.Errorf("Expected 3 requests, got %d", finalCount)
		}
	})
}

// TestHelperFunctions_Integration tests helper functions used in integration scenarios
func TestHelperFunctions_Integration(t *testing.T) {
	t.Run("CaptureOutput", func(t *testing.T) {
		// Save original stdout
		oldStdout := os.Stdout

		// Create pipe to capture output
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Generate and print output
		password := randomPassword(12)
		fmt.Println(password)

		// Restore stdout
		if err := w.Close(); err != nil {
			t.Fatalf("Failed to close pipe writer: %v", err)
		}
		os.Stdout = oldStdout

		// Read captured output
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			t.Fatalf("Failed to copy output: %v", err)
		}
		output := buf.String()

		// Verify output
		if !strings.Contains(output, password) {
			t.Errorf("Expected output to contain password %s, got %s", password, output)
		}
	})

	t.Run("UserUpdateHelper", func(t *testing.T) {
		user := &admin.User{
			PrimaryEmail: "test@example.com",
		}

		updateUser(user, "personal@example.com", "Test", "User")

		// Verify user was updated correctly
		if user.Name == nil {
			t.Fatal("Expected Name to be set")
		}
		if user.Name.FullName != "Test User" {
			t.Errorf("Expected FullName 'Test User', got %s", user.Name.FullName)
		}

		emails, ok := user.Emails.([]admin.UserEmail)
		if !ok || len(emails) != 2 {
			t.Errorf("Expected 2 emails, got %d", len(emails))
		}
	})
}
