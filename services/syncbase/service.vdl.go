// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Source: service.vdl

// Package syncbase defines the wire API for a structured store that supports
// peer-to-peer synchronization.
//
// TODO(sadovsky): Write a detailed package description.
package syncbase

import (
	// VDL system imports
	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/i18n"
	"v.io/v23/rpc"
	"v.io/v23/vdl"
	"v.io/v23/verror"

	// VDL user imports
	"v.io/v23/security/access"
	"v.io/v23/services/permissions"
)

var (
	ErrInvalidName = verror.Register("v.io/syncbase/v23/services/syncbase.InvalidName", verror.NoRetry, "{1:}{2:} invalid name: {3}")
)

func init() {
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrInvalidName.ID), "{1:}{2:} invalid name: {3}")
}

// NewErrInvalidName returns an error with the ErrInvalidName ID.
func NewErrInvalidName(ctx *context.T, name string) error {
	return verror.New(ErrInvalidName, ctx, name)
}

// ServiceClientMethods is the client interface
// containing Service methods.
//
// Service represents a Vanadium Syncbase service.
// Service.Glob operates over App names.
type ServiceClientMethods interface {
	// Object provides access control for Vanadium objects.
	//
	// Vanadium services implementing dynamic access control would typically embed
	// this interface and tag additional methods defined by the service with one of
	// Admin, Read, Write, Resolve etc. For example, the VDL definition of the
	// object would be:
	//
	//   package mypackage
	//
	//   import "v.io/v23/security/access"
	//   import "v.io/v23/services/permissions"
	//
	//   type MyObject interface {
	//     permissions.Object
	//     MyRead() (string, error) {access.Read}
	//     MyWrite(string) error    {access.Write}
	//   }
	//
	// If the set of pre-defined tags is insufficient, services may define their
	// own tag type and annotate all methods with this new type.
	//
	// Instead of embedding this Object interface, define SetPermissions and
	// GetPermissions in their own interface. Authorization policies will typically
	// respect annotations of a single type. For example, the VDL definition of an
	// object would be:
	//
	//  package mypackage
	//
	//  import "v.io/v23/security/access"
	//
	//  type MyTag string
	//
	//  const (
	//    Blue = MyTag("Blue")
	//    Red  = MyTag("Red")
	//  )
	//
	//  type MyObject interface {
	//    MyMethod() (string, error) {Blue}
	//
	//    // Allow clients to change access via the access.Object interface:
	//    SetPermissions(perms access.Permissions, version string) error         {Red}
	//    GetPermissions() (perms access.Permissions, version string, err error) {Blue}
	//  }
	permissions.ObjectClientMethods
}

// ServiceClientStub adds universal methods to ServiceClientMethods.
type ServiceClientStub interface {
	ServiceClientMethods
	rpc.UniversalServiceMethods
}

// ServiceClient returns a client stub for Service.
func ServiceClient(name string) ServiceClientStub {
	return implServiceClientStub{name, permissions.ObjectClient(name)}
}

type implServiceClientStub struct {
	name string

	permissions.ObjectClientStub
}

// ServiceServerMethods is the interface a server writer
// implements for Service.
//
// Service represents a Vanadium Syncbase service.
// Service.Glob operates over App names.
type ServiceServerMethods interface {
	// Object provides access control for Vanadium objects.
	//
	// Vanadium services implementing dynamic access control would typically embed
	// this interface and tag additional methods defined by the service with one of
	// Admin, Read, Write, Resolve etc. For example, the VDL definition of the
	// object would be:
	//
	//   package mypackage
	//
	//   import "v.io/v23/security/access"
	//   import "v.io/v23/services/permissions"
	//
	//   type MyObject interface {
	//     permissions.Object
	//     MyRead() (string, error) {access.Read}
	//     MyWrite(string) error    {access.Write}
	//   }
	//
	// If the set of pre-defined tags is insufficient, services may define their
	// own tag type and annotate all methods with this new type.
	//
	// Instead of embedding this Object interface, define SetPermissions and
	// GetPermissions in their own interface. Authorization policies will typically
	// respect annotations of a single type. For example, the VDL definition of an
	// object would be:
	//
	//  package mypackage
	//
	//  import "v.io/v23/security/access"
	//
	//  type MyTag string
	//
	//  const (
	//    Blue = MyTag("Blue")
	//    Red  = MyTag("Red")
	//  )
	//
	//  type MyObject interface {
	//    MyMethod() (string, error) {Blue}
	//
	//    // Allow clients to change access via the access.Object interface:
	//    SetPermissions(perms access.Permissions, version string) error         {Red}
	//    GetPermissions() (perms access.Permissions, version string, err error) {Blue}
	//  }
	permissions.ObjectServerMethods
}

