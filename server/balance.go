package server

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/gofiber/fiber/v2"
)

// GetBalances returns the balances for a specified account
func GetBalances(c *fiber.Ctx) error {
	// Get account from query parameter, default to core funding if not provided
	exchangeId := oc.ExchangeId(c.Params("exchange"))
	accountMaybe := c.Params("account")
	accountType := c.Query("type")
	conf := unwrapConfig(c)

	_ = accountMaybe

	exchangeCfg, ok := conf.GetExchange(exchangeId)
	if !ok {
		return NotFoundf(c, "exchange not found: %s", exchangeId)
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, &exchangeCfg.MultiSecret)
	if err != nil {
		return InternalErrorf(c, "failed to create client: %s", err)
	}

	// Create balance args
	if accountType == "" {
		// accountMaybe = string(client.CoreFunding)
	}
	balanceArgs := client.NewGetBalanceArgs(oc.AccountType(accountType))

	// Get balances
	assets, err := cli.ListBalances(balanceArgs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get balances: " + err.Error(),
		})
	}

	// TODO define API types
	return c.JSON(assets)
}
