package main

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cobra.OnInitialize(func() {
		viper.AutomaticEnv()
	})

	cmd := &cobra.Command{
		Run: run,
	}

	viper.BindPFlags(cmd.Flags())

	cmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	log.Println("REPL")

	ctx := cmd.Context()
	<-ctx.Done()
	_, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

}
