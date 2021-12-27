package github

// go test -v ./... -args -access-token={ACCESS_TOKEN}

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"
	_ "github.com/whosonfirst/go-writer-github"
)

var access_token = flag.String("access-token", "", "A valid GitHub access token")

func TestEnsureRepoForCurrentYear(t *testing.T) {

	if *access_token == "" {
		t.Fatalf("Missing -args -access-token={TOKEN} flag")
	}

	writer_uri := fmt.Sprintf("githubapi://sfomuseum-data/sfomuseum-data-test-{YYYY}?access_token=%s", *access_token)

	license_fh, err := os.Open("fixtures/LICENSE")

	if err != nil {
		t.Fatalf("Failed to open LICENSE, %v", err)
	}

	defer license_fh.Close()

	readme_fh, err := os.Open("fixtures/README.md.txt")

	if err != nil {
		t.Fatalf("Failed to open README, %v", err)
	}

	defer readme_fh.Close()

	opts := &EnsureRepoForCurrentYearOptions{
		Description: "{YYYY} data for testing.",
		Private:     false,
		License:     license_fh,
		Readme:      readme_fh,
	}

	ctx := context.Background()

	created, name, err := EnsureRepoForCurrentYear(ctx, writer_uri, opts)

	if err != nil {
		t.Fatalf("Failed to create repo, %v", err)
	}

	now := time.Now()
	yyyy := now.Year()

	expected_name := fmt.Sprintf("sfomuseum-data-test-%d", yyyy)

	if name != expected_name {
		t.Fatalf("Unexpected repo name. Expected '%s' but got '%s'", expected_name, name)
	}

	fmt.Printf("Created %s: %t\n", name, created)
}
