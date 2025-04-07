package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

// AccountTransferRequest represents the request body for account transfers
type AccountTransferRequest struct {
	Symbol   string `json:"symbol"`
	Amount   string `json:"amount"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	FromType string `json:"from_type,omitempty"`
	ToType   string `json:"to_type,omitempty"`
}

// AccountTransfer handles transfers between accounts on an exchange
func AccountTransfer(c *fiber.Ctx) error {
	exchangeCfg, secrets, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Parse request body
	var req AccountTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return servererrors.BadRequestf(c, "invalid request body: %s", err)
	}

	// Validate required fields
	if req.Symbol == "" {
		return servererrors.BadRequestf(c, "symbol is required")
	}

	if req.Amount == "" {
		return servererrors.BadRequestf(c, "amount is required")
	}

	// Parse amount
	amount, err := oc.NewAmountFromString(req.Amount)
	if err != nil {
		return servererrors.BadRequestf(c, "invalid amount: %s", err)
	}

	if amount.IsZero() {
		return servererrors.BadRequestf(c, "amount must be greater than 0")
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, secrets)
	if err != nil {
		return servererrors.InternalErrorf(c, "failed to create client: %s", err)
	}

	// Prepare transfer arguments
	transferArgs := client.NewAccountTransferArgs(
		oc.SymbolId(req.Symbol),
		amount,
	)

	// Handle from type
	if req.FromType != "" {
		resolvedFromType, ok, message := exchangeCfg.ResolveAccountType(req.FromType)
		if !ok {
			return servererrors.BadRequestf(c, message)
		}
		transferArgs.SetFromType(resolvedFromType.Type)
	} else {
		// Use the first account type by default
		first, ok := exchangeCfg.FirstAccountType()
		if ok {
			transferArgs.SetFromType(first.Type)
		}
	}

	// Handle to type
	if req.ToType != "" {
		resolvedToType, ok, message := exchangeCfg.ResolveAccountType(req.ToType)
		if !ok {
			return servererrors.BadRequestf(c, message)
		}
		transferArgs.SetToType(resolvedToType.Type)
	} else {
		// Use the first account type by default
		first, ok := exchangeCfg.FirstAccountType()
		if ok {
			transferArgs.SetToType(first.Type)
		}
	}

	// Handle from account
	if req.From != "" {
		fromSubaccount, ok := exchangeCfg.ResolveSubAccount(req.From)
		if !ok {
			return servererrors.BadRequestf(c, "invalid from account")
		}
		transferArgs.SetFrom(fromSubaccount.Id)
	}

	// Handle to account
	if req.To != "" {
		toSubaccount, ok := exchangeCfg.ResolveSubAccount(req.To)
		if !ok {
			return servererrors.BadRequestf(c, "invalid to account")
		}
		transferArgs.SetTo(toSubaccount.Id)
	}

	// Ensure at least one transfer parameter is specified
	if (req.To + req.From + req.FromType + req.ToType) == "" {
		return servererrors.BadRequestf(c, "must specify at least one of to, from, from_type, to_type")
	}

	// Execute transfer
	resp, err := cli.CreateAccountTransfer(transferArgs)
	if err != nil {
		return servererrors.Conflictf(c, "failed to create account transfer: %s", err)
	}

	return c.JSON(resp)
}
