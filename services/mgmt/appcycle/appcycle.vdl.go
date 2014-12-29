// This file was auto-generated by the veyron vdl tool.
// Source: appcycle.vdl

// Package appcycle supports managing the application process.
package appcycle

import (
	// The non-user imports are prefixed with "__" to prevent collisions.
	__io "io"
	__veyron2 "v.io/core/veyron2"
	__context "v.io/core/veyron2/context"
	__ipc "v.io/core/veyron2/ipc"
	__vdl "v.io/core/veyron2/vdl"
	__vdlutil "v.io/core/veyron2/vdl/vdlutil"
	__wiretype "v.io/core/veyron2/wiretype"
)

// TODO(toddw): Remove this line once the new signature support is done.
// It corrects a bug where __wiretype is unused in VDL pacakges where only
// bootstrap types are used on interfaces.
const _ = __wiretype.TypeIDInvalid

// Task is streamed by Stop to provide the client with a sense of the progress
// of the shutdown.
// The meaning of Progress and Goal are up to the developer (the server provides
// the framework with values for these).  The recommended meanings are:
// - Progress: how far along the shutdown sequence the server is.  This should
//   be a monotonically increasing number.
// - Goal: when Progress reaches this value, the shutdown is expected to
//   complete.  This should not change during a stream, but could change if
//   e.g. new shutdown tasks are triggered that were not forseen at the outset
//   of the shutdown.
type Task struct {
	Progress int32
	Goal     int32
}

func (Task) __VDLReflect(struct {
	Name string "v.io/core/veyron2/services/mgmt/appcycle.Task"
}) {
}

func init() {
	__vdl.Register(Task{})
}

// AppCycleClientMethods is the client interface
// containing AppCycle methods.
//
// AppCycle interfaces with the process running a veyron runtime.
type AppCycleClientMethods interface {
	// Stop initiates shutdown of the server.  It streams back periodic
	// updates to give the client an idea of how the shutdown is
	// progressing.
	Stop(*__context.T, ...__ipc.CallOpt) (AppCycleStopCall, error)
	// ForceStop tells the server to shut down right away.  It can be issued
	// while a Stop is outstanding if for example the client does not want
	// to wait any longer.
	ForceStop(*__context.T, ...__ipc.CallOpt) error
}

// AppCycleClientStub adds universal methods to AppCycleClientMethods.
type AppCycleClientStub interface {
	AppCycleClientMethods
	__ipc.UniversalServiceMethods
}

// AppCycleClient returns a client stub for AppCycle.
func AppCycleClient(name string, opts ...__ipc.BindOpt) AppCycleClientStub {
	var client __ipc.Client
	for _, opt := range opts {
		if clientOpt, ok := opt.(__ipc.Client); ok {
			client = clientOpt
		}
	}
	return implAppCycleClientStub{name, client}
}

type implAppCycleClientStub struct {
	name   string
	client __ipc.Client
}

func (c implAppCycleClientStub) c(ctx *__context.T) __ipc.Client {
	if c.client != nil {
		return c.client
	}
	return __veyron2.RuntimeFromContext(ctx).Client()
}

func (c implAppCycleClientStub) Stop(ctx *__context.T, opts ...__ipc.CallOpt) (ocall AppCycleStopCall, err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "Stop", nil, opts...); err != nil {
		return
	}
	ocall = &implAppCycleStopCall{Call: call}
	return
}

func (c implAppCycleClientStub) ForceStop(ctx *__context.T, opts ...__ipc.CallOpt) (err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "ForceStop", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&err); ierr != nil {
		err = ierr
	}
	return
}

func (c implAppCycleClientStub) Signature(ctx *__context.T, opts ...__ipc.CallOpt) (o0 __ipc.ServiceSignature, err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "Signature", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&o0, &err); ierr != nil {
		err = ierr
	}
	return
}

