-- POSTGRES DB

CREATE DATABASE order_service;

CREATE TABLE orders(
    id SERIAL PRIMARY KEY,
    status VARCHAR(255) NOT NULL,
    vendor varchar(255) NOT NULL,
    location varchar(255) NOT NULL,
    line_items JSON NOT NULL
);

CREATE TABLE sagas(
    id VARCHAR(2555) PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    current_step INT,
    last_command_id VARCHAR(255),
    end_state BOOLEAN,
    compensating BOOLEAN,
    data JSON
);

CREATE TABLE processed_messages(
    id VARCHAR(255) PRIMARY KEY,
    message JSON NOT NULL
);
