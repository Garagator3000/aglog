package server

import (
	"aglog/internal/log"
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
)

type TCP struct {
	ip       string
	port     string
	Listener *net.TCPListener
	logger   log.Logger
}

func NewTCPServer(ip string, port int, logger log.Logger) Server {
	return &TCP{
		ip:     ip,
		port:   fmt.Sprintf(":%d", port),
		logger: logger,
	}
}

func (t *TCP) GetAddr() string {
	return t.ip + t.port
}

func (t *TCP) GetProto() string {
	return "tcp"
}

func (t *TCP) Start() {
	addr, err := net.ResolveTCPAddr("tcp", t.ip+t.port)
	if err != nil {
		panic(fmt.Errorf("failed to resolve tcp address: %w", err))
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(fmt.Errorf("failed to listen tcp: %w", err))
	}

	t.Listener = listener
}
func (t *TCP) Stop() {
	t.logger.Info("stopping TCP server...")
	if err := t.Listener.Close(); err != nil {
		t.logger.Error("failed to close tcp connection", log.Error(err))
	}
}

func (t *TCP) Listen(ctx context.Context, buffer chan string) {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			t.logger.Error("failed to accept tcp connection", log.Error(err))
			return
		}

		go t.handleConn(conn, buffer)
	}
}

func (t *TCP) handleConn(conn net.Conn, buffer chan string) {
	defer handleCloser(conn, t.logger)

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			t.logger.Error("failed to read tcp message", log.Error(err))
		}

		buffer <- message
	}
}

func handleCloser(closer io.Closer, logger log.Logger) {
	if err := closer.Close(); err != nil {
		logger.Error("failed to close", log.Error(err))
	}
}
