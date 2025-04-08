package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/client/api"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

func exportBalances(resp []*client.BalanceDetail) []*api.Balance {
	balances := make([]*api.Balance, len(resp))
	for i, b := range resp {
		balances[i] = &api.Balance{
			// TODO name
			Asset:       nil,
			Network:     string(b.NetworkId),
			Symbol:      string(b.SymbolId),
			Available:   b.Available.String(),
			Unavailable: b.Unavailable.String(),
		}
	}
	return balances
}

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
		return servererrors.InternalErrorf("failed to create client: %s", err)
	}

	// Create balance args
	balanceArgs := client.NewGetBalanceArgs(oc.AccountType(accountType))

	// Get balances
	assets, err := cli.ListBalances(balanceArgs)
	if err != nil {
		return servererrors.Conflictf("failed to get balances: %s", err)
	}

	return c.JSON(exportBalances(assets))
}
