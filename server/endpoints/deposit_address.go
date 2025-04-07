package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

func exportDepositAddress(resp oc.Address) string {
	return string(resp)
}

// GetDepositAddress returns a deposit address for a specified symbol and network
func GetDepositAddress(c *fiber.Ctx) error {
	exchangeCfg, account, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Get required query parameters
	symbol := c.Query("symbol")
	if symbol == "" {
		return servererrors.BadRequestf(c, "symbol query parameter is required")
	}

	network := c.Query("network")
	if network == "" {
		return servererrors.BadRequestf(c, "network query parameter is required")
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, account)
	if err != nil {
		return servererrors.InternalErrorf(c, "failed to create client: %s", err)
	}

	// Prepare options
	options := []client.GetDepositAddressOption{}

	// Check if a specific subaccount was requested
	subaccount := c.Query("subaccount")
	if subaccount != "" {
		options = append(options, client.WithSubaccount(oc.AccountId(subaccount)))
	}

	// Get deposit address
	resp, err := cli.GetDepositAddress(client.NewGetDepositAddressArgs(
		oc.SymbolId(symbol),
		oc.NetworkId(network),
		options...,
	))

	if err != nil {
		return servererrors.Conflictf(c, "failed to get deposit address: %s", err)
	}

	return c.JSON(exportDepositAddress(resp))
}
