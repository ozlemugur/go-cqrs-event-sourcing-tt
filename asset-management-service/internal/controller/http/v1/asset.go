package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/usecase"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type assetRoutes struct {
	t usecase.AssetHandler
	l logger.Interface
}

// NewAssetRoutes sets up the asset-related API endpoints.
func newAssetRoutes(handler *gin.RouterGroup, t usecase.AssetHandler, l logger.Interface) {
	r := &assetRoutes{t, l}

	h := handler.Group("/assets")
	{
		h.POST("/withdraw", r.Withdraw)            // Withdraw funds
		h.POST("/deposit", r.Deposit)              // Deposit funds
		h.POST("/transfer", r.Transfer)            // Transfer between wallets
		h.POST("/schedule", r.ScheduleTransaction) // Schedule a future transaction
	}
}

// Generic response structure for asset operations
type assetResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// @Summary     Withdraw funds
// @Description Withdraw funds from a wallet
// @ID          withdraw
// @Tags        assets
// @Accept      json
// @Produce     json
// @Param       request body entity.TransactionRequest true "Withdraw request"
// @Success     200 {object} assetResponse
// @Failure     400 {object} assetResponse
// @Failure     500 {object} assetResponse
// @Router      /assets/withdraw [post]
func (r *assetRoutes) Withdraw(c *gin.Context) {
	var req entity.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - Withdraw - invalid input")
		errorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Include asset_name in the use case call
	if err := r.t.Withdraw(c.Request.Context(), req.WalletID, req.AssetName, req.Amount); err != nil {
		r.l.Error(err, "http - v1 - Withdraw - use case error")
		errorResponse(c, http.StatusInternalServerError, "Withdraw failed")
		return
	}

	c.JSON(http.StatusOK, assetResponse{Status: "success"})
}

// @Summary     Deposit funds
// @Description Deposit funds into a wallet
// @ID          deposit
// @Tags        assets
// @Accept      json
// @Produce     json
// @Param       request body entity.TransactionRequest true "Deposit request"
// @Success     200 {object} assetResponse
// @Failure     400 {object} assetResponse
// @Failure     500 {object} assetResponse
// @Router      /assets/deposit [post]
func (r *assetRoutes) Deposit(c *gin.Context) {
	var req entity.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - Deposit - invalid input")
		errorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Include asset_name in the use case call
	if err := r.t.Deposit(c.Request.Context(), req.WalletID, req.AssetName, req.Amount); err != nil {
		r.l.Error(err, "http - v1 - Deposit - use case error")
		errorResponse(c, http.StatusInternalServerError, "Deposit failed")
		return
	}

	c.JSON(http.StatusOK, assetResponse{Status: "success"})
}

// @Summary     Transfer funds
// @Description Transfer funds between wallets
// @ID          transfer
// @Tags        assets
// @Accept      json
// @Produce     json
// @Param       request body entity.TransferRequest true "Transfer request"
// @Success     200 {object} assetResponse
// @Failure     400 {object} assetResponse
// @Failure     500 {object} assetResponse
// @Router      /assets/transfer [post]
func (r *assetRoutes) Transfer(c *gin.Context) {
	var req entity.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - Transfer - invalid input")
		errorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Include asset_name in the use case call
	if err := r.t.Transfer(c.Request.Context(), req.FromWalletID, req.ToWalletID, req.AssetName, req.Amount); err != nil {
		r.l.Error(err, "http - v1 - Transfer - use case error")
		errorResponse(c, http.StatusInternalServerError, "Transfer failed")
		return
	}

	c.JSON(http.StatusOK, assetResponse{Status: "success"})
}

// @Summary     Schedule a transaction
// @Description Schedule a future transaction
// @ID          schedule-transaction
// @Tags        assets
// @Accept      json
// @Produce     json
// @Param       request body entity.ScheduledTransactionRequest true "Scheduled transaction request"
// @Success     200 {object} assetResponse
// @Failure     400 {object} assetResponse
// @Failure     500 {object} assetResponse
// @Router      /assets/schedule [post]
func (r *assetRoutes) ScheduleTransaction(c *gin.Context) {
	var req entity.ScheduledTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - ScheduleTransaction - invalid input")
		errorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// Include asset_name in the use case call
	//	if err := r.t.ScheduleTransaction(c.Request.Context(), req); err != nil {
	//		r.l.Error(err, "http - v1 - ScheduleTransaction - use case error")
	//		errorResponse(c, http.StatusInternalServerError, "Scheduling failed")
	//		return
	//	}

	c.JSON(http.StatusOK, assetResponse{Status: "success"})
}
