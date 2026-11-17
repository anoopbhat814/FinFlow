package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/finflow/ledger/internal/config"
	"github.com/finflow/ledger/internal/handlers"
	"github.com/finflow/ledger/internal/middleware"
	"github.com/finflow/ledger/internal/repository"
	"github.com/finflow/ledger/internal/service"
	"github.com/finflow/ledger/pkg/idempotency"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if printed, err := parseFlags(os.Args[1:]); err != nil {
		os.Exit(2)
	} else if printed {
		return
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		log.Error("config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	accounts := repository.NewAccountRepo(pool)
	ledger := repository.NewLedgerRepo(pool)
	fraud := service.NewHTTPFraudClient(cfg.FraudURL, cfg.FraudTimeout)
	transfers := service.NewTransferService(accounts, ledger, fraud)
	statements := service.NewLedgerService(accounts, ledger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health(pool))
	handlers.NewAccountHandler(accounts, statements).Register(mux)
	handlers.NewTransferHandler(transfers, idempotency.NewStore(24*time.Hour)).Register(mux)

	var h http.Handler = mux
	for _, mw := range []func(http.Handler) http.Handler{
		middleware.RateLimit(cfg.RateLimitRPS),
		middleware.Auth(cfg.JWTSecret),
		middleware.Logging(log),
		middleware.Recovery(log),
	} {
		h = mw(h)
	}

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: h, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()

	log.Info("ledger listening", "port", cfg.Port)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("server", "err", err)
		os.Exit(1)
	}
}
