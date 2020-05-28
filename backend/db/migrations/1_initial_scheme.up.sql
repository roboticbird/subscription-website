CREATE TABLE users (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE products (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  price INT NOT NULL,
  duration INT NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE subscriptions (
  id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  product_id INT NOT NULL,
  start_date DATETIME NOT NULL,
  status VARCHAR(255) NOT NULL,
  PRIMARY KEY (id),
  INDEX(user_id),
  INDEX(product_id),
  FOREIGN KEY (user_id)
    REFERENCES users(id),
  FOREIGN KEY (product_id)
    REFERENCES products(id)
);

