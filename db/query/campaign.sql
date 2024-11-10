-- name: IsCampaignValid :one
SELECT EXISTS (
            SELECT 1 
            FROM Campaigns 
            WHERE code = $1 AND end_date > $2 AND start_date < $2) isValid;

-- name: CreateCampaign :one
INSERT INTO Campaigns (name, description, code, start_date, end_date, max_vouchers, voucher_lifetime, discount_percentage)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;


-- name: GetCampaignForUpdate :one
SELECT 
    name, 
    code,
    end_date, 
    max_vouchers, 
    voucher_lifetime, 
    discount_percentage,
    start_date <= NOW() AND end_date > NOW() AND redeemed_vouchers < max_vouchers AS isValid
FROM 
    Campaigns
WHERE 
    code = $1
FOR NO KEY UPDATE;


--- Doing update directly will lock the record for update, if concurrent update happens (when two user finish register at the same time and there only one voucher left to be generated) this will keep the logic correctly
-- name: IncreaseRedeemedVoucher :one 
UPDATE Campaigns
SET redeemed_vouchers = redeemed_vouchers + 1
WHERE code = $1
AND end_date > NOW()
AND redeemed_vouchers < max_vouchers
RETURNING 
    start_date, 
    end_date, 
    voucher_lifetime, 
    discount_percentage;