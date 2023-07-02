
CREATE TABLE customers (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email VARCHAR(255),

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
  updated_at TIMESTAMP  ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE orders ADD COLUMN customer_id INTEGER NOT NULL;
ALTER TABLE orders ADD  CONSTRAINT `FK_order_customer_id`  FOREIGN KEY (customer_id)
 REFERENCES customers(id)   ON UPDATE NO ACTION; 