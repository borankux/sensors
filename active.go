package serial

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func triggerRestart(ch chan struct{}) {
	select {
	case ch <- struct{}{}:
	default:
	}
}

type Active struct {
	BaseSensor
	frequency      time.Duration
	ticker         *time.Ticker
	restartTimeout time.Duration
}

func NewActive(id, projectId int, host, address string, port int, frequency time.Duration, deviceType int, restartTimeout time.Duration) *Active {
	return &Active{
		BaseSensor: BaseSensor{
			id:         id,
			host:       host,
			projectId:  projectId,
			port:       port,
			deviceAddr: address,
			stopChan:   make(chan struct{}),
			mu:         sync.Mutex{},
			isRunning:  false,
			deviceType: deviceType,
		},
		frequency:      frequency,
		restartTimeout: restartTimeout,
	}
}

func (a *Active) startActiveDemon(ctx context.Context, restartChan chan struct{}) {
	for {
		select {
		case <-a.stopChan:
			return
		case <-ctx.Done():
			return
		case <-restartChan:
			time.Sleep(a.restartTimeout)
			go a.startIt(ctx, restartChan)
		}
	}
}

func (a *Active) startIt(ctx context.Context, restartChan chan struct{}) {
	address := fmt.Sprintf("%s:%d", a.host, a.port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		log.Printf("Dial error: %v", err)
		if ctx.Err() == nil {
			triggerRestart(restartChan)
		}
		return
	}

	a.mu.Lock()
	a.conn = conn
	a.isRunning = true
	a.stopChan = make(chan struct{})
	a.mu.Unlock()

	a.ticker = time.NewTicker(a.frequency)

	go func() {
		defer func() {
			conn.Close()
			a.mu.Lock()
			wasRunning := a.isRunning
			a.isRunning = false
			a.mu.Unlock()
			if wasRunning && ctx.Err() == nil {
				triggerRestart(restartChan)
			}
		}()
		for {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				return
			}
			data := buffer[:n]
			parsed, err := Parse(data, a.deviceType)
			if err != nil {
				continue
			}
			a.messageHandler(data, a.id, a.deviceType, parsed)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			conn.Close()
			return
		case <-a.stopChan:
			conn.Close()
			return
		case <-a.ticker.C:
			request, err := GetRequest(a.deviceType, a.deviceAddr)
			if err != nil {
				log.Printf("GetRequest error: %v", err)
				continue
			}
			a.mu.Lock()
			running := a.isRunning && a.conn != nil
			a.mu.Unlock()
			if !running {
				return
			}
			_, err = conn.Write(request)
			if err != nil {
				log.Printf("Write error: %v", err)
				conn.Close()
				a.mu.Lock()
				a.isRunning = false
				a.mu.Unlock()
				if ctx.Err() == nil {
					triggerRestart(restartChan)
				}
				return
			}
		}
	}
}

func (a *Active) Start(ctx context.Context) {
	restartChan := make(chan struct{}, 1)
	go a.startIt(ctx, restartChan)
	go a.startActiveDemon(ctx, restartChan)
}

func (a *Active) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.isRunning {
		close(a.stopChan)
		if a.ticker != nil {
			a.ticker.Stop()
		}
		a.isRunning = false
	}
}

func (a *Active) ID() int {
	return a.id
}

func (a *Active) Info() string {
	return fmt.Sprintf("id:%d host:%s port:%d address:%s frequency:%v, isRunning:%t, conn:%+v",
		a.id, a.host, a.port, a.deviceAddr, a.frequency, a.isRunning, a.conn)
}
