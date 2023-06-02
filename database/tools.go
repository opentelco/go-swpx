package database

import "github.com/hashicorp/go-uuid"

// used by the repo to generate a new ID for a device or configuration
func NewID() string {
	guid, _ := uuid.GenerateUUID()
	return guid
}
