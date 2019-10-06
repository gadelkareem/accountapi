package accountapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	Version        = "v1"
	defaultTimeout = 15 * time.Second
	userAgent      = "account_api_client/" + Version + " " + runtime.GOOS + " " + runtime.GOARCH
)

type client struct {
	*http.Client
	endpoint *url.URL
}

func newClient(endpoint *url.URL, timeout time.Duration) *client {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	cl := http.DefaultClient
	cl.Timeout = timeout
	return &client{Client: cl, endpoint: endpoint}
}

func (c *client) buildURL(resourcePath string) string {
	u := c.endpoint
	u.Path = path.Join(Version, resourcePath)
	return u.String()
}

func (c *client) post(resourcePath string, body io.Reader) (*http.Response, error) {
	return c.retryRequest(http.MethodPost, resourcePath, body)
}

func (c *client) get(resourcePath string) (*http.Response, error) {
	return c.retryRequest(http.MethodGet, resourcePath, nil)
}

func (c *client) delete(resourcePath string) (*http.Response, error) {
	return c.retryRequest(http.MethodDelete, resourcePath, nil)
}

func (c *client) retryRequest(method, resourcePath string, body io.Reader) (resp *http.Response, err error) {
	backoff.Retry(func() error {
		resp, err = c.request(method, resourcePath, body)
		if resp.StatusCode >= http.StatusTooManyRequests {
			return errors.New("try again")
		}
		return nil
	}, backoff.NewExponentialBackOff())
	return
}

func (c *client) request(method, resourcePath string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, c.buildURL(resourcePath), body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", userAgent)
	r.Header.Set("Accept", "application/vnd.api+json")
	if method == http.MethodPost {
		r.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return resp, formatError(resp)
	}

	return resp, nil
}

func formatError(r *http.Response) error {
	defer r.Body.Close()

	pl := &struct {
		ErrorMessage string `json:"error_message"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(pl); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	} else if pl.ErrorMessage != "" {
		return errors.New(pl.ErrorMessage)
	}

	return fmt.Errorf("error status code %d: %s", r.StatusCode, http.StatusText(r.StatusCode))
}
