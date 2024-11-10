package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/dltruongdev/TrinityApp/db/sqlc"
	"github.com/gin-gonic/gin"
)

// @Summary Get a voucher by code
// @Description Retrieve a voucher using its unique code
// @Tags vouchers
// @Produce json
// @Param code path string true "Voucher Code"
// @Success 200 {object} db.Voucher
// @Failure 404 {object} gin.H{"error": "Voucher not found"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /vouchers/{code} [get]
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

// @Summary Delete a voucher by code
// @Description Remove a voucher from the system using its code
// @Tags vouchers
// @Param code path string true "Voucher Code"
// @Success 204 "No Content"
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /vouchers/{code} [delete]
func (server *Server) deleteVoucherByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	err := server.store.DeleteVoucherByCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Redeem a voucher
// @Description Redeem a voucher for a user
// @Tags vouchers
// @Accept json
// @Produce json
// @Param voucher body redeemVoucherRequest true "Redeem Voucher data"
// @Success 200 {object} gin.H{"message": "Voucher redeemed successfully"}
// @Failure 400 {object} gin.H{"error": "Cannot apply voucher"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /vouchers/redeem [post]
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
