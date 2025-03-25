package main

import (
	"log"

	"geoip"

	"go.uber.org/zap"
)

func run(logger *zap.Logger) error {
	client := geoip.NewHTTPClient()
	service := geoip.NewService(client, "")

	ip := "217.150.32.5"
	resp, err := service.Get(ip)
	if err != nil {
		return err
	}
	logger.Info("Single IP lookup", zap.String("ip", ip), zap.Any("response", resp))

	ips := []string{"217.150.32.5", "8.8.8.8", "1.1.1.1"}
	batchRes, err := service.GetBatch(ips)
	if err != nil {
		logger.Warn("Batch lookup completed with errors", zap.Error(err))
	}

	for ip, res := range batchRes {
		logger.Info("Batch IP lookup", zap.String("ip", ip), zap.Any("response", res))
	}

	return nil
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	if err := run(logger); err != nil {
		logger.Fatal("Application failed", zap.Error(err))
	}
}
