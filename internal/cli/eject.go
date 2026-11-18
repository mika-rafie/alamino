package cli

import (
	"fmt"

	"github.com/mika-rafie/alamino/internal/gate"
	"github.com/mika-rafie/alamino/internal/safeurl"
	"github.com/spf13/cobra"
)

func newEjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "eject",
		Short: "Remove alamino gate from the current repository",
		Long: `Removes the "alamino" git remote when its only literal URL is inside this
home's repos directory (any other remote is left in place), deletes the bare
repo and worktrees, and removes the repo record from the database.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, d, err := openResources()
			if err != nil {
				return err
			}
			defer d.Close()

			repo, err := gate.Eject(cmd.Context(), d, p, ".")
			if err != nil {
				return fmt.Errorf("eject: %w", err)
			}

			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "  %s Gate removed\n", sGreen.Render("✓"))
			fmt.Fprintln(w)
			fmt.Fprintf(w, "  %s  %s\n", sDim.Render("  repo"), repo.WorkingPath)
			remoteURL := repo.UpstreamURL
			if repo.ForkURL != "" {
				remoteURL = safeurl.Redact(remoteURL)
			}
			fmt.Fprintf(w, "  %s  %s\n", sDim.Render("remote"), remoteURL)
			if repo.ForkURL != "" {
				fmt.Fprintf(w, "  %s  %s\n", sDim.Render("  fork"), safeurl.Redact(repo.ForkURL))
			}
			return nil
		},
	}
}
