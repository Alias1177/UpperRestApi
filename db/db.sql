ALTER SEQUENCE users_id_seq RESTART WITH 1;
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(100),
                       email VARCHAR(150) UNIQUE,
                       password TEXT NOT NULL,
                       age INT
);

