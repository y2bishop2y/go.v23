// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: syncbase

// Package syncbase defines the wire API for a structured store that supports
// peer-to-peer synchronization.
//
// TODO(sadovsky): Write a detailed package description, or provide a reference
// to the Syncbase documentation.
package syncbase

import (
	"fmt"
	"time"
	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/i18n"
	"v.io/v23/rpc"
	"v.io/v23/security/access"
	"v.io/v23/services/permissions"
	"v.io/v23/vdl"
	time_2 "v.io/v23/vdlroot/time"
	"v.io/v23/verror"
)

// DevModeUpdateVClockOpts specifies what DevModeUpdateVClock should do, as
// described below.
type DevModeUpdateVClockOpts struct {
	// If specified, sets the NTP host to talk to for subsequent NTP requests.
	NtpHost string
	// If Now is specified, the fake system clock is updated to the given values
	// of Now and ElapsedTime. If Now is not specified (i.e. takes the zero
	// value), the system clock is not touched by DevModeUpdateVClock.
	Now         time.Time
	ElapsedTime time.Duration
	// If specified, the clock daemon's local and/or NTP update code is triggered
	// after applying the updates specified by the fields above. (Helpful because
	// otherwise these only run periodically.) These functions work even if the
	// clock daemon hasn't been started.
	DoNtpUpdate   bool
	DoLocalUpdate bool
}

func (DevModeUpdateVClockOpts) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/syncbase.DevModeUpdateVClockOpts"`
}) {
}

func (m *DevModeUpdateVClockOpts) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	__VDLEnsureNativeBuilt()
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("NtpHost")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget3.FromString(string(m.NtpHost), vdl.StringType); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	var wireValue4 time_2.Time
	if err := time_2.TimeFromNative(&wireValue4, m.Now); err != nil {
		return err
	}

	keyTarget5, fieldTarget6, err := fieldsTarget1.StartField("Now")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		if err := wireValue4.FillVDLTarget(fieldTarget6, __VDLType_time_Time); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget5, fieldTarget6); err != nil {
			return err
		}
	}
	var wireValue7 time_2.Duration
	if err := time_2.DurationFromNative(&wireValue7, m.ElapsedTime); err != nil {
		return err
	}

	keyTarget8, fieldTarget9, err := fieldsTarget1.StartField("ElapsedTime")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		if err := wireValue7.FillVDLTarget(fieldTarget9, __VDLType_time_Duration); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget8, fieldTarget9); err != nil {
			return err
		}
	}
	keyTarget10, fieldTarget11, err := fieldsTarget1.StartField("DoNtpUpdate")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget11.FromBool(bool(m.DoNtpUpdate), vdl.BoolType); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget10, fieldTarget11); err != nil {
			return err
		}
	}
	keyTarget12, fieldTarget13, err := fieldsTarget1.StartField("DoLocalUpdate")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget13.FromBool(bool(m.DoLocalUpdate), vdl.BoolType); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget12, fieldTarget13); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *DevModeUpdateVClockOpts) MakeVDLTarget() vdl.Target {
	return &DevModeUpdateVClockOptsTarget{Value: m}
}

