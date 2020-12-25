CREATE TABLE IF NOT EXISTS item (
    id serial PRIMARY KEY,
    store_id integer REFERENCES store (id) ON DELETE CASCADE,
    name VARCHAR(160) NOT NULL,
    description TEXT NOT NULL,
    price bigint NOT NULL
);