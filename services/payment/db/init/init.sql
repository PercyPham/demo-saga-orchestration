-- POSTGRES DB

CREATE DATABASE payment_service;

-- TODO: implement
CREATE TABLE payments(
);

CREATE TABLE processed_messages(
    id VARCHAR(255) PRIMARY KEY,
    message JSON NOT NULL
);