type DevModeUpdateVClockOptsTarget struct {
	Value               *DevModeUpdateVClockOpts
	ntpHostTarget       vdl.StringTarget
	nowTarget           time_2.TimeTarget
	elapsedTimeTarget   time_2.DurationTarget
	doNtpUpdateTarget   vdl.BoolTarget
	doLocalUpdateTarget vdl.BoolTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *DevModeUpdateVClockOptsTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if !vdl.Compatible(tt, __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts)
	}
	return t, nil
}
func (t *DevModeUpdateVClockOptsTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "NtpHost":
		t.ntpHostTarget.Value = &t.Value.NtpHost
		target, err := &t.ntpHostTarget, error(nil)
		return nil, target, err
	case "Now":
		t.nowTarget.Value = &t.Value.Now
		target, err := &t.nowTarget, error(nil)
		return nil, target, err
	case "ElapsedTime":
		t.elapsedTimeTarget.Value = &t.Value.ElapsedTime
		target, err := &t.elapsedTimeTarget, error(nil)
		return nil, target, err
	case "DoNtpUpdate":
		t.doNtpUpdateTarget.Value = &t.Value.DoNtpUpdate
		target, err := &t.doNtpUpdateTarget, error(nil)
		return nil, target, err
	case "DoLocalUpdate":
		t.doLocalUpdateTarget.Value = &t.Value.DoLocalUpdate
		target, err := &t.doLocalUpdateTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct %v", name, __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts)
	}
}
func (t *DevModeUpdateVClockOptsTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *DevModeUpdateVClockOptsTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

func init() {
	vdl.Register((*DevModeUpdateVClockOpts)(nil))
}

var __VDLType0 *vdl.Type

func __VDLType0_gen() *vdl.Type {
	__VDLType0Builder := vdl.TypeBuilder{}

	__VDLType01 := __VDLType0Builder.Optional()
	__VDLType02 := __VDLType0Builder.Struct()
	__VDLType03 := __VDLType0Builder.Named("v.io/v23/services/syncbase.DevModeUpdateVClockOpts").AssignBase(__VDLType02)
	__VDLType04 := vdl.StringType
	__VDLType02.AppendField("NtpHost", __VDLType04)
	__VDLType05 := __VDLType0Builder.Struct()
	__VDLType06 := __VDLType0Builder.Named("time.Time").AssignBase(__VDLType05)
	__VDLType07 := vdl.Int64Type
	__VDLType05.AppendField("Seconds", __VDLType07)
	__VDLType08 := vdl.Int32Type
	__VDLType05.AppendField("Nanos", __VDLType08)
	__VDLType02.AppendField("Now", __VDLType06)
	__VDLType09 := __VDLType0Builder.Struct()
	__VDLType010 := __VDLType0Builder.Named("time.Duration").AssignBase(__VDLType09)
	__VDLType09.AppendField("Seconds", __VDLType07)
	__VDLType09.AppendField("Nanos", __VDLType08)
	__VDLType02.AppendField("ElapsedTime", __VDLType010)
	__VDLType011 := vdl.BoolType
	__VDLType02.AppendField("DoNtpUpdate", __VDLType011)
	__VDLType02.AppendField("DoLocalUpdate", __VDLType011)
	__VDLType01.AssignElem(__VDLType03)
	__VDLType0Builder.Build()
	__VDLType0v, err := __VDLType01.Built()
	if err != nil {
		panic(err)
	}
	return __VDLType0v
}
func init() {
	__VDLType0 = __VDLType0_gen()
}

var __VDLType_time_Duration *vdl.Type

func __VDLType_time_Duration_gen() *vdl.Type {
	__VDLType_time_DurationBuilder := vdl.TypeBuilder{}

	__VDLType_time_Duration1 := __VDLType_time_DurationBuilder.Struct()
	__VDLType_time_Duration2 := __VDLType_time_DurationBuilder.Named("time.Duration").AssignBase(__VDLType_time_Duration1)
	__VDLType_time_Duration3 := vdl.Int64Type
	__VDLType_time_Duration1.AppendField("Seconds", __VDLType_time_Duration3)
	__VDLType_time_Duration4 := vdl.Int32Type
	__VDLType_time_Duration1.AppendField("Nanos", __VDLType_time_Duration4)
	__VDLType_time_DurationBuilder.Build()
	__VDLType_time_Durationv, err := __VDLType_time_Duration2.Built()
	if err != nil {
		panic(err)
	}
	return __VDLType_time_Durationv
}
func init() {
	__VDLType_time_Duration = __VDLType_time_Duration_gen()
}

var __VDLType_time_Time *vdl.Type

func __VDLType_time_Time_gen() *vdl.Type {
	__VDLType_time_TimeBuilder := vdl.TypeBuilder{}

	__VDLType_time_Time1 := __VDLType_time_TimeBuilder.Struct()
	__VDLType_time_Time2 := __VDLType_time_TimeBuilder.Named("time.Time").AssignBase(__VDLType_time_Time1)
	__VDLType_time_Time3 := vdl.Int64Type
	__VDLType_time_Time1.AppendField("Seconds", __VDLType_time_Time3)
	__VDLType_time_Time4 := vdl.Int32Type
	__VDLType_time_Time1.AppendField("Nanos", __VDLType_time_Time4)
	__VDLType_time_TimeBuilder.Build()
	__VDLType_time_Timev, err := __VDLType_time_Time2.Built()
	if err != nil {
		panic(err)
	}
	return __VDLType_time_Timev
}
func init() {
	__VDLType_time_Time = __VDLType_time_Time_gen()
}

var __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts *vdl.Type

func __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts_gen() *vdl.Type {
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsBuilder := vdl.TypeBuilder{}

	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts1 := __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsBuilder.Struct()
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts2 := __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsBuilder.Named("v.io/v23/services/syncbase.DevModeUpdateVClockOpts").AssignBase(__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts1)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts3 := vdl.StringType
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts1.AppendField("NtpHost", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts3)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts4 := __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsBuilder.Struct()
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts5 := __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsBuilder.Named("time.Time").AssignBase(__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts4)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts6 := vdl.Int64Type
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts4.AppendField("Seconds", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts6)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts7 := vdl.Int32Type
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts4.AppendField("Nanos", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts7)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts1.AppendField("Now", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts5)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts8 := __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsBuilder.Struct()
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts9 := __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsBuilder.Named("time.Duration").AssignBase(__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts8)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts8.AppendField("Seconds", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts6)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts8.AppendField("Nanos", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts7)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts1.AppendField("ElapsedTime", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts9)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts10 := vdl.BoolType
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts1.AppendField("DoNtpUpdate", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts10)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts1.AppendField("DoLocalUpdate", __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts10)
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsBuilder.Build()
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsv, err := __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts2.Built()
	if err != nil {
		panic(err)
	}
	return __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOptsv
}
func init() {
	__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts = __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts_gen()
}
func __VDLEnsureNativeBuilt() {
	if __VDLType0 == nil {
		__VDLType0 = __VDLType0_gen()
	}
	if __VDLType_time_Duration == nil {
		__VDLType_time_Duration = __VDLType_time_Duration_gen()
	}
	if __VDLType_time_Time == nil {
		__VDLType_time_Time = __VDLType_time_Time_gen()
	}
	if __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts == nil {
		__VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts = __VDLType_v_io_v23_services_syncbase_DevModeUpdateVClockOpts_gen()
	}
}

var (
	ErrNotInDevMode    = verror.Register("v.io/v23/services/syncbase.NotInDevMode", verror.NoRetry, "{1:}{2:} not running with --dev=true")
	ErrInvalidName     = verror.Register("v.io/v23/services/syncbase.InvalidName", verror.NoRetry, "{1:}{2:} invalid name: {3}")
	ErrCorruptDatabase = verror.Register("v.io/v23/services/syncbase.CorruptDatabase", verror.NoRetry, "{1:}{2:} database corrupt, moved to {3}; client must create a new database")
	ErrUnknownBatch    = verror.Register("v.io/v23/services/syncbase.UnknownBatch", verror.NoRetry, "{1:}{2:} unknown batch, perhaps the server restarted")
)

func init() {
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNotInDevMode.ID), "{1:}{2:} not running with --dev=true")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrInvalidName.ID), "{1:}{2:} invalid name: {3}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrCorruptDatabase.ID), "{1:}{2:} database corrupt, moved to {3}; client must create a new database")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrUnknownBatch.ID), "{1:}{2:} unknown batch, perhaps the server restarted")
}

