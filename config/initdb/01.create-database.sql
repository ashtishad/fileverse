\c fileverse

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS file_metadata
(
    id        SERIAL PRIMARY KEY,
    file_id   UUID         not null default uuid_generate_v4(),
    file_name VARCHAR(255) NOT NULL,
    size      BIGINT       NOT NULL,
    timestamp timestamptz  not null default now(),
    ipfs_hash VARCHAR(255) NOT NULL
);
