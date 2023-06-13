package monitor

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var mode string
var verbose bool

var delivery_agent_uri string
var delivery_webhook_uri string
var github_token_uri string

var deliver_to string
var deliver_from string
var deliver_subject string

var minimum_remaining int

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("monitor")

	fs.StringVar(&mode, "mode", "cli", "Valid options are: cli, lambda")
	fs.BoolVar(&verbose, "verbose", false, "")

	fs.StringVar(&delivery_agent_uri, "delivery-agent-uri", "stdout://", "")
	fs.StringVar(&delivery_webhook_uri, "delivery-webhook-uri", "", "A valid gocloud.dev/runtimevar URI")

	fs.StringVar(&github_token_uri, "github-token-uri", "", "A valid gocloud.dev/runtimevar URI")

	fs.StringVar(&deliver_to, "deliver-to", "", "")
	fs.StringVar(&deliver_from, "deliver-from", "", "")
	fs.StringVar(&deliver_subject, "deliver-subject", "", "")

	fs.IntVar(&minimum_remaining, "minimum-remaining", 500, "")
	return fs
}
