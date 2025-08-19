CREATE TABLE books (
    id UUID PRIMARY KEY,                       -- maps to string, UUID is common for IDs
    name VARCHAR(200) NOT NULL,                -- validate: required, min=2, max=200
    author VARCHAR(255) NOT NULL,              -- validate: required
    price_cents BIGINT NOT NULL CHECK (price_cents > 0), -- validate: gt=0
    isbn VARCHAR(50),                          -- optional, can hold ISBN-10/13
    published_at TIMESTAMP NOT NULL,           -- published date
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
