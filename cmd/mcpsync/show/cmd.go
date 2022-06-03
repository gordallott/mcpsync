package show

import (
	"context"
	"fmt"
	"os"

	"github.com/gordallott/mcpsync/pkg/api"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "show",
	Short: "show will show information about the target MemCard Pro",
	Run:   run,
}

var (
	targetIP string
)

func init() {
	Cmd.Flags().StringVar(&targetIP, "ip", "", "IP address of the target MemCard Pro")
	if err := Cmd.MarkFlagRequired("ip"); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	version, err := api.GetVersion(ctx, targetIP)
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), err)
		os.Exit(1)
		return
	}

	fmt.Println("Version:", version.Version)

	cards, err := api.GetCards(ctx, targetIP)
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), err)
		os.Exit(1)
		return
	}

	if len(cards) > 0 {
		fmt.Println("Cards:\n--------")
	} else {
		fmt.Println("No cards found")
	}

	for _, card := range cards {
		fmt.Printf("\tGameID:   %s\n", card.GameID)
		fmt.Printf("\tName:     %s\n", card.Name)
		fmt.Printf("\tFullPath: %s\n", card.FullPath)
		fmt.Printf("--------\n")
	}

}
