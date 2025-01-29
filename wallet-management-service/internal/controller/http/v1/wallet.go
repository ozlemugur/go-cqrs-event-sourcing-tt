package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/internal/usecase"
)

type walletRoutes struct {
	t usecase.WalletHandler
	l logger.Interface
}

func newWalletRoutes(handler *gin.RouterGroup, t usecase.WalletHandler, l logger.Interface) {
	r := &walletRoutes{t, l}

	h := handler.Group("/wallets")
	{
		h.GET("", r.GetAllWallets)       // Retrieve all wallets
		h.GET("/:id", r.GetWalletByID)   // Retrieve a wallet by ID
		h.POST("", r.CreateWallet)       // Create a new wallet
		h.PUT("/:id", r.UpdateWallet)    // Update a wallet by ID
		h.DELETE("/:id", r.DeleteWallet) // Delete a wallet by ID
	}
}

// Response format for wallet-related responses
type walletResponse struct {
	Wallets []entity.Wallet `json:"wallets"`
	Count   int             `json:"count"`
	Status  string          `json:"status"`
	Error   string          `json:"error,omitempty"`
}

// @Summary     Retrieve all wallets
// @Description Get a list of all wallets from the database
// @ID          get-all-wallets
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Success     200 {array} walletResponse
// @Failure     500 {object} response
// @Router      /wallets [get]
func (r *walletRoutes) GetAllWallets(c *gin.Context) {
	wallets, err := r.t.GetAllWallets(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - GetAllWallets")
		errorResponse(c, http.StatusInternalServerError, "Failed to retrieve wallets")
		return
	}
	response := walletResponse{
		Wallets: wallets,
		Count:   len(wallets),
		Status:  "success",
	}
	if len(wallets) == 0 {
		response.Status = "empty"
		response.Wallets = []entity.Wallet{}
	}
	c.JSON(http.StatusOK, response)
}

// @Summary     Retrieve a wallet by ID
// @Description Get details of a specific wallet by its ID
// @ID          get-wallet-by-id
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       id path int true "Wallet ID"
// @Success     200 {object} entity.Wallet
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /wallets/{id} [get]
func (r *walletRoutes) GetWalletByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // Convert string to int
	if err != nil {
		r.l.Error(err, "http - v1 - GetWalletByID - invalid ID format")
		errorResponse(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	wallet, err := r.t.GetWalletByID(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - GetWalletByID")
		errorResponse(c, http.StatusInternalServerError, "Failed to retrieve wallet")
		return
	}
	if wallet == nil {
		errorResponse(c, http.StatusNotFound, "Wallet not found")
		return
	}
	c.JSON(http.StatusOK, wallet)
}

// @Summary     Create a new wallet
// @Description Add a new wallet to the database
// @ID          create-wallet
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       wallet body entity.Wallet true "Wallet details"
// @Success     201 {object} response
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /wallets [post]
func (r *walletRoutes) CreateWallet(c *gin.Context) {
	var wallet entity.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
		r.l.Error(err, "http - v1 - CreateWallet - invalid input")
		errorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := r.t.CreateWallet(c.Request.Context(), wallet); err != nil {
		r.l.Error(err, "http - v1 - CreateWallet - use case error")
		errorResponse(c, http.StatusInternalServerError, "Failed to create wallet")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

// @Summary     Update a wallet by ID
// @Description Update details of a specific wallet
// @ID          update-wallet
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       id path int true "Wallet ID"
// @Param       wallet body entity.Wallet true "Updated wallet details"
// @Success     200 {object} response
// @Failure     400 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /wallets/{id} [put]
func (r *walletRoutes) UpdateWallet(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // Convert string to int
	if err != nil {
		r.l.Error(err, "http - v1 - UpdateWallet - invalid ID format")
		errorResponse(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	var wallet entity.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
		r.l.Error(err, "http - v1 - UpdateWallet - invalid input")
		errorResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}
	if err := r.t.UpdateWallet(c.Request.Context(), id, wallet); err != nil {
		r.l.Error(err, "http - v1 - UpdateWallet - use case error")
		errorResponse(c, http.StatusInternalServerError, "Failed to update wallet")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// @Summary     Delete a wallet by ID
// @Description Remove a specific wallet from the database
// @ID          delete-wallet
// @Tags        wallets
// @Accept      json
// @Produce     json
// @Param       id path int true "Wallet ID"
// @Success     200 {object} response
// @Failure     404 {object} response
// @Failure     500 {object} response
// @Router      /wallets/{id} [delete]
func (r *walletRoutes) DeleteWallet(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam) // Convert string to int
	if err != nil {
		r.l.Error(err, "http - v1 - DeleteWallet - invalid ID format")
		errorResponse(c, http.StatusBadRequest, "Invalid wallet ID")
		return
	}

	if err := r.t.DeleteWallet(c.Request.Context(), id); err != nil {
		r.l.Error(err, "http - v1 - DeleteWallet - use case error")
		errorResponse(c, http.StatusInternalServerError, "Failed to delete wallet")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
