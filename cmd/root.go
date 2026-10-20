/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	"github.com/HenriqueOliveira29/cliToDoList/internal/config"
  	"github.com/HenriqueOliveira29/cliToDoList/internal/service"
	"github.com/HenriqueOliveira29/cliToDoList/internal/repository"

)

var (
	cfgPath string

	pool    *pgxpool.Pool
	svc     *service.TodoService
)

// rootCmd represents the base command when called without any subcommands
func initApp(cmd *cobra.Command) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	pool, err = pgxpool.New(ctx, cfg.Database.DSN())
	if err != nil {
		return fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	svc = service.NewToDoService(repository.NewPostgresRepo(pool))
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "clitodolist",
	Short: "Todo list CLI and API server",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initApp(cmd)
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if pool != nil {
			pool.Close()
		}
	},
}


func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "config.yaml", "path to config file")
}