// AppCycleStopClientStream is the client stream for AppCycle.Stop.
type AppCycleStopClientStream interface {
	// RecvStream returns the receiver side of the AppCycle.Stop client stream.
	RecvStream() interface {
		// Advance stages an item so that it may be retrieved via Value.  Returns
		// true iff there is an item to retrieve.  Advance must be called before
		// Value is called.  May block if an item is not available.
		Advance() bool
		// Value returns the item that was staged by Advance.  May panic if Advance
		// returned false or was not called.  Never blocks.
		Value() Task
		// Err returns any error encountered by Advance.  Never blocks.
		Err() error
	}
}

// AppCycleStopCall represents the call returned from AppCycle.Stop.
type AppCycleStopCall interface {
	AppCycleStopClientStream
	// Finish blocks until the server is done, and returns the positional return
	// values for call.
	//
	// Finish returns immediately if the call has been canceled; depending on the
	// timing the output could either be an error signaling cancelation, or the
	// valid positional return values from the server.
	//
	// Calling Finish is mandatory for releasing stream resources, unless the call
	// has been canceled or any of the other methods return an error.  Finish should
	// be called at most once.
	Finish() error
}

type implAppCycleStopCall struct {
	__ipc.Call
	valRecv Task
	errRecv error
}

func (c *implAppCycleStopCall) RecvStream() interface {
	Advance() bool
	Value() Task
	Err() error
} {
	return implAppCycleStopCallRecv{c}
}

type implAppCycleStopCallRecv struct {
	c *implAppCycleStopCall
}

func (c implAppCycleStopCallRecv) Advance() bool {
	c.c.valRecv = Task{}
	c.c.errRecv = c.c.Recv(&c.c.valRecv)
	return c.c.errRecv == nil
}
func (c implAppCycleStopCallRecv) Value() Task {
	return c.c.valRecv
}
func (c implAppCycleStopCallRecv) Err() error {
	if c.c.errRecv == __io.EOF {
		return nil
	}
	return c.c.errRecv
}
func (c *implAppCycleStopCall) Finish() (err error) {
	if ierr := c.Call.Finish(&err); ierr != nil {
		err = ierr
	}
	return
}

// AppCycleServerMethods is the interface a server writer
// implements for AppCycle.
//
// AppCycle interfaces with the process running a veyron runtime.
type AppCycleServerMethods interface {
	// Stop initiates shutdown of the server.  It streams back periodic
	// updates to give the client an idea of how the shutdown is
	// progressing.
	Stop(AppCycleStopContext) error
	// ForceStop tells the server to shut down right away.  It can be issued
	// while a Stop is outstanding if for example the client does not want
	// to wait any longer.
	ForceStop(__ipc.ServerContext) error
}

// AppCycleServerStubMethods is the server interface containing
// AppCycle methods, as expected by ipc.Server.
// The only difference between this interface and AppCycleServerMethods
// is the streaming methods.
type AppCycleServerStubMethods interface {
	// Stop initiates shutdown of the server.  It streams back periodic
	// updates to give the client an idea of how the shutdown is
	// progressing.
	Stop(*AppCycleStopContextStub) error
	// ForceStop tells the server to shut down right away.  It can be issued
	// while a Stop is outstanding if for example the client does not want
	// to wait any longer.
	ForceStop(__ipc.ServerContext) error
}

// AppCycleServerStub adds universal methods to AppCycleServerStubMethods.
type AppCycleServerStub interface {
	AppCycleServerStubMethods
	// Describe the AppCycle interfaces.
	Describe__() []__ipc.InterfaceDesc
	// Signature will be replaced with Describe__.
	Signature(ctx __ipc.ServerContext) (__ipc.ServiceSignature, error)
}

