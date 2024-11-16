package cmd

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
	"github.com/spf13/cobra"
)

var restCmd = &cobra.Command{
	Use:   "serve-rest-api",
	Short: "API for Apexa Application",
	Long:  `The functionality is to be exposed as REST APIs as per the documentation provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		ServeRestAPI()
	},
}

func init() {
	rootCmd.AddCommand(restCmd)
}

func ServeRestAPI() {
	ctx := context.Background()
	appConfigs, err := config.GetConfig()
	if err != nil {
		app_log.Fatalf("Fail to GetConfig: %s\n", err)
	}

	drivers := dependency.InitDrivers(appConfigs)
	adapters := dependency.InitAdapters(appConfigs, drivers)
	commonRepos := dependency.InitCommonRepos(adapters, drivers, appConfigs)
	appRepos := dependency.InitRepos(drivers.Db, drivers.Redis)
	appUsecase := dependency.InitUsecase(appConfigs, appRepos, commonRepos)
	appHandlers := dependency.InitHandlers(ctx, appConfigs, appUsecase)

	appRouter := router.Router{
		ApplicationHandler: appHandlers,
		HttpConfig:         appConfigs.Http,
	}
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := router.Run(&appRouter); err != nil && err != http.ErrServerClosed {
			app_log.Fatalf("listen: %s\n", err)
		}
	}()
	<-sig
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	defer func() {
		//extra handling here
		err := drivers.Db.Close()
		if err != nil {
			app_log.Fatalf("failed close db %s", err.Error())
		}
		serverStopCtx()
		<-serverCtx.Done()
	}()
	app_log.Info("Server Exited Properly")
}
