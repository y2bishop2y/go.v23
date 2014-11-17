// This file was auto-generated by the veyron vdl tool.
// Source: base.vdl

// Package base is a simple single-file test of vdl functionality.
package base

import (
	// The non-user imports are prefixed with "__" to prevent collisions.
	__io "io"
	__veyron2 "veyron.io/veyron/veyron2"
	__context "veyron.io/veyron/veyron2/context"
	__ipc "veyron.io/veyron/veyron2/ipc"
	__vdl "veyron.io/veyron/veyron2/vdl"
	__vdlutil "veyron.io/veyron/veyron2/vdl/vdlutil"
	__verror "veyron.io/veyron/veyron2/verror"
	__wiretype "veyron.io/veyron/veyron2/wiretype"
)

// TODO(toddw): Remove this line once the new signature support is done.
// It corrects a bug where __wiretype is unused in VDL pacakges where only
// bootstrap types are used on interfaces.
const _ = __wiretype.TypeIDInvalid

type NamedBool bool

type NamedByte byte

type NamedUint16 uint16

type NamedUint32 uint32

type NamedUint64 uint64

type NamedInt16 int16

type NamedInt32 int32

type NamedInt64 int64

type NamedFloat32 float32

type NamedFloat64 float64

type NamedComplex64 complex64

type NamedComplex128 complex128

type NamedString string

//NamedEnum       enum{A;B;C}
type NamedArray [2]bool

type NamedList []uint32

type NamedSet map[string]struct{}

type NamedMap map[string]float32

type NamedStruct struct {
	A bool
	B string
	C int32
}

type Scalars struct {
	A0  bool
	A1  byte
	A2  uint16
	A3  uint32
	A4  uint64
	A5  int16
	A6  int32
	A7  int64
	A8  float32
	A9  float64
	A10 complex64
	A11 complex128
	A12 string
	A13 error
	A14 __vdlutil.Any
	A15 *__vdl.Type
	B0  NamedBool
	B1  NamedByte
	B2  NamedUint16
	B3  NamedUint32
	B4  NamedUint64
	B5  NamedInt16
	B6  NamedInt32
	B7  NamedInt64
	B8  NamedFloat32
	B9  NamedFloat64
	B10 NamedComplex64
	B11 NamedComplex128
	B12 NamedString
}

// These are all scalars that may be used as map or set keys.
type KeyScalars struct {
	A0  bool
	A1  byte
	A2  uint16
	A3  uint32
	A4  uint64
	A5  int16
	A6  int32
	A7  int64
	A8  float32
	A9  float64
	A10 complex64
	A11 complex128
	A12 string
	A13 error
	B0  NamedBool
	B1  NamedByte
	B2  NamedUint16
	B3  NamedUint32
	B4  NamedUint64
	B5  NamedInt16
	B6  NamedInt32
	B7  NamedInt64
	B8  NamedFloat32
	B9  NamedFloat64
	B10 NamedComplex64
	B11 NamedComplex128
	B12 NamedString
}

type Composites struct {
	A0 Scalars
	A1 [2]Scalars
	A2 []Scalars
	A3 map[KeyScalars]struct{}
	A4 map[string]Scalars
	A5 map[KeyScalars][]map[string]complex128
}

type CompComp struct {
	A0 Composites
	A1 [2]Composites
	A2 []Composites
	A3 map[string]Composites
	A4 map[KeyScalars][]map[string]Composites
}

// NestedArgs is defined before Args; that's allowed in regular Go, and also
// allowed in our vdl files.  The compiler will re-order dependent types to ease
// code generation in other languages.
type NestedArgs struct {
	Args Args
}

// Args will be reordered to show up before NestedArgs in the generated output.
type Args struct {
	A int32
	B int32
}

const Cbool = true

const Cbyte = byte(1)

const Cint32 = int32(2)

const Cint64 = int64(3)

const Cuint32 = uint32(4)

const Cuint64 = uint64(5)

const Cfloat32 = float32(6)

const Cfloat64 = float64(7)

const Ccomplex64 = complex64(8 + 9i)

const Ccomplex128 = complex128(10 + 11i)

const Cstring = "foo"

const True = true

const Foo = "foo"

const Five = int32(5)

const Six = uint64(6)

const SixSquared = uint64(36)

const FiveSquared = int32(25)

const ErrIDFoo = __verror.ID("veyron.io/veyron/veyron2/vdl/testdata/base.ErrIDFoo")

const ErrIDBar = __verror.ID("some/path.ErrIdOther")

// ServiceAClientMethods is the client interface
// containing ServiceA methods.
type ServiceAClientMethods interface {
	MethodA1(__context.T, ...__ipc.CallOpt) error
	MethodA2(ctx __context.T, a int32, b string, opts ...__ipc.CallOpt) (s string, err error)
	MethodA3(ctx __context.T, a int32, opts ...__ipc.CallOpt) (ServiceAMethodA3Call, error)
	MethodA4(ctx __context.T, a int32, opts ...__ipc.CallOpt) (ServiceAMethodA4Call, error)
}

// ServiceAClientStub adds universal methods to ServiceAClientMethods.
type ServiceAClientStub interface {
	ServiceAClientMethods
	__ipc.UniversalServiceMethods
}