// ServiceServerStubMethods is the server interface containing
// Service methods, as expected by rpc.Server.
// There is no difference between this interface and ServiceServerMethods
// since there are no streaming methods.
type ServiceServerStubMethods ServiceServerMethods

// ServiceServerStub adds universal methods to ServiceServerStubMethods.
type ServiceServerStub interface {
	ServiceServerStubMethods
	// Describe the Service interfaces.
	Describe__() []rpc.InterfaceDesc
}

// ServiceServer returns a server stub for Service.
// It converts an implementation of ServiceServerMethods into
// an object that may be used by rpc.Server.
func ServiceServer(impl ServiceServerMethods) ServiceServerStub {
	stub := implServiceServerStub{
		impl:             impl,
		ObjectServerStub: permissions.ObjectServer(impl),
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := rpc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := rpc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implServiceServerStub struct {
	impl ServiceServerMethods
	permissions.ObjectServerStub
	gs *rpc.GlobState
}

func (s implServiceServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implServiceServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{ServiceDesc, permissions.ObjectDesc}
}

// ServiceDesc describes the Service interface.
var ServiceDesc rpc.InterfaceDesc = descService

// descService hides the desc to keep godoc clean.
var descService = rpc.InterfaceDesc{
	Name:    "Service",
	PkgPath: "v.io/syncbase/v23/services/syncbase",
	Doc:     "// Service represents a Vanadium Syncbase service.\n// Service.Glob operates over App names.",
	Embeds: []rpc.EmbedDesc{
		{"Object", "v.io/v23/services/permissions", "// Object provides access control for Vanadium objects.\n//\n// Vanadium services implementing dynamic access control would typically embed\n// this interface and tag additional methods defined by the service with one of\n// Admin, Read, Write, Resolve etc. For example, the VDL definition of the\n// object would be:\n//\n//   package mypackage\n//\n//   import \"v.io/v23/security/access\"\n//   import \"v.io/v23/services/permissions\"\n//\n//   type MyObject interface {\n//     permissions.Object\n//     MyRead() (string, error) {access.Read}\n//     MyWrite(string) error    {access.Write}\n//   }\n//\n// If the set of pre-defined tags is insufficient, services may define their\n// own tag type and annotate all methods with this new type.\n//\n// Instead of embedding this Object interface, define SetPermissions and\n// GetPermissions in their own interface. Authorization policies will typically\n// respect annotations of a single type. For example, the VDL definition of an\n// object would be:\n//\n//  package mypackage\n//\n//  import \"v.io/v23/security/access\"\n//\n//  type MyTag string\n//\n//  const (\n//    Blue = MyTag(\"Blue\")\n//    Red  = MyTag(\"Red\")\n//  )\n//\n//  type MyObject interface {\n//    MyMethod() (string, error) {Blue}\n//\n//    // Allow clients to change access via the access.Object interface:\n//    SetPermissions(perms access.Permissions, version string) error         {Red}\n//    GetPermissions() (perms access.Permissions, version string, err error) {Blue}\n//  }"},
	},
}

// AppClientMethods is the client interface
// containing App methods.
//
// App represents the data for a specific app instance (possibly a combination
// of user, device, and app).
// App.Glob operates over Database names.
type AppClientMethods interface {
	// Object provides access control for Vanadium objects.
	//
	// Vanadium services implementing dynamic access control would typically embed
	// this interface and tag additional methods defined by the service with one of
	// Admin, Read, Write, Resolve etc. For example, the VDL definition of the
	// object would be:
	//
	//   package mypackage
	//
	//   import "v.io/v23/security/access"
	//   import "v.io/v23/services/permissions"
	//
	//   type MyObject interface {
	//     permissions.Object
	//     MyRead() (string, error) {access.Read}
	//     MyWrite(string) error    {access.Write}
	//   }
	//
	// If the set of pre-defined tags is insufficient, services may define their
	// own tag type and annotate all methods with this new type.
	//
	// Instead of embedding this Object interface, define SetPermissions and
	// GetPermissions in their own interface. Authorization policies will typically
	// respect annotations of a single type. For example, the VDL definition of an
	// object would be:
	//
	//  package mypackage
	//
	//  import "v.io/v23/security/access"
	//
	//  type MyTag string
	//
	//  const (
	//    Blue = MyTag("Blue")
	//    Red  = MyTag("Red")
	//  )
	//
	//  type MyObject interface {
	//    MyMethod() (string, error) {Blue}
	//
	//    // Allow clients to change access via the access.Object interface:
	//    SetPermissions(perms access.Permissions, version string) error         {Red}
	//    GetPermissions() (perms access.Permissions, version string, err error) {Blue}
	//  }
	permissions.ObjectClientMethods
	// Create creates this App.
	// If perms is nil, we inherit (copy) the Service perms.
	// Create requires the caller to have Write permission at the Service.
	Create(ctx *context.T, perms access.Permissions, opts ...rpc.CallOpt) error
	// Delete deletes this App.
	Delete(*context.T, ...rpc.CallOpt) error
}

// AppClientStub adds universal methods to AppClientMethods.
type AppClientStub interface {
	AppClientMethods
	rpc.UniversalServiceMethods
}

// AppClient returns a client stub for App.
func AppClient(name string) AppClientStub {
	return implAppClientStub{name, permissions.ObjectClient(name)}
}

type implAppClientStub struct {
	name string

	permissions.ObjectClientStub
}

func (c implAppClientStub) Create(ctx *context.T, i0 access.Permissions, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Create", []interface{}{i0}, nil, opts...)
	return
}

func (c implAppClientStub) Delete(ctx *context.T, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Delete", nil, nil, opts...)
	return
}

// AppServerMethods is the interface a server writer
// implements for App.
//
// App represents the data for a specific app instance (possibly a combination
// of user, device, and app).
// App.Glob operates over Database names.
type AppServerMethods interface {
	// Object provides access control for Vanadium objects.
	//
	// Vanadium services implementing dynamic access control would typically embed
	// this interface and tag additional methods defined by the service with one of
	// Admin, Read, Write, Resolve etc. For example, the VDL definition of the
	// object would be:
	//
	//   package mypackage
	//
	//   import "v.io/v23/security/access"
	//   import "v.io/v23/services/permissions"
	//
	//   type MyObject interface {
	//     permissions.Object
	//     MyRead() (string, error) {access.Read}
	//     MyWrite(string) error    {access.Write}
	//   }
	//
	// If the set of pre-defined tags is insufficient, services may define their
	// own tag type and annotate all methods with this new type.
	//
	// Instead of embedding this Object interface, define SetPermissions and
	// GetPermissions in their own interface. Authorization policies will typically
	// respect annotations of a single type. For example, the VDL definition of an
	// object would be:
	//
	//  package mypackage
	//
	//  import "v.io/v23/security/access"
	//
	//  type MyTag string
	//
	//  const (
	//    Blue = MyTag("Blue")
	//    Red  = MyTag("Red")
	//  )
	//
	//  type MyObject interface {
	//    MyMethod() (string, error) {Blue}
	//
	//    // Allow clients to change access via the access.Object interface:
	//    SetPermissions(perms access.Permissions, version string) error         {Red}
	//    GetPermissions() (perms access.Permissions, version string, err error) {Blue}
	//  }
	permissions.ObjectServerMethods
	// Create creates this App.
	// If perms is nil, we inherit (copy) the Service perms.
	// Create requires the caller to have Write permission at the Service.
	Create(ctx *context.T, call rpc.ServerCall, perms access.Permissions) error
	// Delete deletes this App.
	Delete(*context.T, rpc.ServerCall) error
}

// AppServerStubMethods is the server interface containing
// App methods, as expected by rpc.Server.
// There is no difference between this interface and AppServerMethods
// since there are no streaming methods.
type AppServerStubMethods AppServerMethods

// AppServerStub adds universal methods to AppServerStubMethods.
type AppServerStub interface {
	AppServerStubMethods
	// Describe the App interfaces.
	Describe__() []rpc.InterfaceDesc
}

// AppServer returns a server stub for App.
// It converts an implementation of AppServerMethods into
// an object that may be used by rpc.Server.
func AppServer(impl AppServerMethods) AppServerStub {
	stub := implAppServerStub{
		impl:             impl,
		ObjectServerStub: permissions.ObjectServer(impl),
	}
	// Initialize GlobState; always check the stub itself first, to handle the
	// case where the user has the Glob method defined in their VDL source.
	if gs := rpc.NewGlobState(stub); gs != nil {
		stub.gs = gs
	} else if gs := rpc.NewGlobState(impl); gs != nil {
		stub.gs = gs
	}
	return stub
}

type implAppServerStub struct {
	impl AppServerMethods
	permissions.ObjectServerStub
	gs *rpc.GlobState
}

func (s implAppServerStub) Create(ctx *context.T, call rpc.ServerCall, i0 access.Permissions) error {
	return s.impl.Create(ctx, call, i0)
}

func (s implAppServerStub) Delete(ctx *context.T, call rpc.ServerCall) error {
	return s.impl.Delete(ctx, call)
}

func (s implAppServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implAppServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{AppDesc, permissions.ObjectDesc}
}

// AppDesc describes the App interface.
var AppDesc rpc.InterfaceDesc = descApp

// descApp hides the desc to keep godoc clean.
var descApp = rpc.InterfaceDesc{
	Name:    "App",
	PkgPath: "v.io/syncbase/v23/services/syncbase",
	Doc:     "// App represents the data for a specific app instance (possibly a combination\n// of user, device, and app).\n// App.Glob operates over Database names.",
	Embeds: []rpc.EmbedDesc{
		{"Object", "v.io/v23/services/permissions", "// Object provides access control for Vanadium objects.\n//\n// Vanadium services implementing dynamic access control would typically embed\n// this interface and tag additional methods defined by the service with one of\n// Admin, Read, Write, Resolve etc. For example, the VDL definition of the\n// object would be:\n//\n//   package mypackage\n//\n//   import \"v.io/v23/security/access\"\n//   import \"v.io/v23/services/permissions\"\n//\n//   type MyObject interface {\n//     permissions.Object\n//     MyRead() (string, error) {access.Read}\n//     MyWrite(string) error    {access.Write}\n//   }\n//\n// If the set of pre-defined tags is insufficient, services may define their\n// own tag type and annotate all methods with this new type.\n//\n// Instead of embedding this Object interface, define SetPermissions and\n// GetPermissions in their own interface. Authorization policies will typically\n// respect annotations of a single type. For example, the VDL definition of an\n// object would be:\n//\n//  package mypackage\n//\n//  import \"v.io/v23/security/access\"\n//\n//  type MyTag string\n//\n//  const (\n//    Blue = MyTag(\"Blue\")\n//    Red  = MyTag(\"Red\")\n//  )\n//\n//  type MyObject interface {\n//    MyMethod() (string, error) {Blue}\n//\n//    // Allow clients to change access via the access.Object interface:\n//    SetPermissions(perms access.Permissions, version string) error         {Red}\n//    GetPermissions() (perms access.Permissions, version string, err error) {Blue}\n//  }"},
	},
	Methods: []rpc.MethodDesc{
		{
			Name: "Create",
			Doc:  "// Create creates this App.\n// If perms is nil, we inherit (copy) the Service perms.\n// Create requires the caller to have Write permission at the Service.",
			InArgs: []rpc.ArgDesc{
				{"perms", ``}, // access.Permissions
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
		{
			Name: "Delete",
			Doc:  "// Delete deletes this App.",
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
	},
}
