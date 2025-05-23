create table if not exists orders (
    id BIGSERIAL,
    accrual BIGINT NOT NULL CHECK ( accrual >= 0 ),
    number TEXT UNIQUE NOT NULL,
    -- status can be a enum
    status SMALLINT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    user_id BIGINT NOT NULL,

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
