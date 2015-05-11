// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package v23 defines the runtime interface of Vanadium, and its subdirectories
// define the entire Vanadium public API.
//
// Once Vanadium reaches version 1.0 these public APIs will be stable over an
// extended period.  Changes to APIs will be managed to ensure backwards
// compatibility, using the same policy as http://golang.org/doc/go1compat.
//
// This is version 0.1 - we will do our best to maintain backwards
// compatibility, but there's no guarantee until version 1.0.
//
// For more details about the Vanadium project, please visit https://v.io.
package v23

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"

	"v.io/v23/context"
	"v.io/v23/namespace"
	"v.io/v23/naming"
	"v.io/v23/rpc"
	"v.io/v23/security"
)

const (
	// LocalStop is the message received on AppCycle.WaitForStop when the stop was
	// initiated by the process itself.
	LocalStop = "localstop"
	// RemoteStop is the message received on AppCycle.WaitForStop when the stop was
	// initiated via an RPC call (AppCycle.Stop).
	RemoteStop            = "remotestop"
	UnhandledStopExitCode = 1
	ForceStopExitCode     = 1
)

// Task is streamed to channels registered using TrackTask to provide a sense of
// the progress of the application's shutdown sequence.  For a description of
// the fields, see the Task struct in the v23/services/appcycle package, which
// it mirrors.
type Task struct {
	Progress, Goal int32
}

// AppCycle is the interface for managing the shutdown of a runtime
// remotely and locally. An appropriate instance of this is provided by
// the RuntimeFactory to the runtime implementation which in turn arranges to
// serve it on an appropriate network address.
type AppCycle interface {
	// Stop causes all the channels returned by WaitForStop to return the
	// LocalStop message, to give the application a chance to shut down.
	// Stop does not block.  If any of the channels are not receiving,
	// the message is not sent on them.
	// If WaitForStop had never been called, Stop acts like ForceStop.
	Stop()

	// ForceStop causes the application to exit immediately with an error
	// code.
	ForceStop()

	// WaitForStop takes in a channel on which a stop event will be
	// conveyed.  The stop event is represented by a string identifying the
	// source of the event.  For example, when Stop is called locally, the
	// LocalStop message will be received on the channel.  If the channel is
	// not being received on, or is full, no message is sent on it.
	//
	// The channel is assumed to remain open while messages could be sent on
	// it.  The channel will be automatically closed during the call to
	// Cleanup.
	WaitForStop(chan<- string)

	// AdvanceGoal extends the goal value in the shutdown task tracker.
	// Non-positive delta is ignored.
	AdvanceGoal(delta int32)
	// AdvanceProgress advances the progress value in the shutdown task
	// tracker.  Non-positive delta is ignored.
	AdvanceProgress(delta int32)
	// TrackTask registers a channel to receive task updates (a Task will be
	// sent on the channel if either the goal or progress values of the
	// task have changed).  If the channel is not being received on, or is
	// full, no Task is sent on it.
	//
	// The channel is assumed to remain open while Tasks could be sent on
	// it.
	TrackTask(chan<- Task)

	// Remote returns an object to serve the remotely accessible AppCycle
	// interface (as defined in v23/services/appcycle)
	Remote() interface{}
}

