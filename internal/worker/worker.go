package worker

import (
	"aglog/internal/log"
	"aglog/internal/loki"
	"aglog/internal/queue"
	"aglog/internal/server"
	"context"
	"regexp"
	"time"
)

// Listen log messages from the all servers and send them to queue.
func ListenMessages(
	ctx context.Context,
	srv []server.Server,
	logger log.Logger,
	messageQueue queue.Queue,
) {
	messageChan := make(chan string, 1024)

	// Startring all servers.
	for _, s := range srv {
		s.Start()
		logger.Info("server listening", log.String("address", s.GetAddr()), log.String("protocol", s.GetProto()))
		go s.Listen(ctx, messageChan)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("logreceiver stopped by context")
			return
		case newlog := <-messageChan:
			messageQueue.Enqueue(newlog)
		}
	}
}

// Handle the log messages from the queue.
func Work(
	ctx context.Context,
	lokiClient *loki.Client,
	messageQueue queue.Queue,
	logger log.Logger,
	formats map[string]*regexp.Regexp,
) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("worker stopped by context")
			return
		default:
			handleLog(lokiClient, messageQueue, logger, formats)
		}
	}
}

// Handle one log message from queue.
func handleLog(
	lokiClient *loki.Client,
	messageQueue queue.Queue,
	logger log.Logger,
	formats map[string]*regexp.Regexp,
) {
	ts, message := messageQueue.Dequeue()
	if message == "" {
		return
	}

	labels := make(map[string]string)

	err := parseLog(labels, message, formats)
	if err != nil {
		logger.Error("failed to parse log", log.Error(err))
		return
	}

	labels["service_name"] = "aglog"
	labels["ts"] = time.Unix(0, ts).Format(time.RFC3339)

	lokiLog := loki.MapToLokiFormat(labels, ts, message)

	err = lokiClient.Push(lokiLog)
	if err != nil {
		logger.Error("failed to send log to loki", log.Error(err))
	}
}

// Add labels for Loki from log message by regex pattern.
func parseLog(labels map[string]string, logMessage string, formats map[string]*regexp.Regexp) error {
	for _, re := range formats {
		if !re.MatchString(logMessage) {
			continue
		}

		subexpNames := re.SubexpNames()
		matches := re.FindStringSubmatch(logMessage)

		if len(matches) != len(subexpNames)+1 {
			continue
		}

		for i, name := range subexpNames {
			if i != 0 && name != "" {
				labels[name] = matches[i]
			}
		}
		return nil
	}

	return nil
}
