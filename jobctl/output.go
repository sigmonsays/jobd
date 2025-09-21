package main

import (
	"context"
	"fmt"

	"github.com/sigmonsays/jobd/api"
	"github.com/spf13/cobra"
)

func init() {
	outputCmd := &cobra.Command{
		Use:   "output [job-id]",
		Short: "get job output",
		Long:  "return job output and output",
		Args:  cobra.ExactArgs(1),
		RunE:  outputJob,
	}
	rootCmd.AddCommand(outputCmd)
}
func outputJob(cmd *cobra.Command, args []string) error {
	appCtx, err := GetAppContext(cmd)
	if err != nil {
		return err
	}
	jid := args[0]

	ctx := context.Background()

	req := api.JobOutputParams{
		Jid: api.NewOptString(jid),
	}

	res, err := appCtx.ApiClient.JobOutput(ctx, req)
	if err != nil {
		return err
	}

	out, ok := res.(*api.JobOutputResponse)
	if !ok {
		return fmt.Errorf("Invalid response type %T", res)
	}

	fmt.Printf("%s\n", out.Output.Value)

	return nil
}
