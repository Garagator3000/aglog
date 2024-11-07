package loki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(baseURL string, opts ...func(params *lokiClientParams)) *Client {
	params := defaultParams()
	for _, opt := range opts {
		opt(&params)
	}

	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: params.timeout,
		},
	}
}

func (c *Client) Push(log Log) error {
	body, err := json.Marshal(log)
	if err != nil {
		return err
	}

	resp, err := c.post(pushEP, body)
	if resp == nil {
		return fmt.Errorf("nil response from loki server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to push to loki: %s", resp.Status)
	}

	return err
}