// Runtime is the interface that concrete Vanadium implementations must
// implement.  It will not be used directly by application builders.
// They will instead use the package level functions that mirror these
// factories.
type Runtime interface {

	// Init is a chance to initialize state in the runtime implementation
	// after the runtime has been registered in the v23 package.
	// Code that runs in this routine, unlike the code in the runtimes
	// constructor, can use the v23.Get/With methods.
	Init(ctx *context.T) error

	// NewEndpoint returns an Endpoint by parsing the supplied endpoint
	// string as per the format described above. It can be used to test
	// a string to see if it's in valid endpoint format.
	//
	// NewEndpoint will accept strings both in the @ format described
	// above and in internet host:port format.
	//
	// All implementations of NewEndpoint should provide appropriate
	// defaults for any endpoint subfields not explicitly provided as
	// follows:
	// - a missing protocol will default to a protocol appropriate for the
	//   implementation hosting NewEndpoint
	// - a missing host:port will default to :0 - i.e. any port on all
	//   interfaces
	// - a missing routing id should default to the null routing id
	// - a missing codec version should default to AnyCodec
	// - a missing RPC version should default to the highest version
	//   supported by the runtime implementation hosting NewEndpoint
	NewEndpoint(ep string) (naming.Endpoint, error)

	// NewServer creates a new Server instance.
	//
	// It accepts at least the following options:
	// ServesMountTable and ServerBlessings.
	NewServer(ctx *context.T, opts ...rpc.ServerOpt) (rpc.Server, error)

	// WithNewStreamManager creates a new StreamManager instance and context
	// with that StreamManager attached.
	WithNewStreamManager(ctx *context.T) (*context.T, error)

	// WithPrincipal attaches 'principal' to the returned context.
	WithPrincipal(ctx *context.T, principal security.Principal) (*context.T, error)

	// GetPrincipal returns the Principal in 'ctx'.
	GetPrincipal(ctx *context.T) security.Principal

	// WithNewClient creates a new Client instance and attaches it to a
	// new context.
	WithNewClient(ctx *context.T, opts ...rpc.ClientOpt) (*context.T, rpc.Client, error)

	// GetClient returns the Client in 'ctx'.
	GetClient(ctx *context.T) rpc.Client

	// WithNewNamespace creates a new Namespace instance and attaches it to the
	// returned context.
	WithNewNamespace(ctx *context.T, roots ...string) (*context.T, namespace.T, error)

	// GetNamespace returns the Namespace in 'ctx'.
	GetNamespace(ctx *context.T) namespace.T

	// GetAppCycle returns the AppCycle in 'ctx'.
	GetAppCycle(ctx *context.T) AppCycle

	// GetListenSpec returns the ListenSpec in 'ctx'.
	GetListenSpec(ctx *context.T) rpc.ListenSpec

	// WithBackgroundContext creates a new context derived from 'ctx'
	// with the given context set as the background context.
	WithBackgroundContext(ctx *context.T) *context.T

	// GetBackgroundContext returns a background context. This context can be used
	// for general background activities.
	GetBackgroundContext(ctx *context.T) *context.T

	// WithReservedNameDispatcher returns a context that uses the
	// provided dispatcher to control access to the framework managed
	// portion of the namespace.
	WithReservedNameDispatcher(ctx *context.T, d rpc.Dispatcher) *context.T

	// GetReservedNameDispatcher returns the dispatcher used for
	// reserved names.
	GetReservedNameDispatcher(ctx *context.T) rpc.Dispatcher
}

// NewEndpoint returns an Endpoint by parsing the supplied endpoint
// string as per the format described above. It can be used to test
// a string to see if it's in valid endpoint format.
//
// NewEndpoint will accept strings both in the @ format described
// above and in internet host:port format.
//
// All implementations of NewEndpoint should provide appropriate
// defaults for any endpoint subfields not explicitly provided as
// follows:
// - a missing protocol will default to a protocol appropriate for the
//   implementation hosting NewEndpoint
// - a missing host:port will default to :0 - i.e. any port on all
//   interfaces
// - a missing routing id should default to the null routing id
// - a missing codec version should default to AnyCodec
// - a missing RPC version should default to the highest version
//   supported by the runtime implementation hosting NewEndpoint
func NewEndpoint(ep string) (naming.Endpoint, error) {
	return initState.currentRuntime().NewEndpoint(ep)
}

// NewServer creates a new Server instance.
//
// It accepts at least the following options:
// ServesMountTable and ServerBlessings.
//
// ServerBlessings defaults to v23.GetPrincipal(ctx).BlessingStore().Default().
// These Blessings are the Server's Blessings for its lifetime.
func NewServer(ctx *context.T, opts ...rpc.ServerOpt) (rpc.Server, error) {
	return initState.currentRuntime().NewServer(ctx, opts...)
}

// WithNewStreamManager creates a new StreamManager instance and context
// with that StreamManager attached.
func WithNewStreamManager(ctx *context.T) (*context.T, error) {
	return initState.currentRuntime().WithNewStreamManager(ctx)
}

// WithPrincipal attaches 'principal' to the returned context.
func WithPrincipal(ctx *context.T, principal security.Principal) (*context.T, error) {
	return initState.currentRuntime().WithPrincipal(ctx, principal)
}

