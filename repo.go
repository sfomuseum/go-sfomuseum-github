package github

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-github/organizations"
	"github.com/whosonfirst/go-writer"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// EnsureCurrentYearWithURI will replace any instances of `{YYYY}` in the net/url path component
// of 'uri' with the current year.
func EnsureCurrentYearWithURI(uri string) (string, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return "", err
	}

	path, err := ensureCurrentYearWithPath(u.Path)

	if err != nil {
		return "", err
	}

	u.Path = path
	uri = u.String()

	return uri, nil
}

// EnsureCurrentYearWithURI will replace any instances of `{YYYY}` in 'uri' with the current year.
// It is assumed that 'uri' is an escaped net/url path component.
func ensureCurrentYearWithPath(uri string) (string, error) {

	if strings.Contains(uri, "{YYYY}") {

		now := time.Now()
		this_year := fmt.Sprintf("%04d", now.Year())

		uri = strings.Replace(uri, "{YYYY}", this_year, -1)
	}

	return uri, nil
}

// EnsureRepoForCurrentYearOptions contains specific properties to assign when creating new repositories.
type EnsureRepoForCurrentYearOptions struct {
	// The description of the new repository being created.
	Description string
	// Whether or not the new new repository being created is public or private.
	Private bool
	// An optional `io.Reader` instance containing the body of a `LICENSE` file to create in the new repository.
	License io.ReadSeeker
	// An optional `io.Reader` instance containing the body of a `README.md` file to create in the new repository.
	Readme io.ReadSeeker
}

// EnsureRepoForCurrentYear checks whether 'writer_uri' contains a "{YYYY}" replacement string for a GitHub repository.
// If present it will replace "{YYYY}" with the current year and ensure that the corresponding repository exists on
// GitHub, creating it if not.
func EnsureRepoForCurrentYear(ctx context.Context, writer_uri string, opts *EnsureRepoForCurrentYearOptions) (bool, string, error) {

	wr_u, err := url.Parse(writer_uri)

	if err != nil {
		return false, "", fmt.Errorf("Failed to parse URI, %v", err)
	}

	re_yyyy, err := regexp.Compile(`sfomuseum-data-[a-z0-9\-]+-{YYYY}.*`)

	if err != nil {
		return false, "", fmt.Errorf("Failed to compile YYYY regular expression")
	}

	if !re_yyyy.MatchString(wr_u.Path) {
		return false, "", nil
	}

	t := time.Now()
	yyyy := t.Year()

	str_yyyy := strconv.Itoa(yyyy)

	wr_u.Path = strings.Replace(wr_u.Path, "{YYYY}", str_yyyy, 1)

	create_repo := false

	var repo_owner string
	var repo_name string
	var repo_token string

	// As in https://github.com/whosonfirst/go-writer-github/blob/main/api.go

	if wr_u.Scheme == "githubapi" {

		repo_owner = wr_u.Host

		path := strings.TrimLeft(wr_u.Path, "/")
		parts := strings.Split(path, "/")

		if len(parts) != 1 {
			return false, "", fmt.Errorf("Invalid path for repo")
		}

		repo_name = parts[0]

		q := wr_u.Query()
		repo_token = q.Get("access_token")

		list_opts := &organizations.ListOptions{
			Prefix: []string{repo_name},
		}

		// This has a known-known bug listing private repos
		// https://github.com/whosonfirst/go-whosonfirst-github/issues/13

		repos, err := organizations.ListRepos(repo_owner, list_opts)

		if err != nil {
			return false, repo_name, fmt.Errorf("Failed to list repos, %v", err)
		}

		if len(repos) == 0 {
			create_repo = true
		}
	}

	if create_repo {

		desc := strings.Replace(opts.Description, "{YYYY}", str_yyyy, 1)

		create_opts := &organizations.CreateOptions{
			AccessToken: repo_token,
			Name:        repo_name,
			Description: desc,
			Private:     opts.Private,
		}

		err := organizations.CreateRepo(repo_owner, create_opts)

		if err != nil {
			return true, repo_name, fmt.Errorf("Failed to create repo '%s', %v", repo_name, err)
		}

		// Create a new writer with no prefix

		q := wr_u.Query()
		q.Del("prefix")

		wr_u.RawQuery = q.Encode()

		wr_uri := wr_u.String()
		repo_wr, err := writer.NewWriter(ctx, wr_uri)

		if err != nil {
			return true, repo_name, fmt.Errorf("Failed to create new writer for '%s', %v", wr_uri, err)
		}

		if opts.License != nil {

			_, err = repo_wr.Write(ctx, "LICENSE", opts.License)

			if err != nil {
				return true, repo_name, fmt.Errorf("Failed to write LICENSE file, %v", err)
			}
		}

		if opts.Readme != nil {

			readme_t, err := io.ReadAll(opts.Readme)

			if err != nil {
				return true, repo_name, fmt.Errorf("Failed to create README file, %v", err)
			}

			t, err := template.New("README").Parse(string(readme_t))

			if err != nil {
				return true, repo_name, fmt.Errorf("Failed to parse README template, %v", err)
			}

			vars := struct {
				Year int
			}{
				Year: yyyy,
			}

			var buf bytes.Buffer
			buf_wr := bufio.NewWriter(&buf)

			err = t.Execute(buf_wr, vars)

			if err != nil {
				return true, repo_name, fmt.Errorf("Failed to render README template, %v", err)
			}

			buf_wr.Flush()

			readme_r := bytes.NewReader(buf.Bytes())

			_, err = repo_wr.Write(ctx, "README.md", readme_r)

			if err != nil {
				return true, repo_name, fmt.Errorf("Failed to write README.md file, %v", err)
			}
		}
	}

	return create_repo, repo_name, nil
}
