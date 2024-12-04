package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/gomajido/hospital-cms-golang/internal/module/appointment/constant"
	"github.com/gomajido/hospital-cms-golang/internal/module/appointment/domain"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) domain.AppointmentRepository {
	return &AppointmentRepository{
		db: db,
	}
}

// Create creates a new appointment
func (r *AppointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	query := `INSERT INTO appointments (
		id, user_id, doctor_id, doctor_schedule_id, appointment_date,
		appointment_time, status, reason, notes, reschedule_count,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	appointment.CreatedAt = now
	appointment.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		appointment.ID, appointment.UserID, appointment.DoctorID,
		appointment.ScheduleID, appointment.AppointmentDate,
		appointment.AppointmentTime, appointment.Status,
		appointment.Reason, appointment.Notes,
		appointment.RescheduleCount, appointment.CreatedAt,
		appointment.UpdatedAt,
	)

	return err
}

// GetByID gets an appointment by ID with related data
func (r *AppointmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	query := `
		SELECT 
			a.id, a.user_id, a.doctor_id, a.doctor_schedule_id,
			a.appointment_date, a.appointment_time, a.status,
			a.reason, a.notes, a.reschedule_count,
			a.created_at, a.updated_at,
			u.name as user_name, u.email as user_email,
			d.name as doctor_name, d.specialization as doctor_specialization,
			d.service_id as doctor_service_id,
			ds.day as schedule_day, ds.start_time as schedule_start_time,
			ds.end_time as schedule_end_time
		FROM appointments a
		LEFT JOIN users u ON a.user_id = u.id
		LEFT JOIN doctors d ON a.doctor_id = d.id
		LEFT JOIN doctor_schedules ds ON a.doctor_schedule_id = ds.id
		WHERE a.id = ?`

	appointment := &domain.Appointment{}
	var userName, userEmail string
	var doctorName, doctorSpecialization string
	var doctorServiceID uuid.UUID
	var scheduleDay, scheduleStartTime, scheduleEndTime string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&appointment.ID, &appointment.UserID, &appointment.DoctorID,
		&appointment.ScheduleID, &appointment.AppointmentDate,
		&appointment.AppointmentTime, &appointment.Status,
		&appointment.Reason, &appointment.Notes,
		&appointment.RescheduleCount, &appointment.CreatedAt,
		&appointment.UpdatedAt,
		&userName, &userEmail,
		&doctorName, &doctorSpecialization, &doctorServiceID,
		&scheduleDay, &scheduleStartTime, &scheduleEndTime,
	)

	if err == sql.ErrNoRows {
		return nil, constant.ErrAppointmentNotFound
	}
	if err != nil {
		return nil, err
	}

	// Set related data
	appointment.User = &domain.User{
		ID:    appointment.UserID,
		Name:  userName,
		Email: userEmail,
	}

	appointment.Doctor = &domain.Doctor{
		ID:             appointment.DoctorID,
		Name:           doctorName,
		Specialization: doctorSpecialization,
		ServiceID:      doctorServiceID,
	}

	appointment.Schedule = &domain.DoctorSchedule{
		ID:        appointment.ScheduleID,
		DoctorID:  appointment.DoctorID,
		Day:       scheduleDay,
		StartTime: scheduleStartTime,
		EndTime:   scheduleEndTime,
	}

	return appointment, nil
}

// GetByUserID gets appointments for a user with pagination
func (r *AppointmentRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.Appointment, int64, error) {
	var appointments []domain.Appointment
	var total int64
	offset := (page - 1) * limit

	// Get total count
	countQuery := "SELECT COUNT(*) FROM appointments WHERE user_id = ?"
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get appointments with related data
	query := `
		SELECT 
			a.id, a.user_id, a.doctor_id, a.doctor_schedule_id,
			a.appointment_date, a.appointment_time, a.status,
			a.reason, a.notes, a.reschedule_count,
			a.created_at, a.updated_at,
			u.name as user_name, u.email as user_email,
			d.name as doctor_name, d.specialization as doctor_specialization,
			d.service_id as doctor_service_id,
			ds.day as schedule_day, ds.start_time as schedule_start_time,
			ds.end_time as schedule_end_time
		FROM appointments a
		LEFT JOIN users u ON a.user_id = u.id
		LEFT JOIN doctors d ON a.doctor_id = d.id
		LEFT JOIN doctor_schedules ds ON a.doctor_schedule_id = ds.id
		WHERE a.user_id = ?
		ORDER BY a.appointment_date DESC, a.appointment_time DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var appointment domain.Appointment
		var userName, userEmail string
		var doctorName, doctorSpecialization string
		var doctorServiceID uuid.UUID
		var scheduleDay, scheduleStartTime, scheduleEndTime string

		err := rows.Scan(
			&appointment.ID, &appointment.UserID, &appointment.DoctorID,
			&appointment.ScheduleID, &appointment.AppointmentDate,
			&appointment.AppointmentTime, &appointment.Status,
			&appointment.Reason, &appointment.Notes,
			&appointment.RescheduleCount, &appointment.CreatedAt,
			&appointment.UpdatedAt,
			&userName, &userEmail,
			&doctorName, &doctorSpecialization, &doctorServiceID,
			&scheduleDay, &scheduleStartTime, &scheduleEndTime,
		)
		if err != nil {
			return nil, 0, err
		}

		// Set related data
		appointment.User = &domain.User{
			ID:    appointment.UserID,
			Name:  userName,
			Email: userEmail,
		}

		appointment.Doctor = &domain.Doctor{
			ID:             appointment.DoctorID,
			Name:           doctorName,
			Specialization: doctorSpecialization,
			ServiceID:      doctorServiceID,
		}

		appointment.Schedule = &domain.DoctorSchedule{
			ID:        appointment.ScheduleID,
			DoctorID:  appointment.DoctorID,
			Day:       scheduleDay,
			StartTime: scheduleStartTime,
			EndTime:   scheduleEndTime,
		}

		appointments = append(appointments, appointment)
	}

	return appointments, total, nil
}

