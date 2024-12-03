package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service represents medical services offered
type Service struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// Doctor represents the doctor entity
type Doctor struct {
	ID             uuid.UUID        `json:"id"`
	Name           string           `json:"name"`
	ServiceID      uuid.UUID        `json:"service_id"`
	Description    string           `json:"description"`
	Specialization string           `json:"specialization"`
	Degree         string           `json:"degree"`
	Experience     string           `json:"experience"`
	Service        *Service         `json:"service"`
	Schedules      []DoctorSchedule `json:"schedules,omitempty"`
}

// DoctorSchedule represents the doctor's regular schedule
type DoctorSchedule struct {
	ID          uuid.UUID          `json:"id"`
	DoctorID    uuid.UUID          `json:"doctor_id"`
	Day         string             `json:"day"`
	StartTime   string             `json:"start_time"`
	EndTime     string             `json:"end_time"`
	Reschedules []DoctorReschedule `json:"reschedules,omitempty"`
}

// DoctorReschedule represents schedule changes or cancellations
type DoctorReschedule struct {
	ID               uuid.UUID       `json:"id"`
	DoctorScheduleID uuid.UUID       `json:"doctor_schedule_id"`
	Date             time.Time       `json:"date"`
	StartTime        string          `json:"start_time"`
	EndTime          string          `json:"end_time"`
	Status           string          `json:"status"`
	Description      string          `json:"description"`
	Schedule         *DoctorSchedule `json:"schedule"`
}

// DoctorRepository defines the interface for doctor data operations
type DoctorRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Doctor, error)
	List(ctx context.Context, page, limit int) ([]Doctor, int64, error)
	Create(ctx context.Context, doctor *Doctor) error
	Update(ctx context.Context, doctor *Doctor) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Schedule operations
	CreateSchedule(ctx context.Context, schedule *DoctorSchedule) error
	GetSchedulesByDoctorID(ctx context.Context, doctorID uuid.UUID) ([]DoctorSchedule, error)
	UpdateSchedule(ctx context.Context, schedule *DoctorSchedule) error
	DeleteSchedule(ctx context.Context, id uuid.UUID) error

	// Reschedule operations
	CreateReschedule(ctx context.Context, reschedule *DoctorReschedule) error
	GetReschedulesByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]DoctorReschedule, error)
	UpdateReschedule(ctx context.Context, reschedule *DoctorReschedule) error
	DeleteReschedule(ctx context.Context, id uuid.UUID) error
}

// DoctorUsecase defines the interface for doctor business logic
type DoctorUsecase interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Doctor, error)
	List(ctx context.Context, page, limit int) ([]Doctor, int64, error)
	Create(ctx context.Context, req CreateDoctorRequest) (*Doctor, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateDoctorRequest) (*Doctor, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Schedule operations
	CreateSchedule(ctx context.Context, doctorID uuid.UUID, req CreateScheduleRequest) (*DoctorSchedule, error)
	GetSchedulesByDoctorID(ctx context.Context, doctorID uuid.UUID) ([]DoctorSchedule, error)
	UpdateSchedule(ctx context.Context, id uuid.UUID, req UpdateScheduleRequest) (*DoctorSchedule, error)
	DeleteSchedule(ctx context.Context, id uuid.UUID) error

	// Reschedule operations
	CreateReschedule(ctx context.Context, scheduleID uuid.UUID, req CreateRescheduleRequest) (*DoctorReschedule, error)
	GetReschedulesByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]DoctorReschedule, error)
	UpdateReschedule(ctx context.Context, id uuid.UUID, req UpdateRescheduleRequest) (*DoctorReschedule, error)
	DeleteReschedule(ctx context.Context, id uuid.UUID) error
}
