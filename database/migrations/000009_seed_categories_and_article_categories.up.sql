-- Insert health and hospital related categories if they don't exist
INSERT IGNORE INTO categories (id, name, slug, description, created_at, updated_at) VALUES
(UUID(), 'Medical News', 'medical-news', 'Latest medical news and breakthroughs in healthcare', NOW(), NOW()),
(UUID(), 'Health Tips', 'health-tips', 'General health tips and wellness advice', NOW(), NOW()),
(UUID(), 'Hospital Services', 'hospital-services', 'Information about various hospital services and departments', NOW(), NOW()),
(UUID(), 'Patient Care', 'patient-care', 'Guidelines and information about patient care', NOW(), NOW()),
(UUID(), 'Medical Research', 'medical-research', 'Latest research and studies in medicine', NOW(), NOW()),
(UUID(), 'Healthcare Technology', 'healthcare-technology', 'Innovations and technology in healthcare', NOW(), NOW()),
(UUID(), 'Mental Health', 'mental-health', 'Mental health awareness and resources', NOW(), NOW()),
(UUID(), 'Nutrition', 'nutrition', 'Diet and nutrition information for better health', NOW(), NOW()),
(UUID(), 'Emergency Care', 'emergency-care', 'Emergency medical services and first aid information', NOW(), NOW()),
(UUID(), 'Preventive Care', 'preventive-care', 'Preventive healthcare measures and screenings', NOW(), NOW());

-- Insert seeded articles into categories
INSERT IGNORE INTO article_categories (id, article_id, category_id) VALUES
('550e8400-e29b-41d4-a716-446655440000','123e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'medical-news')),
('550e8400-e29b-41d4-a716-446655440001','223e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'health-tips')),
('550e8400-e29b-41d4-a716-446655440002','323e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'hospital-services')),
('550e8400-e29b-41d4-a716-446655440003','123e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'patient-care')),
('550e8400-e29b-41d4-a716-446655440004','223e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'medical-research')),
('550e8400-e29b-41d4-a716-446655440005','323e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'healthcare-technology')),
('550e8400-e29b-41d4-a716-446655440006','123e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'mental-health')),
('550e8400-e29b-41d4-a716-446655440007','223e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'nutrition')),
('550e8400-e29b-41d4-a716-446655440008','323e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'emergency-care')),
('550e8400-e29b-41d4-a716-446655440009','123e4567-e89b-12d3-a456-426614174000', (SELECT id FROM categories WHERE slug = 'preventive-care'));
