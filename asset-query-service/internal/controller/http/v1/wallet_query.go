package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-service/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-service/internal/usecase"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type walletQueryRoutes struct {
	t usecase.WalletQueryHandler
	l logger.Interface
}

func newWalletQueryRoutes(handler *gin.RouterGroup, t usecase.WalletQueryHandler, l logger.Interface) {
	r := &walletQueryRoutes{t, l}

	h := handler.Group("/wallets")
	{
		h.GET("/:id", r.GetWalletByID)
		h.GET("/:id/balance", r.GetWalletBalance)           // Retrieve wallet balance
		h.GET("/:id/transactions", r.GetTransactionHistory) // Retrieve transaction history
	}
}

// Swagger nesnelerindeki isimleri büyük harfle başladık

type WalletResponse struct {
	Wallet entity.Wallet `json:"wallet"`
	Status string        `json:"status"`
	Error  string        `json:"error,omitempty"`
}

type BalanceResponse struct {
	WalletID int     `json:"wallet_id"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
	Error    string  `json:"error,omitempty"`
}

type TransactionHistoryResponse struct {
	Transactions []entity.Transaction `json:"transactions"`
	Count        int                  `json:"count"`
	Status       string               `json:"status"`
	Error        string               `json:"error,omitempty"`
}

// @Summary     Retrieve a wallet by ID
// @Description Get details of a specific wallet by its ID
// @ID          get-wallet-by-id
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       id path int true "Wallet ID"
// @Success     200 {object} WalletResponse
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /wallets/{id} [get]
func (r *walletQueryRoutes) GetWalletByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	wallet, err := r.t.GetWalletByID(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, http.StatusInternalServerError, "Failed to retrieve wallet")
		return
	}

	if wallet == nil {
		r.handleError(c, http.StatusNotFound, "Wallet not found")
		return
	}

	c.JSON(http.StatusOK, WalletResponse{Wallet: *wallet, Status: "success"})
}

// @Summary     Retrieve wallet balance
// @Description Get the balance of a specific wallet
// @ID          get-wallet-balance
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       id path int true "Wallet ID"
// @Success     200 {object} BalanceResponse
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /wallets/{id}/balance [get]
func (r *walletQueryRoutes) GetWalletBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	balance, err := r.t.GetBalance(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, http.StatusInternalServerError, "Failed to retrieve balance")
		return
	}

	c.JSON(http.StatusOK, BalanceResponse{WalletID: id, Balance: balance, Status: "success"})
}

// @Summary     Retrieve transaction history
// @Description Get the transaction history for a wallet
// @ID          get-wallet-transaction-history
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       id path int true "Wallet ID"
// @Success     200 {object} TransactionHistoryResponse
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /wallets/{id}/transactions [get]
func (r *walletQueryRoutes) GetTransactionHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	transactions, err := r.t.GetTransactionHistory(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, http.StatusInternalServerError, "Failed to retrieve transaction history")
		return
	}

	c.JSON(http.StatusOK, TransactionHistoryResponse{
		Transactions: transactions,
		Count:        len(transactions),
		Status:       "success",
	})
}

// handleError - Hataları daha iyi yönetmek için bir yardımcı fonksiyon
func (r *walletQueryRoutes) handleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"status": "error", "message": message})
}
