package domain

import (
	"github.com/google/uuid"
)

// CreateAppointmentRequest represents the request to create a new appointment
type CreateAppointmentRequest struct {
	UserID          uuid.UUID `json:"user_id" validate:"required"`
	DoctorID        uuid.UUID `json:"doctor_id" validate:"required"`
	ScheduleID      uuid.UUID `json:"doctor_schedule_id" validate:"required"`
	AppointmentDate string    `json:"appointment_date" validate:"required"`
	AppointmentTime string    `json:"appointment_time" validate:"required"`
	Reason          string    `json:"reason" validate:"required"`
	Notes           string    `json:"notes"`
}

// CancelAppointmentRequest represents the request to cancel an appointment
type CancelAppointmentRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Reason string    `json:"reason" validate:"required"`
	Notes  string    `json:"notes"`
}

// RescheduleAppointmentRequest represents the request to reschedule an appointment
type RescheduleAppointmentRequest struct {
	UserID          uuid.UUID `json:"user_id" validate:"required"`
	ScheduleID      uuid.UUID `json:"doctor_schedule_id" validate:"required"`
	AppointmentDate string    `json:"appointment_date" validate:"required"`
	AppointmentTime string    `json:"appointment_time" validate:"required"`
	Reason          string    `json:"reason" validate:"required"`
	Notes           string    `json:"notes"`
}

// CheckAvailabilityRequest represents the request to check doctor's availability
type CheckAvailabilityRequest struct {
	DoctorID        uuid.UUID `json:"doctor_id" validate:"required"`
	ScheduleID      uuid.UUID `json:"doctor_schedule_id" validate:"required"`
	AppointmentDate string    `json:"appointment_date" validate:"required"`
	AppointmentTime string    `json:"appointment_time" validate:"required"`
}
