-- POSTGRES DB

CREATE DATABASE order_service;

CREATE TABLE orders(
  id SERIAL PRIMARY KEY,
  state VARCHAR(255) NOT NULL,
  vendor varchar(255) NOT NULL,
  location varchar(255) NOT NULL,
  line_items JSON NOT NULL
);
