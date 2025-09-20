package main

import (
	"fmt"
	"os"

	"github.com/sigmonsays/jobd/api"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jobd",
	Short: "A CLI client for jobd API",
	Long:  "A command line interface for interacting with the jobd CI/CD job daemon",
}

type Options struct {
	BaseUrl string
	ApiKey  string
	Timeout int
}

func main() {

	opts := &Options{
		Timeout: 30,
		BaseUrl: "http://localhost:8093",
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&opts.BaseUrl, "url", opts.BaseUrl, "Base URL of the jobd server")
	rootCmd.PersistentFlags().StringVar(&opts.ApiKey, "api-key", "", "API key for authentication (can also use JOBD_API_KEY env var)")
	rootCmd.PersistentFlags().IntVar(&opts.Timeout, "timeout", opts.Timeout, "HTTP request timeout in seconds")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type AppContext struct {
	Options   *Options
	ApiClient *api.Client
}

func GetAppContext(cmd *cobra.Command) (*AppContext, error) {
	opts := &Options{}
	opts.BaseUrl, _ = cmd.Flags().GetString("url")
	opts.ApiKey, _ = cmd.Flags().GetString("api-key")
	opts.Timeout, _ = cmd.Flags().GetInt("timeout")

	ret := &AppContext{
		Options: opts,
	}

	apiClient, err := api.NewClient(opts.BaseUrl)
	if err != nil {
		return ret, err
	}
	ret.ApiClient = apiClient

	return ret, nil
}
