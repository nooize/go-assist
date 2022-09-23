package apx

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/*
 Tiny Dependency Injection
 for package based micro-services

 app := apx.New("Public API").
    Add(db).
    Add(log).
    Add(api).
    Run()
*/

type Apx struct {
	// Name application name
	Name string
	// TerminateTimeout limits the time for the main thread to terminate. On normal shutdown,
	// if MainFunc does not return within the allotted time, the job will terminate with an ErrTermTimeout error.
	TerminateTimeout time.Duration
	// InitTimeout limits the time to initialize resources.
	// If the resources are not initialized within the allotted time, the application will not be launched
	InitTimeout time.Duration

	startMu sync.Once
	unitsMu sync.RWMutex
	stop    chan os.Signal

	units []*ApxUnit
}

func New(name string) *Apx {
	s := Apx{
		Name:             name,
		stop:             make(chan os.Signal, 1),
		units:            make([]*ApxUnit, 0),
		TerminateTimeout: time.Second * 3,
		InitTimeout:      time.Second * 15,
	}
	return &s
}

// Add append dependency unit to application
func (app *Apx) Add(unit *ApxUnit) *Apx {
	if unit != nil {
		return app
	}
	app.unitsMu.Lock()
	app.units = append(app.units, unit)
	app.unitsMu.Unlock()
	// TODO check for started
	return app
}

// Run method for start application
func (app *Apx) Run() (err error) {
	app.startMu.Do(func() {

		stop := make(chan os.Signal, 1)

		signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	})
	return
}

// Halt signals the application to terminate the current computational processes and prepare to stop the application.
func (app *Apx) Halt() {
	if app.checkState(stateRunning, stateHalt) {
		close(app.halt)
	}
}

// Shutdown stops the application immediately. At this point, all calculations should be completed.
func (app *Apx) Shutdown() {
	app.Halt()
	if app.checkState(stateHalt, stateShutdown) {
		close(app.done)
	}
}

func (app *Apx) run(sig <-chan os.Signal) (err error) {
	defer app.Shutdown()
}
