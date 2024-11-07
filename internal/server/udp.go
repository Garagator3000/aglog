package server

import (
	"aglog/internal/log"
	"bytes"
	"context"
	"fmt"
	"net"
)

type UDP struct {
	ip         string
	port       string
	Connection *net.UDPConn
	logger     log.Logger
}

func NewUdpServer(ip string, port int, logger log.Logger) Server {
	return &UDP{
		ip:     ip,
		port:   fmt.Sprintf(":%d", port),
		logger: logger,
	}
}

func (u *UDP) GetAddr() string {
	return u.ip + u.port
}

func (u *UDP) Start() {
	addr, err := net.ResolveUDPAddr("udp", u.ip+u.port)
	if err != nil {
		panic(fmt.Errorf("failed to resolve udp address: %w", err))
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(fmt.Errorf("failed to listen udp: %w", err))
	}

	u.Connection = conn
}

func (u *UDP) Stop() {
	err := u.Connection.Close()
	if err != nil {
		u.logger.Error("failed to close udp connection", log.Error(err))
	}
}

func (u *UDP) Listen(ctx context.Context, buffer chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buffer <- u.read()
		}
	}
}

func (u *UDP) read() string {
	buffer := make([]byte, 2048)

	_, _, err := u.Connection.ReadFromUDP(buffer)
	if err != nil {
		u.logger.Error("failed to read udp message", log.Error(err))
	}

	buffer = bytes.ReplaceAll(buffer, []byte{0}, []byte{})
	buffer = bytes.TrimSuffix(buffer, []byte{'\n'})

	return string(buffer)
}
