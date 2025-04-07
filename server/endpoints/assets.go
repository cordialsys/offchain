package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/client/api"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

func exportAssets(resp []*oc.Asset) []*api.Asset {
	assets := make([]*api.Asset, len(resp))
	for i, a := range resp {
		assets[i] = &api.Asset{
			// TODO name
			Name:    nil,
			Network: string(a.NetworkId),
			Symbol:  string(a.SymbolId),
		}
	}
	return assets
}

// GetAssets returns the list of assets for a specified exchange
func GetAssets(c *fiber.Ctx) error {
	exchangeCfg, account, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, account)
	if err != nil {
		return servererrors.InternalErrorf(c, "failed to create client: %s", err)
	}

	// Get assets
	assets, err := cli.ListAssets()
	if err != nil {
		return servererrors.Conflictf(c, "failed to get assets: %s", err)
	}

	return c.JSON(exportAssets(assets))
}
