package domain

import (
	"fmt"
	"time"

	"github.com/gomajido/hospital-cms-golang/internal/constant"
	"github.com/gomajido/hospital-cms-golang/internal/response"
)

const (
	// Field names for validation messages
	DOCTOR_ID_FIELD        = "doctor_id"
	SCHEDULE_ID_FIELD      = "doctor_schedule_id"
	APPOINTMENT_DATE_FIELD = "appointment_date"
	APPOINTMENT_TIME_FIELD = "appointment_time"
	REASON_FIELD           = "reason"
	STATUS_FIELD           = "status"
	PAGE_FIELD             = "page"
	LIMIT_FIELD            = "limit"
)

// Validate validates CreateAppointmentRequest
func (r *CreateAppointmentRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if r.DoctorID.String() == "00000000-0000-0000-0000-000000000000" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DOCTOR_ID_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, DOCTOR_ID_FIELD),
		})
	}

	if r.ScheduleID.String() == "00000000-0000-0000-0000-000000000000" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        SCHEDULE_ID_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, SCHEDULE_ID_FIELD),
		})
	}

	if r.AppointmentDate == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        APPOINTMENT_DATE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, APPOINTMENT_DATE_FIELD),
		})
	} else {
		_, err := time.Parse("2006-01-02", r.AppointmentDate)
		if err != nil {
			errorInfo = append(errorInfo, response.ErrorInfo{
				Field:        APPOINTMENT_DATE_FIELD,
				ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, APPOINTMENT_DATE_FIELD, "YYYY-MM-DD"),
			})
		}
	}

	if r.AppointmentTime == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        APPOINTMENT_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, APPOINTMENT_TIME_FIELD),
		})
	} else {
		_, err := time.Parse("15:04", r.AppointmentTime)
		if err != nil {
			errorInfo = append(errorInfo, response.ErrorInfo{
				Field:        APPOINTMENT_TIME_FIELD,
				ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, APPOINTMENT_TIME_FIELD, "HH:mm"),
			})
		}
	}

	if r.Reason == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        REASON_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, REASON_FIELD),
		})
	}

	return errorInfo
}

// Validate validates RescheduleAppointmentRequest
func (r *RescheduleAppointmentRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if r.AppointmentDate == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        APPOINTMENT_DATE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, APPOINTMENT_DATE_FIELD),
		})
	} else {
		_, err := time.Parse("2006-01-02", r.AppointmentDate)
		if err != nil {
			errorInfo = append(errorInfo, response.ErrorInfo{
				Field:        APPOINTMENT_DATE_FIELD,
				ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, APPOINTMENT_DATE_FIELD, "YYYY-MM-DD"),
			})
		}
	}

	if r.AppointmentTime == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        APPOINTMENT_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, APPOINTMENT_TIME_FIELD),
		})
	} else {
		_, err := time.Parse("15:04", r.AppointmentTime)
		if err != nil {
			errorInfo = append(errorInfo, response.ErrorInfo{
				Field:        APPOINTMENT_TIME_FIELD,
				ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, APPOINTMENT_TIME_FIELD, "HH:mm"),
			})
		}
	}

	if r.Reason == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        REASON_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, REASON_FIELD),
		})
	}

	return errorInfo
}

// Validate validates CancelAppointmentRequest
func (r *CancelAppointmentRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if r.Reason == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        REASON_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, REASON_FIELD),
		})
	}

	return errorInfo
}

// Validate validates CheckAvailabilityRequest
func (r *CheckAvailabilityRequest) Validate() []response.ErrorInfo {
	var errorInfo []response.ErrorInfo

	if r.DoctorID.String() == "00000000-0000-0000-0000-000000000000" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        DOCTOR_ID_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, DOCTOR_ID_FIELD),
		})
	}

	if r.ScheduleID.String() == "00000000-0000-0000-0000-000000000000" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        SCHEDULE_ID_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, SCHEDULE_ID_FIELD),
		})
	}

	if r.AppointmentDate == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        APPOINTMENT_DATE_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, APPOINTMENT_DATE_FIELD),
		})
	} else {
		_, err := time.Parse("2006-01-02", r.AppointmentDate)
		if err != nil {
			errorInfo = append(errorInfo, response.ErrorInfo{
				Field:        APPOINTMENT_DATE_FIELD,
				ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, APPOINTMENT_DATE_FIELD, "YYYY-MM-DD"),
			})
		}
	}

	if r.AppointmentTime == "" {
		errorInfo = append(errorInfo, response.ErrorInfo{
			Field:        APPOINTMENT_TIME_FIELD,
			ErrorMessage: fmt.Sprintf(constant.VALIDATION_REQUIRED, APPOINTMENT_TIME_FIELD),
		})
	} else {
		_, err := time.Parse("15:04", r.AppointmentTime)
		if err != nil {
			errorInfo = append(errorInfo, response.ErrorInfo{
				Field:        APPOINTMENT_TIME_FIELD,
				ErrorMessage: fmt.Sprintf(constant.VALIDATION_INVALID_FORMAT, APPOINTMENT_TIME_FIELD, "HH:mm"),
			})
		}
	}

	return errorInfo
}
