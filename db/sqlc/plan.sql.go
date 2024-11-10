// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: plan.sql

package db

import (
	"context"
)

const createPlan = `-- name: CreatePlan :one
INSERT INTO Plans (plan_name, price)
VALUES ($1, $2)
RETURNING plan_id, plan_name, price
`

type CreatePlanParams struct {
	PlanName string `json:"plan_name"`
	Price    string `json:"price"`
}

func (q *Queries) CreatePlan(ctx context.Context, arg CreatePlanParams) (Plan, error) {
	row := q.db.QueryRowContext(ctx, createPlan, arg.PlanName, arg.Price)
	var i Plan
	err := row.Scan(&i.PlanID, &i.PlanName, &i.Price)
	return i, err
}

const deletePlan = `-- name: DeletePlan :exec
DELETE FROM Plans
WHERE plan_id = $1
`

func (q *Queries) DeletePlan(ctx context.Context, planID int32) error {
	_, err := q.db.ExecContext(ctx, deletePlan, planID)
	return err
}

const getPlanByID = `-- name: GetPlanByID :one
SELECT plan_id, plan_name, price
FROM Plans
WHERE plan_id = $1
`

func (q *Queries) GetPlanByID(ctx context.Context, planID int32) (Plan, error) {
	row := q.db.QueryRowContext(ctx, getPlanByID, planID)
	var i Plan
	err := row.Scan(&i.PlanID, &i.PlanName, &i.Price)
	return i, err
}

const listPlans = `-- name: ListPlans :many
SELECT plan_id, plan_name, price
FROM Plans
`

func (q *Queries) ListPlans(ctx context.Context) ([]Plan, error) {
	rows, err := q.db.QueryContext(ctx, listPlans)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Plan{}
	for rows.Next() {
		var i Plan
		if err := rows.Scan(&i.PlanID, &i.PlanName, &i.Price); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePlan = `-- name: UpdatePlan :one
UPDATE Plans
SET plan_name = $2, price = $3
WHERE plan_id = $1
RETURNING plan_id, plan_name, price
`

type UpdatePlanParams struct {
	PlanID   int32  `json:"plan_id"`
	PlanName string `json:"plan_name"`
	Price    string `json:"price"`
}

func (q *Queries) UpdatePlan(ctx context.Context, arg UpdatePlanParams) (Plan, error) {
	row := q.db.QueryRowContext(ctx, updatePlan, arg.PlanID, arg.PlanName, arg.Price)
	var i Plan
	err := row.Scan(&i.PlanID, &i.PlanName, &i.Price)
	return i, err
}