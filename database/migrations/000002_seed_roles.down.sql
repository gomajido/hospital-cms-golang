-- Remove seeded roles
DELETE FROM roles WHERE name IN ('admin', 'doctor', 'nurse', 'receptionist', 'patient');
