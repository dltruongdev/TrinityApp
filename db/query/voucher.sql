-- name: CreateVoucher :one
INSERT INTO Vouchers (user_id, campaign_id, code, valid_until)
VALUES ($1, $2, $3, $4)
RETURNING voucher_id, user_id, campaign_id, code, valid_until;

-- name: GetVoucherByID :one
SELECT voucher_id, user_id, campaign_id, code, valid_until
FROM Vouchers
WHERE voucher_id = $1;

-- name: UpdateVoucher :one
UPDATE Vouchers
SET user_id = $2, campaign_id = $3, code = $4, valid_until = $5
WHERE voucher_id = $1
RETURNING voucher_id, user_id, campaign_id, code, valid_until;

-- name: DeleteVoucher :exec
DELETE FROM Vouchers
WHERE voucher_id = $1;