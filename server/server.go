package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"slices"
	"strings"
	"syscall"
	"time"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/pkg/httpsignature"
	"github.com/cordialsys/offchain/pkg/httpsignature/verifier"
	"github.com/cordialsys/offchain/server/endpoints"
	"github.com/cordialsys/offchain/server/servererrors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server represents the HTTP server
type Server struct {
	app    *fiber.App
	ocConf *oc.Config
	ServerArgs
}

type ServerArgs struct {
	Listen    string
	AnyOrigin bool
	Origins   []string

	BearerTokens        []string
	PublicKeys          []verifier.VerifierI
	PublicReadEndpoints bool
}

// New creates a new server instance
func New(ocConf *oc.Config, args ServerArgs) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if apiErr, ok := err.(*servererrors.ErrorResponse); ok {
				return apiErr.Send(c)
			} else {
				return servererrors.InternalErrorf("%v", err)
			}
		},
	})

	// Add middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			fmt.Printf("[PANIC] %v\n%s\n", e, debug.Stack())
		},
	}))
	app.Use(logger.New())

	allowOrigins := []string{"http://localhost:3000"}
	if args.AnyOrigin {
		allowOrigins = []string{"*"}
	} else if len(args.Origins) > 0 {
		allowOrigins = args.Origins
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(allowOrigins, ","),
	}))
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("conf", ocConf)

		// read subaccount if used in header or query
		subaccount := c.Get("sub-account")
		if subaccount == "" {
			subaccount = c.Query("sub-account")
		}
		c.Locals("sub-account", subaccount)

		return c.Next()
	})

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	httpSigAuth := func(c *fiber.Ctx) error {
		verified := false
		lastErr := fmt.Errorf("invalid signature")
		requiredHeaders := []string{}
		if c.Get("sub-account") != "" {
			// ensure sub-account is in the signature input if it is used
			requiredHeaders = append(requiredHeaders, "sub-account")
		}
		for _, key := range args.PublicKeys {
			_, lastErr := httpsignature.VerifyFiber(c, key, requiredHeaders...)
			if lastErr == nil {
				verified = true
				break
			}
		}
		if !verified {
			return servererrors.Unauthorizedf("%v", lastErr)
		}
		return c.Next()
	}

	bearerOrHttpSigAuth := func(c *fiber.Ctx) error {
		if args.PublicReadEndpoints {
			return c.Next()
		}
		// expect Authorization: Bearer <token>
		auth := strings.TrimSpace(c.Get("Authorization"))
		parts := strings.Split(auth, " ")
		if len(parts) != 2 {
			// check for http-signature
			if c.Get(httpsignature.HeaderSignature) != "" &&
				c.Get(httpsignature.HeaderSignatureInput) != "" &&
				c.Get(httpsignature.HeaderContentDigest) != "" {
				return httpSigAuth(c)
			}

			return servererrors.BadRequestf("no authorization header or http-signature")
		}
		if parts[0] != "Bearer" {
			return servererrors.BadRequestf("expected Bearer token in authorization header")
		}
		token := parts[1]
		if slices.Contains(args.BearerTokens, token) {
			return c.Next()
		}
		return servererrors.Unauthorizedf("invalid bearer token")
	}

	// API routes
	v1 := app.Group("/v1")
	// public
	v1.Get("/exchanges/:exchange/assets", endpoints.GetAssets)
	v1.Get("/exchanges/:exchange/account-types", endpoints.GetAccountTypes)

	// bearer or http sig auth
	v1.Get("/exchanges/:exchange/balances", bearerOrHttpSigAuth, endpoints.GetBalances)
	v1.Get("/exchanges/:exchange/deposit-address", bearerOrHttpSigAuth, endpoints.GetDepositAddress)
	v1.Get("/exchanges/:exchange/subaccounts", bearerOrHttpSigAuth, endpoints.ListSubaccounts)
	v1.Get("/exchanges/:exchange/withdrawal-history", bearerOrHttpSigAuth, endpoints.ListWithdrawalHistory)

	// http sig auth only
	v1.Post("/exchanges/:exchange/account-transfer", httpSigAuth, endpoints.AccountTransfer)
	v1.Post("/exchanges/:exchange/withdrawal", httpSigAuth, endpoints.CreateWithdrawal)

	return &Server{
		app:        app,
		ocConf:     ocConf,
		ServerArgs: args,
	}
}

// Start begins listening for requests
func (s *Server) Start() error {
	// Start server in a goroutine so we can handle graceful shutdown
	go func() {
		if err := s.app.Listen(s.Listen); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	fmt.Printf("Server started on %s\n", s.Listen)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}

	return nil
}
