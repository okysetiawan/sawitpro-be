/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE IF NOT EXISTS users
(
    id           serial PRIMARY KEY,
    name         VARCHAR(50)        NOT NULL,
    phone_number VARCHAR(50) UNIQUE NOT NULL,
    password     text               NOT NULL,
    salt         VARCHAR(255)       NOT NULL,

    created_at   timestamptz default current_timestamp,
    updated_at   timestamptz default current_timestamp
);
