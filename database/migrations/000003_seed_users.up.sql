-- Get role IDs
SET @admin_role_id = (SELECT id FROM roles WHERE name = 'admin' LIMIT 1);
SET @member_role_id = (SELECT id FROM roles WHERE name = 'member' LIMIT 1);
SET @doctor_role_id = (SELECT id FROM roles WHERE name = 'doctor' LIMIT 1);
SET @nurse_role_id = (SELECT id FROM roles WHERE name = 'nurse' LIMIT 1);
SET @receptionist_role_id = (SELECT id FROM roles WHERE name = 'receptionist' LIMIT 1);
SET @patient_role_id = (SELECT id FROM roles WHERE name = 'patient' LIMIT 1);

-- Insert users with hashed password 'password123' ($2a$10$your_salt_here)
INSERT INTO users (id, name, email, password, phone, status, created_at, updated_at) VALUES
-- Admin users
(UUID(), 'Admin One', 'admin1@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567890', 'active', NOW(), NOW()),
(UUID(), 'Admin Two', 'admin2@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567891', 'active', NOW(), NOW()),

-- Member users
(UUID(), 'Member One', 'member1@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567892', 'active', NOW(), NOW()),
(UUID(), 'Member Two', 'member2@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567893', 'active', NOW(), NOW()),

-- Doctor users
(UUID(), 'Doctor One', 'doctor1@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567894', 'active', NOW(), NOW()),
(UUID(), 'Doctor Two', 'doctor2@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567895', 'active', NOW(), NOW()),

-- Nurse users
(UUID(), 'Nurse One', 'nurse1@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567896', 'active', NOW(), NOW()),
(UUID(), 'Nurse Two', 'nurse2@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567897', 'active', NOW(), NOW()),

-- Receptionist users
(UUID(), 'Receptionist One', 'receptionist1@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567898', 'active', NOW(), NOW()),
(UUID(), 'Receptionist Two', 'receptionist2@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567899', 'active', NOW(), NOW()),

-- Patient users
(UUID(), 'Patient One', 'patient1@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567900', 'active', NOW(), NOW()),
(UUID(), 'Patient Two', 'patient2@example.com', '$2a$10$NqPVxUcKGBxeVZvXUVvVeO9kcIQ1f.6GQTiZeQH/YVlqBvYl8LYhO', '+1234567901', 'active', NOW(), NOW());

-- Assign roles to users
INSERT INTO user_roles (id, user_id, role_id) 
SELECT UUID(), id, @admin_role_id FROM users WHERE email LIKE 'admin%@example.com'
UNION ALL
SELECT UUID(), id, @member_role_id FROM users WHERE email LIKE 'member%@example.com'
UNION ALL
SELECT UUID(), id, @doctor_role_id FROM users WHERE email LIKE 'doctor%@example.com'
UNION ALL
SELECT UUID(), id, @nurse_role_id FROM users WHERE email LIKE 'nurse%@example.com'
UNION ALL
SELECT UUID(), id, @receptionist_role_id FROM users WHERE email LIKE 'receptionist%@example.com'
UNION ALL
SELECT UUID(), id, @patient_role_id FROM users WHERE email LIKE 'patient%@example.com';
