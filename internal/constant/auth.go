package constant

import "time"

// Authentication Constants
const (
	AUTH_BEARER_PREFIX     = "Bearer "
	AUTH_BASIC_PREFIX      = "Basic "
	AUTH_TOKEN_TYPE_ACCESS = "access"
	AUTH_TOKEN_TYPE_REFRESH = "refresh"
)

// Token Durations
const (
	ACCESS_TOKEN_DURATION  = time.Hour * 24     // 1 day
	REFRESH_TOKEN_DURATION = time.Hour * 24 * 30 // 30 days
)

// User Roles
const (
	ROLE_ADMIN     = "admin"
	ROLE_DOCTOR    = "doctor"
	ROLE_PATIENT   = "patient"
	ROLE_STAFF     = "staff"
)
