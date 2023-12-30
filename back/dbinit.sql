CREATE ROLE portfolioapp WITH LOGIN;

CREATE DATABASE portfolio_stats WITH OWNER portfolioapp;

\c portfolio_stats;

CREATE TABLE guests (
    guest_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL,
    salt TEXT,
    secret TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE visits (
    visit_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    guest_id INT REFERENCES guests(guest_id) ON DELETE CASCADE,
    visited_at TIMESTAMP NOT NULL
);

