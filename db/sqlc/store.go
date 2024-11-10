package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

type CreateCampaignTxParams struct {
	Name               string         `json:"name"`
	Description        sql.NullString `json:"description"`
	Code               string         `json:"code"`
	StartDate          time.Time      `json:"start_date"`
	EndDate            time.Time      `json:"end_date"`
	MaxVouchers        int32          `json:"max_vouchers"`
	VoucherLifetime    int32          `json:"voucher_lifetime"`
	DiscountPercentage int32          `json:"discount_percentage"`
	PlanID             int32          `json:"plan_id"`
}

type CreateCampaignTxResult struct {
	Campaign
	PlanID int32 `json:"plan_id"`
}

// CreateCampaignTx perform Campaign insertion into both tables:
// campaigns and campaignplans in a single db transaction
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
	CampaignCode string    `json:"campaign_code"`
	VoucherID    int32     `json:"voucher_id"`
	UserID       int32     `json:"user_id"`
	CampaignID   int32     `json:"campaign_id"`
	ValidUntil   time.Time `json:"valid_until"`
}

func (store *Store) GenerateVoucher(ctx context.Context, arg GenerateVoucherParams) (CreateVoucherRow, error) {
	var result CreateVoucherRow

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		campaign, err := q.IncreaseRedeemedVoucher(ctx, arg.CampaignCode)
		if err != nil {
			return err
		}

		validUntil := time.Now().Add(time.Duration(campaign.VoucherLifetime))
		if campaign.EndDate.Before(validUntil) {
			validUntil = campaign.EndDate
		}

		// Update the many-to-many table campaignplans
		result, err = q.CreateVoucher(ctx, CreateVoucherParams{
			UserID:     arg.UserID,
			CampaignID: arg.CampaignID,
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
