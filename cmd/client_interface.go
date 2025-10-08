package cmd

import (
	admin "google.golang.org/api/admin/directory/v1"
)

// adminClientInterface defines the interface for admin client operations
// This allows us to inject mock clients for testing
type adminClientInterface interface {
	InsertUser(user *admin.User) (*admin.User, error)
	GetUser(email string) (*admin.User, error)
	ListUsers() (*admin.Users, error)
	InsertMember(groupEmail string, member *admin.Member) (*admin.Member, error)
	ListMembers(groupEmail string) (*admin.Members, error)
}

// realAdminClientAdapter adapts the real admin.Service to our interface
type realAdminClientAdapter struct {
	service *admin.Service
}

// newRealAdminClientAdapter creates an adapter for the real admin service
func newRealAdminClientAdapter(service *admin.Service) adminClientInterface {
	return &realAdminClientAdapter{service: service}
}

func (a *realAdminClientAdapter) InsertUser(user *admin.User) (*admin.User, error) {
	return a.service.Users.Insert(user).Do()
}

func (a *realAdminClientAdapter) GetUser(email string) (*admin.User, error) {
	return a.service.Users.Get(email).Do()
}

func (a *realAdminClientAdapter) ListUsers() (*admin.Users, error) {
	return a.service.Users.List().Customer("my_customer").Do()
}

func (a *realAdminClientAdapter) InsertMember(groupEmail string, member *admin.Member) (*admin.Member, error) {
	return a.service.Members.Insert(groupEmail, member).Do()
}

func (a *realAdminClientAdapter) ListMembers(groupEmail string) (*admin.Members, error) {
	return a.service.Members.List(groupEmail).Do()
}
