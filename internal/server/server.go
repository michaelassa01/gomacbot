package server

import (
	"fmt"

	db "github.com/michaelassa01/gomacbot/internal/database"
	// h "github.com/michaelassa01/gomacbot/internal/server/api"
	// wsh "github.com/michaelassa01/gomacbot/server/ws"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/michaelassa01/gomacbot/pkg/token"
	"github.com/michaelassa01/gomacbot/utils"
)

// Server serves HTTP request for klusta service
type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	// handler    *h.Handler
	// wsHandler  *wsh.WSHandler
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {

	// TokenMaker is set to use pasetoMaker and can be changed to JWTMaker
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// wsHandler := wsh.NewHandler(store, config)

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		// handler:    h.NewHandler(store, config),
		// wsHandler:  wsHandler,
	}

	// router func imported
	server.setupRouter(config)

	// to register validator with gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	return server, nil
}

// Start runs the HTTP server on a specification address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// func (ws *Server) WSHandler() *wsh.WSHandler {
// 	return ws.wsHandler
// }

func errorResponse(err error) gin.H {
	// // Split the error message into individual validation errors
	// lines := strings.Split(err.Error(), "\n")

	// // Prepare a slice of maps to hold individual field errors
	// var details []map[string]string

	// for _, line := range lines {
	// 	// Example line: "Key: 'tansferRequest.FromAccountID' Error:Field validation for 'FromAccountID' failed on the 'required' tag"
	// 	parts := strings.SplitN(line, "Error:", 2)
	// 	if len(parts) == 2 {
	// 		fieldPart := strings.TrimSpace(parts[0])
	// 		errorMsg := strings.TrimSpace(parts[1])

	// 		details = append(details, map[string]string{
	// 			"key":   fieldPart,
	// 			"error": errorMsg,
	// 		})
	// 	}
	// }

	// return gin.H{
	// 	"errors": details,
	// }
	return gin.H{"error": err.Error()}

}