// GetPrincipal returns the Principal in 'ctx'.
func GetPrincipal(ctx *context.T) security.Principal {
	return initState.currentRuntime().GetPrincipal(ctx)
}

// WithNewClient creates a new Client instance and attaches it to a
// new context.
func WithNewClient(ctx *context.T, opts ...rpc.ClientOpt) (*context.T, rpc.Client, error) {
	return initState.currentRuntime().WithNewClient(ctx, opts...)
}

// GetClient returns the Client in 'ctx'.
func GetClient(ctx *context.T) rpc.Client {
	return initState.currentRuntime().GetClient(ctx)
}

// WithNewNamespace creates a new Namespace instance and attaches it to the
// returned context.
func WithNewNamespace(ctx *context.T, roots ...string) (*context.T, namespace.T, error) {
	return initState.currentRuntime().WithNewNamespace(ctx, roots...)
}

// GetNamespace returns the Namespace in 'ctx'.
func GetNamespace(ctx *context.T) namespace.T {
	return initState.currentRuntime().GetNamespace(ctx)
}

// GetAppCycle returns the AppCycle in 'ctx'.
func GetAppCycle(ctx *context.T) AppCycle {
	return initState.currentRuntime().GetAppCycle(ctx)
}

// GetListenSpec returns the ListenSpec in 'ctx'.
func GetListenSpec(ctx *context.T) rpc.ListenSpec {
	return initState.currentRuntime().GetListenSpec(ctx)
}

// WithBackgroundContext creates a new context derived from 'ctx'
// with the given context set as the background context.
func WithBackgroundContext(ctx *context.T) *context.T {
	return initState.runtime.WithBackgroundContext(ctx)
}

// GetBackgroundContext returns a background context. This context can be used
// for general background activities.
func GetBackgroundContext(ctx *context.T) *context.T {
	return initState.runtime.GetBackgroundContext(ctx)
}

// WithReservedNameDispatcher returns a context that uses the
// provided dispatcher to handle reserved names in particular
// __debug.
func WithReservedNameDispatcher(ctx *context.T, d rpc.Dispatcher) *context.T {
	return initState.currentRuntime().WithReservedNameDispatcher(ctx, d)
}

// GetReservedNameDispatcher returns the dispatcher used for
// reserved names.
func GetReservedNameDispatcher(ctx *context.T) rpc.Dispatcher {
	return initState.currentRuntime().GetReservedNameDispatcher(ctx)
}

var initState = &initStateData{}

type initStateData struct {
	mu                  sync.RWMutex
	runtime             Runtime
	runtimeStack        string
	runtimeFactory      RuntimeFactory
	runtimeFactoryStack string
}

func (i *initStateData) currentRuntime() Runtime {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if i.runtimeStack == "" {
		panic(`Calling v23 method before initializing the runtime with Init().
You should call Init from your main or test function before calling
other v23 operations.`)
	}
	if i.runtime == nil {
		panic(`Calling v23 method during runtime initialization.  You cannot
call v23 methods until after the runtime has been constructed.  You may
be able to move the offending caller to the Runtime.Init() method of your
runtime implementation.`)
	}

	return i.runtime
}

// A RuntimeFactory represents the combination of hardware, operating system,
// compiler and libraries available to the application. The RuntimeFactory
// creates a runtime implementation with the required hardware, operating system
// and library specific dependencies included.
//
// The implementations of the RuntimeFactory are intended to capture all of
// the dependencies implied by that RuntimeFactory. For example, if a RuntimeFactory requires
// a particular hardware specific library (say Bluetooth support), then the
// implementation of the RuntimeFactory should include that dependency and
// the resulting runtime instance; the package implementing
// the RuntimeFactory should expose the additional APIs needed to use the
// functionality.
//
// RuntimeFactories range from the generic to the very specific (e.g. "linux" or
// "my-sprinkler-controller-v2". Applications should, in general, use
// as generic a RuntimeFactory as possbile.
//
// RuntimeFactories are registered using v23.RegisterRuntimeFactory and subsequent
// registrations will panic. Packages that implement RuntimeFactories will typically
// call RegisterRuntimeFactory in their init functions so importing a RuntimeFactory will
// be sufficient to register it. Only one RuntimeFactory can be registered in any
// program, and subsequent registrations will panic.  Typically a program's main
// package will be the only place to import a RuntimeFactory.
//
// This scheme allows applications to use a pre-supplied RuntimeFactory as well
// as for developers to create their own RuntimeFactories (to represent their
// hardware and software system).
//
// At a minimum a RuntimeFactory must do the following:
//   - Initialize a Runtime implementation (providing the flags to it)
//   - Return a Runtime implemenation, initial context, Shutdown func.
//
// See the v.io/x/ref/runtime/factories package for a complete description of the
// precanned RuntimeFactories and how to use them.
type RuntimeFactory func(ctx *context.T) (Runtime, *context.T, Shutdown, error)

