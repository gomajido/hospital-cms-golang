package cmd

import (
	"os"

	"github.com/gomajido/hospital-cms-golang/internal/dependency"
	"github.com/spf13/cobra"
)

type CoreDependency struct {
	Drivers     *dependency.Drivers
	CommonRepos *dependency.CommonRepositories
	AppRepos    *dependency.AppRepositories
	AppUsecase  *dependency.AppUsecase
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nabitu-service",
	Short: "Nabitu Core Application",
	Long:  "Nabitu Core Application",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
