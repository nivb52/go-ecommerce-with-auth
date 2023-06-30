
CREATE TABLE IF NOT EXISTS statuses (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT  NOT NULL CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP  ON UPDATE DEFAULT CURRENT_TIMESTAMP
)

INSERT INTO statuses (
    name
)
VALUES
('Cleared'),
('Refunded'),
('Cancelled')
;


CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    widget_id INTEGER NOT NULL,
    FORGIEN KEY (widget_id) REFERENCES widgets(id),

    transaction_id INTEGER NOT NULL,
    FORGIEN KEY (transaction_id) REFERENCES transactions(id) 
        ON DELETE CASCADE, 

    status_id INTEGER NOT NULL,
    FORGIEN KEY (status_id) REFERENCES statuses(id) 
        ON DELETE CASCADE,

    quantity INTEGER,
    amount INTEGER,

    created_at TIMESTAMP DEFAULT NOT NULL CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP ON UPDATE DEFAULT CURRENT_TIMESTAMP
);

