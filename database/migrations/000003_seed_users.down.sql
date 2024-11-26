-- Remove role assignments
DELETE FROM user_roles 
WHERE user_id IN (
    SELECT id FROM users 
    WHERE email IN (
        'admin1@example.com', 'admin2@example.com',
        'member1@example.com', 'member2@example.com',
        'doctor1@example.com', 'doctor2@example.com',
        'nurse1@example.com', 'nurse2@example.com',
        'receptionist1@example.com', 'receptionist2@example.com',
        'patient1@example.com', 'patient2@example.com'
    )
);

-- Remove seeded users
DELETE FROM users 
WHERE email IN (
    'admin1@example.com', 'admin2@example.com',
    'member1@example.com', 'member2@example.com',
    'doctor1@example.com', 'doctor2@example.com',
    'nurse1@example.com', 'nurse2@example.com',
    'receptionist1@example.com', 'receptionist2@example.com',
    'patient1@example.com', 'patient2@example.com'
);
