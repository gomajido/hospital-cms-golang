-- Insert default roles
INSERT INTO roles (id, name, description) VALUES
    (UUID(), 'admin', 'Administrator with full access'),
    (UUID(), 'member', 'Member'),
    (UUID(), 'doctor', 'Medical doctor with patient access'),
    (UUID(), 'nurse', 'Nurse with limited patient access'),
    (UUID(), 'receptionist', 'Front desk staff'),
    (UUID(), 'patient', 'Patient user')
ON DUPLICATE KEY UPDATE
    description = VALUES(description);
