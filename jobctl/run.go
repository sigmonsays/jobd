package main

import (
	"context"

	"github.com/sigmonsays/jobd/api"
	"github.com/spf13/cobra"
)

func init() {
	runCmd := &cobra.Command{
		Use:   "run [job-id]",
		Short: "Run a job",
		Long:  "Execute a job by its ID using the jobd API",
		Args:  cobra.ExactArgs(1),
		RunE:  runJob,
	}
	rootCmd.AddCommand(runCmd)
}
func runJob(cmd *cobra.Command, args []string) error {
	appCtx, err := GetAppContext(cmd)
	if err != nil {
		return err
	}
	jid := args[0]

	ctx := context.Background()

	req := &api.RunJobRequest{
		Jobid: api.NewOptString(jid),
	}

	res, err := appCtx.ApiClient.RunJob(ctx, req)
	if err != nil {
		return err
	}

	PrintJson(res)

	return nil
}
