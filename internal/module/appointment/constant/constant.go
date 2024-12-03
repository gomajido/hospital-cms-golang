package constant

import "errors"

const (
	AppointmentStatusScheduled = "scheduled"
	AppointmentStatusCompleted = "completed"
	AppointmentStatusCancelled = "cancelled"
	MaxRescheduleCount        = 3
)

// Common errors for appointment module
var (
	ErrTimeSlotNotAvailable   = errors.New("time slot is not available")
	ErrMaxReschedulesExceeded = errors.New("maximum number of reschedules exceeded")
	ErrAppointmentNotFound    = errors.New("appointment not found")
	ErrInvalidAppointmentDate = errors.New("invalid appointment date")
	ErrInvalidAppointmentTime = errors.New("invalid appointment time")
)
