
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    amount INSERT NOT NULL,
    currency VARCHAR(6),
    last_four string VARCHAR(4),
    bank_return_code VARCHAR(255),
    -- transaction_status_id FORGIEN KEY TO transaction table
    transaction_status_id INTEGER,
    FOREIGN KEY (transaction_status_id) REFERENCES transaction_statuses(id) 
        ON DELETE CASCADE, 
    created_at TIMESTAMP DEFAULT NOT NULL CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP ON UPDATE DEFAULT CURRENT_TIMESTAMP
);
