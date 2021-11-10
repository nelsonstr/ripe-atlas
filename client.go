package atlas // import "github.com/keltia/ripe-atlas"

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/keltia/proxy"
	"github.com/pkg/errors"
)

// NewClient is the first function to call.
// Yes, it does take multiple config
// and the last one wins.
func NewClient(cfgs ...Config) (*Client, error) {
	c := &Client{}
	for _, cfg := range cfgs {
		c.config = cfg
	}

	// This holds the global options
	c.opts = make(map[string]string)

	// If no log output is specified, use the default one
	if c.config.Log == nil {
		c.log = log.New(os.Stderr, "", log.LstdFlags|log.LUTC)
	} else {
		c.log = c.config.Log
	}

	// Set log levels
	if c.config.Verbose {
		c.level = 1
	}

	if c.config.Level != 0 {
		c.level = c.config.Level
	}

	// Final check
	if c.config.Level > 2 {
		c.level = 2
	}

	// Ensure this is not empty
	if c.config.endpoint == "" {
		c.config.endpoint = apiEndpoint
	}
	c.verbose("c.config=%#v", c.config)

	// Create and save the http.Client
	return c.addHTTPClient()
}

// HasAPIKey returns whether an API key is stored
func (c *Client) HasAPIKey() (string, bool) {
	if c.config.APIKey == "" {
		return "", false
	}
	return c.config.APIKey, true
}

// call is s shortcut
func (c *Client) call(req *http.Request) (*http.Response, error) {
	c.verbose("Full URL:\n%v", req.URL)

	myurl, _ := url.Parse(apiEndpoint)
	req.Header.Set("Host", myurl.Host)
	req.Header.Set("User-Agent", fmt.Sprintf("ripe-atlas/%s", ourVersion))

	return c.client.Do(req)
}

func (c *Client) addHTTPClient() (*Client, error) {
	_, transport := proxy.SetupTransport(apiEndpoint)
	if transport == nil {
		return c, errors.New("addhttpclient")
	}
	c.client = &http.Client{Transport: transport, Timeout: 20 * time.Second}
	return c, nil
}

// SetOption sets a global option
func (c *Client) SetOption(name, value string) *Client {
	if value != "" {
		c.opts[name] = value
	}
	return c
}

// GetVersion returns the API wrapper version
func GetVersion() string {
	return ourVersion
}
