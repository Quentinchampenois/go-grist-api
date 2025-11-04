package grist

import (
	"encoding/json"
	"fmt"
)

// AccessRole represents the role types from the API
type AccessRole string

// Define the enum values
const (
	AccessRoleOwner  AccessRole = "owners"
	AccessRoleEditor AccessRole = "editors"
	AccessRoleViewer AccessRole = "viewers"
)

// UnmarshalJSON allows parsing JSON into AccessRole with validation
func (r *AccessRole) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	role := AccessRole(s)
	switch role {
	case AccessRoleOwner, AccessRoleEditor, AccessRoleViewer:
		*r = role
		return nil
	default:
		return fmt.Errorf("invalid AccessRole: %s", s)
	}
}
