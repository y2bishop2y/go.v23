// This file was auto-generated by the veyron vdl tool.
// Source: build.vdl

// Package build supports building and describing Veyron binaries.
package build

import (
	"veyron2/services/mgmt/binary"

	// The non-user imports are prefixed with "_gen_" to prevent collisions.
	_gen_veyron2 "veyron2"
	_gen_context "veyron2/context"
	_gen_ipc "veyron2/ipc"
	_gen_naming "veyron2/naming"
	_gen_rt "veyron2/rt"
	_gen_vdl "veyron2/vdl"
	_gen_wiretype "veyron2/wiretype"
)

type File struct {
	Name     string
	Contents []byte
}

// Build describes an interface for building binaries from source.
// Build is the interface the client binds and uses.
// Build_ExcludingUniversal is the interface without internal framework-added methods
// to enable embedding without method collisions.  Not to be used directly by clients.
type Build_ExcludingUniversal interface {
	// Build streams sources to the build server, which then attempts to
	// build the sources and returns the output.
	Build(ctx _gen_context.T, opts ..._gen_ipc.CallOpt) (reply BuildBuildStream, err error)
	// Describe generates a description for a binary identified by
	// the given Object name.
	Describe(ctx _gen_context.T, Name string, opts ..._gen_ipc.CallOpt) (reply binary.Description, err error)
}
type Build interface {
	_gen_ipc.UniversalServiceMethods
	Build_ExcludingUniversal
}

// BuildService is the interface the server implements.
type BuildService interface {

	// Build streams sources to the build server, which then attempts to
	// build the sources and returns the output.
	Build(context _gen_ipc.ServerContext, stream BuildServiceBuildStream) (reply []byte, err error)
	// Describe generates a description for a binary identified by
	// the given Object name.
	Describe(context _gen_ipc.ServerContext, Name string) (reply binary.Description, err error)
}

// BuildBuildStream is the interface for streaming responses of the method
// Build in the service interface Build.
type BuildBuildStream interface {

	// Send places the item onto the output stream, blocking if there is no buffer
	// space available.
	Send(item File) error

	// CloseSend indicates to the server that no more items will be sent; server
	// Recv calls will receive io.EOF after all sent items.  Subsequent calls to
	// Send on the client will fail.  This is an optional call - it's used by
	// streaming clients that need the server to receive the io.EOF terminator.
	CloseSend() error

	// Finish closes the stream and returns the positional return values for
	// call.
	Finish() (reply []byte, err error)

	// Cancel cancels the RPC, notifying the server to stop processing.
	Cancel()
}

// Implementation of the BuildBuildStream interface that is not exported.
type implBuildBuildStream struct {
	clientCall _gen_ipc.Call
}

func (c *implBuildBuildStream) Send(item File) error {
	return c.clientCall.Send(item)
}

func (c *implBuildBuildStream) CloseSend() error {
	return c.clientCall.CloseSend()
}

func (c *implBuildBuildStream) Finish() (reply []byte, err error) {
	if ierr := c.clientCall.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

func (c *implBuildBuildStream) Cancel() {
	c.clientCall.Cancel()
}

// BuildServiceBuildStream is the interface for streaming responses of the method
// Build in the service interface Build.
type BuildServiceBuildStream interface {

	// Recv fills itemptr with the next item in the input stream, blocking until
	// an item is available.  Returns io.EOF to indicate graceful end of input.
	Recv() (item File, err error)
}

// Implementation of the BuildServiceBuildStream interface that is not exported.
type implBuildServiceBuildStream struct {
	serverCall _gen_ipc.ServerCall
}

func (s *implBuildServiceBuildStream) Recv() (item File, err error) {
	err = s.serverCall.Recv(&item)
	return
}

// BindBuild returns the client stub implementing the Build
// interface.
//
// If no _gen_ipc.Client is specified, the default _gen_ipc.Client in the
// global Runtime is used.
func BindBuild(name string, opts ..._gen_ipc.BindOpt) (Build, error) {
	var client _gen_ipc.Client
	switch len(opts) {
	case 0:
		client = _gen_rt.R().Client()
	case 1:
		switch o := opts[0].(type) {
		case _gen_veyron2.Runtime:
			client = o.Client()
		case _gen_ipc.Client:
			client = o
		default:
			return nil, _gen_vdl.ErrUnrecognizedOption
		}
	default:
		return nil, _gen_vdl.ErrTooManyOptionsToBind
	}
	stub := &clientStubBuild{client: client, name: name}

	return stub, nil
}

// NewServerBuild creates a new server stub.
//
// It takes a regular server implementing the BuildService
// interface, and returns a new server stub.
func NewServerBuild(server BuildService) interface{} {
	return &ServerStubBuild{
		service: server,
	}
}

// clientStubBuild implements Build.
type clientStubBuild struct {
	client _gen_ipc.Client
	name   string
}

func (__gen_c *clientStubBuild) Build(ctx _gen_context.T, opts ..._gen_ipc.CallOpt) (reply BuildBuildStream, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "Build", nil, opts...); err != nil {
		return
	}
	reply = &implBuildBuildStream{clientCall: call}
	return
}

