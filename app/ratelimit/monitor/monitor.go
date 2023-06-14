// Package monitor provides a command-line application to query the GitHub API for API rate limits and deliver a
// notification message if the number of remaining API calls is <= (n).
package monitor

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/v48/github"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-messenger"
	"github.com/sfomuseum/runtimevar"
	"golang.org/x/oauth2"
)

// Run executes the command-line application to query the GitHub API for API rate limits
// using the default `flag.FlagSet` instance.
func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

// Run executes the command-line application to query the GitHub API for API rate limits
// using a custom `flag.FlagSet` instance.
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SFOMUSEUM")

	if err != nil {
		return fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	if delivery_webhook_uri != "" {

		webhook_uri, err := runtimevar.StringVar(ctx, delivery_webhook_uri)

		if err != nil {
			return fmt.Errorf("Failed to resolve webhook URI, %w", err)
		}

		delivery_agent_uri = strings.Replace(delivery_agent_uri, "{webhook}", webhook_uri, 1)
	}

	m, err := messenger.NewDeliveryAgent(ctx, delivery_agent_uri)

	if err != nil {
		return fmt.Errorf("Failed to create delivery agent, %w", err)
	}

	github_token, err := runtimevar.StringVar(ctx, github_token_uri)

	if err != nil {
		return fmt.Errorf("Failed to resolve github_token_uri, %w", err)
	}

	token_source := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: github_token},
	)

	oauth2_client := oauth2.NewClient(ctx, token_source)
	github_client := github.NewClient(oauth2_client)

	runMonitor := func(ctx context.Context) error {

		limits, _, err := github_client.RateLimits(ctx)

		if err != nil {
			return fmt.Errorf("Failed to derive rate limits, %w", err)
		}

		if verbose {
			logger.Printf("%d of %d remaining, resets at %v", limits.Core.Remaining, limits.Core.Limit, limits.Core.Reset)
		}

		if limits.Core.Remaining > minimum_remaining {
			return nil
		}

		body := fmt.Sprintf("%d of %d, resets at %v", limits.Core.Remaining, limits.Core.Limit, limits.Core.Reset)

		msg := &messenger.Message{
			To:      deliver_to,
			From:    deliver_from,
			Subject: deliver_subject,
			Body:    body,
		}

		err = m.DeliverMessage(ctx, msg)

		if err != nil {
			return fmt.Errorf("Failed to deliver message, %w", err)
		}

		return nil
	}

	switch mode {
	case "cli":
		return runMonitor(ctx)
	case "lambda":
		lambda.Start(runMonitor)
		return nil
	default:
		return fmt.Errorf("Invalid mode")
	}

}