// NewErrNotInDevMode returns an error with the ErrNotInDevMode ID.
func NewErrNotInDevMode(ctx *context.T) error {
	return verror.New(ErrNotInDevMode, ctx)
}

// NewErrInvalidName returns an error with the ErrInvalidName ID.
func NewErrInvalidName(ctx *context.T, name string) error {
	return verror.New(ErrInvalidName, ctx, name)
}

// NewErrCorruptDatabase returns an error with the ErrCorruptDatabase ID.
func NewErrCorruptDatabase(ctx *context.T, path string) error {
	return verror.New(ErrCorruptDatabase, ctx, path)
}

// NewErrUnknownBatch returns an error with the ErrUnknownBatch ID.
func NewErrUnknownBatch(ctx *context.T) error {
	return verror.New(ErrUnknownBatch, ctx)
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
	// DevModeUpdateVClock updates various bits of Syncbase virtual clock and clock
	// daemon state based on the specified options.
	// Requires --dev flag to be set (in addition to Admin check).
	DevModeUpdateVClock(_ *context.T, uco DevModeUpdateVClockOpts, _ ...rpc.CallOpt) error
	// DevModeGetTime returns the current time per the Syncbase clock.
	// Requires --dev flag to be set (in addition to Admin check).
	DevModeGetTime(*context.T, ...rpc.CallOpt) (time.Time, error)
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

func (c implServiceClientStub) DevModeUpdateVClock(ctx *context.T, i0 DevModeUpdateVClockOpts, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "DevModeUpdateVClock", []interface{}{i0}, nil, opts...)
	return
}

func (c implServiceClientStub) DevModeGetTime(ctx *context.T, opts ...rpc.CallOpt) (o0 time.Time, err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "DevModeGetTime", nil, []interface{}{&o0}, opts...)
	return
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
	// DevModeUpdateVClock updates various bits of Syncbase virtual clock and clock
	// daemon state based on the specified options.
	// Requires --dev flag to be set (in addition to Admin check).
	DevModeUpdateVClock(_ *context.T, _ rpc.ServerCall, uco DevModeUpdateVClockOpts) error
	// DevModeGetTime returns the current time per the Syncbase clock.
	// Requires --dev flag to be set (in addition to Admin check).
	DevModeGetTime(*context.T, rpc.ServerCall) (time.Time, error)
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