// ServiceAClient returns a client stub for ServiceA.
func ServiceAClient(name string, opts ...__ipc.BindOpt) ServiceAClientStub {
	var client __ipc.Client
	for _, opt := range opts {
		if clientOpt, ok := opt.(__ipc.Client); ok {
			client = clientOpt
		}
	}
	return implServiceAClientStub{name, client}
}

type implServiceAClientStub struct {
	name   string
	client __ipc.Client
}

func (c implServiceAClientStub) c(ctx __context.T) __ipc.Client {
	if c.client != nil {
		return c.client
	}
	return __veyron2.RuntimeFromContext(ctx).Client()
}

func (c implServiceAClientStub) MethodA1(ctx __context.T, opts ...__ipc.CallOpt) (err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "MethodA1", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&err); ierr != nil {
		err = ierr
	}
	return
}

func (c implServiceAClientStub) MethodA2(ctx __context.T, i0 int32, i1 string, opts ...__ipc.CallOpt) (o0 string, err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "MethodA2", []interface{}{i0, i1}, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&o0, &err); ierr != nil {
		err = ierr
	}
	return
}

func (c implServiceAClientStub) MethodA3(ctx __context.T, i0 int32, opts ...__ipc.CallOpt) (ocall ServiceAMethodA3Call, err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "MethodA3", []interface{}{i0}, opts...); err != nil {
		return
	}
	ocall = &implServiceAMethodA3Call{Call: call}
	return
}

func (c implServiceAClientStub) MethodA4(ctx __context.T, i0 int32, opts ...__ipc.CallOpt) (ocall ServiceAMethodA4Call, err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "MethodA4", []interface{}{i0}, opts...); err != nil {
		return
	}
	ocall = &implServiceAMethodA4Call{Call: call}
	return
}

func (c implServiceAClientStub) Signature(ctx __context.T, opts ...__ipc.CallOpt) (o0 __ipc.ServiceSignature, err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "Signature", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&o0, &err); ierr != nil {
		err = ierr
	}
	return
}

// ServiceAMethodA3ClientStream is the client stream for ServiceA.MethodA3.
type ServiceAMethodA3ClientStream interface {
	// RecvStream returns the receiver side of the ServiceA.MethodA3 client stream.
	RecvStream() interface {
		// Advance stages an item so that it may be retrieved via Value.  Returns
		// true iff there is an item to retrieve.  Advance must be called before
		// Value is called.  May block if an item is not available.
		Advance() bool
		// Value returns the item that was staged by Advance.  May panic if Advance
		// returned false or was not called.  Never blocks.
		Value() Scalars
		// Err returns any error encountered by Advance.  Never blocks.
		Err() error
	}
}

// ServiceAMethodA3Call represents the call returned from ServiceA.MethodA3.
type ServiceAMethodA3Call interface {
	ServiceAMethodA3ClientStream
	// Finish blocks until the server is done, and returns the positional return
	// values for call.
	//
	// Finish returns immediately if Cancel has been called; depending on the
	// timing the output could either be an error signaling cancelation, or the
	// valid positional return values from the server.
	//
	// Calling Finish is mandatory for releasing stream resources, unless Cancel
	// has been called or any of the other methods return an error.  Finish should
	// be called at most once.
	Finish() (s string, err error)
	// Cancel cancels the RPC, notifying the server to stop processing.  It is
	// safe to call Cancel concurrently with any of the other stream methods.
	// Calling Cancel after Finish has returned is a no-op.
	Cancel()
}

type implServiceAMethodA3Call struct {
	__ipc.Call
	valRecv Scalars
	errRecv error
}

func (c *implServiceAMethodA3Call) RecvStream() interface {
	Advance() bool
	Value() Scalars
	Err() error
} {
	return implServiceAMethodA3CallRecv{c}
}

type implServiceAMethodA3CallRecv struct {
	c *implServiceAMethodA3Call
}

func (c implServiceAMethodA3CallRecv) Advance() bool {
	c.c.valRecv = Scalars{}
	c.c.errRecv = c.c.Recv(&c.c.valRecv)
	return c.c.errRecv == nil
}
func (c implServiceAMethodA3CallRecv) Value() Scalars {
	return c.c.valRecv
}
func (c implServiceAMethodA3CallRecv) Err() error {
	if c.c.errRecv == __io.EOF {
		return nil
	}
	return c.c.errRecv
}
func (c *implServiceAMethodA3Call) Finish() (o0 string, err error) {
	if ierr := c.Call.Finish(&o0, &err); ierr != nil {
		err = ierr
	}
	return
}

