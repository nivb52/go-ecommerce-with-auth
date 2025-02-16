
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
    FOREIGN  KEY (transaction_id) REFERENCES transactions(id),

    status_id INTEGER NOT NULL,
    FOREIGN  KEY (status_id) REFERENCES statuses(id),

    customer_id INTEGER NOT NULL,
    quantity INTEGER,
    amount INTEGER,
    request_id VARCHAR(255),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE orders ADD INDEX request_id_index (request_id);
ALTER TABLE orders ADD INDEX customer_id_index (customer_id);
ALTER TABLE orders ADD INDEX status_id_index (status_id);

