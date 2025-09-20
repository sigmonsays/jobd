package main

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	// Add list command for completeness
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all jobs",
		Long:  "List all available jobs from the jobd server",
		RunE:  listJobs,
	}
	rootCmd.AddCommand(listCmd)
}

func listJobs(cmd *cobra.Command, args []string) error {
	appCtx, err := GetAppContext(cmd)
	if err != nil {
		return err
	}

	ctx := context.Background()

	jobs, err := appCtx.ApiClient.ListJob(ctx)
	if err != nil {
		return err
	}

	PrintJson(jobs)

	return nil
}
