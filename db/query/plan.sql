-- name: CreatePlan :one
INSERT INTO Plans (plan_name, price)
VALUES ($1, $2)
RETURNING plan_id, plan_name, price;

-- name: GetPlanByID :one
SELECT plan_id, plan_name, price
FROM Plans
WHERE plan_id = $1;

-- name: UpdatePlan :one
UPDATE Plans
SET plan_name = $2, price = $3
WHERE plan_id = $1
RETURNING plan_id, plan_name, price;

-- name: DeletePlan :exec
DELETE FROM Plans
WHERE plan_id = $1;

-- name: ListPlans :many
SELECT plan_id, plan_name, price
FROM Plans;