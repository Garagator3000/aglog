package loki

import "fmt"

type Log struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	Stream any        `json:"stream"`
	Values [][]string `json:"values"`
}

func MapToLokiFormat(labels map[string]string, timestamp int64, log string) Log {
	return Log{
		Streams: []Stream{
			{
				Stream: labels,
				Values: [][]string{
					{
						fmt.Sprintf("%d", timestamp), log,
					},
				},
			},
		},
	}
}
