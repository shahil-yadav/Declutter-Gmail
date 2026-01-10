-- Create the initial table of mails
CREATE TABLE mails (
    id VARCHAR(255) PRIMARY KEY,
    user_id int,
    account_email VARCHAR(255),
    sender_email VARCHAR(255),
    snippet VARCHAR(255),
    date DATETIME,

    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Create the table users for managin my google account
CREATE TABLE users (
    -- Using AUTO_INCREMENT, cannot use uuid due to skill issues
    user_id INT AUTO_INCREMENT PRIMARY KEY,

    -- Basic User Info
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    cover_photo TEXT,

    -- OAuth Credentials
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    expiry DATETIME,
    expires_in INT NOT NULL,
    token_type VARCHAR(30) DEFAULT 'Bearer',

    -- Audit timestamps
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


-- Create the table recording the download jobs made by user
CREATE TABLE scan_jobs (
    job_id VARCHAR(255) PRIMARY KEY,
    user_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(user_id)
)

-- Create the table recording the trash jobs made by user
CREATE TABLE trash_jobs (
    job_id VARCHAR(255) PRIMARY KEY,
    user_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(user_id) REFERENCES users(user_id)
)