package bybit

import (
	"fmt"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

func Validate(exchange *oc.ExchangeConfig) error {
	if exchange.Id == "" {
		return fmt.Errorf("bybit main account id is required; use the UID shown on the website")
	}
	_, err := strconv.ParseUint(string(exchange.Id), 10, 64)
	if err != nil {
		return fmt.Errorf("bybit main account id must be a numeric (not %s); use the UID shown on the website", exchange.Id)
	}
	for _, sub := range exchange.SubAccounts {
		_, err := strconv.ParseUint(string(sub.Id), 10, 64)
		if err != nil {
			return fmt.Errorf("bybit subaccount id must be a numeric (not %s); use the UID shown on the website", sub.Id)
		}
	}
	return nil
}