// ServiceAMethodA4ClientStream is the client stream for ServiceA.MethodA4.
type ServiceAMethodA4ClientStream interface {
	// RecvStream returns the receiver side of the ServiceA.MethodA4 client stream.
	RecvStream() interface {
		// Advance stages an item so that it may be retrieved via Value.  Returns
		// true iff there is an item to retrieve.  Advance must be called before
		// Value is called.  May block if an item is not available.
		Advance() bool
		// Value returns the item that was staged by Advance.  May panic if Advance
		// returned false or was not called.  Never blocks.
		Value() string
		// Err returns any error encountered by Advance.  Never blocks.
		Err() error
	}
	// SendStream returns the send side of the ServiceA.MethodA4 client stream.
	SendStream() interface {
		// Send places the item onto the output stream.  Returns errors encountered
		// while sending, or if Send is called after Close or Cancel.  Blocks if
		// there is no buffer space; will unblock when buffer space is available or
		// after Cancel.
		Send(item int32) error
		// Close indicates to the server that no more items will be sent; server
		// Recv calls will receive io.EOF after all sent items.  This is an optional
		// call - e.g. a client might call Close if it needs to continue receiving
		// items from the server after it's done sending.  Returns errors
		// encountered while closing, or if Close is called after Cancel.  Like
		// Send, blocks if there is no buffer space available.
		Close() error
	}
}

// ServiceAMethodA4Call represents the call returned from ServiceA.MethodA4.
type ServiceAMethodA4Call interface {
	ServiceAMethodA4ClientStream
	// Finish performs the equivalent of SendStream().Close, then blocks until
	// the server is done, and returns the positional return values for the call.
	//
	// Finish returns immediately if Cancel has been called; depending on the
	// timing the output could either be an error signaling cancelation, or the
	// valid positional return values from the server.
	//
	// Calling Finish is mandatory for releasing stream resources, unless Cancel
	// has been called or any of the other methods return an error.  Finish should
	// be called at most once.
	Finish() error
	// Cancel cancels the RPC, notifying the server to stop processing.  It is
	// safe to call Cancel concurrently with any of the other stream methods.
	// Calling Cancel after Finish has returned is a no-op.
	Cancel()
}

type implServiceAMethodA4Call struct {
	__ipc.Call
	valRecv string
	errRecv error
}

func (c *implServiceAMethodA4Call) RecvStream() interface {
	Advance() bool
	Value() string
	Err() error
} {
	return implServiceAMethodA4CallRecv{c}
}

type implServiceAMethodA4CallRecv struct {
	c *implServiceAMethodA4Call
}

func (c implServiceAMethodA4CallRecv) Advance() bool {
	c.c.errRecv = c.c.Recv(&c.c.valRecv)
	return c.c.errRecv == nil
}
func (c implServiceAMethodA4CallRecv) Value() string {
	return c.c.valRecv
}
func (c implServiceAMethodA4CallRecv) Err() error {
	if c.c.errRecv == __io.EOF {
		return nil
	}
	return c.c.errRecv
}
func (c *implServiceAMethodA4Call) SendStream() interface {
	Send(item int32) error
	Close() error
} {
	return implServiceAMethodA4CallSend{c}
}

type implServiceAMethodA4CallSend struct {
	c *implServiceAMethodA4Call
}

func (c implServiceAMethodA4CallSend) Send(item int32) error {
	return c.c.Send(item)
}
func (c implServiceAMethodA4CallSend) Close() error {
	return c.c.CloseSend()
}
func (c *implServiceAMethodA4Call) Finish() (err error) {
	if ierr := c.Call.Finish(&err); ierr != nil {
		err = ierr
	}
	return
}

// ServiceAServerMethods is the interface a server writer
// implements for ServiceA.
type ServiceAServerMethods interface {
	MethodA1(__ipc.ServerContext) error
	MethodA2(ctx __ipc.ServerContext, a int32, b string) (s string, err error)
	MethodA3(ctx ServiceAMethodA3Context, a int32) (s string, err error)
	MethodA4(ctx ServiceAMethodA4Context, a int32) error
}

// ServiceAServerStubMethods is the server interface containing
// ServiceA methods, as expected by ipc.Server.
// The only difference between this interface and ServiceAServerMethods
// is the streaming methods.
type ServiceAServerStubMethods interface {
	MethodA1(__ipc.ServerContext) error
	MethodA2(ctx __ipc.ServerContext, a int32, b string) (s string, err error)
	MethodA3(ctx *ServiceAMethodA3ContextStub, a int32) (s string, err error)
	MethodA4(ctx *ServiceAMethodA4ContextStub, a int32) error
}

// ServiceAServerStub adds universal methods to ServiceAServerStubMethods.
type ServiceAServerStub interface {
	ServiceAServerStubMethods
	// Describe the ServiceA interfaces.
	Describe__() []__ipc.InterfaceDesc
	// Signature will be replaced with Describe__.
	Signature(ctx __ipc.ServerContext) (__ipc.ServiceSignature, error)
}

