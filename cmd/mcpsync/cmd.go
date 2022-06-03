package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gordallott/mcpsync/cmd/mcpsync/show"
	"github.com/gordallott/mcpsync/pkg/sync"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sync",
	Short: "mcp sync will sync a wifi enabled MemCard Pro to destination directory",
	Long:  `mcp sync will sync a wifi enabled MemCard Pro to destination directory`,
	Run:   run,
}

var (
	targetIP string
	syncDest string

	pollFrequency time.Duration
)

func init() {
	rootCmd.Flags().StringVar(&targetIP, "ip", "", "IP address of the target MemCard Pro")
	if err := rootCmd.MarkFlagRequired("ip"); err != nil {
		panic(err)
	}

	rootCmd.Flags().StringVar(&syncDest, "dest", "", "Destination directory to sync to")
	if err := rootCmd.MarkFlagRequired("dest"); err != nil {
		panic(err)
	}

	rootCmd.Flags().DurationVar(&pollFrequency, "pollFrequency", time.Minute, "frequency to poll for saves at")
	rootCmd.AddCommand(show.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	err := sync.Sync(ctx, targetIP, syncDest, pollFrequency)
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), err)
		os.Exit(1)
		return
	}
}
