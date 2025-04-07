package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

// GetBalances returns the balances for a specified account
func GetBalances(c *fiber.Ctx) error {
	// Get account from query parameter, default to core funding if not provided
	accountType := c.Query("type")

	exchangeCfg, account, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, account)
	if err != nil {
		return servererrors.InternalErrorf(c, "failed to create client: %s", err)
	}

	// Create balance args
	balanceArgs := client.NewGetBalanceArgs(oc.AccountType(accountType))

	// Get balances
	assets, err := cli.ListBalances(balanceArgs)
	if err != nil {
		return servererrors.Conflictf(c, "failed to get balances: %s", err)
	}

	// TODO define API types
	return c.JSON(assets)
}
