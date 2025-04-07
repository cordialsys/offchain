package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/server/endpoints"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server represents the HTTP server
type Server struct {
	app    *fiber.App
	listen string
	ocConf *oc.Config
}

// New creates a new server instance
func New(listen string, ocConf *oc.Config, serverConfig *Config) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(logger.New())

	allowOrigins := []string{"localhost:3000"}
	if serverConfig.AnyOrigin {
		allowOrigins = []string{"*"}
	} else if len(serverConfig.Origins) > 0 {
		allowOrigins = serverConfig.Origins
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(allowOrigins, ","),
	}))
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("conf", ocConf)

		// read subaccount if used in header or query
		subaccount := c.Get("sub-account")
		if subaccount != "" {
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

	// API routes
	v1 := app.Group("/v1")
	v1.Get("/exchanges/:exchange/balances", endpoints.GetBalances)
	v1.Get("/exchanges/:exchange/assets", endpoints.GetAssets)
	v1.Get("/exchanges/:exchange/deposit-address", endpoints.GetDepositAddress)
	v1.Get("/exchanges/:exchange/account-types", endpoints.GetAccountTypes)
	v1.Get("/exchanges/:exchange/subaccounts", endpoints.ListSubaccounts)
	v1.Get("/exchanges/:exchange/withdrawal-history", endpoints.ListWithdrawalHistory)

	v1.Post("/exchanges/:exchange/account-transfer", endpoints.AccountTransfer)
	v1.Post("/exchanges/:exchange/withdrawal", endpoints.CreateWithdrawal)

	return &Server{
		app:    app,
		listen: listen,
		ocConf: ocConf,
	}
}

// Start begins listening for requests
func (s *Server) Start() error {
	// Start server in a goroutine so we can handle graceful shutdown
	go func() {
		if err := s.app.Listen(s.listen); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	fmt.Printf("Server started on %s\n", s.listen)

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
