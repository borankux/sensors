package sensors

import (
	"context"
	"log"
	"sync"
	"time"
)

type ManagerStatus struct {
	Manager bool
	Devices map[int]bool
}

type ManagerStatusHandler func(ms *ManagerStatus)

type Manager struct {
	isRunning      bool
	sensors        map[int]Sensor
	mu             sync.RWMutex
	handler        MessageHandler
	ctx            context.Context
	onStatusChange ManagerStatusHandler
}

func (m *Manager) SetStatusHandler(handler ManagerStatusHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onStatusChange = handler
}

func (m *Manager) getStatus() *ManagerStatus {
	deviceStatus := make(map[int]bool)
	for _, sensor := range m.sensors {
		deviceStatus[sensor.ID()] = sensor.IsRunning()
	}
	return &ManagerStatus{
		Manager: m.isRunning,
		Devices: deviceStatus,
	}
}

func (m *Manager) AddSensor(sensor Sensor) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s, existed := m.sensors[sensor.ID()]; existed {
		m.RemoveSensor(s.ID())
	}

	sensor.SetHandler(m.handler)
	if m.isRunning {
		sensor.Start(m.ctx)
	}
	m.sensors[sensor.ID()] = sensor
	m.onStatusChange(m.getStatus())
}

func (m *Manager) GetSensor(id int) Sensor {
	return m.sensors[id]
}

func (m *Manager) RemoveSensor(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sensor, exists := m.sensors[id]; exists {
		sensor.Stop()
		delete(m.sensors, id)
	}
	if len(m.sensors) == 0 {
		m.isRunning = false
	}
	m.onStatusChange(m.getStatus())
}

func (m *Manager) StartAllSensors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sensor := range m.sensors {
		sensor.Start(m.ctx)
	}
	m.isRunning = true
	m.onStatusChange(m.getStatus())
}

func (m *Manager) StopAllSensors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sensor := range m.sensors {
		go sensor.Stop()
	}
	m.isRunning = false
	m.onStatusChange(m.getStatus())
}

func (m *Manager) GetInfo() {
	if len(m.sensors) == 0 {
		log.Println("No serial")
	}
	for _, sensor := range m.sensors {
		info := sensor.Info()
		log.Println(info)
	}
}

func NewManager(handler MessageHandler, ctx context.Context) *Manager {
	manager := &Manager{
		sensors: make(map[int]Sensor),
		handler: handler,
		ctx:     ctx,
	}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			manager.onStatusChange(manager.getStatus())
		}
	}()
	return manager
}
