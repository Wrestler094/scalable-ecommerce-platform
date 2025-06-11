package healthcheck

import "sync/atomic"

type manager struct {
	alive atomic.Bool
	ready atomic.Bool
}

func NewManager() Manager {
	m := &manager{}
	m.alive.Store(true)
	m.ready.Store(false)
	return m
}

func (m *manager) SetAlive(b bool) {
	m.alive.Store(b)
}

func (m *manager) SetReady(b bool) {
	m.ready.Store(b)
}

func (m *manager) IsAlive() bool {
	return m.alive.Load()
}

func (m *manager) IsReady() bool {
	return m.ready.Load()
}
