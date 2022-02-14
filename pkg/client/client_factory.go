package client

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
	"k8s.io/test-infra/prow/bugzilla"
	"net/http"
	totbz "tot/pkg/bugzilla"
)

type TotClientFactory interface {
	NewBugzillaClient() (totbz.TotBugzillaClient, error)
	NewJiraClient() (*jira.Client, error)
	NewGitHubClient() (*github.Client, error)
}

func NewTotClientFactory() TotClientFactory {
	return &totClientFactory{}
}

type totClientFactory struct {
}

func (t *totClientFactory) NewBugzillaClient() (totbz.TotBugzillaClient, error) {
	apiKey := func() []byte {
		return []byte("Your BZ API key")
	}
	bzClient := bugzilla.NewClient(apiKey, "https://bugzilla.redhat.com/", 1)
	return totbz.WrapBugzillaClient(bzClient), nil
}

func (t *totClientFactory) NewJiraClient() (*jira.Client, error) {
	transport := TokenAuthTransport{Token: "your jira token"}
	return jira.NewClient(transport.Client(), "https://issues.redhat.com")
}

func (t *totClientFactory) NewGitHubClient() (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "your gh token"},
	)
	return github.NewClient(oauth2.NewClient(ctx, ts)), nil
}

type TokenAuthTransport struct {
	Token string

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

func (t *TokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := cloneRequest(req) // per RoundTripper contract
	req2.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.Token))
	return t.transport().RoundTrip(req2)
}

func (t *TokenAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *TokenAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