func (__gen_c *clientStubBuild) Describe(ctx _gen_context.T, Name string, opts ..._gen_ipc.CallOpt) (reply binary.Description, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "Describe", []interface{}{Name}, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

func (__gen_c *clientStubBuild) UnresolveStep(ctx _gen_context.T, opts ..._gen_ipc.CallOpt) (reply []string, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "UnresolveStep", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

func (__gen_c *clientStubBuild) Signature(ctx _gen_context.T, opts ..._gen_ipc.CallOpt) (reply _gen_ipc.ServiceSignature, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "Signature", nil, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

func (__gen_c *clientStubBuild) GetMethodTags(ctx _gen_context.T, method string, opts ..._gen_ipc.CallOpt) (reply []interface{}, err error) {
	var call _gen_ipc.Call
	if call, err = __gen_c.client.StartCall(ctx, __gen_c.name, "GetMethodTags", []interface{}{method}, opts...); err != nil {
		return
	}
	if ierr := call.Finish(&reply, &err); ierr != nil {
		err = ierr
	}
	return
}

// ServerStubBuild wraps a server that implements
// BuildService and provides an object that satisfies
// the requirements of veyron2/ipc.ReflectInvoker.
type ServerStubBuild struct {
	service BuildService
}

func (__gen_s *ServerStubBuild) GetMethodTags(call _gen_ipc.ServerCall, method string) ([]interface{}, error) {
	// TODO(bprosnitz) GetMethodTags() will be replaces with Signature().
	// Note: This exhibits some weird behavior like returning a nil error if the method isn't found.
	// This will change when it is replaced with Signature().
	switch method {
	case "Build":
		return []interface{}{}, nil
	case "Describe":
		return []interface{}{}, nil
	default:
		return nil, nil
	}
}

func (__gen_s *ServerStubBuild) Signature(call _gen_ipc.ServerCall) (_gen_ipc.ServiceSignature, error) {
	result := _gen_ipc.ServiceSignature{Methods: make(map[string]_gen_ipc.MethodSignature)}
	result.Methods["Build"] = _gen_ipc.MethodSignature{
		InArgs: []_gen_ipc.MethodArgument{},
		OutArgs: []_gen_ipc.MethodArgument{
			{Name: "", Type: 66},
			{Name: "", Type: 67},
		},
		InStream: 68,
	}
	result.Methods["Describe"] = _gen_ipc.MethodSignature{
		InArgs: []_gen_ipc.MethodArgument{
			{Name: "Name", Type: 3},
		},
		OutArgs: []_gen_ipc.MethodArgument{
			{Name: "", Type: 70},
			{Name: "", Type: 67},
		},
	}

	result.TypeDefs = []_gen_vdl.Any{
		_gen_wiretype.NamedPrimitiveType{Type: 0x32, Name: "byte", Tags: []string(nil)}, _gen_wiretype.SliceType{Elem: 0x41, Name: "", Tags: []string(nil)}, _gen_wiretype.NamedPrimitiveType{Type: 0x1, Name: "error", Tags: []string(nil)}, _gen_wiretype.StructType{
			[]_gen_wiretype.FieldType{
				_gen_wiretype.FieldType{Type: 0x3, Name: "Name"},
				_gen_wiretype.FieldType{Type: 0x42, Name: "Contents"},
			},
			"veyron2/services/mgmt/build.File", []string(nil)},
		_gen_wiretype.MapType{Key: 0x3, Elem: 0x2, Name: "", Tags: []string(nil)}, _gen_wiretype.StructType{
			[]_gen_wiretype.FieldType{
				_gen_wiretype.FieldType{Type: 0x3, Name: "Name"},
				_gen_wiretype.FieldType{Type: 0x45, Name: "Profiles"},
			},
			"veyron2/services/mgmt/binary.Description", []string(nil)},
	}

	return result, nil
}

func (__gen_s *ServerStubBuild) UnresolveStep(call _gen_ipc.ServerCall) (reply []string, err error) {
	if unresolver, ok := __gen_s.service.(_gen_ipc.Unresolver); ok {
		return unresolver.UnresolveStep(call)
	}
	if call.Server() == nil {
		return
	}
	var published []string
	if published, err = call.Server().Published(); err != nil || published == nil {
		return
	}
	reply = make([]string, len(published))
	for i, p := range published {
		reply[i] = _gen_naming.Join(p, call.Name())
	}
	return
}

func (__gen_s *ServerStubBuild) Build(call _gen_ipc.ServerCall) (reply []byte, err error) {
	stream := &implBuildServiceBuildStream{serverCall: call}
	reply, err = __gen_s.service.Build(call, stream)
	return
}

func (__gen_s *ServerStubBuild) Describe(call _gen_ipc.ServerCall, Name string) (reply binary.Description, err error) {
	reply, err = __gen_s.service.Describe(call, Name)
	return
}
