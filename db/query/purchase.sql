-- name: CreatePurchase :one
INSERT INTO Purchases (user_id, plan_id, purchase_date, amount)
VALUES ($1, $2, $3, $4)
RETURNING purchase_id, user_id, plan_id, purchase_date, amount;

-- name: GetPurchaseByID :one
SELECT purchase_id, user_id, plan_id, purchase_date, amount
FROM Purchases
WHERE purchase_id = $1;

-- name: UpdatePurchase :one
UPDATE Purchases
SET user_id = $2, plan_id = $3, purchase_date = $4, amount = $5
WHERE purchase_id = $1
RETURNING purchase_id, user_id, plan_id, purchase_date, amount;

-- name: DeletePurchase :exec
DELETE FROM Purchases
WHERE purchase_id = $1;

-- name: ListPurchases :many
SELECT purchase_id, user_id, plan_id, purchase_date, amount
FROM Purchases;