// ServiceAServer returns a server stub for ServiceA.
// It converts an implementation of ServiceAServerMethods into
// an object that may be used by ipc.Server.
func ServiceAServer(impl ServiceAServerMethods) ServiceAServerStub {
	stub := implServiceAServerStub{
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

type implServiceAServerStub struct {
	impl ServiceAServerMethods
	gs   *__ipc.GlobState
}

func (s implServiceAServerStub) MethodA1(ctx __ipc.ServerContext) error {
	return s.impl.MethodA1(ctx)
}

func (s implServiceAServerStub) MethodA2(ctx __ipc.ServerContext, i0 int32, i1 string) (string, error) {
	return s.impl.MethodA2(ctx, i0, i1)
}

func (s implServiceAServerStub) MethodA3(ctx *ServiceAMethodA3ContextStub, i0 int32) (string, error) {
	return s.impl.MethodA3(ctx, i0)
}

func (s implServiceAServerStub) MethodA4(ctx *ServiceAMethodA4ContextStub, i0 int32) error {
	return s.impl.MethodA4(ctx, i0)
}

func (s implServiceAServerStub) VGlob() *__ipc.GlobState {
	return s.gs
}

func (s implServiceAServerStub) Describe__() []__ipc.InterfaceDesc {
	return []__ipc.InterfaceDesc{ServiceADesc}
}

// ServiceADesc describes the ServiceA interface.
var ServiceADesc __ipc.InterfaceDesc = descServiceA

// descServiceA hides the desc to keep godoc clean.
var descServiceA = __ipc.InterfaceDesc{
	Name:    "ServiceA",
	PkgPath: "veyron.io/veyron/veyron2/vdl/testdata/base",
	Methods: []__ipc.MethodDesc{
		{
			Name: "MethodA1",
			OutArgs: []__ipc.ArgDesc{
				{"", ``}, // error
			},
		},
		{
			Name: "MethodA2",
			InArgs: []__ipc.ArgDesc{
				{"a", ``}, // int32
				{"b", ``}, // string
			},
			OutArgs: []__ipc.ArgDesc{
				{"s", ``},   // string
				{"err", ``}, // error
			},
		},
		{
			Name: "MethodA3",
			InArgs: []__ipc.ArgDesc{
				{"a", ``}, // int32
			},
			OutArgs: []__ipc.ArgDesc{
				{"s", ``},   // string
				{"err", ``}, // error
			},
			Tags: []__vdlutil.Any{"tag", uint64(6)},
		},
		{
			Name: "MethodA4",
			InArgs: []__ipc.ArgDesc{
				{"a", ``}, // int32
			},
			OutArgs: []__ipc.ArgDesc{
				{"", ``}, // error
			},
		},
	},
}

func (s implServiceAServerStub) Signature(ctx __ipc.ServerContext) (__ipc.ServiceSignature, error) {
	// TODO(toddw): Replace with new Describe__ implementation.
	result := __ipc.ServiceSignature{Methods: make(map[string]__ipc.MethodSignature)}
	result.Methods["MethodA1"] = __ipc.MethodSignature{
		InArgs: []__ipc.MethodArgument{},
		OutArgs: []__ipc.MethodArgument{
			{Name: "", Type: 65},
		},
	}
	result.Methods["MethodA2"] = __ipc.MethodSignature{
		InArgs: []__ipc.MethodArgument{
			{Name: "a", Type: 36},
			{Name: "b", Type: 3},
		},
		OutArgs: []__ipc.MethodArgument{
			{Name: "s", Type: 3},
			{Name: "err", Type: 65},
		},
	}
	result.Methods["MethodA3"] = __ipc.MethodSignature{
		InArgs: []__ipc.MethodArgument{
			{Name: "a", Type: 36},
		},
		OutArgs: []__ipc.MethodArgument{
			{Name: "s", Type: 3},
			{Name: "err", Type: 65},
		},

		OutStream: 82,
	}
	result.Methods["MethodA4"] = __ipc.MethodSignature{
		InArgs: []__ipc.MethodArgument{
			{Name: "a", Type: 36},
		},
		OutArgs: []__ipc.MethodArgument{
			{Name: "", Type: 65},
		},
		InStream:  36,
		OutStream: 3,
	}

	result.TypeDefs = []__vdlutil.Any{
		__wiretype.NamedPrimitiveType{Type: 0x1, Name: "error", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x32, Name: "byte", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x1, Name: "anydata", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x7, Name: "TypeID", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x2, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedBool", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x32, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedByte", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x33, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedUint16", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x34, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedUint32", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x35, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedUint64", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x23, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedInt16", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x24, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedInt32", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x25, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedInt64", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x19, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedFloat32", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x1a, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedFloat64", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x38, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedComplex64", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x39, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedComplex128", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x3, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedString", Tags: []string(nil)}, __wiretype.StructType{
			[]__wiretype.FieldType{
				__wiretype.FieldType{Type: 0x2, Name: "A0"},
				__wiretype.FieldType{Type: 0x42, Name: "A1"},
				__wiretype.FieldType{Type: 0x33, Name: "A2"},
				__wiretype.FieldType{Type: 0x34, Name: "A3"},
				__wiretype.FieldType{Type: 0x35, Name: "A4"},
				__wiretype.FieldType{Type: 0x23, Name: "A5"},
				__wiretype.FieldType{Type: 0x24, Name: "A6"},
				__wiretype.FieldType{Type: 0x25, Name: "A7"},
				__wiretype.FieldType{Type: 0x19, Name: "A8"},
				__wiretype.FieldType{Type: 0x1a, Name: "A9"},
				__wiretype.FieldType{Type: 0x38, Name: "A10"},
				__wiretype.FieldType{Type: 0x39, Name: "A11"},
				__wiretype.FieldType{Type: 0x3, Name: "A12"},
				__wiretype.FieldType{Type: 0x41, Name: "A13"},
				__wiretype.FieldType{Type: 0x43, Name: "A14"},
				__wiretype.FieldType{Type: 0x44, Name: "A15"},
				__wiretype.FieldType{Type: 0x45, Name: "B0"},
				__wiretype.FieldType{Type: 0x46, Name: "B1"},
				__wiretype.FieldType{Type: 0x47, Name: "B2"},
				__wiretype.FieldType{Type: 0x48, Name: "B3"},
				__wiretype.FieldType{Type: 0x49, Name: "B4"},
				__wiretype.FieldType{Type: 0x4a, Name: "B5"},
				__wiretype.FieldType{Type: 0x4b, Name: "B6"},
				__wiretype.FieldType{Type: 0x4c, Name: "B7"},
				__wiretype.FieldType{Type: 0x4d, Name: "B8"},
				__wiretype.FieldType{Type: 0x4e, Name: "B9"},
				__wiretype.FieldType{Type: 0x4f, Name: "B10"},
				__wiretype.FieldType{Type: 0x50, Name: "B11"},
				__wiretype.FieldType{Type: 0x51, Name: "B12"},
			},
			"veyron.io/veyron/veyron2/vdl/testdata/base.Scalars", []string(nil)},
	}

	return result, nil
}

// ServiceAMethodA3ServerStream is the server stream for ServiceA.MethodA3.
type ServiceAMethodA3ServerStream interface {
	// SendStream returns the send side of the ServiceA.MethodA3 server stream.
	SendStream() interface {
		// Send places the item onto the output stream.  Returns errors encountered
		// while sending.  Blocks if there is no buffer space; will unblock when
		// buffer space is available.
		Send(item Scalars) error
	}
}

// ServiceAMethodA3Context represents the context passed to ServiceA.MethodA3.
type ServiceAMethodA3Context interface {
	__ipc.ServerContext
	ServiceAMethodA3ServerStream
}

// ServiceAMethodA3ContextStub is a wrapper that converts ipc.ServerCall into
// a typesafe stub that implements ServiceAMethodA3Context.
type ServiceAMethodA3ContextStub struct {
	__ipc.ServerCall
}

// Init initializes ServiceAMethodA3ContextStub from ipc.ServerCall.
func (s *ServiceAMethodA3ContextStub) Init(call __ipc.ServerCall) {
	s.ServerCall = call
}

// SendStream returns the send side of the ServiceA.MethodA3 server stream.
func (s *ServiceAMethodA3ContextStub) SendStream() interface {
	Send(item Scalars) error
} {
	return implServiceAMethodA3ContextSend{s}
}

type implServiceAMethodA3ContextSend struct {
	s *ServiceAMethodA3ContextStub
}

func (s implServiceAMethodA3ContextSend) Send(item Scalars) error {
	return s.s.Send(item)
}

// ServiceAMethodA4ServerStream is the server stream for ServiceA.MethodA4.
type ServiceAMethodA4ServerStream interface {
	// RecvStream returns the receiver side of the ServiceA.MethodA4 server stream.
	RecvStream() interface {
		// Advance stages an item so that it may be retrieved via Value.  Returns
		// true iff there is an item to retrieve.  Advance must be called before
		// Value is called.  May block if an item is not available.
		Advance() bool
		// Value returns the item that was staged by Advance.  May panic if Advance
		// returned false or was not called.  Never blocks.
		Value() int32
		// Err returns any error encountered by Advance.  Never blocks.
		Err() error
	}
	// SendStream returns the send side of the ServiceA.MethodA4 server stream.
	SendStream() interface {
		// Send places the item onto the output stream.  Returns errors encountered
		// while sending.  Blocks if there is no buffer space; will unblock when
		// buffer space is available.
		Send(item string) error
	}
}

// ServiceAMethodA4Context represents the context passed to ServiceA.MethodA4.
type ServiceAMethodA4Context interface {
	__ipc.ServerContext
	ServiceAMethodA4ServerStream
}

// ServiceAMethodA4ContextStub is a wrapper that converts ipc.ServerCall into
// a typesafe stub that implements ServiceAMethodA4Context.
type ServiceAMethodA4ContextStub struct {
	__ipc.ServerCall
	valRecv int32
	errRecv error
}

// Init initializes ServiceAMethodA4ContextStub from ipc.ServerCall.
func (s *ServiceAMethodA4ContextStub) Init(call __ipc.ServerCall) {
	s.ServerCall = call
}

// RecvStream returns the receiver side of the ServiceA.MethodA4 server stream.
func (s *ServiceAMethodA4ContextStub) RecvStream() interface {
	Advance() bool
	Value() int32
	Err() error
} {
	return implServiceAMethodA4ContextRecv{s}
}

type implServiceAMethodA4ContextRecv struct {
	s *ServiceAMethodA4ContextStub
}

func (s implServiceAMethodA4ContextRecv) Advance() bool {
	s.s.errRecv = s.s.Recv(&s.s.valRecv)
	return s.s.errRecv == nil
}
func (s implServiceAMethodA4ContextRecv) Value() int32 {
	return s.s.valRecv
}
func (s implServiceAMethodA4ContextRecv) Err() error {
	if s.s.errRecv == __io.EOF {
		return nil
	}
	return s.s.errRecv
}

// SendStream returns the send side of the ServiceA.MethodA4 server stream.
func (s *ServiceAMethodA4ContextStub) SendStream() interface {
	Send(item string) error
} {
	return implServiceAMethodA4ContextSend{s}
}

type implServiceAMethodA4ContextSend struct {
	s *ServiceAMethodA4ContextStub
}

func (s implServiceAMethodA4ContextSend) Send(item string) error {
	return s.s.Send(item)
}

// ServiceBClientMethods is the client interface
// containing ServiceB methods.
type ServiceBClientMethods interface {
	ServiceAClientMethods
	MethodB1(ctx __context.T, a Scalars, b Composites, opts ...__ipc.CallOpt) (c CompComp, err error)
}

// ServiceBClientStub adds universal methods to ServiceBClientMethods.
type ServiceBClientStub interface {
	ServiceBClientMethods
	__ipc.UniversalServiceMethods
}

// ServiceBClient returns a client stub for ServiceB.
func ServiceBClient(name string, opts ...__ipc.BindOpt) ServiceBClientStub {
	var client __ipc.Client
	for _, opt := range opts {
		if clientOpt, ok := opt.(__ipc.Client); ok {
			client = clientOpt
		}
	}
	return implServiceBClientStub{name, client, ServiceAClient(name, client)}
}

type implServiceBClientStub struct {
	name   string
	client __ipc.Client

	ServiceAClientStub
}

func (c implServiceBClientStub) c(ctx __context.T) __ipc.Client {
	if c.client != nil {
		return c.client
	}
	return __veyron2.RuntimeFromContext(ctx).Client()
}

func (c implServiceBClientStub) MethodB1(ctx __context.T, i0 Scalars, i1 Composites, opts ...__ipc.CallOpt) (o0 CompComp, err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "MethodB1", []interface{}{i0, i1}, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&o0, &err); ierr != nil {
		err = ierr
	}
	return
}

func (c implServiceBClientStub) Signature(ctx __context.T, opts ...__ipc.CallOpt) (o0 __ipc.ServiceSignature, err error) {
	var call __ipc.Call
	if call, err = c.c(ctx).StartCall(ctx, c.name, "Signature", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&o0, &err); ierr != nil {
		err = ierr
	}
	return
}

// ServiceBServerMethods is the interface a server writer
// implements for ServiceB.
type ServiceBServerMethods interface {
	ServiceAServerMethods
	MethodB1(ctx __ipc.ServerContext, a Scalars, b Composites) (c CompComp, err error)
}

// ServiceBServerStubMethods is the server interface containing
// ServiceB methods, as expected by ipc.Server.
// The only difference between this interface and ServiceBServerMethods
// is the streaming methods.
type ServiceBServerStubMethods interface {
	ServiceAServerStubMethods
	MethodB1(ctx __ipc.ServerContext, a Scalars, b Composites) (c CompComp, err error)
}

// ServiceBServerStub adds universal methods to ServiceBServerStubMethods.
type ServiceBServerStub interface {
	ServiceBServerStubMethods
	// Describe the ServiceB interfaces.
	Describe__() []__ipc.InterfaceDesc
	// Signature will be replaced with Describe__.
	Signature(ctx __ipc.ServerContext) (__ipc.ServiceSignature, error)
}

// ServiceBServer returns a server stub for ServiceB.
// It converts an implementation of ServiceBServerMethods into
// an object that may be used by ipc.Server.
func ServiceBServer(impl ServiceBServerMethods) ServiceBServerStub {
	stub := implServiceBServerStub{
		impl:               impl,
		ServiceAServerStub: ServiceAServer(impl),
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

type implServiceBServerStub struct {
	impl ServiceBServerMethods
	ServiceAServerStub
	gs *__ipc.GlobState
}

func (s implServiceBServerStub) MethodB1(ctx __ipc.ServerContext, i0 Scalars, i1 Composites) (CompComp, error) {
	return s.impl.MethodB1(ctx, i0, i1)
}

func (s implServiceBServerStub) VGlob() *__ipc.GlobState {
	return s.gs
}

func (s implServiceBServerStub) Describe__() []__ipc.InterfaceDesc {
	return []__ipc.InterfaceDesc{ServiceBDesc, ServiceADesc}
}

// ServiceBDesc describes the ServiceB interface.
var ServiceBDesc __ipc.InterfaceDesc = descServiceB

// descServiceB hides the desc to keep godoc clean.
var descServiceB = __ipc.InterfaceDesc{
	Name:    "ServiceB",
	PkgPath: "veyron.io/veyron/veyron2/vdl/testdata/base",
	Embeds: []__ipc.EmbedDesc{
		{"ServiceA", "veyron.io/veyron/veyron2/vdl/testdata/base", ``},
	},
	Methods: []__ipc.MethodDesc{
		{
			Name: "MethodB1",
			InArgs: []__ipc.ArgDesc{
				{"a", ``}, // Scalars
				{"b", ``}, // Composites
			},
			OutArgs: []__ipc.ArgDesc{
				{"c", ``},   // CompComp
				{"err", ``}, // error
			},
		},
	},
}

func (s implServiceBServerStub) Signature(ctx __ipc.ServerContext) (__ipc.ServiceSignature, error) {
	// TODO(toddw): Replace with new Describe__ implementation.
	result := __ipc.ServiceSignature{Methods: make(map[string]__ipc.MethodSignature)}
	result.Methods["MethodB1"] = __ipc.MethodSignature{
		InArgs: []__ipc.MethodArgument{
			{Name: "a", Type: 82},
			{Name: "b", Type: 91},
		},
		OutArgs: []__ipc.MethodArgument{
			{Name: "c", Type: 97},
			{Name: "err", Type: 66},
		},
	}

	result.TypeDefs = []__vdlutil.Any{
		__wiretype.NamedPrimitiveType{Type: 0x32, Name: "byte", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x1, Name: "error", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x1, Name: "anydata", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x7, Name: "TypeID", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x2, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedBool", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x32, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedByte", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x33, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedUint16", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x34, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedUint32", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x35, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedUint64", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x23, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedInt16", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x24, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedInt32", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x25, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedInt64", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x19, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedFloat32", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x1a, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedFloat64", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x38, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedComplex64", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x39, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedComplex128", Tags: []string(nil)}, __wiretype.NamedPrimitiveType{Type: 0x3, Name: "veyron.io/veyron/veyron2/vdl/testdata/base.NamedString", Tags: []string(nil)}, __wiretype.StructType{
			[]__wiretype.FieldType{
				__wiretype.FieldType{Type: 0x2, Name: "A0"},
				__wiretype.FieldType{Type: 0x41, Name: "A1"},
				__wiretype.FieldType{Type: 0x33, Name: "A2"},
				__wiretype.FieldType{Type: 0x34, Name: "A3"},
				__wiretype.FieldType{Type: 0x35, Name: "A4"},
				__wiretype.FieldType{Type: 0x23, Name: "A5"},
				__wiretype.FieldType{Type: 0x24, Name: "A6"},
				__wiretype.FieldType{Type: 0x25, Name: "A7"},
				__wiretype.FieldType{Type: 0x19, Name: "A8"},
				__wiretype.FieldType{Type: 0x1a, Name: "A9"},
				__wiretype.FieldType{Type: 0x38, Name: "A10"},
				__wiretype.FieldType{Type: 0x39, Name: "A11"},
				__wiretype.FieldType{Type: 0x3, Name: "A12"},
				__wiretype.FieldType{Type: 0x42, Name: "A13"},
				__wiretype.FieldType{Type: 0x43, Name: "A14"},
				__wiretype.FieldType{Type: 0x44, Name: "A15"},
				__wiretype.FieldType{Type: 0x45, Name: "B0"},
				__wiretype.FieldType{Type: 0x46, Name: "B1"},
				__wiretype.FieldType{Type: 0x47, Name: "B2"},
				__wiretype.FieldType{Type: 0x48, Name: "B3"},
				__wiretype.FieldType{Type: 0x49, Name: "B4"},
				__wiretype.FieldType{Type: 0x4a, Name: "B5"},
				__wiretype.FieldType{Type: 0x4b, Name: "B6"},
				__wiretype.FieldType{Type: 0x4c, Name: "B7"},
				__wiretype.FieldType{Type: 0x4d, Name: "B8"},
				__wiretype.FieldType{Type: 0x4e, Name: "B9"},
				__wiretype.FieldType{Type: 0x4f, Name: "B10"},
				__wiretype.FieldType{Type: 0x50, Name: "B11"},
				__wiretype.FieldType{Type: 0x51, Name: "B12"},
			},
			"veyron.io/veyron/veyron2/vdl/testdata/base.Scalars", []string(nil)},
		__wiretype.ArrayType{Elem: 0x52, Len: 0x2, Name: "", Tags: []string(nil)}, __wiretype.SliceType{Elem: 0x52, Name: "", Tags: []string(nil)}, __wiretype.StructType{
			[]__wiretype.FieldType{
				__wiretype.FieldType{Type: 0x2, Name: "A0"},
				__wiretype.FieldType{Type: 0x41, Name: "A1"},
				__wiretype.FieldType{Type: 0x33, Name: "A2"},
				__wiretype.FieldType{Type: 0x34, Name: "A3"},
				__wiretype.FieldType{Type: 0x35, Name: "A4"},
				__wiretype.FieldType{Type: 0x23, Name: "A5"},
				__wiretype.FieldType{Type: 0x24, Name: "A6"},
				__wiretype.FieldType{Type: 0x25, Name: "A7"},
				__wiretype.FieldType{Type: 0x19, Name: "A8"},
				__wiretype.FieldType{Type: 0x1a, Name: "A9"},
				__wiretype.FieldType{Type: 0x38, Name: "A10"},
				__wiretype.FieldType{Type: 0x39, Name: "A11"},
				__wiretype.FieldType{Type: 0x3, Name: "A12"},
				__wiretype.FieldType{Type: 0x42, Name: "A13"},
				__wiretype.FieldType{Type: 0x45, Name: "B0"},
				__wiretype.FieldType{Type: 0x46, Name: "B1"},
				__wiretype.FieldType{Type: 0x47, Name: "B2"},
				__wiretype.FieldType{Type: 0x48, Name: "B3"},
				__wiretype.FieldType{Type: 0x49, Name: "B4"},
				__wiretype.FieldType{Type: 0x4a, Name: "B5"},
				__wiretype.FieldType{Type: 0x4b, Name: "B6"},
				__wiretype.FieldType{Type: 0x4c, Name: "B7"},
				__wiretype.FieldType{Type: 0x4d, Name: "B8"},
				__wiretype.FieldType{Type: 0x4e, Name: "B9"},
				__wiretype.FieldType{Type: 0x4f, Name: "B10"},
				__wiretype.FieldType{Type: 0x50, Name: "B11"},
				__wiretype.FieldType{Type: 0x51, Name: "B12"},
			},
			"veyron.io/veyron/veyron2/vdl/testdata/base.KeyScalars", []string(nil)},
		__wiretype.MapType{Key: 0x55, Elem: 0x2, Name: "", Tags: []string(nil)}, __wiretype.MapType{Key: 0x3, Elem: 0x52, Name: "", Tags: []string(nil)}, __wiretype.MapType{Key: 0x3, Elem: 0x39, Name: "", Tags: []string(nil)}, __wiretype.SliceType{Elem: 0x58, Name: "", Tags: []string(nil)}, __wiretype.MapType{Key: 0x55, Elem: 0x59, Name: "", Tags: []string(nil)}, __wiretype.StructType{
			[]__wiretype.FieldType{
				__wiretype.FieldType{Type: 0x52, Name: "A0"},
				__wiretype.FieldType{Type: 0x53, Name: "A1"},
				__wiretype.FieldType{Type: 0x54, Name: "A2"},
				__wiretype.FieldType{Type: 0x56, Name: "A3"},
				__wiretype.FieldType{Type: 0x57, Name: "A4"},
				__wiretype.FieldType{Type: 0x5a, Name: "A5"},
			},
			"veyron.io/veyron/veyron2/vdl/testdata/base.Composites", []string(nil)},
		__wiretype.ArrayType{Elem: 0x5b, Len: 0x2, Name: "", Tags: []string(nil)}, __wiretype.SliceType{Elem: 0x5b, Name: "", Tags: []string(nil)}, __wiretype.MapType{Key: 0x3, Elem: 0x5b, Name: "", Tags: []string(nil)}, __wiretype.SliceType{Elem: 0x5e, Name: "", Tags: []string(nil)}, __wiretype.MapType{Key: 0x55, Elem: 0x5f, Name: "", Tags: []string(nil)}, __wiretype.StructType{
			[]__wiretype.FieldType{
				__wiretype.FieldType{Type: 0x5b, Name: "A0"},
				__wiretype.FieldType{Type: 0x5c, Name: "A1"},
				__wiretype.FieldType{Type: 0x5d, Name: "A2"},
				__wiretype.FieldType{Type: 0x5e, Name: "A3"},
				__wiretype.FieldType{Type: 0x60, Name: "A4"},
			},
			"veyron.io/veyron/veyron2/vdl/testdata/base.CompComp", []string(nil)},
	}
	var ss __ipc.ServiceSignature
	var firstAdded int
	ss, _ = s.ServiceAServerStub.Signature(ctx)
	firstAdded = len(result.TypeDefs)
	for k, v := range ss.Methods {
		for i, _ := range v.InArgs {
			if v.InArgs[i].Type >= __wiretype.TypeIDFirst {
				v.InArgs[i].Type += __wiretype.TypeID(firstAdded)
			}
		}
		for i, _ := range v.OutArgs {
			if v.OutArgs[i].Type >= __wiretype.TypeIDFirst {
				v.OutArgs[i].Type += __wiretype.TypeID(firstAdded)
			}
		}
		if v.InStream >= __wiretype.TypeIDFirst {
			v.InStream += __wiretype.TypeID(firstAdded)
		}
		if v.OutStream >= __wiretype.TypeIDFirst {
			v.OutStream += __wiretype.TypeID(firstAdded)
		}
		result.Methods[k] = v
	}
	//TODO(bprosnitz) combine type definitions from embeded interfaces in a way that doesn't cause duplication.
	for _, d := range ss.TypeDefs {
		switch wt := d.(type) {
		case __wiretype.SliceType:
			if wt.Elem >= __wiretype.TypeIDFirst {
				wt.Elem += __wiretype.TypeID(firstAdded)
			}
			d = wt
		case __wiretype.ArrayType:
			if wt.Elem >= __wiretype.TypeIDFirst {
				wt.Elem += __wiretype.TypeID(firstAdded)
			}
			d = wt
		case __wiretype.MapType:
			if wt.Key >= __wiretype.TypeIDFirst {
				wt.Key += __wiretype.TypeID(firstAdded)
			}
			if wt.Elem >= __wiretype.TypeIDFirst {
				wt.Elem += __wiretype.TypeID(firstAdded)
			}
			d = wt
		case __wiretype.StructType:
			for i, fld := range wt.Fields {
				if fld.Type >= __wiretype.TypeIDFirst {
					wt.Fields[i].Type += __wiretype.TypeID(firstAdded)
				}
			}
			d = wt
			// NOTE: other types are missing, but we are upgrading anyways.
		}
		result.TypeDefs = append(result.TypeDefs, d)
	}

	return result, nil
}
