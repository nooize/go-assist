package di

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
 Tiny Dependency Injection
 for micro service


 web.After(db)
 web.After(log)

*/

type Service struct {
	Unit
	Name string
	quit chan os.Signal
	shutdown chan error
}

func (s *Service) Start() {
	s.Unit.run()
}

func (s *Service) After(d *Unit) *Service {
	s.Unit.After(d)
	return s
}

func (s *Service) run() error {

	signal.Notify(s.quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGHUP, syscall.SIGQUIT)

	log.Printf(s.Name + " is up.")
	// logs.Logs().Info().Msgf("%v is up", core.ProductName)

	defer func() {
		log.Printf(s.Name + " is stop.")
		// logs.Logs().Info().Msgf("%v stop", core.ProductName)
	}()

	_, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	for {
		select {
		case sig := <- s.quit:
			switch sig {
			case syscall.SIGINT:
				// s.Shutdown()
				os.Exit(0)
			case syscall.SIGTERM:
				// s.Shutdown()
				os.Exit(143)
			case syscall.SIGUSR1:
				// File log re-open for rotating file logs.
				// s.log.ReopenLogFile()
			case syscall.SIGHUP:
				//s.mu.Lock()
				//ns := s.natsServer
				//s.mu.Unlock()
				//if ns != nil {
				//	if err := ns.Reload(); err != nil {
				//		s.log.Errorf("Reload: %v", err)
				//	}
				//} else {
				//	s.log.Warnf("Reload supported only for embedded NATS Server's configuration")
				//}
			}
		case <- s.shutdown:
			return nil
		}
	}

	os.Exit(0)

	return nil
}

func NewService(name string) *Service {
	s := Service{
		quit: make(chan os.Signal, 2),
		shutdown: make(chan error, 2),
		Name: name,
	}
	s.Unit = *NewUnit(&unitHandler{
		Start: func(state chan *UnitState) {
			err := s.run()
			state <- &UnitState{ Error: err }
		},
		Stop: func(state chan *UnitState) {
			os.Exit(0)
		},
	})
	return &s
}
