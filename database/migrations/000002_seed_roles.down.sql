-- Remove seeded roles
DELETE FROM roles WHERE name IN ('admin', 'member','doctor', 'nurse', 'receptionist', 'patient');
