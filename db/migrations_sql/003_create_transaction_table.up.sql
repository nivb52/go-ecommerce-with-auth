
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    amount INTEGER NOT NULL,
    currency VARCHAR(6),
    last_four VARCHAR(4),
    bank_return_code VARCHAR(255),
    -- transaction_status_id FOREIGN  KEY TO transactions table
    transaction_status_id INTEGER,
     CONSTRAINT `FK_transactions_transaction_status_id` 
    FOREIGN KEY (transaction_status_id) REFERENCES transaction_statuses(id) 
        ON DELETE SET NULL  ON UPDATE NO ACTION, 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP  ON UPDATE CURRENT_TIMESTAMP
);
