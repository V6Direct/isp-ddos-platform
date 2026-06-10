package main

import (
	"flag"
	"log"
	"os"

	"github.com/V6Direct/isp-ddos-platform/controller/api"
	"github.com/V6Direct/isp-ddos-platform/controller/db"
	"github.com/V6Direct/isp-ddos-platform/controller/config"
)

func main() {
	cfgPath := flag.String("config", "/etc/ddos/controller.yaml", "Path to config file")
	listen := flag.String("listen", "", "Override listen address")
	database := flag.String("db", "", "Override database path")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Printf("Warning: could not load config file: %v, using defaults", err)
		cfg = config.Default()
	}
	if *listen != "" {
		cfg.Listen = *listen
	}
	if *database != "" {
		cfg.Database = *database
	}

	if cfg.AdminToken == "" {
		log.Fatal("admin_token must be set in config or ADMIN_TOKEN env var")
	}

	database_conn, err := db.Open(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	if err := db.Migrate(database_conn); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Printf("DDoS Controller starting on %s", cfg.Listen)
	log.Printf("Database: %s", cfg.Database)

	server := api.NewServer(cfg, database_conn)
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