// GetByDoctorID gets appointments for a doctor with pagination
func (r *AppointmentRepository) GetByDoctorID(ctx context.Context, doctorID uuid.UUID, page, limit int) ([]domain.Appointment, int64, error) {
	var appointments []domain.Appointment
	var total int64
	offset := (page - 1) * limit

	// Get total count
	countQuery := "SELECT COUNT(*) FROM appointments WHERE doctor_id = ?"
	err := r.db.QueryRowContext(ctx, countQuery, doctorID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get appointments with related data
	query := `
		SELECT 
			a.id, a.user_id, a.doctor_id, a.doctor_schedule_id,
			a.appointment_date, a.appointment_time, a.status,
			a.reason, a.notes, a.reschedule_count,
			a.created_at, a.updated_at,
			u.name as user_name, u.email as user_email,
			d.name as doctor_name, d.specialization as doctor_specialization,
			d.service_id as doctor_service_id,
			ds.day as schedule_day, ds.start_time as schedule_start_time,
			ds.end_time as schedule_end_time
		FROM appointments a
		LEFT JOIN users u ON a.user_id = u.id
		LEFT JOIN doctors d ON a.doctor_id = d.id
		LEFT JOIN doctor_schedules ds ON a.doctor_schedule_id = ds.id
		WHERE a.doctor_id = ?
		ORDER BY a.appointment_date DESC, a.appointment_time DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, doctorID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var appointment domain.Appointment
		var userName, userEmail string
		var doctorName, doctorSpecialization string
		var doctorServiceID uuid.UUID
		var scheduleDay, scheduleStartTime, scheduleEndTime string

		err := rows.Scan(
			&appointment.ID, &appointment.UserID, &appointment.DoctorID,
			&appointment.ScheduleID, &appointment.AppointmentDate,
			&appointment.AppointmentTime, &appointment.Status,
			&appointment.Reason, &appointment.Notes,
			&appointment.RescheduleCount, &appointment.CreatedAt,
			&appointment.UpdatedAt,
			&userName, &userEmail,
			&doctorName, &doctorSpecialization, &doctorServiceID,
			&scheduleDay, &scheduleStartTime, &scheduleEndTime,
		)
		if err != nil {
			return nil, 0, err
		}

		// Set related data
		appointment.User = &domain.User{
			ID:    appointment.UserID,
			Name:  userName,
			Email: userEmail,
		}

		appointment.Doctor = &domain.Doctor{
			ID:             appointment.DoctorID,
			Name:           doctorName,
			Specialization: doctorSpecialization,
			ServiceID:      doctorServiceID,
		}

		appointment.Schedule = &domain.DoctorSchedule{
			ID:        appointment.ScheduleID,
			DoctorID:  appointment.DoctorID,
			Day:       scheduleDay,
			StartTime: scheduleStartTime,
			EndTime:   scheduleEndTime,
		}

		appointments = append(appointments, appointment)
	}

	return appointments, total, nil
}

// Update updates an appointment
func (r *AppointmentRepository) Update(ctx context.Context, appointment *domain.Appointment) error {
	query := `UPDATE appointments SET
		appointment_date = ?, appointment_time = ?, status = ?,
		reason = ?, notes = ?, reschedule_count = ?, updated_at = ?
		WHERE id = ?`

	appointment.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		appointment.AppointmentDate, appointment.AppointmentTime,
		appointment.Status, appointment.Reason, appointment.Notes,
		appointment.RescheduleCount, appointment.UpdatedAt,
		appointment.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return constant.ErrAppointmentNotFound
	}

	return nil
}

// Cancel cancels an appointment
func (r *AppointmentRepository) Cancel(ctx context.Context, id uuid.UUID, req *domain.CancelAppointmentRequest) (*domain.Appointment, error) {
	// Get the appointment first to verify it exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the appointment status``
	query := `UPDATE appointments SET
		status = ?, reason = ?, notes = ?, updated_at = ?
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		constant.AppointmentStatusCancelled,
		req.Reason,
		req.Notes,
		time.Now(),
		id,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, constant.ErrAppointmentNotFound
	}

	// Get the updated appointment with all relations
	updatedAppointment, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedAppointment, nil
}

