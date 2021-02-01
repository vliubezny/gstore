BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS product (
    id serial PRIMARY KEY,
    category_id integer REFERENCES category (id),
    name VARCHAR(160) NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS position (
    product_id integer REFERENCES product (id) ON DELETE CASCADE,
    store_id integer REFERENCES store (id) ON DELETE CASCADE,
    price numeric NOT NULL CHECK (price > 0),
    PRIMARY KEY (product_id, store_id)
);

DROP TABLE IF EXISTS item;

COMMIT TRANSACTION;