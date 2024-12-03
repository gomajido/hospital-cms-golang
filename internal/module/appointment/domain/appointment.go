package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents minimal user information needed for appointments
type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// Doctor represents minimal doctor information needed for appointments
type Doctor struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Specialization string    `json:"specialization"`
	ServiceID      uuid.UUID `json:"service_id"`
}

// DoctorSchedule represents minimal schedule information needed for appointments
type DoctorSchedule struct {
	ID        uuid.UUID `json:"id"`
	DoctorID  uuid.UUID `json:"doctor_id"`
	Day       string    `json:"day"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
}

// Appointment represents an appointment between a user and a doctor
type Appointment struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	DoctorID        uuid.UUID       `json:"doctor_id"`
	ScheduleID      uuid.UUID       `json:"doctor_schedule_id"`
	AppointmentDate time.Time       `json:"appointment_date"`
	AppointmentTime string          `json:"appointment_time"`
	Status          string          `json:"status"`
	Reason          string          `json:"reason"`
	Notes           string          `json:"notes,omitempty"`
	RescheduleCount int             `json:"reschedule_count"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	User            *User           `json:"user,omitempty"`
	Doctor          *Doctor         `json:"doctor,omitempty"`
	Schedule        *DoctorSchedule `json:"schedule,omitempty"`
}

// AppointmentRepository defines the interface for appointment data operations
type AppointmentRepository interface {
	Create(ctx context.Context, appointment *Appointment) error
	GetByID(ctx context.Context, id uuid.UUID) (*Appointment, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]Appointment, int64, error)
	GetByDoctorID(ctx context.Context, doctorID uuid.UUID, page, limit int) ([]Appointment, int64, error)
	Update(ctx context.Context, appointment *Appointment) error
	Cancel(ctx context.Context, id uuid.UUID, req *CancelAppointmentRequest) (*Appointment, error)
	Reschedule(ctx context.Context, id uuid.UUID, date time.Time, timeSlot string) error
	CheckAvailability(ctx context.Context, req *CheckAvailabilityRequest) (bool, error)
}

// AppointmentUsecase defines the interface for appointment business logic
type AppointmentUsecase interface {
	Create(ctx context.Context, req CreateAppointmentRequest) (*Appointment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Appointment, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]Appointment, int64, error)
	GetByDoctorID(ctx context.Context, doctorID uuid.UUID, page, limit int) ([]Appointment, int64, error)
	Cancel(ctx context.Context, id uuid.UUID, req CancelAppointmentRequest) (*Appointment, error)
	Reschedule(ctx context.Context, id uuid.UUID, req RescheduleAppointmentRequest) (*Appointment, error)
	CheckAvailability(ctx context.Context, req CheckAvailabilityRequest) (bool, error)
}
