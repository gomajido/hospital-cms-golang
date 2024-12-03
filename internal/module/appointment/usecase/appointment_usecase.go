package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gomajido/hospital-cms-golang/internal/module/appointment/constant"
	"github.com/gomajido/hospital-cms-golang/internal/module/appointment/domain"
)

type appointmentUsecase struct {
	appointmentRepo domain.AppointmentRepository
}

// NewAppointmentUsecase creates a new instance of appointmentUsecase
func NewAppointmentUsecase(ar domain.AppointmentRepository) domain.AppointmentUsecase {
	return &appointmentUsecase{
		appointmentRepo: ar,
	}
}

func (u *appointmentUsecase) Create(ctx context.Context, req domain.CreateAppointmentRequest) (*domain.Appointment, error) {
	// Parse appointment date
	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid appointment date format: %v", err)
	}

	// Check if doctor schedule is available
	available, err := u.appointmentRepo.CheckAvailability(ctx, &domain.CheckAvailabilityRequest{
		DoctorID:        req.DoctorID,
		ScheduleID:      req.ScheduleID,
		AppointmentDate: req.AppointmentDate, // Use the parsed appointmentDate,
		AppointmentTime: req.AppointmentTime,
	})
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, constant.ErrTimeSlotNotAvailable
	}

	// Create appointment
	appointment := &domain.Appointment{
		ID:              uuid.New(),
		UserID:          req.UserID,
		DoctorID:        req.DoctorID,
		ScheduleID:      req.ScheduleID,
		AppointmentDate: appointmentDate,
		AppointmentTime: req.AppointmentTime,
		Status:          constant.AppointmentStatusScheduled,
		Reason:          req.Reason,
		Notes:           req.Notes,
		RescheduleCount: 0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := u.appointmentRepo.Create(ctx, appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

func (u *appointmentUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	return u.appointmentRepo.GetByID(ctx, id)
}

func (u *appointmentUsecase) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.Appointment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return u.appointmentRepo.GetByUserID(ctx, userID, page, limit)
}

func (u *appointmentUsecase) GetByDoctorID(ctx context.Context, doctorID uuid.UUID, page, limit int) ([]domain.Appointment, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return u.appointmentRepo.GetByDoctorID(ctx, doctorID, page, limit)
}

func (u *appointmentUsecase) Cancel(ctx context.Context, id uuid.UUID, req domain.CancelAppointmentRequest) (*domain.Appointment, error) {
	// Get appointment
	appointment, err := u.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the appointment
	if appointment.UserID != req.UserID {
		return nil, fmt.Errorf("unauthorized: only the appointment owner can cancel")
	}

	// Check if appointment can be cancelled
	if appointment.Status != constant.AppointmentStatusScheduled {
		return nil, fmt.Errorf("appointment cannot be cancelled: invalid status")
	}

	// Update appointment status
	appointment.Status = constant.AppointmentStatusCancelled
	appointment.Notes = req.Notes
	appointment.UpdatedAt = time.Now()

	if err := u.appointmentRepo.Update(ctx, appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

func (u *appointmentUsecase) Reschedule(ctx context.Context, id uuid.UUID, req domain.RescheduleAppointmentRequest) (*domain.Appointment, error) {
	// Get appointment
	appointment, err := u.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user owns the appointment
	if appointment.UserID != req.UserID {
		return nil, fmt.Errorf("unauthorized: only the appointment owner can reschedule")
	}

	// Check if appointment can be rescheduled
	if appointment.Status != constant.AppointmentStatusScheduled {
		return nil, fmt.Errorf("appointment cannot be rescheduled: invalid status")
	}

	// Check reschedule count
	if appointment.RescheduleCount >= constant.MaxRescheduleCount {
		return nil, constant.ErrMaxReschedulesExceeded
	}

	// Parse new appointment date
	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return nil, fmt.Errorf("invalid appointment date format: %v", err)
	}

	// Check if new schedule is available
	available, err := u.appointmentRepo.CheckAvailability(ctx, &domain.CheckAvailabilityRequest{
		DoctorID:        appointment.DoctorID,
		ScheduleID:      req.ScheduleID,
		AppointmentDate: req.AppointmentDate,
		AppointmentTime: req.AppointmentTime,
	})
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, constant.ErrTimeSlotNotAvailable
	}

	// Update appointment
	appointment.ScheduleID = req.ScheduleID
	appointment.AppointmentDate = appointmentDate
	appointment.AppointmentTime = req.AppointmentTime
	appointment.Notes = req.Notes
	appointment.RescheduleCount++
	appointment.UpdatedAt = time.Now()

	if err := u.appointmentRepo.Update(ctx, appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

func (u *appointmentUsecase) CheckAvailability(ctx context.Context, req domain.CheckAvailabilityRequest) (bool, error) {
	// Parse appointment date
	_, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return false, fmt.Errorf("invalid appointment date format: %v", err)
	}

	return u.appointmentRepo.CheckAvailability(ctx, &domain.CheckAvailabilityRequest{
		DoctorID:        req.DoctorID,
		ScheduleID:      req.ScheduleID,
		AppointmentDate: req.AppointmentDate,
		AppointmentTime: req.AppointmentTime,
	})
}
