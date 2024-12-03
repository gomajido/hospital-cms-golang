package usecase

import (
	"context"
	"fmt"

	"github.com/gomajido/hospital-cms-golang/internal/module/doctor/domain"
	"github.com/google/uuid"
)

type doctorUsecase struct {
	doctorRepo domain.DoctorRepository
}

func NewDoctorUsecase(dr domain.DoctorRepository) domain.DoctorUsecase {
	return &doctorUsecase{
		doctorRepo: dr,
	}
}

func (u *doctorUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Doctor, error) {
	return u.doctorRepo.GetByID(ctx, id)
}

func (u *doctorUsecase) List(ctx context.Context, page, limit int) ([]domain.Doctor, int64, error) {
	return u.doctorRepo.List(ctx, page, limit)
}

func (u *doctorUsecase) Create(ctx context.Context, req domain.CreateDoctorRequest) (*domain.Doctor, error) {
	doctor := &domain.Doctor{
		ID:             uuid.New(),
		Name:           req.Name,
		ServiceID:      req.ServiceID,
		Description:    req.Description,
		Specialization: req.Specialization,
		Degree:         req.Degree,
		Experience:     req.Experience,
	}

	if err := u.doctorRepo.Create(ctx, doctor); err != nil {
		return nil, err
	}

	return doctor, nil
}

func (u *doctorUsecase) Update(ctx context.Context, id uuid.UUID, req domain.UpdateDoctorRequest) (*domain.Doctor, error) {
	doctor, err := u.doctorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	doctor.Name = req.Name
	doctor.ServiceID = req.ServiceID
	doctor.Description = req.Description
	doctor.Specialization = req.Specialization
	doctor.Degree = req.Degree
	doctor.Experience = req.Experience

	if err := u.doctorRepo.Update(ctx, doctor); err != nil {
		return nil, err
	}

	return doctor, nil
}

func (u *doctorUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.doctorRepo.Delete(ctx, id)
}

// Schedule operations
func (u *doctorUsecase) CreateSchedule(ctx context.Context, doctorID uuid.UUID, req domain.CreateScheduleRequest) (*domain.DoctorSchedule, error) {
	schedule := &domain.DoctorSchedule{
		ID:        uuid.New(),
		DoctorID:  doctorID,
		Day:       req.Day,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := u.doctorRepo.CreateSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (u *doctorUsecase) GetSchedulesByDoctorID(ctx context.Context, doctorID uuid.UUID) ([]domain.DoctorSchedule, error) {
	return u.doctorRepo.GetSchedulesByDoctorID(ctx, doctorID)
}

func (u *doctorUsecase) UpdateSchedule(ctx context.Context, id uuid.UUID, req domain.UpdateScheduleRequest) (*domain.DoctorSchedule, error) {
	schedule := &domain.DoctorSchedule{
		ID:        id,
		Day:       req.Day,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := u.doctorRepo.UpdateSchedule(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return schedule, nil
}

func (u *doctorUsecase) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	return u.doctorRepo.DeleteSchedule(ctx, id)
}

// Reschedule operations
func (u *doctorUsecase) CreateReschedule(ctx context.Context, scheduleID uuid.UUID, req domain.CreateRescheduleRequest) (*domain.DoctorReschedule, error) {
	reschedule := &domain.DoctorReschedule{
		ID:               uuid.New(),
		DoctorScheduleID: scheduleID,
		Date:             req.Date,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Status:           req.Status,
		Description:      req.Description,
	}

	if err := u.doctorRepo.CreateReschedule(ctx, reschedule); err != nil {
		return nil, err
	}

	return reschedule, nil
}

func (u *doctorUsecase) GetReschedulesByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]domain.DoctorReschedule, error) {
	return u.doctorRepo.GetReschedulesByScheduleID(ctx, scheduleID)
}

func (u *doctorUsecase) UpdateReschedule(ctx context.Context, id uuid.UUID, req domain.UpdateRescheduleRequest) (*domain.DoctorReschedule, error) {
	reschedule := &domain.DoctorReschedule{
		ID:          id,
		Date:        req.Date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      req.Status,
		Description: req.Description,
	}

	if err := u.doctorRepo.UpdateReschedule(ctx, reschedule); err != nil {
		return nil, fmt.Errorf("failed to update reschedule: %w", err)
	}

	return reschedule, nil
}

func (u *doctorUsecase) DeleteReschedule(ctx context.Context, id uuid.UUID) error {
	return u.doctorRepo.DeleteReschedule(ctx, id)
}
