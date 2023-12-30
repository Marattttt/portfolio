CREATE DATABASE portfolio_stats;

\c portfolio_stats;

CREATE TABLE guests (
    guest_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL,
    message TEXT NOT NULL
);

CREATE TABLE visits (
    visit_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    visited_at TIMESTAMP NOT NULL
);

CREATE TABLE guest_visit (
    guest_id integer REFERENCES guests(guest_id) ON DELETE CASCADE,
    visit_data_id integer REFERENCES visits(visit_id) ON DELETE RESTRICT,

    PRIMARY KEY (guest_id, visit_data_id)
);

