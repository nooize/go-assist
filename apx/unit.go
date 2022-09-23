package apx

import "context"

/*
 Simple Dependency Injection
 for micro service


 web.After(db)
 web.After(log)

*/

type ApxUnit interface {
	// Start tries to perform the initial initialization of the service, the logic of the function must make sure
	// that all created connections to remote services are in working order and are pinging. Otherwise, the
	// application will need additional error handling.
	Start(ctx context.Context) error
	// IsReady will be called by the service controller at regular intervals, it is important that a response with
	// any error will be regarded as an unrecoverable state of the service and will lead to an emergency stop of
	// the application. If the service is not critical for the application, like a memcached, then try to implement
	// the logic of self-diagnosis and service recovery inside Ping, and return the nil as a response even if the
	// recovery failed.
	IsReady(ctx context.Context) error
	// Close will be executed when the service controller receives a stop command. Normally, this happens after the
	// main thread of the application has already finished. That is, no more requests from the outside are expected.
	Stop() error
}
