package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/dltruongdev/TrinityApp/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (server *Server) getVoucherByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	voucher, err := server.store.GetVoucherByCode(ctx, code)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Voucher not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, voucher)
}

func (server *Server) deleteVoucherByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	err := server.store.DeleteVoucherByCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

type redeemVoucherRequest struct {
	CampageCode string `json:"campaign_code" binding:"required"`
	VoucherCode string `json:"voucher_code" binding:"required"`
	UserID      int32  `json:"user_id" binding:"required"`
}

func (server *Server) redeemVoucher(ctx *gin.Context) {
	var req redeemVoucherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	arg := db.RedeemVoucherTxParams{
		CampaignCode: req.CampageCode,
		VoucherCode:  req.VoucherCode,
		UserID:       req.UserID,
	}

	err := server.store.RedeemVoucherTx(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot apply voucher"})
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Voucher redeemed successfully"})
}
