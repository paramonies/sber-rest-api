CREATE TABLE user_types
(
    id      SERIAL NOT NULL PRIMARY KEY,
    name    VARCHAR(255) NOT NULL
);

CREATE TABLE users
(
    id              SERIAL NOT NULL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    age             INTEGER NOT NULL,
    user_type_id    INTEGER ,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_type_id) REFERENCES user_types (id) ON DELETE SET NULL
);

CREATE TABLE items
(
    id          SERIAL NOT NULL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    user_id     INTEGER NOT NULL, 
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);