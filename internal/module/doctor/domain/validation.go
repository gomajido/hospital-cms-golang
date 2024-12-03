package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/gomajido/hospital-cms-golang/internal/constant"
	"github.com/gomajido/hospital-cms-golang/internal/response"
	"github.com/google/uuid"
)

const (
	// Field names for validation messages
	NAME_FIELD           = "name"
	SERVICE_ID_FIELD     = "service_id"
	DESCRIPTION_FIELD    = "description"
	SPECIALIZATION_FIELD = "specialization"
	DEGREE_FIELD         = "degree"
	EXPERIENCE_FIELD     = "experience"
	DAY_FIELD           = "day"
	START_TIME_FIELD    = "start_time"
	END_TIME_FIELD      = "end_time"
	DATE_FIELD          = "date"
	STATUS_FIELD        = "status"
)

// CreateDoctorRequest represents the request to create a doctor
type CreateDoctorRequest struct {
	Name           string    `json:"name"`
	ServiceID      uuid.UUID `json:"service_id"`
	Description    string    `json:"description"`
	Specialization string    `json:"specialization"`
	Degree         string    `json:"degree"`
	Experience     string    `json:"experience"`
}

// UpdateDoctorRequest represents the request to update a doctor
type UpdateDoctorRequest struct {
	Name           string    `json:"name"`
	ServiceID      uuid.UUID `json:"service_id"`
	Description    string    `json:"description"`
	Specialization string    `json:"specialization"`
	Degree         string    `json:"degree"`
	Experience     string    `json:"experience"`
}

// CreateScheduleRequest represents the request to create a doctor schedule
type CreateScheduleRequest struct {
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// UpdateScheduleRequest represents the request to update a doctor schedule
type UpdateScheduleRequest struct {
	Day       string `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// CreateRescheduleRequest represents the request to create a schedule change
type CreateRescheduleRequest struct {
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
}

// UpdateRescheduleRequest represents the request to update a schedule change
type UpdateRescheduleRequest struct {
	Date        time.Time `json:"date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
}

func (c *CreateDoctorRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if c.Name == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        NAME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, NAME_FIELD),
		})
	}

	if c.ServiceID == uuid.Nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        SERVICE_ID_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, SERVICE_ID_FIELD),
		})
	}

	if c.Specialization == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        SPECIALIZATION_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, SPECIALIZATION_FIELD),
		})
	}

	if c.Degree == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DEGREE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, DEGREE_FIELD),
		})
	}

	if c.Experience == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        EXPERIENCE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, EXPERIENCE_FIELD),
		})
	}

	return errorInfo
}

func (u *UpdateDoctorRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if u.Name == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        NAME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, NAME_FIELD),
		})
	}

	if u.ServiceID == uuid.Nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        SERVICE_ID_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, SERVICE_ID_FIELD),
		})
	}

	if u.Specialization == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        SPECIALIZATION_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, SPECIALIZATION_FIELD),
		})
	}

	if u.Degree == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DEGREE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, DEGREE_FIELD),
		})
	}

	if u.Experience == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        EXPERIENCE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, EXPERIENCE_FIELD),
		})
	}

	return errorInfo
}

func (c *CreateScheduleRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	validDays := map[string]bool{
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
		"Sunday":    true,
	}

	if c.Day == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DAY_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, DAY_FIELD),
		})
	} else if !validDays[c.Day] {
		validDaysList := strings.Join([]string{
			"Monday", "Tuesday", "Wednesday", "Thursday",
			"Friday", "Saturday", "Sunday",
		}, ", ")
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DAY_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_VALUE, DAY_FIELD, validDaysList),
		})
	}

	if c.StartTime == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        START_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, START_TIME_FIELD),
		})
	} else if _, err := time.Parse("15:04:05", c.StartTime); err != nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        START_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, START_TIME_FIELD, "HH:MM:SS"),
		})
	}

	if c.EndTime == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        END_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, END_TIME_FIELD),
		})
	} else if _, err := time.Parse("15:04:05", c.EndTime); err != nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        END_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, END_TIME_FIELD, "HH:MM:SS"),
		})
	}

	return errorInfo
}

func (u *UpdateScheduleRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	validDays := map[string]bool{
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
		"Sunday":    true,
	}

	if u.Day == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DAY_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, DAY_FIELD),
		})
	} else if !validDays[u.Day] {
		validDaysList := strings.Join([]string{
			"Monday", "Tuesday", "Wednesday", "Thursday",
			"Friday", "Saturday", "Sunday",
		}, ", ")
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DAY_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_VALUE, DAY_FIELD, validDaysList),
		})
	}

	if u.StartTime == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        START_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, START_TIME_FIELD),
		})
	} else if _, err := time.Parse("15:04:05", u.StartTime); err != nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        START_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, START_TIME_FIELD, "HH:MM:SS"),
		})
	}

	if u.EndTime == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        END_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, END_TIME_FIELD),
		})
	} else if _, err := time.Parse("15:04:05", u.EndTime); err != nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        END_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, END_TIME_FIELD, "HH:MM:SS"),
		})
	}

	return errorInfo
}

func (c *CreateRescheduleRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	validStatus := map[string]bool{
		"changed":   true,
		"cancelled": true,
	}

	if c.Date.IsZero() {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DATE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, DATE_FIELD),
		})
	} else if c.Date.Before(time.Now()) {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DATE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_FUTURE_DATE, DATE_FIELD),
		})
	}

	if c.StartTime == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        START_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, START_TIME_FIELD),
		})
	} else if _, err := time.Parse("15:04:05", c.StartTime); err != nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        START_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, START_TIME_FIELD, "HH:MM:SS"),
		})
	}

	if c.EndTime == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        END_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, END_TIME_FIELD),
		})
	} else if _, err := time.Parse("15:04:05", c.EndTime); err != nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        END_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, END_TIME_FIELD, "HH:MM:SS"),
		})
	}

	if c.Status == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        STATUS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, STATUS_FIELD),
		})
	} else if !validStatus[strings.ToLower(c.Status)] {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        STATUS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_VALUE, STATUS_FIELD, "changed, cancelled"),
		})
	}

	return errorInfo
}

func (u *UpdateRescheduleRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if u.Date.IsZero() {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DATE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, DATE_FIELD),
		})
	} else if u.Date.Before(time.Now().Truncate(24*time.Hour)) {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DATE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_FUTURE_DATE, DATE_FIELD),
		})
	}

	if u.StartTime == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        START_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, START_TIME_FIELD),
		})
	} else if _, err := time.Parse("15:04:05", u.StartTime); err != nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        START_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, START_TIME_FIELD, "HH:MM:SS"),
		})
	}

	if u.EndTime == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        END_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, END_TIME_FIELD),
		})
	} else if _, err := time.Parse("15:04:05", u.EndTime); err != nil {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        END_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, END_TIME_FIELD, "HH:MM:SS"),
		})
	}

	if u.Status == constant.EMPTY_STRING {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        STATUS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, STATUS_FIELD),
		})
	} else if !isValidStatus(u.Status) {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        STATUS_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_VALUE, STATUS_FIELD, "active, inactive"),
		})
	}

	return errorInfo
}

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
	}
	return validStatuses[status]
}
