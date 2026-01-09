package constants

type UserRole string

const (
	UserRoleUser    UserRole = "user"
	UserRoleManager UserRole = "manager"
)

// user login context
type contextKey string

const (
	ContextUserID contextKey = "user_id"
	ContextRole   contextKey = "role"
)
