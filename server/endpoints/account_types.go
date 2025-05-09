package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/pkg/secret"
	"github.com/cordialsys/offchain/server/client/api"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

// NopAccount is a placeholder account with dummy credentials
var NopAccount = &oc.Account{
	Id:         "",
	SubAccount: false,
	MultiSecret: oc.MultiSecret{
		ApiKeyRef:     secret.Secret("raw:AAA"),
		SecretKeyRef:  secret.Secret("raw:AAA"),
		PassphraseRef: secret.Secret("raw:AAA"),
	},
}

func exportAccountTypes(resp []*oc.AccountTypeConfig) []*api.AccountType {
	types := make([]*api.AccountType, len(resp))
	for i, t := range resp {
		types[i] = &api.AccountType{
			Type:        api.AccountTypeID(t.Type),
			Description: api.As(t.Description),
			Aliases:     t.Aliases,
		}
	}
	return types
}

// GetAccountTypes returns the list of valid account types for an exchange
func GetAccountTypes(c *fiber.Ctx) error {
	// Get configuration from context
	conf, ok := c.Locals("conf").(*oc.Config)
	if !ok {
		return servererrors.InternalErrorf("configuration not found in context")
	}

	// Get exchange ID from path parameter
	exchangeId := oc.ExchangeId(c.Params("exchange"))

	// Get exchange configuration
	exchangeConfig, ok := conf.GetExchange(exchangeId)
	if !ok {
		return servererrors.NotFoundf("exchange not found. options are: %v", oc.ValidExchangeIds)
	}

	// Create client with NopAccount (we don't need real credentials for this operation)
	cli, err := loader.NewClient(exchangeConfig, NopAccount)
	if err != nil {
		return servererrors.InternalErrorf("failed to create client: %s", err)
	}

	// Get account types
	accountTypes, err := cli.ListAccountTypes()
	if err != nil {
		return servererrors.Conflictf("failed to get account types: %s", err)
	}

	return c.JSON(exportAccountTypes(accountTypes))
}
