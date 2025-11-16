package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"pr-manager-service/config"

	httpadapter "pr-manager-service/internal/adapters/httpadapter"
	metricsadapter "pr-manager-service/internal/adapters/metricsadapter"
	repo "pr-manager-service/internal/repository"
	uc "pr-manager-service/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	kitlogger "github.com/nikitadev-work/avito-test-task-internship-autumn-2025/common/kit/logger"
	"github.com/nikitadev-work/avito-test-task-internship-autumn-2025/common/kit/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(ctx context.Context, cfg *config.Config) error {
	// logger
	l := kitlogger.NewLogger(
		cfg.Log.Level,
		map[string]any{
			"service": cfg.App.Name,
			"version": cfg.App.Version,
		},
	)

	l.Info("start configuration", nil)

	// metrics
	metrics.InitMetrics()

	// business metrics adapter
	businessMetrics := metricsadapter.NewMetrics(cfg.App.Name)

	// postgresql
	sslMode := "require"
	if !cfg.PostgreSQL.SslEnabled {
		sslMode = "disable"
	}
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgreSQL.User, cfg.PostgreSQL.Password, cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port, cfg.PostgreSQL.Name, sslMode)
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		l.Error("unable to create connection pool", map[string]any{
			"error": err.Error(),
		})
		return err
	}
	defer pool.Close()

	// repositories
	teamRepo := repo.NewTeamRepository(pool)
	userRepo := repo.NewUserRepository(pool)
	prRepo := repo.NewPullRequestRepository(pool)

	// usecase
	usecase := uc.NewService(teamRepo, userRepo, prRepo, l, businessMetrics)

	// http
	httpMux := httpadapter.NewRouter(usecase)
	httpMux.Handle("/metrics", promhttp.Handler())

	handlerWithMetrics := metrics.HTTPMiddleware(cfg.App.Name, httpMux)

	httpAddr := ":" + cfg.HTTP.Port
	httpServer := httpadapter.NewServer(httpAddr, handlerWithMetrics)

	httpErrCh := make(chan error, 1)
	go func() {
		l.Info("start http server", map[string]any{
			"http.addr": httpAddr,
		})
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			httpErrCh <- err
		}
	}()

	l.Info("pr-manager-service service started", map[string]any{
		"http.port": cfg.HTTP.Port,
		"log.level": cfg.Log.Level,
		"db.name":   cfg.PostgreSQL.Name,
		"db.host":   cfg.PostgreSQL.Host,
		"db.port":   cfg.PostgreSQL.Port,
	})

	// gracefull shutdown
	select {
	case <-ctx.Done():
		l.Info("starting graceful shutdown", nil)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		done := make(chan struct{})

		go func() {
			var wg sync.WaitGroup
			wg.Add(1)

			// http server
			go func() {
				defer wg.Done()
				if err := httpServer.Shutdown(shutdownCtx); err != nil {
					l.Error("http server graceful shutdown error", map[string]any{
						"error": err.Error(),
					})
				}
			}()

			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// successfully finished
			l.Info("gracefully finished", nil)
			return nil
		case <-shutdownCtx.Done():
			if err := httpServer.Close(); err != nil {
				l.Error("failed to close http server", map[string]any{
					"error": err.Error(),
				})
			}
			err := errors.New("graceful shutdown timeout")
			l.Error("graceful shutdown error", map[string]any{
				"error": err.Error(),
			})
			return err
		}
	case err := <-httpErrCh:
		l.Error("http server error", map[string]any{
			"error": err.Error(),
		})
		return err
	}
}
