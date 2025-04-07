package endpoints

import (
	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

func UnwrapConfig(c *fiber.Ctx) *oc.Config {
	return c.Locals("conf").(*oc.Config)
}
func WrapConfig(c *fiber.Ctx, conf *oc.Config) {
	c.Locals("conf", conf)
}

func loadAccount(c *fiber.Ctx, exchangeId string) (*oc.ExchangeConfig, *oc.Account, error) {
	conf := UnwrapConfig(c)
	exchangeConfig, ok := conf.GetExchange(oc.ExchangeId(exchangeId))
	if !ok {
		return nil, nil, servererrors.NotFoundf(c, "exchange not found: %s", exchangeId)
	}
	subaccountId := c.Locals("sub-account").(string)
	if subaccountId == "" {
		// done
		acc := exchangeConfig.AsAccount()
		err := acc.LoadSecrets()
		if err != nil {
			return nil, nil, servererrors.InternalErrorf(c, "failed to load secrets for main account of %s: %s", exchangeId, err)
		}
		return exchangeConfig, acc, nil
	}

	var account *oc.Account
	for _, subaccount := range exchangeConfig.SubAccounts {
		if string(subaccount.Id) == subaccountId || subaccount.Alias == subaccountId {
			account = subaccount.AsAccount()
		}
	}
	if account == nil {
		return nil, nil, servererrors.NotFoundf(c, "subaccount %s for exchange %s not found", subaccountId, exchangeId)
	}
	err := account.LoadSecrets()
	if err != nil {
		return nil, nil, servererrors.InternalErrorf(c, "failed to load secrets for %s subaccount %s: %v", exchangeId, subaccountId, err)
	}
	return exchangeConfig, account, nil
}
