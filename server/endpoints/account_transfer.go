package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/client/api"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

func exportAccountTransfer(resp *client.TransferStatus) *api.TransferResponse {
	return &api.TransferResponse{
		Id:     resp.ID,
		Status: api.OperationStatus(resp.Status),
	}
}

// AccountTransfer handles transfers between accounts on an exchange
func AccountTransfer(c *fiber.Ctx) error {
	exchangeCfg, secrets, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Parse request body
	var req api.Transfer
	if err := c.BodyParser(&req); err != nil {
		return servererrors.BadRequestf("invalid request body: %s", err)
	}

	// Validate required fields
	symbol := api.DerefOrZero(req.Symbol)
	if symbol == "" {
		return servererrors.BadRequestf("symbol is required")
	}

	if req.Amount == "" {
		return servererrors.BadRequestf("amount is required")
	}

	// Parse amount
	amount, err := oc.NewAmountFromString(req.Amount)
	if err != nil {
		return servererrors.BadRequestf("invalid amount: %s", err)
	}

	if amount.IsZero() {
		return servererrors.BadRequestf("amount must be greater than 0")
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, secrets)
	if err != nil {
		return servererrors.InternalErrorf("failed to create client: %s", err)
	}

	// Prepare transfer arguments
	transferArgs := client.NewAccountTransferArgs(
		oc.SymbolId(symbol),
		amount,
	)

	// Handle from type
	if fromType := api.DerefOrZero(req.FromType); fromType != "" {
		resolvedFromType, ok, message := exchangeCfg.ResolveAccountType(fromType)
		if !ok {
			return servererrors.BadRequestf(message)
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
	if toType := api.DerefOrZero(req.ToType); toType != "" {
		resolvedToType, ok, message := exchangeCfg.ResolveAccountType(toType)
		if !ok {
			return servererrors.BadRequestf(message)
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
	if from := api.DerefOrZero(req.From); from != "" {
		fromSubaccount, ok := exchangeCfg.ResolveSubAccount(from)
		if !ok {
			return servererrors.BadRequestf("invalid from account")
		}
		transferArgs.SetFrom(fromSubaccount.Id)
	}

	// Handle to account
	if to := api.DerefOrZero(req.To); to != "" {
		toSubaccount, ok := exchangeCfg.ResolveSubAccount(to)
		if !ok {
			return servererrors.BadRequestf("invalid to account")
		}
		transferArgs.SetTo(toSubaccount.Id)
	}

	// Ensure at least one transfer parameter is specified
	if (api.DerefOrZero(req.To) + api.DerefOrZero(req.From) + api.DerefOrZero(req.FromType) + api.DerefOrZero(req.ToType)) == "" {
		return servererrors.BadRequestf("must specify at least one of to, from, from_type, to_type")
	}

	// Execute transfer
	resp, err := cli.CreateAccountTransfer(transferArgs)
	if err != nil {
		return servererrors.Conflictf("failed to create account transfer: %s", err)
	}

	return c.JSON(exportAccountTransfer(resp))
}
