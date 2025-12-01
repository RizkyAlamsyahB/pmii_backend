-- Rollback seed admin user
DELETE FROM tbl_user WHERE user_email IN ('admin@pmii.or.id', 'author@pmii.or.id');
