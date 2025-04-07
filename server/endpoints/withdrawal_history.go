package endpoints

import (
	"strconv"

	"github.com/cordialsys/offchain/client"
	"github.com/cordialsys/offchain/loader"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
)

// ListWithdrawalHistory returns the withdrawal history for an exchange account
func ListWithdrawalHistory(c *fiber.Ctx) error {
	exchangeCfg, account, err := loadAccount(c, c.Params("exchange"))
	if err != nil {
		return err
	}

	// Create client
	cli, err := loader.NewClient(exchangeCfg, account)
	if err != nil {
		return servererrors.InternalErrorf(c, "failed to create client: %s", err)
	}

	// Create withdrawal history arguments
	args := client.NewWithdrawalHistoryArgs()

	// Handle limit parameter
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return servererrors.BadRequestf(c, "invalid limit parameter: must be a number")
		}
		if limit <= 0 {
			return servererrors.BadRequestf(c, "invalid limit parameter: must be greater than 0")
		}
		args.SetLimit(limit)
	}

	// Handle page token parameter
	if pageToken := c.Query("page_token"); pageToken != "" {
		args.SetPageToken(pageToken)
	}

	// Get withdrawal history
	resp, err := cli.ListWithdrawalHistory(args)
	if err != nil {
		return servererrors.Conflictf(c, "failed to get withdrawal history: %s", err)
	}

	return c.JSON(resp)
}
