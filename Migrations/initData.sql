INSERT INTO abonement (id, title, validity, visiting_time, photo, price)
VALUES
    (gen_random_uuid(), 'Monthly Gym Pass', '1 Month', '8am - 10pm', 'gym_pass_photo.jpg', 50),
    (gen_random_uuid(), 'Annual Gym Membership', '1 Year', '6am - 10pm', 'annual_pass_photo.jpg', 500);

INSERT INTO service (id, title, photo)
VALUES
    (gen_random_uuid(), 'Personal Training', 'personal_training.jpg'),
    (gen_random_uuid(), 'Yoga Classes', 'yoga_classes.jpg');

INSERT INTO "user" (id, name, email, role, password_hash, photo, created_time)
VALUES
    (gen_random_uuid(), 'John', 'john.doe@example.com', 'admin', 'hashpassword1', 'john_photo.jpg', NOW()),
    (gen_random_uuid(), 'Jane', 'jane.smith@example.com', 'user', 'hashpassword2', 'jane_photo.jpg', NOW());

INSERT INTO coach (id, name, description, photo)
VALUES
    (gen_random_uuid(), 'Mike Johnson', 'Experienced fitness coach', 'mike_photo.jpg'),
    (gen_random_uuid(), 'Sarah Lee', 'Yoga and wellness instructor', 'sarah_photo.jpg');

INSERT INTO abonement_service (abonement_id, service_id)
VALUES
    ((SELECT id FROM abonement WHERE title = 'Monthly Gym Pass'), (SELECT id FROM service WHERE title = 'Personal Training')),
    ((SELECT id FROM abonement WHERE title = 'Annual Gym Membership'), (SELECT id FROM service WHERE title = 'Yoga Classes'));

INSERT INTO coach_service (coach_id, service_id)
VALUES
    ((SELECT id FROM coach WHERE name = 'Mike Johnson'), (SELECT id FROM service WHERE title = 'Personal Training')),
    ((SELECT id FROM coach WHERE name = 'Sarah Lee'), (SELECT id FROM service WHERE title = 'Yoga Classes'));

INSERT INTO comment (id, comment_body, user_id, coach_id, create_date)
VALUES
    (gen_random_uuid(), 'Great personal training session!', (SELECT id FROM "user" WHERE name = 'John'), (SELECT id FROM coach WHERE name = 'Mike Johnson'), NOW()),
    (gen_random_uuid(), 'Yoga class was very relaxing!', (SELECT id FROM "user" WHERE name = 'Jane'), (SELECT id FROM coach WHERE name = 'Sarah Lee'), NOW());

INSERT INTO "order" (id, abonement_id, user_id, status)
VALUES
    (gen_random_uuid(), (SELECT id FROM abonement WHERE title = 'Monthly Gym Pass'), (SELECT id FROM "user" WHERE name = 'John'), 1),
    (gen_random_uuid(), (SELECT id FROM abonement WHERE title = 'Annual Gym Membership'), (SELECT id FROM "user" WHERE name = 'Jane'), 1);