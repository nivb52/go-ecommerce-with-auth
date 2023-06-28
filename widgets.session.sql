CREATE TABLE IF NOT EXISTS widgets (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  price INTEGER NOT NULL,
  inventory_level INTEGER NOT NULL DEFAULT 100,
  description TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,    
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