// Reschedule updates appointment date and time, and increments reschedule count
func (r *AppointmentRepository) Reschedule(ctx context.Context, id uuid.UUID, date time.Time, timeSlot string) error {
	query := `UPDATE appointments SET
		appointment_date = ?, appointment_time = ?,
		reschedule_count = reschedule_count + 1,
		updated_at = ?
		WHERE id = ? AND reschedule_count < ?`

	result, err := r.db.ExecContext(ctx, query,
		date, timeSlot, time.Now(),
		id, constant.MaxRescheduleCount,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Check if appointment exists
		var exists bool
		err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM appointments WHERE id = ?)", id).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			return constant.ErrAppointmentNotFound
		}
		return constant.ErrMaxReschedulesExceeded
	}

	return nil
}

// CheckAvailability checks if a time slot is available for a doctor
func (r *AppointmentRepository) CheckAvailability(ctx context.Context, req *domain.CheckAvailabilityRequest) (bool, error) {
	// Parse appointment date
	appointmentDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	if err != nil {
		return false, err
	}

	var count int
	query := `SELECT COUNT(*) FROM appointments 
		WHERE doctor_id = ? 
		AND doctor_schedule_id = ? 
		AND appointment_date = ? 
		AND appointment_time = ?
		AND status = ?`

	err = r.db.QueryRowContext(ctx, query,
		req.DoctorID,
		req.ScheduleID,
		appointmentDate,
		req.AppointmentTime,
		constant.AppointmentStatusScheduled,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
