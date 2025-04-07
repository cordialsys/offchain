package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

// WithdrawalRequest represents the request body for withdrawals
type WithdrawalRequest struct {
	To      oc.Address   `json:"to"`
	Symbol  oc.SymbolId  `json:"symbol"`
	Network oc.NetworkId `json:"network"`
	Amount  oc.Amount    `json:"amount"`
}

// CreateWithdrawal handles withdrawal requests from an exchange
func CreateWithdrawal(c *fiber.Ctx) error {
	exchangeCfg, secrets, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Parse request body
	var req WithdrawalRequest
	if err := c.BodyParser(&req); err != nil {
		return servererrors.BadRequestf(c, "invalid request body: %s", err)
	}

	// Validate required fields
	if req.To == "" {
		return servererrors.BadRequestf(c, "to address is required")
	}

	if req.Symbol == "" {
		return servererrors.BadRequestf(c, "symbol is required")
	}

	if req.Network == "" {
		return servererrors.BadRequestf(c, "network is required")
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, secrets)
	if err != nil {
		return servererrors.InternalErrorf(c, "failed to create client: %s", err)
	}

	// Create withdrawal
	resp, err := cli.CreateWithdrawal(client.NewWithdrawalArgs(
		req.To,
		req.Symbol,
		req.Network,
		req.Amount,
	))

	if err != nil {
		return servererrors.Conflictf(c, "failed to create withdrawal: %s", err)
	}

	return c.JSON(resp)
}
