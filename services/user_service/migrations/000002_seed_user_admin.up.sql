-- Seed admin user
-- Default password: admin123 (hashed with bcrypt DefaultCost)
-- Idempotent: Only insert if not exists
INSERT INTO tbl_user (user_name, user_email, user_password, user_level, user_status, user_photo)
SELECT * FROM (
    SELECT 'Admin PMII' as user_name, 'admin@pmii.or.id' as user_email, '$2a$10$qVYhLsqJvauH6.5iOddLyeyhZSdgM5wJxV0Z/ul5tjhVvx5/2Uo3K' as user_password, '1' as user_level, '1' as user_status, '' as user_photo
    UNION ALL
    SELECT 'Author PMII', 'author@pmii.or.id', '$2a$10$qVYhLsqJvauH6.5iOddLyeyhZSdgM5wJxV0Z/ul5tjhVvx5/2Uo3K', '2', '1', ''
) AS tmp
WHERE NOT EXISTS (
    SELECT 1 FROM tbl_user WHERE user_email IN ('admin@pmii.or.id', 'author@pmii.or.id')
);

-- Note: Default password for both users is 'admin123'
-- You should change this in production!
