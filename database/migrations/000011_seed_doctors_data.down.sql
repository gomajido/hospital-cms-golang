-- Delete data in reverse order to handle foreign key constraints
DELETE FROM doctor_reschedules;
DELETE FROM doctor_schedules;
DELETE FROM doctors;
