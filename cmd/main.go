package main

import (
	"codematic-backend/internal/config"
	"codematic-backend/internal/infrastructure/db"
	"codematic-backend/internal/infrastructure/search"
)

func main() {
	cfg := config.LoadAppConfig()
	zapLogger := config.InitLogger()
	defer zapLogger.Close()

	store := db.InitDB(cfg, zapLogger.Logger)
	search := search.InitES(cfg)

	_ = store
	_ = search

}
