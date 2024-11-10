-- name: CreateVoucher :one
INSERT INTO Vouchers (user_id, campaign_id, code, valid_until)
VALUES ($1, $2, $3, $4)
RETURNING voucher_id, user_id, campaign_id, code, valid_until;

-- name: GetVoucherByCode :one
SELECT voucher_id, user_id, campaign_id, code, valid_until
FROM Vouchers
WHERE code = $1;

-- name: DeleteVoucherByCode :exec
DELETE FROM Vouchers
WHERE code = $1;

-- name: RedeemVoucher :execrows
UPDATE Vouchers
SET is_redeemed = true, updated_at = NOW()
WHERE code = $1 AND user_id = $2 AND valid_until > NOW();