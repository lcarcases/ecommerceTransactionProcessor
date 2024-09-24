-- init.sql
CREATE TABLE transactions (
  transaction_id SERIAL PRIMARY KEY,
  date TIMESTAMP NOT NULL,
  product_id INT NOT NULL,
  quantity INT NOT NULL,
  price DECIMAL(10, 2) NOT NULL
);
