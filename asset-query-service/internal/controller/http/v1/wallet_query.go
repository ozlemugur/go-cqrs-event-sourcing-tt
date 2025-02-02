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
	t usecase.WalletQueryUseCaseHandler
	l logger.Interface
}

func newWalletQueryRoutes(handler *gin.RouterGroup, t usecase.WalletQueryUseCaseHandler, l logger.Interface) {
	r := &walletQueryRoutes{t, l}

	h := handler.Group("/wallets")
	{
		h.GET("/:id/assets", r.GetAllAssets)                // Retrieve all assets of a wallet
		h.GET("/:id/assets/:asset", r.GetAssetBalance)      // Retrieve balance of a specific asset
		h.GET("/:id/transactions", r.GetTransactionHistory) // Retrieve transaction history
	}
}

// **Response Structs**

type AssetResponse struct {
	Assets []entity.WalletAsset `json:"assets"`
	Status string               `json:"status"`
	Error  string               `json:"error,omitempty"`
}

type AssetBalanceResponse struct {
	WalletID  int     `json:"wallet_id"`
	AssetName string  `json:"asset_name"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Error     string  `json:"error,omitempty"`
}

type TransactionHistoryResponse struct {
	Transactions []entity.Transaction `json:"transactions"`
	Count        int                  `json:"count"`
	Status       string               `json:"status"`
	Error        string               `json:"error,omitempty"`
}

// **Route Handlers**

// @Summary     Retrieve all assets of a wallet
// @Description Get all assets for a specific wallet by its ID
// @ID          get-all-assets
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       id path int true "Wallet ID"
// @Success     200 {object} AssetResponse
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /wallets/{id}/assets [get]
func (r *walletQueryRoutes) GetAllAssets(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	assets, err := r.t.GetAllAssets(c.Request.Context(), id)
	if err != nil {
		r.handleError(c, http.StatusInternalServerError, "Failed to retrieve assets")
		return
	}

	c.JSON(http.StatusOK, AssetResponse{Assets: assets, Status: "success"})
}

// @Summary     Retrieve balance of a specific asset
// @Description Get the balance of a specific asset in a wallet
// @ID          get-asset-balance
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       id path int true "Wallet ID"
// @Param       asset path string true "Asset Name"
// @Success     200 {object} AssetBalanceResponse
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /wallets/{id}/assets/{asset} [get]
func (r *walletQueryRoutes) GetAssetBalance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	assetName := c.Param("asset")
	asset, err := r.t.GetAssetBalance(c.Request.Context(), id, assetName)
	if err != nil {
		r.handleError(c, http.StatusInternalServerError, "Failed to retrieve asset balance")
		return
	}

	c.JSON(http.StatusOK, AssetBalanceResponse{
		WalletID:  id,
		AssetName: asset.AssetName,
		Amount:    asset.Amount,
		Status:    "success",
	})
}

// @Summary     Retrieve transaction history
// @Description Get the transaction history for a wallet
// @ID          get-transaction-history
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

	transactions, err := r.t.GetTransactionHistory(c.Request.Context(), id, "")
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

// **Helper Function for Error Handling**
func (r *walletQueryRoutes) handleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"status": "error", "message": message})
}
