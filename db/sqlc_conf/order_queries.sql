-- name: CreateOrder :one
INSERT INTO orders (
    accrual, number, status, uploaded_at, user_id 
) VALUES ( 0, $1, 'NEW', NOW(), $2 )
RETURNING *;

-- name: UserOrders :many
SELECT * FROM orders
where user_id = $1;

-- name: DeleteOrderById :exec
DELETE FROM orders
    WHERE id = $1;

-- name: DeleteOrderByNumber :exec
DELETE FROM orders
    WHERE number = $1;
