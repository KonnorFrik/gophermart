create table if not exists orders (
    id BIGSERIAL PRIMARY KEY,
    accrual BIGINT NOT NULL CHECK ( accrual >= 0 ),
    number TEXT UNIQUE NOT NULL,
    -- status can be a enum
    status order_status NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,
    user_id BIGINT NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
