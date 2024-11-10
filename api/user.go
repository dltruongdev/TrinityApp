package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/dltruongdev/TrinityApp/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Name       string `json:"name" binding:"required,lte=50"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	UserTypeID int32  `json:"user_type_id" binding:"required,oneof=1 2"`
}

// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body createUserRequest true "User data"
// @Param promoCode query string false "Promotional code"
// @Success 201 {object} db.User
// @Failure 400 {object} gin.H{"error": "Bad Request"}
// @Failure 409 {object} gin.H{"error": "Conflict"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /users [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	// check if user exists
	exist, err := server.store.IsUserExist(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("something went wrong")))
		return
	}

	if exist {
		ctx.JSON(http.StatusConflict, errorResponse(errors.New("email is already existed")))
		return
	}

	// TODO verify if this is admin user
	// if true keep UserTypeID

	// else set usertype to 2 (User)
	req.UserTypeID = 2

	arg := db.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password, // TODO hashed password
		UserTypeID:   req.UserTypeID,
		PlanID: sql.NullInt32{
			Int32: 1,
			Valid: true,
		},
	}

	// create user
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get promo code in the querystring
	promotionalCode := ctx.Query("promoCode")
	if promotionalCode != "" {
		// generate new voucher and update redeemed_voucher to 1
		voucher, err := server.store.GenerateVoucherTx(ctx, db.GenerateVoucherParams{
			CampaignCode: promotionalCode,
			UserID:       user.UserID,
		})

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.JSON(http.StatusCreated, user)
				return
			}
			ctx.JSON(http.StatusCreated, user)
			return
			// do some logging
		}

		redirectURL := "/promotion?promoCode=" + promotionalCode + "&voucher=" + voucher.Code
		ctx.Redirect(http.StatusFound, redirectURL)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}
