-- Remove seeded articles
DELETE FROM articles 
WHERE slug IN (
    'understanding-heart-health-prevention-and-treatment',
    'mental-health-and-wellness-in-modern-times',
    'managing-diabetes-comprehensive-guide'
);
