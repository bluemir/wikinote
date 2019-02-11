package backend

import (
	"github.com/bluemir/go-utils/auth"
)

const (
	defaultRule = `
rules:
  admin:  [ "view", "edit", "user", "search" ]
  editor: [ "view', "edit", "attach", "search" ]
  viewer: [ "view", "search" ]
  guest:  [ "view" ]
`
)

const (
	// pre-defined role
	RoleAdmin  auth.Role = "admin"
	RoleEditor auth.Role = "editor"
	RoleViewer auth.Role = "editor"
	RoleGuest  auth.Role = "guest"
	// pre-defined action
	ActionView auth.Action = "view"
	ActionEdit auth.Action = "edit"
)
