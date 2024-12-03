package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gomajido/hospital-cms-golang/internal/module/doctor/domain"
	"github.com/google/uuid"
)

type doctorRepository struct {
	db *sql.DB
}

func NewDoctorRepository(db *sql.DB) domain.DoctorRepository {
	return &doctorRepository{
		db: db,
	}
}

func (r *doctorRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Doctor, error) {
	doctor := &domain.Doctor{}
	service := &domain.Service{}

	// Get doctor and service data
	query := `
		SELECT 
			d.id, d.name, d.service_id, d.description, d.specialization,
			d.degree, d.experience,
			s.id, s.name, s.description
		FROM doctors d
		LEFT JOIN services s ON d.service_id = s.id
		WHERE d.id = ?`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&doctor.ID, &doctor.Name, &doctor.ServiceID, &doctor.Description,
		&doctor.Specialization, &doctor.Degree, &doctor.Experience,
		&service.ID, &service.Name, &service.Description,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("doctor not found")
	}
	if err != nil {
		return nil, err
	}

	doctor.Service = service

	// Get doctor schedules
	schedules, err := r.GetSchedulesByDoctorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting doctor schedules: %v", err)
	}

	// For each schedule, get future reschedules
	for i := range schedules {
		query = `
			SELECT id, doctor_schedule_id, date, start_time, end_time, status, description
			FROM doctor_reschedules
			WHERE doctor_schedule_id = ? AND date >= CURDATE()
			ORDER BY date ASC`

		rows, err := r.db.QueryContext(ctx, query, schedules[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error getting schedule reschedules: %v", err)
		}
		defer rows.Close()

		var reschedules []domain.DoctorReschedule
		for rows.Next() {
			var reschedule domain.DoctorReschedule
			err := rows.Scan(
				&reschedule.ID,
				&reschedule.DoctorScheduleID,
				&reschedule.Date,
				&reschedule.StartTime,
				&reschedule.EndTime,
				&reschedule.Status,
				&reschedule.Description,
			)
			if err != nil {
				return nil, fmt.Errorf("error scanning reschedule: %v", err)
			}
			reschedules = append(reschedules, reschedule)
		}

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating reschedules: %v", err)
		}

		schedules[i].Reschedules = reschedules
	}

	doctor.Schedules = schedules
	return doctor, nil
}

func (r *doctorRepository) List(ctx context.Context, page, limit int) ([]domain.Doctor, int64, error) {
	var total int64
	var doctors []domain.Doctor

	// Get total count
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM doctors").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get doctors with services
	offset := (page - 1) * limit
	query := `
		SELECT 
			d.id, d.name, d.service_id, d.description, d.specialization,
			d.degree, d.experience,
			s.id, s.name, s.description
		FROM doctors d
		LEFT JOIN services s ON d.service_id = s.id
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		doctor := domain.Doctor{}
		service := domain.Service{}

		err := rows.Scan(
			&doctor.ID, &doctor.Name, &doctor.ServiceID, &doctor.Description,
			&doctor.Specialization, &doctor.Degree, &doctor.Experience,
			&service.ID, &service.Name, &service.Description,
		)
		if err != nil {
			return nil, 0, err
		}

		doctor.Service = &service
		doctors = append(doctors, doctor)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return doctors, total, nil
}

func (r *doctorRepository) Create(ctx context.Context, doctor *domain.Doctor) error {
	query := `
		INSERT INTO doctors (
			id, name, service_id, description, specialization,
			degree, experience, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	doctor.ID = uuid.New()

	_, err := r.db.ExecContext(ctx, query,
		doctor.ID, doctor.Name, doctor.ServiceID, doctor.Description,
		doctor.Specialization, doctor.Degree, doctor.Experience,
	)

	return err
}

func (r *doctorRepository) Update(ctx context.Context, doctor *domain.Doctor) error {
	query := `
		UPDATE doctors SET
			name = ?, service_id = ?, description = ?, specialization = ?,
			degree = ?, experience = ?, updated_at = NOW()
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		doctor.Name, doctor.ServiceID, doctor.Description,
		doctor.Specialization, doctor.Degree, doctor.Experience,
		doctor.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("doctor not found")
	}

	return nil
}

func (r *doctorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM doctors WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("doctor not found")
	}

	return nil
}

// Schedule operations
func (r *doctorRepository) CreateSchedule(ctx context.Context, schedule *domain.DoctorSchedule) error {
	query := `
		INSERT INTO doctor_schedules (
			id, doctor_id, day, start_time, end_time,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`

	schedule.ID = uuid.New()

	_, err := r.db.ExecContext(ctx, query,
		schedule.ID, schedule.DoctorID, schedule.Day,
		schedule.StartTime, schedule.EndTime,
	)

	return err
}

func (r *doctorRepository) GetSchedulesByDoctorID(ctx context.Context, doctorID uuid.UUID) ([]domain.DoctorSchedule, error) {
	query := `
		SELECT id, doctor_id, day, start_time, end_time
		FROM doctor_schedules
		WHERE doctor_id = ?
		ORDER BY CASE day
			WHEN 'Monday' THEN 1
			WHEN 'Tuesday' THEN 2
			WHEN 'Wednesday' THEN 3
			WHEN 'Thursday' THEN 4
			WHEN 'Friday' THEN 5
			WHEN 'Saturday' THEN 6
			WHEN 'Sunday' THEN 7
		END`

	rows, err := r.db.QueryContext(ctx, query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.DoctorSchedule
	for rows.Next() {
		var schedule domain.DoctorSchedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.DoctorID,
			&schedule.Day,
			&schedule.StartTime,
			&schedule.EndTime,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *doctorRepository) UpdateSchedule(ctx context.Context, schedule *domain.DoctorSchedule) error {
	query := `
		UPDATE doctor_schedules SET
			day = ?, start_time = ?, end_time = ?, updated_at = NOW()
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		schedule.Day, schedule.StartTime, schedule.EndTime,
		schedule.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("schedule not found")
	}

	return nil
}

func (r *doctorRepository) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM doctor_schedules WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("schedule not found")
	}

	return nil
}

