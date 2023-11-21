package main

import (
	"blog/config"
	"blog/server"
	"flag"
	"github.com/sirupsen/logrus"
)

var (
	appConfig = flag.String("config", "config/app.yaml", "application config path")
)

func main() {
	flag.Args()

	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.JSONFormatter{})

	config, err := config.Parse(*appConfig)
	if err != nil {
		logger.Fatalf("Failed to parse config: %v", err)
	}

	s, err := server.New(config, logger)
	if err != nil {
		logger.Fatalf("Init server failed: %v", err)
	}
	if err := s.Run(); err != nil {
		logger.Fatalf("Start server failed: %v", err)
	}
}
