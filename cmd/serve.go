package main

import (
	"context"
	"lila_games/internal/database"
	"lila_games/pkg/config"
	"lila_games/pkg/server"
	"lila_games/pkg/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// initialize app's context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// load config from ENVs
	cfg, err := config.Get()
	if err != nil {
		log.Fatal("config load failed, err:" + err.Error())
	}
	// allow supporting service to boot up, only if asked
	if cfg.TimeToWaitForSupportingServicesToComeUp > 0 {
		log.Printf("allowing %v to supporting services for boot up", cfg.TimeToWaitForSupportingServicesToComeUp)
		time.Sleep(cfg.TimeToWaitForSupportingServicesToComeUp)
	}
	// create db connection
	conn, err := database.NewConnection(cfg.DB)
	if err != nil {
		log.Fatal("database connection failed, err: " + err.Error())
	}
	// run migration if asked
	if cfg.DB.RunMigration {
		err = conn.Migrate(ctx)
		if err != nil {
			log.Fatal("database migration failed, err: " + err.Error())
		}
	}
	// create game service
	svc := service.New(conn)
	// create server
	s := server.New(":1234", svc)
	// start server
	s.Start()
	//await interruption
	deathRay := make(chan os.Signal, 1)
	signal.Notify(deathRay, syscall.SIGABRT, syscall.SIGTERM)
	<-deathRay
	// stop server
	s.Stop(ctx)
}