// Reschedule operations
func (r *doctorRepository) CreateReschedule(ctx context.Context, reschedule *domain.DoctorReschedule) error {
	query := `
		INSERT INTO doctor_reschedules (
			id, doctor_schedule_id, date, start_time, end_time,
			status, description, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	reschedule.ID = uuid.New()

	_, err := r.db.ExecContext(ctx, query,
		reschedule.ID, reschedule.DoctorScheduleID, reschedule.Date,
		reschedule.StartTime, reschedule.EndTime, reschedule.Status,
		reschedule.Description,
	)

	return err
}

func (r *doctorRepository) GetReschedulesByScheduleID(ctx context.Context, scheduleID uuid.UUID) ([]domain.DoctorReschedule, error) {
	query := `
		SELECT 
			dr.id, dr.doctor_schedule_id, dr.date, dr.start_time,
			dr.end_time, dr.status, dr.description,
			ds.id, ds.doctor_id, ds.day, ds.start_time, ds.end_time
		FROM doctor_reschedules dr
		LEFT JOIN doctor_schedules ds ON dr.doctor_schedule_id = ds.id
		WHERE dr.doctor_schedule_id = ?`

	rows, err := r.db.QueryContext(ctx, query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reschedules []domain.DoctorReschedule

	for rows.Next() {
		reschedule := domain.DoctorReschedule{}
		schedule := domain.DoctorSchedule{}

		err := rows.Scan(
			&reschedule.ID, &reschedule.DoctorScheduleID, &reschedule.Date,
			&reschedule.StartTime, &reschedule.EndTime, &reschedule.Status,
			&reschedule.Description,
			&schedule.ID, &schedule.DoctorID, &schedule.Day,
			&schedule.StartTime, &schedule.EndTime,
		)
		if err != nil {
			return nil, err
		}

		reschedule.Schedule = &schedule
		reschedules = append(reschedules, reschedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reschedules, nil
}

func (r *doctorRepository) UpdateReschedule(ctx context.Context, reschedule *domain.DoctorReschedule) error {
	query := `
		UPDATE doctor_reschedules SET
			date = ?, start_time = ?, end_time = ?,
			status = ?, description = ?, updated_at = NOW()
		WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query,
		reschedule.Date, reschedule.StartTime, reschedule.EndTime,
		reschedule.Status, reschedule.Description, reschedule.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reschedule not found")
	}

	return nil
}

func (r *doctorRepository) DeleteReschedule(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM doctor_reschedules WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reschedule not found")
	}

	return nil
}
