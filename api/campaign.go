package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/dltruongdev/TrinityApp/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createCampaignRequest struct {
	Name               string    `json:"name" binding:"required"`
	Description        string    `json:"description"`
	Code               string    `json:"code" binding:"required"`
	StartDate          time.Time `json:"start_date" binding:"required"`
	EndDate            time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
	MaxVouchers        int32     `json:"max_vouchers" binding:"required"`
	VoucherLifetime    int32     `json:"voucher_lifetime" binding:"required"`
	DiscountPercentage int32     `json:"discount_percentage" binding:"required"`
	PlanID             int32     `json:"plan_id" binding:"required"`
}

func (server *Server) createCampaign(ctx *gin.Context) {
	var req createCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// check if user exists
	exist, err := server.store.IsCampaginExist(ctx, req.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if exist {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email is already existed"})
		return
	}

	arg := db.CreateCampaignTxParams{
		Name:               req.Name,
		Description:        req.Description,
		Code:               req.Code,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		MaxVouchers:        req.MaxVouchers,
		VoucherLifetime:    req.VoucherLifetime,
		DiscountPercentage: req.DiscountPercentage,
		PlanID:             req.PlanID,
	}

	// create campaign
	campaign, err := server.store.CreateCampaignTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, campaign)
}

func (server *Server) getCampgaignByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	campaign, err := server.store.GetCompaignByCode(ctx, code)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, campaign)
}

func (server *Server) deleteCampgaignByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	err := server.store.DeleteCompaignByCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
