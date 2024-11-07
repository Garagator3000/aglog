package loki

import "time"

type lokiClientParams struct {
	timeout time.Duration
}

func defaultParams() lokiClientParams {
	return lokiClientParams{
		timeout: 10 * time.Second,
	}
}

func WithTimeout(timeout string) func(p *lokiClientParams) {
	return func(p *lokiClientParams) {
		p.timeout, _ = time.ParseDuration(timeout)
	}
}
