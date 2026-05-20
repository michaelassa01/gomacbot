package transactions

import "github.com/gin-gonic/gin"

func RegisterTransactionsRoutes(r *gin.RouterGroup, h *Handler) {
	group := r.Group("/transactions")
	{
		group.POST("", h.CreateTransaction)
		group.GET("", h.ListTransactions)
		group.GET("/reference/:reference", h.GetTransactionByReference)
		group.GET("/:id", h.GetTransaction)
		group.PATCH("/:id/status", h.UpdateTransactionStatus)
		group.DELETE("/:id", h.DeleteTransaction)
	}
}
