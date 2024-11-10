-- name: CreateCampaignPlan :exec
INSERT INTO CampaignPlans (campaign_id, plan_id)
VALUES ($1, $2);