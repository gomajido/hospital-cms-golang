-- Create services table
CREATE TABLE IF NOT EXISTS services (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create doctors table
CREATE TABLE IF NOT EXISTS doctors (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    service_id CHAR(36) NOT NULL,
    description TEXT,
    specialization VARCHAR(255) NOT NULL,
    degree VARCHAR(255) NOT NULL,
    experience VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

-- Create doctor schedules table
CREATE TABLE IF NOT EXISTS doctor_schedules (
    id CHAR(36) PRIMARY KEY,
    doctor_id CHAR(36) NOT NULL,
    day VARCHAR(20) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
    CHECK (day IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'))
);

-- Create doctor reschedules table
CREATE TABLE IF NOT EXISTS doctor_reschedules (
    id CHAR(36) PRIMARY KEY,
    doctor_schedule_id CHAR(36) NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status ENUM('changed', 'cancelled') NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_schedule_id) REFERENCES doctor_schedules(id) ON DELETE CASCADE
);

-- Insert initial services
INSERT INTO services (id, name, description) VALUES 
    (UUID(), 'General Medicine', 'Primary healthcare services for general medical conditions'),
    (UUID(), 'Cardiology', 'Specialized care for heart and cardiovascular conditions'),
    (UUID(), 'Pediatrics', 'Medical care for infants, children, and adolescents'),
    (UUID(), 'Orthopedics', 'Treatment for musculoskeletal conditions and injuries'),
    (UUID(), 'Dermatology', 'Care for skin, hair, and nail conditions'),
    (UUID(), 'Neurology', 'Treatment for disorders of the nervous system'),
    (UUID(), 'Psychiatry', 'Mental health care and psychological support'),
    (UUID(), 'Obstetrics & Gynecology', 'Women''s health and reproductive care'),
    (UUID(), 'Ophthalmology', 'Eye care and vision services'),
    (UUID(), 'ENT', 'Treatment for ear, nose, and throat conditions');
