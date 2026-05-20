package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/michaelassa01/gomacbot/internal/transactions/domain"
	u "github.com/michaelassa01/gomacbot/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateTransaction(ctx *gin.Context) {
	var req CreateTransactionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := u.ConvertToPgUUIDFromString(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	tx := &domain.Transaction{
		UserID:      userID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Type:        req.Type,
		Reference:   req.Reference,
		Description: req.Description,
		Provider:    req.Provider,
	}
	if len(req.Metadata) > 0 {
		tx.Metadata = req.Metadata
	}

	created, err := h.service.Create(tx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, toTransactionResponse(created))
}

func (h *Handler) GetTransaction(ctx *gin.Context) {
	var req GetTransactionReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.service.GetByID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, toTransactionResponse(tx))
}

func (h *Handler) GetTransactionByReference(ctx *gin.Context) {
	reference := ctx.Param("reference")
	if reference == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "reference is required"})
		return
	}

	tx, err := h.service.GetByReference(reference)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, toTransactionResponse(tx))
}

func (h *Handler) ListTransactions(ctx *gin.Context) {
	var req ListTransactionsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := h.service.ListByUserID(req.UserID, req.Limit, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make([]TransactionRes, 0, len(items))
	for _, item := range items {
		res = append(res, *toTransactionResponse(item))
	}

	ctx.JSON(http.StatusOK, ListTransactionsRes{Transactions: res})
}

func (h *Handler) UpdateTransactionStatus(ctx *gin.Context) {
	var uri GetTransactionReq
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateTransactionStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.service.UpdateStatus(uri.ID, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, toTransactionResponse(updated))
}

func (h *Handler) DeleteTransaction(ctx *gin.Context) {
	var req DeleteTransactionReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
