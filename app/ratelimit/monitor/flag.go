package monitor

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var mode string
var verbose bool

// TBD - Whether or not we want to support multiple DeliveryAgents and whether
// the webhook URI substitution stuff just applies to all the agents URIs or
// requires being clever.

var delivery_agent_uri string
var delivery_webhook_uri string
var github_token_uri string

var deliver_to string
var deliver_from string
var deliver_subject string

var minimum_remaining int

// DefaultFlagSet returns a default `flag.FlagSet` instance configured with the necessary flags to
// execute the command-line application to query the GitHub API for API rate limits.
func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("monitor")

	fs.StringVar(&mode, "mode", "cli", "Valid options are: cli, lambda")
	fs.BoolVar(&verbose, "verbose", false, "Log rate limit statistics.")

	fs.StringVar(&delivery_agent_uri, "delivery-agent-uri", "stdout://", "")
	fs.StringVar(&delivery_webhook_uri, "delivery-webhook-uri", "", "A valid gocloud.dev/runtimevar URI")

	fs.StringVar(&github_token_uri, "github-token-uri", "", "A valid gocloud.dev/runtimevar URI")

	fs.StringVar(&deliver_to, "deliver-to", "", "The address a message will be delivered from.")
	fs.StringVar(&deliver_from, "deliver-from", "", "The address a message will be delivered from.")
	fs.StringVar(&deliver_subject, "deliver-subject", "", "The subject of the message to delivered.")

	fs.IntVar(&minimum_remaining, "minimum-remaining", 500, "The minimum number of remaining API calls to watch for. If this number is greater than the number reported by the GitHub API then a notification message will be dispatched.")
	return fs
}
