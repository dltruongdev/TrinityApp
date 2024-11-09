-- name: IsCampaignValid :one
SELECT EXISTS (
            SELECT 1 
            FROM Campaigns 
            WHERE code = $1 AND end_date >= $2) isValid;

-- name: IsCampaignValidForVoucherGeneration :one
SELECT EXISTS (
            SELECT 1 
            FROM Campaigns 
            WHERE code = $1 
            AND end_date > NOW()
            AND redeemed_vouchers < max_vouchers) isValid;
