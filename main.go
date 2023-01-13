package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"goportscan/core/cmd"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	rootCmd := &cobra.Command{Use: "PortScan"}
	rootCmd.AddCommand(cmd.Consumer())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute failed", err.Error())
	}
}
