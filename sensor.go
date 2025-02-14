package serial

import (
	"context"
	"net"
	"sync"
	"time"
)

type MessageHandler func(data []byte, sensorId int, deviceType int, parsed *ParserComplex)

type Sensor interface {
	Start(ctx context.Context)
	Stop()
	ID() int
	Info() string
	SetHandler(handler MessageHandler)
}

type BaseSensor struct {
	id             int
	host           string
	projectId      int
	port           int
	deviceAddr     string
	messageHandler MessageHandler
	conn           net.Conn
	stopChan       chan struct{}
	mu             sync.Mutex
	isRunning      bool
	deviceType     int
	restartTimeout time.Duration
}

func (s *BaseSensor) SetHandler(handler MessageHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageHandler = handler
}
