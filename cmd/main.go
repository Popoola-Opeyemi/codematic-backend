package main

import "codematic/internal/config"

func main() {
	cfg := config.LoadAppConfig()
	zapLogger := config.InitLogger()
	defer zapLogger.Close()

	_ = cfg

}
