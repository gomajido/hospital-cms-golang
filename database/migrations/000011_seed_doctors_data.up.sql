-- Get service IDs for reference
SET @general_medicine_id = (SELECT id FROM services WHERE name = 'General Medicine' LIMIT 1);
SET @cardiology_id = (SELECT id FROM services WHERE name = 'Cardiology' LIMIT 1);
SET @pediatrics_id = (SELECT id FROM services WHERE name = 'Pediatrics' LIMIT 1);
SET @orthopedics_id = (SELECT id FROM services WHERE name = 'Orthopedics' LIMIT 1);

-- Insert sample doctors
INSERT INTO doctors (id, name, service_id, description, specialization, degree, experience) VALUES
    (UUID(), 'Dr. John Smith', @general_medicine_id, 'Experienced general practitioner with focus on preventive care', 'General Medicine', 'MD, MBBS', '15 years'),
    (UUID(), 'Dr. Sarah Johnson', @cardiology_id, 'Specialized in interventional cardiology', 'Cardiology', 'MD, DM Cardiology', '12 years'),
    (UUID(), 'Dr. Michael Brown', @pediatrics_id, 'Pediatrician with expertise in newborn care', 'Pediatrics', 'MD, DCH', '10 years'),
    (UUID(), 'Dr. Emily Davis', @orthopedics_id, 'Specializes in sports injuries and joint replacements', 'Orthopedics', 'MS Ortho', '8 years');

-- Get doctor IDs for reference
SET @dr_smith_id = (SELECT id FROM doctors WHERE name = 'Dr. John Smith' LIMIT 1);
SET @dr_johnson_id = (SELECT id FROM doctors WHERE name = 'Dr. Sarah Johnson' LIMIT 1);
SET @dr_brown_id = (SELECT id FROM doctors WHERE name = 'Dr. Michael Brown' LIMIT 1);
SET @dr_davis_id = (SELECT id FROM doctors WHERE name = 'Dr. Emily Davis' LIMIT 1);

-- Insert sample schedules for Dr. Smith
INSERT INTO doctor_schedules (id, doctor_id, day, start_time, end_time) VALUES
    (UUID(), @dr_smith_id, 'Monday', '09:00:00', '17:00:00'),
    (UUID(), @dr_smith_id, 'Wednesday', '09:00:00', '17:00:00'),
    (UUID(), @dr_smith_id, 'Friday', '09:00:00', '13:00:00');

-- Insert sample schedules for Dr. Johnson
INSERT INTO doctor_schedules (id, doctor_id, day, start_time, end_time) VALUES
    (UUID(), @dr_johnson_id, 'Tuesday', '10:00:00', '18:00:00'),
    (UUID(), @dr_johnson_id, 'Thursday', '10:00:00', '18:00:00'),
    (UUID(), @dr_johnson_id, 'Saturday', '10:00:00', '14:00:00');

-- Insert sample schedules for Dr. Brown
INSERT INTO doctor_schedules (id, doctor_id, day, start_time, end_time) VALUES
    (UUID(), @dr_brown_id, 'Monday', '08:00:00', '16:00:00'),
    (UUID(), @dr_brown_id, 'Tuesday', '08:00:00', '16:00:00'),
    (UUID(), @dr_brown_id, 'Thursday', '08:00:00', '16:00:00');

-- Insert sample schedules for Dr. Davis
INSERT INTO doctor_schedules (id, doctor_id, day, start_time, end_time) VALUES
    (UUID(), @dr_davis_id, 'Wednesday', '09:00:00', '17:00:00'),
    (UUID(), @dr_davis_id, 'Friday', '09:00:00', '17:00:00'),
    (UUID(), @dr_davis_id, 'Saturday', '09:00:00', '13:00:00');

-- Get schedule IDs for reference
SET @dr_smith_monday_id = (SELECT id FROM doctor_schedules WHERE doctor_id = @dr_smith_id AND day = 'Monday' LIMIT 1);
SET @dr_johnson_tuesday_id = (SELECT id FROM doctor_schedules WHERE doctor_id = @dr_johnson_id AND day = 'Tuesday' LIMIT 1);

-- Insert sample reschedules
INSERT INTO doctor_reschedules (id, doctor_schedule_id, date, start_time, end_time, status, description) VALUES
    (UUID(), @dr_smith_monday_id, DATE_ADD(CURDATE(), INTERVAL 1 WEEK), '10:00:00', '18:00:00', 'changed', 'Extended hours for next week'),
    (UUID(), @dr_smith_monday_id, DATE_ADD(CURDATE(), INTERVAL 2 WEEK), '00:00:00', '00:00:00', 'cancelled', 'Doctor on leave'),
    (UUID(), @dr_johnson_tuesday_id, DATE_ADD(CURDATE(), INTERVAL 1 WEEK), '11:00:00', '19:00:00', 'changed', 'Schedule adjusted for medical conference');
