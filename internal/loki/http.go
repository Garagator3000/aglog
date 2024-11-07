package loki

import (
	"bytes"
	"net/http"
	"strings"
)

const (
	pushEP = "loki/api/v1/push"
)

func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) post(uri string, body []byte) (*http.Response, error) {
	uri = c.baseURL + "/" + strings.TrimPrefix(uri, "/")
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	return c.do(req)
}
