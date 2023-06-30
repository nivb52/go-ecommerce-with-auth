
CREATE TABLE IF NOT EXISTS transaction_statuses (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT  NOT NULL CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP  ON UPDATE DEFAULT CURRENT_TIMESTAMP
)


INSERT INTO transaction_statuses (
    name,
  )
VALUES  
 ('Pending')
 ('Pending'),
 ('Cleared'),
 ('Declined'),
 ('Refunded'),
('Partially refunded');
