
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTO_INCREMENT,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email VARCHAR(255),
  password VARCHAR(60),

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
  updated_at TIMESTAMP  ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO users (
    first_name,
    last_name,
    email,
    password
  )
VALUES  
('Admin',
'User',
'admin@example.com',
'$2a$12$VR1wDmweaF3ZTVgEHiJrNOSi8VcS4j0eamr96A/7iOe8vlum3O3/q'
);
