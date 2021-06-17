-- POSTGRES DB

CREATE DATABASE kitchen_service;

CREATE TABLE tickets(
    order_id INT PRIMARY KEY,
    vendor VARCHAR(255),
    saga_id VARCHAR(255) NOT NULL,
    command_id VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    line_items JSON NOT NULL
);

CREATE TABLE processed_messages(
    id VARCHAR(255) PRIMARY KEY
);
