-- Create users table based on tbl_user design
CREATE TABLE IF NOT EXISTS tbl_user (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(100) NOT NULL,
    user_email VARCHAR(60) NOT NULL UNIQUE,
    user_password VARCHAR(255) NOT NULL,
    user_level VARCHAR(10) NOT NULL DEFAULT '2' COMMENT '1=Admin, 2=Author',
    user_status VARCHAR(10) NOT NULL DEFAULT '1' COMMENT '1=Active, 0=Inactive',
    user_photo VARCHAR(40),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (user_email),
    INDEX idx_status (user_status),
    INDEX idx_level (user_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
