package main

import (
	"aglog/internal/config"
	"aglog/internal/log"
	"aglog/internal/loki"
	"aglog/internal/queue"
	"aglog/internal/server"
	"aglog/internal/worker"
	"context"
	"flag"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

var (
	configPath = flag.String("config", "aglog.yml", "path to configuration file")
)

func main() {
	flag.Parse()

	conf := config.ReadConfig(*configPath)
	logger := log.NewLog(log.WithLevel(conf.Log.Level), log.WithFormat(conf.Log.Format))
	lokiClient := loki.NewClient(conf.Loki.Server, loki.WithTimeout(conf.Loki.Timeout))
	messageQueue := queue.NewSqliteQueue(conf.Storage.PathToStorage, logger)
	udpServer := server.NewUdpServer(conf.Server.IP, conf.Server.Port, logger)

	var servers []server.Server
	servers = append(servers, udpServer)

	formats := generateFormats(conf.Messages.Formats)

	ctx, cancel := context.WithCancel(context.Background())
	exitChan := make(chan struct{}, 1)

	go worker.ListenMessages(ctx, servers, logger, messageQueue)
	go worker.Work(ctx, lokiClient, messageQueue, logger, formats)

	go gracefulShutdown(exitChan, cancel, udpServer.Stop, messageQueue.Close)

	<-exitChan
}

func generateFormats(formats []string) map[string]*regexp.Regexp {
	compiledFormats := make(map[string]*regexp.Regexp)

	for _, format := range formats {
		compiledFormats[format] = regexp.MustCompile(format)
	}

	return compiledFormats
}

func gracefulShutdown(exitChan chan struct{}, cancels ...func()) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	for _, cancel := range cancels {
		cancel()
	}

	exitChan <- struct{}{}
}
