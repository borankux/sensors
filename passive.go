package sensors

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

type Passive struct {
	BaseSensor
	restartTimeout time.Duration
}

func NewPassive(id int, projectId int, host string, address string, port int, deviceType int, restartTimeout time.Duration) *Passive {
	return &Passive{
		BaseSensor: BaseSensor{
			id:         id,
			host:       host,
			projectId:  projectId,
			port:       port,
			deviceAddr: address,
			stopChan:   make(chan struct{}),
			isRunning:  false,
			deviceType: deviceType,
		},
		restartTimeout: restartTimeout,
	}
}

func (p *Passive) startPassiveDemon(ctx context.Context, restartChan chan struct{}) {
	for {
		select {
		case <-p.stopChan:
			return
		case <-ctx.Done():
			return
		case <-restartChan:
			time.Sleep(p.restartTimeout)
			go p.startIt(ctx, restartChan)
		}
	}
}

func (p *Passive) startIt(ctx context.Context, restartChan chan struct{}) {
	go func() {
		address := fmt.Sprintf("%s:%d", p.host, p.port)
		conn, err := net.DialTimeout("tcp", address, time.Second*5)
		if err != nil {
			if ctx.Err() == nil {
				triggerRestart(restartChan)
			}
			return
		}

		defer func() {
			conn.Close()
			p.mu.Lock()
			wasRunning := p.isRunning
			p.isRunning = false
			p.mu.Unlock()
			if wasRunning && ctx.Err() == nil {
				triggerRestart(restartChan)
			}
		}()

		p.mu.Lock()
		p.conn = conn
		p.isRunning = true
		p.stopChan = make(chan struct{})
		p.mu.Unlock()

		for {
			select {
			case <-ctx.Done():
				conn.Close()
				return
			case <-p.stopChan:
				conn.Close()
				return
			default:
				if p.isRunning {
					buffer := make([]byte, 1024)
					n, err := conn.Read(buffer)
					if err != nil {
						return
					}
					data := buffer[:n]
					parsed, err := Parse(data, p.deviceType)
					if err != nil {
						continue
					}
					p.messageHandler(data, p.id, p.deviceType, parsed)
				}
			}
		}
	}()
}

func (p *Passive) Start(ctx context.Context) {
	restartChan := make(chan struct{}, 1)
	go p.startIt(ctx, restartChan)
	go p.startPassiveDemon(ctx, restartChan)
}

func (p *Passive) Stop() {
	fmt.Printf("passive: stopping sensor %d\n", p.id)
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isRunning {
		fmt.Printf("passive: stopping sensor %d\n", p.id)
		close(p.stopChan)
	}
	p.isRunning = false
	log.Printf("passive sensor shutting down, sensor id:%d\n", p.id)
}
func (p *Passive) ID() int {
	return p.id
}
func (p *Passive) Info() string {
	return fmt.Sprintf("id:%d host:%s port:%d address:%s isRunning:%t", p.id, p.host, p.port, p.deviceAddr, p.isRunning)
}
