package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	datatransfer "google.golang.org/api/admin/datatransfer/v1"
	admin "google.golang.org/api/admin/directory/v1"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	clientSecret string
	cacheFile    string

	scopes = []string{
		admin.AdminDirectoryUserReadonlyScope,
		admin.AdminDirectoryUserScope,
		admin.AdminDirectoryGroupReadonlyScope,
		admin.AdminDirectoryGroupMemberReadonlyScope,
		admin.AdminDirectoryGroupMemberScope,
		calendar.CalendarScope,
		calendar.CalendarReadonlyScope,
		calendar.CalendarEventsScope,
		calendar.CalendarEventsReadonlyScope,
		datatransfer.AdminDatatransferScope,
	}
)

// validateCredentialPath validates that a file path is safe to use for credentials
// Prevents directory traversal attacks by ensuring the path is within expected directories
func validateCredentialPath(filePath string) error {
	if filePath == "" {
		return errors.New("file path cannot be empty")
	}

	// Clean the path to resolve any ".." or "." components
	cleanPath := filepath.Clean(filePath)

	// Check for suspicious patterns
	if strings.Contains(cleanPath, "..") {
		return errors.New("path contains directory traversal sequence")
	}

	// Get absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Get user's home directory
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Ensure the path is within the user's home directory or /tmp for safety
	validPrefixes := []string{
		usr.HomeDir,
		"/tmp",
		os.TempDir(),
	}

	isValid := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(absPath, prefix) {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("credential path must be within user home directory or temp directory")
	}

	return nil
}

func newAdminClient() (*admin.Service, error) {
	client, err := newHTTPClient()
	if err != nil {
		return nil, err
	}

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func newCalendarClient() (*calendar.Service, error) {
	client, err := newHTTPClient()
	if err != nil {
		return nil, err
	}

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func newDataTransferClient() (*datatransfer.Service, error) {
	client, err := newHTTPClient()
	if err != nil {
		return nil, err
	}

	srv, err := datatransfer.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

// return an appropriately configured http.Client
func newHTTPClient() (*http.Client, error) {
	if clientSecret == "" {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		clientSecret = filepath.Join(usr.HomeDir, ".credentials", "client_secret.json")
	}

	// Validate credential file path to prevent directory traversal
	if err := validateCredentialPath(clientSecret); err != nil {
		return nil, fmt.Errorf("invalid client secret path: %w", err)
	}

	// #nosec G304 - Path is validated by validateCredentialPath() above
	b, err := os.ReadFile(clientSecret)
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/admin-directory_v1-go-quickstart.json
	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		return nil, err
	}
	cacheFile, err := tokenCacheFile()
	if err != nil {
		return nil, err
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		serr := saveToken(cacheFile, tok)
		if serr != nil {
			log.Fatalf("Unable to save token: %v", err)
		}
	}
	return config.Client(context.Background(), tok), nil

}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	if cacheFile == "" {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
		err = os.MkdirAll(tokenCacheDir, 0700)
		if err != nil {
			return "", err
		}
		return filepath.Join(tokenCacheDir,
			url.QueryEscape("gac.json")), err
	}

	return cacheFile, nil
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	// Validate credential file path to prevent directory traversal
	if err := validateCredentialPath(file); err != nil {
		return nil, fmt.Errorf("invalid token file path: %w", err)
	}

	// #nosec G304 - Path is validated by validateCredentialPath() above
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)

	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) (err error) {
	// Validate credential file path to prevent directory traversal
	if err := validateCredentialPath(file); err != nil {
		return fmt.Errorf("invalid token save path: %w", err)
	}

	fmt.Printf("Saving credential file to: %s\n", file)
	// #nosec G304 - Path is validated by validateCredentialPath() above
	f, err := os.Create(file)
	if err != nil {
		return
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	err = json.NewEncoder(f).Encode(token)

	return
}
