package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/dependency"
	"github.com/gomajido/hospital-cms-golang/internal/router"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
)

func main() {
	ctx := context.Background()
	app_log.Info("Starting application...")

	appConfigs, err := config.GetConfig()
	if err != nil {
		app_log.Fatalf("Fail to GetConfig: %s\n", err)
	}
	app_log.Info("Config loaded successfully")

	app_log.Info("Initializing drivers...")
	drivers := dependency.InitDrivers(appConfigs)
	app_log.Info("Drivers initialized")

	app_log.Info("Initializing adapters...")
	adapters := dependency.InitAdapters(appConfigs, drivers)
	app_log.Info("Adapters initialized")

	app_log.Info("Initializing common repos...")
	commonRepos := dependency.InitCommonRepos(adapters, drivers, appConfigs)
	app_log.Info("Common repos initialized")

	app_log.Info("Initializing app repos...")
	appRepos := dependency.InitRepos(drivers.Db, drivers.Redis)
	app_log.Info("App repos initialized")

	app_log.Info("Initializing app usecase...")
	appUsecase := dependency.InitUsecase(appConfigs, appRepos, commonRepos)
	app_log.Info("App usecase initialized")

	app_log.Info("Initializing app handlers...")
	appHandlers := dependency.InitHandlers(ctx, appConfigs, appUsecase)
	app_log.Info("App handlers initialized")

	appRouter := router.Router{
		ApplicationHandler: appHandlers,
		HttpConfig:         appConfigs.Http,
	}

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		app_log.Info("Starting HTTP server...")
		if err := router.Run(&appRouter); err != nil && err != http.ErrServerClosed {
			app_log.Infof("App handlers initialized : %v", err)
			app_log.Fatalf("Server error: %s\n", err)
		}
	}()

	<-sig
	app_log.Info("Received shutdown signal")

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	defer func() {
		app_log.Info("Closing database connection...")
		err := drivers.Db.Close()
		if err != nil {
			app_log.Fatalf("Failed to close db: %s", err.Error())
		}
		app_log.Info("Database connection closed")
		serverStopCtx()
		<-serverCtx.Done()
	}()

	app_log.Info("Server Exited Properly")
}
