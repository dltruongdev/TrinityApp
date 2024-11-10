package api

import (
	db "github.com/dltruongdev/TrinityApp/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	// config     util.Config
	store *db.Store
	// tokenMaker token.Maker
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store *db.Store) (*Server, error) {
	//tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot create token maker: %w", err)
	// }

	server := &Server{
		//config:     config,
		store: store,
		//tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	/*
	 Setup routes
	*/

	//users
	router.POST("/users/register", server.createUser)

	//campaigns
	router.POST("/campaigns", server.createCampaign)
	router.GET("/campaigns/:code", server.getCampgaignByCode)
	router.DELETE("/campaigns/:code", server.deleteCampgaignByCode)

	//vouchers
	router.GET("/vouchers/:code", server.getVoucherByCode)
	router.DELETE("/vouchers/:code", server.deleteCampgaignByCode)
	router.POST("vouchers/redeem", server.redeemVoucher)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
