package db

import (
	"codematic-backend/internal/config"
	db "codematic-backend/internal/infrastructure/db/sqlc"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DBConn struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
	logger  *zap.Logger
	mu      sync.Mutex
	closed  bool
}

func InitDB(config *config.Config, logger *zap.Logger) *DBConn {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.PostgresUser,
		config.PostgresPass,
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresDB,
	)

	var pool *pgxpool.Pool

	for {
		poolConfig, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			logger.Error("cannot parse database config", zap.Error(err))
			logger.Info("Retrying to parse database config in 5s...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Configure connection pool settings
		poolConfig.MaxConnIdleTime = time.Minute * 15
		poolConfig.MaxConnLifetime = time.Hour * 1
		poolConfig.MaxConns = 20
		poolConfig.MinConns = 5

		enableQueryLogging := config.EnableDBQueryLogging

		// Attach zap tracer with ON/OFF control
		poolConfig.ConnConfig.Tracer = &PgxZapTracer{
			Logger:  logger,
			Enabled: enableQueryLogging,
		}

		pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			logger.Error("cannot open connection to database", zap.Error(err))
			logger.Info("Retrying to connect to database in 5s...")
			time.Sleep(5 * time.Second)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = pool.Ping(ctx)
		if err != nil {
			pool.Close()
			logger.Error("cannot ping database", zap.Error(err))
			logger.Info("Retrying to ping database in 5s...")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	queries := db.New(pool)

	logger.Info("database connection initialized successfully")

	return &DBConn{
		Pool:    pool,
		Queries: queries,
		logger:  logger,
		closed:  false,
	}
}

// Close safely closes the database connection pool
func (db *DBConn) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()

	if !db.closed {
		db.Pool.Close()
		db.closed = true
		db.logger.Info("database connection closed")
	}
}

// IsHealthy checks if the database connection is still operational
func (db *DBConn) IsHealthy(ctx context.Context) bool {
	db.mu.Lock()
	if db.closed {
		db.mu.Unlock()
		return false
	}
	db.mu.Unlock()

	err := db.Pool.Ping(ctx)
	return err == nil
}

// GetPoolStats returns connection pool statistics for monitoring
func (db *DBConn) GetPoolStats() map[string]interface{} {
	stats := db.Pool.Stat()
	return map[string]interface{}{
		"total_connections":    stats.TotalConns(),
		"idle_connections":     stats.IdleConns(),
		"acquired_connections": stats.AcquiredConns(),
	}
}
