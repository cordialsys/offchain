package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/client/api"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

func exportSubaccounts(resp []*oc.SubAccountHeader) []*api.SubAccountHeader {
	subaccounts := make([]*api.SubAccountHeader, len(resp))
	for i, s := range resp {
		subaccounts[i] = &api.SubAccountHeader{
			Id:    string(s.Id),
			Alias: api.As(s.Alias),
		}
	}
	return subaccounts
}

// ListSubaccounts returns the list of configured subaccounts on an exchange
func ListSubaccounts(c *fiber.Ctx) error {
	exchangeCfg, account, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, account)
	if err != nil {
		return servererrors.InternalErrorf(c, "failed to create client: %s", err)
	}

	// Get subaccounts
	subaccounts, err := cli.ListSubaccounts()
	if err != nil {
		return servererrors.Conflictf(c, "failed to list subaccounts: %s", err)
	}

	return c.JSON(exportSubaccounts(subaccounts))
}
