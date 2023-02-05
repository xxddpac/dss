package main

import (
	"dss/core/cmd"
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// @title Distributed Scan Services
// @version 1.0
// @description Distributed Scan Services API DOCS
func main() {
	rootCmd := &cobra.Command{Use: "SecurityScan"}
	rootCmd.AddCommand(cmd.Consumer())
	rootCmd.AddCommand(cmd.Producer())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute failed", err.Error())
	}
}