// RegisterRuntimeFactory register the specified RuntimeFactory.
// It must be called before v23.Init; typically it will be called by an init
// function. It will panic if called more than once.
func RegisterRuntimeFactory(f RuntimeFactory) {
	// Skip 3 frames: runtime.Callers, getStack, RegisterRuntimeFactory.
	stack := getStack(3)
	initState.mu.Lock()
	defer initState.mu.Unlock()
	if initState.runtimeFactory != nil {
		format := `A RuntimeFactory has already been registered.
This is most likely because a library package is
importing a RuntimeFactory.  Look for imports of the form
'v.io/x/ref/runtime/factories/...' and remove them.  RuntimeFactories should
only be imported in your main package.  Previous registration was from:
%s
Current registration is from:
%s
`
		panic(fmt.Sprintf(format, initState.runtimeFactoryStack, stack))
	}
	initState.runtimeFactory = f
	initState.runtimeFactoryStack = stack
}

type Shutdown func()

func getStack(skip int) string {
	var buf bytes.Buffer
	stack := make([]uintptr, 16)
	stack = stack[:runtime.Callers(skip, stack)]
	for _, pc := range stack {
		fnc := runtime.FuncForPC(pc)
		file, line := fnc.FileLine(pc)
		fmt.Fprintf(&buf, "%s:%d: %s\n", file, line, fnc.Name())
	}
	return buf.String()
}

// Init should be called once for each vanadium executable, providing
// the setup of the initial context.T and a Shutdown function that can
// be used to clean up the runtime.  We allow calling Init multiple
// times (useful in tests), but only as long as you call the Shutdown
// returned previously before calling Init the second time.
func Init() (*context.T, Shutdown) {
	initState.mu.Lock()
	runtimeFactory := initState.runtimeFactory
	if initState.runtimeFactory == nil {
		initState.mu.Unlock()
		panic("No RuntimeFactory has been registered nor specified. This is most" +
			" likely because your main package has not imported a RuntimeFactory")
	}

	// Skip 3 stack frames: runtime.Callers, getStack, Init
	stack := getStack(3)
	if initState.runtimeStack != "" {
		initState.mu.Unlock()
		format := `A runtime has already been initialized."
The previous initialization was from:
%s
This registration is from:
%s
`
		panic(fmt.Sprintf(format, initState.runtimeStack, stack))
	}
	initState.runtimeStack = stack
	initState.mu.Unlock()

	rootctx, rootcancel := context.RootContext()
	// Note we derive a second cancelable context here beyond the
	// rootctx.  This allows us to do shutdown in two steps.  First
	// we cancel this initial context to trigger cleanup of all
	// servers, clients, stream managers, etc.  Then after everything
	// is shut down we invoke rootcancel.  This allows the cleanup
	// to perform operations that require uncancelled contexts.
	ctx, cancel := context.WithCancel(rootctx)
	rt, ctx, shutdown, err := runtimeFactory(ctx)
	if err != nil {
		cancel()
		rootcancel()
		panic(err)
	}

	initState.mu.Lock()
	initState.runtime = rt
	initState.mu.Unlock()

	vshutdown := func() {
		// Note we call our own cancel here to ensure that the
		// runtime/runtimeFactory implementor has not attached anything to a
		// non-cancellable context.
		cancel()
		shutdown()
		rootcancel()

		initState.mu.Lock()
		initState.runtime = nil
		initState.runtimeStack = ""
		initState.mu.Unlock()
	}

	if err := rt.Init(ctx); err != nil {
		vshutdown()
		panic(err)
	}

	return ctx, vshutdown
}
