package server

import (
	"aglog/internal/config"
	"aglog/internal/log"
	"context"
)

type Server interface {
	GetAddr() string
	GetProto() string

	Start()
	Stop()

	Listen(ctx context.Context, buffer chan string)
}

func NewServers(conf config.Server, logger log.Logger) []Server {
	var servers []Server

	if conf.UDP.Enabled {
		udpServer := NewUdpServer(conf.UDP.IP, conf.UDP.Port, logger)
		servers = append(servers, udpServer)
	}

	if conf.TCP.Enabled {
		tcpServer := NewTCPServer(conf.TCP.IP, conf.TCP.Port, logger)
		servers = append(servers, tcpServer)
	}

	return servers
}
