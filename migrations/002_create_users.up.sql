CREATE TABLE users (

    -- Unique ID auto generated
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    -- USERNAME that is text
    username TEXT NOT NULL UNIQUE CHECK (length(username) > 0),

    -- PASSWORD HASH
    password_hash TEXT NOT NULL CHECK (length(password_hash) > 0)

);
