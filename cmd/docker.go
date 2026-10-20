package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/HenriqueOliveira29/cliToDoList/internal/config"
)

// composeUp starts Postgres and blocks until its healthcheck passes.
// Credentials are passed to the child process only; nothing is written to disk.
func composeUp(cmd *cobra.Command) error {
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	db := cfg.Database

	c := exec.CommandContext(cmd.Context(), "docker", "compose", "up", "-d", "--wait")
	c.Env = append(os.Environ(),
		"POSTGRES_USER="+db.User,
		"POSTGRES_PASSWORD="+db.Password,
		"POSTGRES_DB="+db.Name,
		"PGADMIN_EMAIL="+cfg.PgAdmin.Email,
		"PGADMIN_PASSWORD="+cfg.PgAdmin.Password,
		fmt.Sprintf("POSTGRES_PORT=%d", db.Port),
	)
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("docker compose up: %w", err)
	}
	return nil
}