package di

import (
	"log"
	"sync"
)

/*
 Simple Dependency Injection
 for micro service


 web.After(db)
 web.After(log)

*/

type unitState struct {
	Error error
}

type unitCreator func() *Unit

type unitHandler struct {
	Start func(chan *unitState)
	Stop  func(chan *unitState)
}

type Unit struct {
	lock      sync.Mutex
	instance  *unitHandler
	state     *unitState
	dependers []*Unit
	after     []*Unit
}

func (m *Unit) After(d *Unit) *Unit {
	if func() bool {
		for _, p := range d.dependers {
			if p == m {
				return false
			}
		}
		return true
	}() {
		d.dependers = append(d.dependers, m)
	}
	if func() bool {
		for _, p := range m.after {
			if p == d {
				return false
			}
		}
		return true
	}() {
		m.after = append(m.after, d)
	}
	return m
}

func (m *Unit) run() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.state != nil {
		return
	}

	for _, u := range m.after {
		u.run()
	}

	stateChannel := make(chan *unitState)
	go m.instance.Start(stateChannel)
	m.state = <-stateChannel
	if m.state.Error != nil {
		log.Printf("run error %v", m.state.Error)
		// TODO log error
	}

}

func (m *Unit) stop() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.state == nil {
		return
	}

	stateChannel := make(chan *unitState)
	go m.instance.Stop(stateChannel)
	m.state = <-stateChannel
	if m.state.Error != nil {
		// TODO log error
	}

	for _, u := range m.after {
		u.stop()
	}

}

func NewUnit(i *unitHandler) *Unit {
	return &Unit{
		instance: i,
	}
}
