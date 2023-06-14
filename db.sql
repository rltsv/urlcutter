CREATE TABLE IF NOT EXISTS shortener (
    link_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    long_url VARCHAR(50) NOT NULL,
    short_url VARCHAR(50) NOT NULL
 );