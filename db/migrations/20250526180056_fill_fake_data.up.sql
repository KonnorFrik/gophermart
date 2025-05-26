INSERT INTO users (
    login, password
) VALUES ( 'test', 'test' );

INSERT INTO orders (
    accrual, number, status, uploaded_at, user_id
)
VALUES
    (123, '12345', 'NEW', NOW(), (select id from users where login = 'test') ),
    (42, '123456', 'NEW', NOW() - interval '1 day', (select id from users where login = 'test') )
;
