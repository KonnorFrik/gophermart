create table if not exists orders (
    id BIGSERIAL,
    accrual BIGINT CHECK ( accrual >= 0 ),
    number TEXT UNIQUE NOT NULL,
    -- status can be a enum
    status SMALLINT,
    uploaded_at TIMESTAMP,
    user_id BIGINT,

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
