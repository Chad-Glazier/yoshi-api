-- this group of tables are used primarily for data that has
-- 1-1 relationships with users, with the exception of `user_
-- sessions`. Each user is represented by their email; foreign 
-- key relationships aren't allowed by PlanetScale (the remote
-- DB that I'm using), so they aren't included here.

CREATE TABLE user_credentials (
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    PRIMARY KEY (email)
);

CREATE TABLE user_data (
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    display_name VARCHAR(100) UNIQUE NOT NULL,
    PRIMARY KEY (email)
);

CREATE TABLE user_preferences (
    email VARCHAR(255) UNIQUE NOT NULL,
    theme VARCHAR(40) NOT NULL DEFAULT 'light',
    kilos BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (email)
);

CREATE TABLE user_sessions (
    email VARCHAR(255) NOT NULL,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    last_renewed TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (uuid)
);
