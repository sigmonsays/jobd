package main

import (
	"context"

	"github.com/sigmonsays/jobd/api"
	"github.com/spf13/cobra"
)

func init() {
	detailCmd := &cobra.Command{
		Use:   "detail [job-id]",
		Short: "get job detail",
		Long:  "return job detail and output",
		Args:  cobra.ExactArgs(1),
		RunE:  detailJob,
	}
	rootCmd.AddCommand(detailCmd)
}
func detailJob(cmd *cobra.Command, args []string) error {
	appCtx, err := GetAppContext(cmd)
	if err != nil {
		return err
	}
	jid := args[0]

	ctx := context.Background()

	req := api.JobDetailParams{
		Jid: api.NewOptString(jid),
	}

	res, err := appCtx.ApiClient.JobDetail(ctx, req)
	if err != nil {
		return err
	}

	PrintJson(res)

	return nil
}
