
CREATE TABLE IF NOT EXISTS statuses (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP  ON UPDATE CURRENT_TIMESTAMP
);

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
    FOREIGN  KEY (widget_id) REFERENCES widgets(id),

    transaction_id INTEGER NOT NULL,
    CONSTRAINT `FK_order_transaction_id`
    FOREIGN  KEY (transaction_id) REFERENCES transactions(id) 
        ON DELETE SET NULL ON UPDATE NO ACTION, 

    status_id INTEGER NOT NULL,
    CONSTRAINT `FK_order_transaction_id`
    FOREIGN  KEY (status_id) REFERENCES statuses(id) 
        ON DELETE SET NULL  ON UPDATE NO ACTION,

    quantity INTEGER,
    amount INTEGER,

    created_at TIMESTAMP DEFAULT NOT NULL CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP ON UPDATE DEFAULT CURRENT_TIMESTAMP
);

