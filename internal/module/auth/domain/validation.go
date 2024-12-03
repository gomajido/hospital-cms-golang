package domain

import (
	"fmt"
	"regexp"

	"github.com/gomajido/hospital-cms-golang/internal/constant"
	"github.com/gomajido/hospital-cms-golang/internal/response"
)

const (
	// Field names for validation messages
	EMAIL_FIELD    = "email"
	PASSWORD_FIELD = "password"
	NAME_FIELD     = "name"
	PHONE_FIELD    = "phone"
	ROLE_IDS_FIELD = "role_ids"
	ROLE_NAMES_FIELD = "role_names"
	USER_ID_FIELD  = "user_id"
)

// UpdateUserRequest represents the request to update user details
type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"-"` // Never expose password in JSON
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

// AssignRolesRequest represents the request to assign roles to a user
type AssignRolesRequest struct {
	RoleNames []string `json:"role_names"`
}

func (r *RegisterRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if r.Email == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        EMAIL_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, EMAIL_FIELD),
		})
	} else if !isValidEmail(r.Email) {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        EMAIL_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, EMAIL_FIELD, "email"),
		})
	}

	if r.Password == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        PASSWORD_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, PASSWORD_FIELD),
		})
	} else if len(r.Password) < 8 {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        PASSWORD_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_MIN_LENGTH, PASSWORD_FIELD, 8),
		})
	}

	if r.Name == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        NAME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, NAME_FIELD),
		})
	}

	return errorInfo
}

func (l *LoginRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if l.Email == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        EMAIL_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, EMAIL_FIELD),
		})
	} else if !isValidEmail(l.Email) {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        EMAIL_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, EMAIL_FIELD, "email"),
		})
	}

	if l.Password == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        PASSWORD_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, PASSWORD_FIELD),
		})
	}

	return errorInfo
}

func (u *UpdateUserRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if u.Name == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        NAME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, NAME_FIELD),
		})
	}

	if u.Password != constant.EMPTY_STRING && len(u.Password) < 8 {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        PASSWORD_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_MIN_LENGTH, PASSWORD_FIELD, 8),
		})
	}

	return errorInfo
}

func (a *AssignRolesRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if len(a.RoleNames) == 0 {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        ROLE_NAMES_FIELD,
			ErrorMessage: "At least one role name is required",
		})
	}

	return errorInfo
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
