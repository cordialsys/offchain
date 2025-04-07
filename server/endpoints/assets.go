package endpoints

import (
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

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

	return c.JSON(assets)
}
