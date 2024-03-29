package main

import (
	"log"

	"github.com/spf13/cobra"
)

var cfgPath string
var env string
var envPath string

const (
	ProdEnv  = "prod"
	DevEnv   = "dev"
	TestEnv  = "test"
	LocalEnv = "local"
)

// rootCmd is the root of all sub commands in the binary
// it doesn't have a Run method as it executes other sub commands
var rootCmd = &cobra.Command{
	Use:     "go-rest-api",
	Short:   "go-rest-api is a http server to serve public facing api",
	Version: "1.0",
}

func init() {
	// Here all other sub commands should be registered to the rootCmd
	rootCmd.AddCommand(srvCmd)
	rootCmd.AddCommand(migrationRoot)

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
