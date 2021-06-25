-- POSTGRES DB

CREATE DATABASE payment_service;

-- TODO: implement
CREATE TABLE payments(
    order_id INT PRIMARY KEY,
    total INT,
    status VARCHAR(255) NOT NULL,
    command_id VARCHAR(255) NOT NULL
);

CREATE TABLE processed_messages(
    id VARCHAR(255) PRIMARY KEY,
    message JSON NOT NULL
);
