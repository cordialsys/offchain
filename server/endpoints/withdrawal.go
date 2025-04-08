package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/client/api"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

func exportWithdrawal(resp *client.WithdrawalResponse) *api.WithdrawalResponse {
	return &api.WithdrawalResponse{
		Id:     resp.ID,
		Status: api.OperationStatus(resp.Status),
	}
}

// CreateWithdrawal handles withdrawal requests from an exchange
func CreateWithdrawal(c *fiber.Ctx) error {
	exchangeCfg, secrets, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Parse request body
	var req api.Withdrawal
	if err := c.BodyParser(&req); err != nil {
		return servererrors.BadRequestf("invalid request body: %s", err)
	}

	// Validate required fields
	if req.Address == "" {
		return servererrors.BadRequestf("to address is required")
	}

	symbol := oc.SymbolId(api.DerefOrZero(req.Symbol))
	network := oc.NetworkId(api.DerefOrZero(req.Network))
	if symbol == "" {
		return servererrors.BadRequestf("symbol is required")
	}

	if network == "" {
		return servererrors.BadRequestf("network is required")
	}

	amount, err := oc.NewAmountFromString(req.Amount)
	if err != nil {
		return servererrors.BadRequestf("invalid amount: %s", err)
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, secrets)
	if err != nil {
		return servererrors.InternalErrorf("failed to create client: %s", err)
	}

	// Create withdrawal
	resp, err := cli.CreateWithdrawal(client.NewWithdrawalArgs(
		oc.Address(req.Address),
		symbol,
		network,
		amount,
	))

	if err != nil {
		return servererrors.Conflictf("failed to create withdrawal: %s", err)
	}

	return c.JSON(exportWithdrawal(resp))
}
