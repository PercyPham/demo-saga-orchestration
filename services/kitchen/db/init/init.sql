-- POSTGRES DB

CREATE DATABASE kitchen_service;

CREATE TABLE tickets(
    order_id INT PRIMARY KEY,
    status VARCHAR(255) NOT NULL,
    command_id VARCHAR(255) NOT NULL,
    vendor VARCHAR(255),
    line_items JSON NOT NULL
);

CREATE TABLE processed_messages(
    id VARCHAR(255) PRIMARY KEY,
    message JSON NOT NULL
);
