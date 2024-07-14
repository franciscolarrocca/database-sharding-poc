CREATE TABLE products (
    id SERIAL NOT NULL PRIMARY KEY,
    product_name TEXT NOT NULL,
    product_code CHARACTER(5) NOT NULL,
    UNIQUE (product_name),
    UNIQUE (product_code)
);