func (s implServiceServerStub) DevModeUpdateVClock(ctx *context.T, call rpc.ServerCall, i0 DevModeUpdateVClockOpts) error {
	return s.impl.DevModeUpdateVClock(ctx, call, i0)
}

func (s implServiceServerStub) DevModeGetTime(ctx *context.T, call rpc.ServerCall) (time.Time, error) {
	return s.impl.DevModeGetTime(ctx, call)
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
	PkgPath: "v.io/v23/services/syncbase",
	Doc:     "// Service represents a Vanadium Syncbase service.\n// Service.Glob operates over App names.",
	Embeds: []rpc.EmbedDesc{
		{"Object", "v.io/v23/services/permissions", "// Object provides access control for Vanadium objects.\n//\n// Vanadium services implementing dynamic access control would typically embed\n// this interface and tag additional methods defined by the service with one of\n// Admin, Read, Write, Resolve etc. For example, the VDL definition of the\n// object would be:\n//\n//   package mypackage\n//\n//   import \"v.io/v23/security/access\"\n//   import \"v.io/v23/services/permissions\"\n//\n//   type MyObject interface {\n//     permissions.Object\n//     MyRead() (string, error) {access.Read}\n//     MyWrite(string) error    {access.Write}\n//   }\n//\n// If the set of pre-defined tags is insufficient, services may define their\n// own tag type and annotate all methods with this new type.\n//\n// Instead of embedding this Object interface, define SetPermissions and\n// GetPermissions in their own interface. Authorization policies will typically\n// respect annotations of a single type. For example, the VDL definition of an\n// object would be:\n//\n//  package mypackage\n//\n//  import \"v.io/v23/security/access\"\n//\n//  type MyTag string\n//\n//  const (\n//    Blue = MyTag(\"Blue\")\n//    Red  = MyTag(\"Red\")\n//  )\n//\n//  type MyObject interface {\n//    MyMethod() (string, error) {Blue}\n//\n//    // Allow clients to change access via the access.Object interface:\n//    SetPermissions(perms access.Permissions, version string) error         {Red}\n//    GetPermissions() (perms access.Permissions, version string, err error) {Blue}\n//  }"},
	},
	Methods: []rpc.MethodDesc{
		{
			Name: "DevModeUpdateVClock",
			Doc:  "// DevModeUpdateVClock updates various bits of Syncbase virtual clock and clock\n// daemon state based on the specified options.\n// Requires --dev flag to be set (in addition to Admin check).",
			InArgs: []rpc.ArgDesc{
				{"uco", ``}, // DevModeUpdateVClockOpts
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Admin"))},
		},
		{
			Name: "DevModeGetTime",
			Doc:  "// DevModeGetTime returns the current time per the Syncbase clock.\n// Requires --dev flag to be set (in addition to Admin check).",
			OutArgs: []rpc.ArgDesc{
				{"", ``}, // time.Time
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Admin"))},
		},
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
	Create(_ *context.T, perms access.Permissions, _ ...rpc.CallOpt) error
	// Destroy destroys this App.
	Destroy(*context.T, ...rpc.CallOpt) error
	// Exists returns true only if this App exists. Insufficient permissions
	// cause Exists to return false instead of an error.
	Exists(*context.T, ...rpc.CallOpt) (bool, error)
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

func (c implAppClientStub) Destroy(ctx *context.T, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Destroy", nil, nil, opts...)
	return
}

func (c implAppClientStub) Exists(ctx *context.T, opts ...rpc.CallOpt) (o0 bool, err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Exists", nil, []interface{}{&o0}, opts...)
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
	Create(_ *context.T, _ rpc.ServerCall, perms access.Permissions) error
	// Destroy destroys this App.
	Destroy(*context.T, rpc.ServerCall) error
	// Exists returns true only if this App exists. Insufficient permissions
	// cause Exists to return false instead of an error.
	Exists(*context.T, rpc.ServerCall) (bool, error)
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

func (s implAppServerStub) Destroy(ctx *context.T, call rpc.ServerCall) error {
	return s.impl.Destroy(ctx, call)
}

func (s implAppServerStub) Exists(ctx *context.T, call rpc.ServerCall) (bool, error) {
	return s.impl.Exists(ctx, call)
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
	PkgPath: "v.io/v23/services/syncbase",
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
			Name: "Destroy",
			Doc:  "// Destroy destroys this App.",
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
		{
			Name: "Exists",
			Doc:  "// Exists returns true only if this App exists. Insufficient permissions\n// cause Exists to return false instead of an error.",
			OutArgs: []rpc.ArgDesc{
				{"", ``}, // bool
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Resolve"))},
		},
	},
}
