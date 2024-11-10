package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreateCampaignTx perform Campaign insertion into both tables:
// campaigns and campaignplans in a single db transaction
type CreateCampaignTxParams struct {
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Code               string    `json:"code"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	MaxVouchers        int32     `json:"max_vouchers"`
	VoucherLifetime    int32     `json:"voucher_lifetime"`
	DiscountPercentage int32     `json:"discount_percentage"`
	PlanID             int32     `json:"plan_id"`
}

type CreateCampaignTxResult struct {
	Campaign
	PlanID int32 `json:"plan_id"`
}

func (store *Store) CreateCampaignTx(ctx context.Context, arg CreateCampaignTxParams) (CreateCampaignTxResult, error) {
	var result CreateCampaignTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Campaign, err = q.CreateCampaign(ctx, CreateCampaignParams{
			Name:               arg.Name,
			Description:        arg.Description,
			Code:               arg.Code,
			StartDate:          arg.StartDate,
			EndDate:            arg.EndDate,
			MaxVouchers:        arg.MaxVouchers,
			VoucherLifetime:    arg.VoucherLifetime,
			DiscountPercentage: arg.DiscountPercentage,
		})

		if err != nil {
			return err
		}

		// Update the many-to-many table campaignplans
		err = q.CreateCampaignPlan(ctx, CreateCampaignPlanParams{
			CampaignID: result.CampaignID,
			PlanID:     arg.PlanID,
		})

		if err != nil {
			return err
		}

		result.PlanID = arg.PlanID

		return nil
	})

	return result, err
}

// Generate new voucher for eligible users
type GenerateVoucherParams struct {
	CampaignCode string `json:"campaign_code"`
	UserID       int32  `json:"user_id"`
}

func (store *Store) GenerateVoucherTx(ctx context.Context, arg GenerateVoucherParams) (CreateVoucherRow, error) {
	var result CreateVoucherRow

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// Validate
		campaign, err := q.GetCampaignForUpdate(ctx, arg.CampaignCode)
		if err != nil {
			return err
		}

		// Check if validaUntil is pass campaign's end_date
		validUntil := time.Now().Add(time.Duration(campaign.VoucherLifetime))
		if campaign.EndDate.Before(validUntil) {
			validUntil = campaign.EndDate
		}

		// Create voucher
		result, err = q.CreateVoucher(ctx, CreateVoucherParams{
			UserID:     arg.UserID,
			CampaignID: campaign.CampaignID,
			Code:       generateVoucherCode(),
			ValidUntil: validUntil,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func generateVoucherCode() string {
	id := uuid.NewString()
	return strings.ReplaceAll(id, "-", "")
}

// Redeeming a voucher
type RedeemVoucherTxParams struct {
	CampaignCode string `json:"campaign_code"`
	UserID       int32  `json:"user_id"`
	VoucherCode  string `json:"voucher_code"`
}

type RedeemVoucherTxResult struct {
}

func (store *Store) RedeemVoucherTx(ctx context.Context, arg RedeemVoucherTxParams) error {

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// Update voucher
		rows, err := q.RedeemVoucher(ctx, RedeemVoucherParams{
			Code:   arg.VoucherCode,
			UserID: arg.UserID, // Todo get it from context
		})
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}

		// Update campaign
		rows, err = q.IncreaseRedeemedVoucher(ctx, arg.CampaignCode)
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}

		return nil
	})

	return err
}
