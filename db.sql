CREATE TABLE IF NOT EXISTS links (
    link_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    original_url VARCHAR(50) NOT NULL,
    short_url VARCHAR(50) NOT NULL,
    correlation_id VARCHAR(50)
 );