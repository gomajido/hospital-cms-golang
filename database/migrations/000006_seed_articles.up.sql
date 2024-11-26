-- Get admin user ID (assuming it exists from previous seeds)
SET @admin_id = (SELECT id FROM users WHERE email = 'admin1@example.com' LIMIT 1);

-- Seed Articles
INSERT INTO articles (
    id,
    title,
    slug,
    content,
    excerpt,
    main_image,
    status,
    author_id,
    published_at,
    meta_title,
    meta_description,
    meta_keywords,
    canonical_url,
    focus_keyphrase,
    og_title,
    og_description,
    og_image
) VALUES
-- Article 1: Heart Health
(
    UUID(),
    'Understanding Heart Health: Prevention and Treatment',
    'understanding-heart-health-prevention-and-treatment',
    '<h2>What is Cardiovascular Health?</h2>
    <p>Cardiovascular health refers to the overall well-being of your heart and blood vessels. Maintaining good heart health is crucial for a long and healthy life.</p>
    
    <h2>Risk Factors</h2>
    <ul>
        <li>High blood pressure</li>
        <li>High cholesterol</li>
        <li>Smoking</li>
        <li>Obesity</li>
        <li>Physical inactivity</li>
        <li>Diabetes</li>
    </ul>
    
    <h2>Prevention Tips</h2>
    <p>Here are some key ways to maintain heart health:</p>
    <ul>
        <li>Regular exercise (at least 30 minutes daily)</li>
        <li>Healthy diet rich in fruits and vegetables</li>
        <li>Limited salt and saturated fat intake</li>
        <li>Regular blood pressure monitoring</li>
        <li>Stress management</li>
        <li>Adequate sleep</li>
    </ul>
    
    <h2>Treatment Options</h2>
    <p>Modern medicine offers various treatment options:</p>
    <ul>
        <li>Medications (statins, blood thinners)</li>
        <li>Lifestyle modifications</li>
        <li>Surgical procedures when necessary</li>
        <li>Cardiac rehabilitation programs</li>
    </ul>',
    'Learn about cardiovascular health, risk factors, prevention strategies, and modern treatment options. Discover how to maintain a healthy heart through lifestyle changes and medical interventions.',
    'https://example.com/images/heart-health.jpg',
    'published',
    @admin_id,
    NOW(),
    'Heart Health Guide: Prevention & Treatment | Hospital CMS',
    'Comprehensive guide to understanding heart health, prevention strategies, and treatment options. Learn how to maintain a healthy heart and prevent cardiovascular disease.',
    'heart health, cardiovascular disease, heart disease prevention, heart treatment, healthy heart',
    'https://example.com/articles/understanding-heart-health-prevention-and-treatment',
    'heart health prevention treatment',
    'Complete Guide to Heart Health and Disease Prevention',
    'Learn about heart health, prevention strategies, and modern treatment options. Expert advice on maintaining cardiovascular wellness.',
    'https://example.com/images/heart-health-social.jpg'
),

-- Article 2: Mental Health
(
    UUID(),
    'Mental Health and Wellness in Modern Times',
    'mental-health-and-wellness-in-modern-times',
    '<h2>Understanding Mental Health</h2>
    <p>Mental health encompasses emotional, psychological, and social well-being. It affects how we think, feel, act, and cope with life''s challenges.</p>
    
    <h2>Common Mental Health Conditions</h2>
    <ul>
        <li>Anxiety disorders</li>
        <li>Depression</li>
        <li>Stress-related conditions</li>
        <li>Sleep disorders</li>
        <li>Eating disorders</li>
    </ul>
    
    <h2>Signs and Symptoms</h2>
    <ul>
        <li>Persistent sadness or anxiety</li>
        <li>Changes in sleep patterns</li>
        <li>Loss of interest in activities</li>
        <li>Difficulty concentrating</li>
        <li>Physical symptoms without clear causes</li>
    </ul>
    
    <h2>Treatment Approaches</h2>
    <p>Modern mental health care offers various treatment options:</p>
    <ul>
        <li>Psychotherapy</li>
        <li>Medication when necessary</li>
        <li>Lifestyle modifications</li>
        <li>Support groups</li>
        <li>Mindfulness and meditation</li>
    </ul>',
    'Explore mental health in the modern world, including common conditions, symptoms, and treatment approaches. Learn about maintaining mental wellness and seeking professional help.',
    'https://example.com/images/mental-health.jpg',
    'published',
    @admin_id,
    NOW(),
    'Mental Health Guide: Understanding and Treatment | Hospital CMS',
    'Comprehensive guide to mental health, including common conditions, symptoms, and modern treatment approaches. Learn about maintaining mental wellness.',
    'mental health, mental wellness, anxiety, depression, mental health treatment',
    'https://example.com/articles/mental-health-and-wellness-in-modern-times',
    'mental health wellness modern times',
    'Complete Guide to Mental Health and Wellness',
    'Understanding mental health in modern times: symptoms, treatments, and wellness strategies for better mental health.',
    'https://example.com/images/mental-health-social.jpg'
),

-- Article 3: Diabetes Care
(
    UUID(),
    'Managing Diabetes: A Comprehensive Guide',
    'managing-diabetes-comprehensive-guide',
    '<h2>Understanding Diabetes</h2>
    <p>Diabetes is a chronic condition that affects how your body processes blood sugar. Understanding and managing diabetes is crucial for maintaining good health.</p>
    
    <h2>Types of Diabetes</h2>
    <ul>
        <li>Type 1 Diabetes</li>
        <li>Type 2 Diabetes</li>
        <li>Gestational Diabetes</li>
        <li>Prediabetes</li>
    </ul>
    
    <h2>Symptoms and Signs</h2>
    <ul>
        <li>Increased thirst and urination</li>
        <li>Unexplained weight loss</li>
        <li>Fatigue</li>
        <li>Blurred vision</li>
        <li>Slow-healing sores</li>
    </ul>
    
    <h2>Management Strategies</h2>
    <p>Effective diabetes management includes:</p>
    <ul>
        <li>Blood sugar monitoring</li>
        <li>Medication management</li>
        <li>Healthy diet planning</li>
        <li>Regular exercise</li>
        <li>Regular medical check-ups</li>
    </ul>',
    'Learn about diabetes types, symptoms, and effective management strategies. Discover how to maintain healthy blood sugar levels through diet, exercise, and medical care.',
    'https://example.com/images/diabetes-care.jpg',
    'published',
    @admin_id,
    NOW(),
    'Diabetes Management Guide: Types, Symptoms & Care | Hospital CMS',
    'Complete guide to understanding and managing diabetes, including types, symptoms, and effective management strategies.',
    'diabetes, diabetes management, blood sugar, diabetes symptoms, diabetes care',
    'https://example.com/articles/managing-diabetes-comprehensive-guide',
    'diabetes management guide',
    'Complete Guide to Diabetes Management and Care',
    'Understanding and managing diabetes: types, symptoms, and effective strategies for blood sugar control.',
    'https://example.com/images/diabetes-care-social.jpg'
);
