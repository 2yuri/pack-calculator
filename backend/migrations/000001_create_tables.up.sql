CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_products_deleted_at ON products(deleted_at);

CREATE TABLE packs (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    size INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_packs_product_id ON packs(product_id);
CREATE INDEX idx_packs_deleted_at ON packs(deleted_at);
CREATE UNIQUE INDEX idx_packs_product_id_size_active ON packs(product_id, size) WHERE deleted_at IS NULL;