// AppCycleServer returns a server stub for AppCycle.
// It converts an implementation of AppCycleServerMethods into
// an object that may be used by ipc.Server.
func AppCycleServer(impl AppCycleServerMethods) AppCycleServerStub {
	stub := implAppCycleServerStub{
		impl: impl,
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := __ipc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := __ipc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implAppCycleServerStub struct {
	impl AppCycleServerMethods
	gs   *__ipc.GlobState
}

func (s implAppCycleServerStub) Stop(ctx *AppCycleStopContextStub) error {
	return s.impl.Stop(ctx)
}

func (s implAppCycleServerStub) ForceStop(ctx __ipc.ServerContext) error {
	return s.impl.ForceStop(ctx)
}

func (s implAppCycleServerStub) Globber() *__ipc.GlobState {
	return s.gs
}

func (s implAppCycleServerStub) Describe__() []__ipc.InterfaceDesc {
	return []__ipc.InterfaceDesc{AppCycleDesc}
}

// AppCycleDesc describes the AppCycle interface.
var AppCycleDesc __ipc.InterfaceDesc = descAppCycle

// descAppCycle hides the desc to keep godoc clean.
var descAppCycle = __ipc.InterfaceDesc{
	Name:    "AppCycle",
	PkgPath: "v.io/core/veyron2/services/mgmt/appcycle",
	Doc:     "// AppCycle interfaces with the process running a veyron runtime.",
	Methods: []__ipc.MethodDesc{
		{
			Name: "Stop",
			Doc:  "// Stop initiates shutdown of the server.  It streams back periodic\n// updates to give the client an idea of how the shutdown is\n// progressing.",
			OutArgs: []__ipc.ArgDesc{
				{"", ``}, // error
			},
		},
		{
			Name: "ForceStop",
			Doc:  "// ForceStop tells the server to shut down right away.  It can be issued\n// while a Stop is outstanding if for example the client does not want\n// to wait any longer.",
			OutArgs: []__ipc.ArgDesc{
				{"", ``}, // error
			},
		},
	},
}

func (s implAppCycleServerStub) Signature(ctx __ipc.ServerContext) (__ipc.ServiceSignature, error) {
	// TODO(toddw): Replace with new Describe__ implementation.
	result := __ipc.ServiceSignature{Methods: make(map[string]__ipc.MethodSignature)}
	result.Methods["ForceStop"] = __ipc.MethodSignature{
		InArgs: []__ipc.MethodArgument{},
		OutArgs: []__ipc.MethodArgument{
			{Name: "", Type: 65},
		},
	}
	result.Methods["Stop"] = __ipc.MethodSignature{
		InArgs: []__ipc.MethodArgument{},
		OutArgs: []__ipc.MethodArgument{
			{Name: "", Type: 65},
		},

		OutStream: 66,
	}

	result.TypeDefs = []__vdlutil.Any{
		__wiretype.NamedPrimitiveType{Type: 0x1, Name: "error", Tags: []string(nil)}, __wiretype.StructType{
			[]__wiretype.FieldType{
				__wiretype.FieldType{Type: 0x24, Name: "Progress"},
				__wiretype.FieldType{Type: 0x24, Name: "Goal"},
			},
			"v.io/core/veyron2/services/mgmt/appcycle.Task", []string(nil)},
	}

	return result, nil
}

// AppCycleStopServerStream is the server stream for AppCycle.Stop.
type AppCycleStopServerStream interface {
	// SendStream returns the send side of the AppCycle.Stop server stream.
	SendStream() interface {
		// Send places the item onto the output stream.  Returns errors encountered
		// while sending.  Blocks if there is no buffer space; will unblock when
		// buffer space is available.
		Send(item Task) error
	}
}

// AppCycleStopContext represents the context passed to AppCycle.Stop.
type AppCycleStopContext interface {
	__ipc.ServerContext
	AppCycleStopServerStream
}

// AppCycleStopContextStub is a wrapper that converts ipc.ServerCall into
// a typesafe stub that implements AppCycleStopContext.
type AppCycleStopContextStub struct {
	__ipc.ServerCall
}

// Init initializes AppCycleStopContextStub from ipc.ServerCall.
func (s *AppCycleStopContextStub) Init(call __ipc.ServerCall) {
	s.ServerCall = call
}

// SendStream returns the send side of the AppCycle.Stop server stream.
func (s *AppCycleStopContextStub) SendStream() interface {
	Send(item Task) error
} {
	return implAppCycleStopContextSend{s}
}

type implAppCycleStopContextSend struct {
	s *AppCycleStopContextStub
}

func (s implAppCycleStopContextSend) Send(item Task) error {
	return s.s.Send(item)
}
