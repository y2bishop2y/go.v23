// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: groups

// Package groups defines interfaces for managing access control groups.  Groups
// can be referenced by BlessingPatterns (e.g. in AccessLists).
package groups

import (
	"fmt"
	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/i18n"
	"v.io/v23/rpc"
	"v.io/v23/security/access"
	"v.io/v23/services/permissions"
	"v.io/v23/vdl"
	"v.io/v23/verror"
)

// BlessingPatternChunk is a substring of a BlessingPattern. As with
// BlessingPatterns, BlessingPatternChunks may contain references to groups.
// However, they may be restricted in other ways. For example, in the future
// BlessingPatterns may support "$" terminators, but these may be disallowed for
// BlessingPatternChunks.
type BlessingPatternChunk string

func (BlessingPatternChunk) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/groups.BlessingPatternChunk"`
}) {
}

func (m *BlessingPatternChunk) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if err := t.FromString(string((*m)), __VDLType_v_io_v23_services_groups_BlessingPatternChunk); err != nil {
		return err
	}
	return nil
}

func (m *BlessingPatternChunk) MakeVDLTarget() vdl.Target {
	return &BlessingPatternChunkTarget{Value: m}
}

type BlessingPatternChunkTarget struct {
	Value *BlessingPatternChunk
	vdl.TargetBase
}

func (t *BlessingPatternChunkTarget) FromString(src string, tt *vdl.Type) error {

	if !vdl.Compatible(tt, __VDLType_v_io_v23_services_groups_BlessingPatternChunk) {
		return fmt.Errorf("type %v incompatible with %v", tt, __VDLType_v_io_v23_services_groups_BlessingPatternChunk)
	}
	*t.Value = BlessingPatternChunk(src)

	return nil
}

type GetRequest struct {
}

func (GetRequest) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/groups.GetRequest"`
}) {
}

func (m *GetRequest) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if __VDLType_v_io_v23_services_groups_GetRequest == nil || __VDLType0 == nil {
		panic("Initialization order error: types generated for FillVDLTarget not initialized. Consider moving caller to an init() block.")
	}
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *GetRequest) MakeVDLTarget() vdl.Target {
	return &GetRequestTarget{Value: m}
}

type GetRequestTarget struct {
	Value *GetRequest
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *GetRequestTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if !vdl.Compatible(tt, __VDLType_v_io_v23_services_groups_GetRequest) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, __VDLType_v_io_v23_services_groups_GetRequest)
	}
	return t, nil
}
func (t *GetRequestTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	default:
		return nil, nil, fmt.Errorf("field %s not in struct %v", name, __VDLType_v_io_v23_services_groups_GetRequest)
	}
}
func (t *GetRequestTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *GetRequestTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

type GetResponse struct {
	Entries map[BlessingPatternChunk]struct{}
}

func (GetResponse) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/groups.GetResponse"`
}) {
}

func (m *GetResponse) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if __VDLType_v_io_v23_services_groups_GetResponse == nil || __VDLType1 == nil {
		panic("Initialization order error: types generated for FillVDLTarget not initialized. Consider moving caller to an init() block.")
	}
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Entries")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {

		setTarget4, err := fieldTarget3.StartSet(__VDLType2, len(m.Entries))
		if err != nil {
			return err
		}
		for key6 := range m.Entries {
			keyTarget5, err := setTarget4.StartKey()
			if err != nil {
				return err
			}

			if err := key6.FillVDLTarget(keyTarget5, __VDLType_v_io_v23_services_groups_BlessingPatternChunk); err != nil {
				return err
			}
			if err := setTarget4.FinishKey(keyTarget5); err != nil {
				return err
			}
		}
		if err := fieldTarget3.FinishSet(setTarget4); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *GetResponse) MakeVDLTarget() vdl.Target {
	return &GetResponseTarget{Value: m}
}

type GetResponseTarget struct {
	Value         *GetResponse
	entriesTarget unnamed_7365745b762e696f2f7632332f73657276696365732f67726f7570732e426c657373696e675061747465726e4368756e6b20737472696e675dTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *GetResponseTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if !vdl.Compatible(tt, __VDLType_v_io_v23_services_groups_GetResponse) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, __VDLType_v_io_v23_services_groups_GetResponse)
	}
	return t, nil
}
func (t *GetResponseTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Entries":
		t.entriesTarget.Value = &t.Value.Entries
		target, err := &t.entriesTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct %v", name, __VDLType_v_io_v23_services_groups_GetResponse)
	}
}
func (t *GetResponseTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *GetResponseTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

// map[BlessingPatternChunk]struct{}
type unnamed_7365745b762e696f2f7632332f73657276696365732f67726f7570732e426c657373696e675061747465726e4368756e6b20737472696e675dTarget struct {
	Value     *map[BlessingPatternChunk]struct{}
	currKey   BlessingPatternChunk
	keyTarget BlessingPatternChunkTarget
	vdl.TargetBase
	vdl.SetTargetBase
}

func (t *unnamed_7365745b762e696f2f7632332f73657276696365732f67726f7570732e426c657373696e675061747465726e4368756e6b20737472696e675dTarget) StartSet(tt *vdl.Type, len int) (vdl.SetTarget, error) {

	if !vdl.Compatible(tt, __VDLType2) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, __VDLType2)
	}
	*t.Value = make(map[BlessingPatternChunk]struct{})
	return t, nil
}
func (t *unnamed_7365745b762e696f2f7632332f73657276696365732f67726f7570732e426c657373696e675061747465726e4368756e6b20737472696e675dTarget) StartKey() (key vdl.Target, _ error) {
	t.currKey = BlessingPatternChunk("")
	t.keyTarget.Value = &t.currKey
	target, err := &t.keyTarget, error(nil)
	return target, err
}
func (t *unnamed_7365745b762e696f2f7632332f73657276696365732f67726f7570732e426c657373696e675061747465726e4368756e6b20737472696e675dTarget) FinishKey(key vdl.Target) error {
	(*t.Value)[t.currKey] = struct{}{}
	return nil
}
func (t *unnamed_7365745b762e696f2f7632332f73657276696365732f67726f7570732e426c657373696e675061747465726e4368756e6b20737472696e675dTarget) FinishSet(list vdl.SetTarget) error {
	if len(*t.Value) == 0 {
		*t.Value = nil
	}

	return nil
}

// ApproximationType defines the type of approximation desired when a Relate
// call encounters an error (inaccessible or undefined group in a blessing
// pattern, cyclic group definitions, storage errors, invalid patterns
// etc). "Under" is used for blessing patterns in "Allow" clauses in an
// AccessList, while "Over" is used for blessing patterns in "Deny" clauses.
type ApproximationType int

const (
	ApproximationTypeUnder ApproximationType = iota
	ApproximationTypeOver
)

// ApproximationTypeAll holds all labels for ApproximationType.
var ApproximationTypeAll = [...]ApproximationType{ApproximationTypeUnder, ApproximationTypeOver}

// ApproximationTypeFromString creates a ApproximationType from a string label.
func ApproximationTypeFromString(label string) (x ApproximationType, err error) {
	err = x.Set(label)
	return
}

// Set assigns label to x.
func (x *ApproximationType) Set(label string) error {
	switch label {
	case "Under", "under":
		*x = ApproximationTypeUnder
		return nil
	case "Over", "over":
		*x = ApproximationTypeOver
		return nil
	}
	*x = -1
	return fmt.Errorf("unknown label %q in groups.ApproximationType", label)
}

// String returns the string label of x.
func (x ApproximationType) String() string {
	switch x {
	case ApproximationTypeUnder:
		return "Under"
	case ApproximationTypeOver:
		return "Over"
	}
	return ""
}

func (ApproximationType) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/groups.ApproximationType"`
	Enum struct{ Under, Over string }
}) {
}

func (m *ApproximationType) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if err := t.FromEnumLabel((*m).String(), __VDLType_v_io_v23_services_groups_ApproximationType); err != nil {
		return err
	}
	return nil
}

func (m *ApproximationType) MakeVDLTarget() vdl.Target {
	return &ApproximationTypeTarget{Value: m}
}

type ApproximationTypeTarget struct {
	Value *ApproximationType
	vdl.TargetBase
}

func (t *ApproximationTypeTarget) FromEnumLabel(src string, tt *vdl.Type) error {

	if !vdl.Compatible(tt, __VDLType_v_io_v23_services_groups_ApproximationType) {
		return fmt.Errorf("type %v incompatible with %v", tt, __VDLType_v_io_v23_services_groups_ApproximationType)
	}
	switch src {
	case "Under":
		*t.Value = 0
	case "Over":
		*t.Value = 1
	default:
		return fmt.Errorf("label %s not in enum %v", src, __VDLType_v_io_v23_services_groups_ApproximationType)
	}

	return nil
}

// Approximation contains information about membership approximations made
// during a Relate call.
type Approximation struct {
	Reason  string
	Details string
}

func (Approximation) __VDLReflect(struct {
	Name string `vdl:"v.io/v23/services/groups.Approximation"`
}) {
}

func (m *Approximation) FillVDLTarget(t vdl.Target, tt *vdl.Type) error {
	if __VDLType_v_io_v23_services_groups_Approximation == nil || __VDLType3 == nil {
		panic("Initialization order error: types generated for FillVDLTarget not initialized. Consider moving caller to an init() block.")
	}
	fieldsTarget1, err := t.StartFields(tt)
	if err != nil {
		return err
	}

	keyTarget2, fieldTarget3, err := fieldsTarget1.StartField("Reason")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget3.FromString(string(m.Reason), vdl.StringType); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget2, fieldTarget3); err != nil {
			return err
		}
	}
	keyTarget4, fieldTarget5, err := fieldsTarget1.StartField("Details")
	if err != vdl.ErrFieldNoExist && err != nil {
		return err
	}
	if err != vdl.ErrFieldNoExist {
		if err := fieldTarget5.FromString(string(m.Details), vdl.StringType); err != nil {
			return err
		}
		if err := fieldsTarget1.FinishField(keyTarget4, fieldTarget5); err != nil {
			return err
		}
	}
	if err := t.FinishFields(fieldsTarget1); err != nil {
		return err
	}
	return nil
}

func (m *Approximation) MakeVDLTarget() vdl.Target {
	return &ApproximationTarget{Value: m}
}

type ApproximationTarget struct {
	Value         *Approximation
	reasonTarget  vdl.StringTarget
	detailsTarget vdl.StringTarget
	vdl.TargetBase
	vdl.FieldsTargetBase
}

func (t *ApproximationTarget) StartFields(tt *vdl.Type) (vdl.FieldsTarget, error) {

	if !vdl.Compatible(tt, __VDLType_v_io_v23_services_groups_Approximation) {
		return nil, fmt.Errorf("type %v incompatible with %v", tt, __VDLType_v_io_v23_services_groups_Approximation)
	}
	return t, nil
}
func (t *ApproximationTarget) StartField(name string) (key, field vdl.Target, _ error) {
	switch name {
	case "Reason":
		t.reasonTarget.Value = &t.Value.Reason
		target, err := &t.reasonTarget, error(nil)
		return nil, target, err
	case "Details":
		t.detailsTarget.Value = &t.Value.Details
		target, err := &t.detailsTarget, error(nil)
		return nil, target, err
	default:
		return nil, nil, fmt.Errorf("field %s not in struct %v", name, __VDLType_v_io_v23_services_groups_Approximation)
	}
}
func (t *ApproximationTarget) FinishField(_, _ vdl.Target) error {
	return nil
}
func (t *ApproximationTarget) FinishFields(_ vdl.FieldsTarget) error {

	return nil
}

func init() {
	vdl.Register((*BlessingPatternChunk)(nil))
	vdl.Register((*GetRequest)(nil))
	vdl.Register((*GetResponse)(nil))
	vdl.Register((*ApproximationType)(nil))
	vdl.Register((*Approximation)(nil))
}

var __VDLType3 *vdl.Type = vdl.TypeOf((*Approximation)(nil))
var __VDLType0 *vdl.Type = vdl.TypeOf((*GetRequest)(nil))
var __VDLType1 *vdl.Type = vdl.TypeOf((*GetResponse)(nil))
var __VDLType2 *vdl.Type = vdl.TypeOf(map[BlessingPatternChunk]struct{}(nil))
var __VDLType_v_io_v23_services_groups_Approximation *vdl.Type = vdl.TypeOf(Approximation{})
var __VDLType_v_io_v23_services_groups_ApproximationType *vdl.Type = vdl.TypeOf(ApproximationTypeUnder)
var __VDLType_v_io_v23_services_groups_BlessingPatternChunk *vdl.Type = vdl.TypeOf(BlessingPatternChunk(""))
var __VDLType_v_io_v23_services_groups_GetRequest *vdl.Type = vdl.TypeOf(GetRequest{})
var __VDLType_v_io_v23_services_groups_GetResponse *vdl.Type = vdl.TypeOf(GetResponse{})

func __VDLEnsureNativeBuilt() {
}

var (
	ErrNoBlessings         = verror.Register("v.io/v23/services/groups.NoBlessings", verror.NoRetry, "{1:}{2:} No blessings recognized; cannot create group Permissions")
	ErrExcessiveContention = verror.Register("v.io/v23/services/groups.ExcessiveContention", verror.RetryBackoff, "{1:}{2:} Gave up after encountering excessive contention; try again later")
	ErrCycleFound          = verror.Register("v.io/v23/services/groups.CycleFound", verror.NoRetry, "{1:}{2:} Found cycle in group definitions{:_}")
)

func init() {
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNoBlessings.ID), "{1:}{2:} No blessings recognized; cannot create group Permissions")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrExcessiveContention.ID), "{1:}{2:} Gave up after encountering excessive contention; try again later")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrCycleFound.ID), "{1:}{2:} Found cycle in group definitions{:_}")
}

// NewErrNoBlessings returns an error with the ErrNoBlessings ID.
func NewErrNoBlessings(ctx *context.T) error {
	return verror.New(ErrNoBlessings, ctx)
}

// NewErrExcessiveContention returns an error with the ErrExcessiveContention ID.
func NewErrExcessiveContention(ctx *context.T) error {
	return verror.New(ErrExcessiveContention, ctx)
}

// NewErrCycleFound returns an error with the ErrCycleFound ID.
func NewErrCycleFound(ctx *context.T) error {
	return verror.New(ErrCycleFound, ctx)
}

// GroupReaderClientMethods is the client interface
// containing GroupReader methods.
//
// GroupReader implements methods to read or query a group's membership
// information.
type GroupReaderClientMethods interface {
	// Relate determines the relationships between the provided blessing
	// names and the members of the group.
	//
	// Given an input set of blessing names and a group defined by a set of
	// blessing patterns S, for each blessing name B in the input, Relate(B)
	// returns a set of "remainders" consisting of every blessing name B"
	// such that there exists some B' for which B = B' B" and B' is in S,
	// and "" if B is a member of S.
	//
	// For example, if a group is defined as S = {n1, n1:n2, n1:n2:n3}, then
	// Relate(n1:n2) = {n2, ""}.
	//
	// reqVersion specifies the expected version of the group's membership
	// information. If this version is set and matches the Group's current
	// version, the response will indicate that fact but will otherwise be
	// empty.
	//
	// visitedGroups is the set of groups already visited in a particular
	// chain of Relate calls, and is used to detect the presence of
	// cycles. When a cycle is detected, it is treated just like any other
	// error, and the result is approximated.
	//
	// Relate also returns information about all the errors encountered that
	// resulted in approximations, if any.
	//
	// TODO(hpucha): scrub "Approximation" for preserving privacy. Flesh
	// versioning out further. Other args we may need: option to Get() the
	// membership set when allowed (to avoid an extra RPC), options related
	// to caching this information.
	Relate(_ *context.T, blessings map[string]struct{}, hint ApproximationType, reqVersion string, visitedGroups map[string]struct{}, _ ...rpc.CallOpt) (remainder map[string]struct{}, approximations []Approximation, version string, _ error)
	// Get returns all entries in the group.
	// TODO(sadovsky): Flesh out this API.
	Get(_ *context.T, req GetRequest, reqVersion string, _ ...rpc.CallOpt) (res GetResponse, version string, _ error)
}

// GroupReaderClientStub adds universal methods to GroupReaderClientMethods.
type GroupReaderClientStub interface {
	GroupReaderClientMethods
	rpc.UniversalServiceMethods
}

// GroupReaderClient returns a client stub for GroupReader.
func GroupReaderClient(name string) GroupReaderClientStub {
	return implGroupReaderClientStub{name}
}

type implGroupReaderClientStub struct {
	name string
}

func (c implGroupReaderClientStub) Relate(ctx *context.T, i0 map[string]struct{}, i1 ApproximationType, i2 string, i3 map[string]struct{}, opts ...rpc.CallOpt) (o0 map[string]struct{}, o1 []Approximation, o2 string, err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Relate", []interface{}{i0, i1, i2, i3}, []interface{}{&o0, &o1, &o2}, opts...)
	return
}

func (c implGroupReaderClientStub) Get(ctx *context.T, i0 GetRequest, i1 string, opts ...rpc.CallOpt) (o0 GetResponse, o1 string, err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Get", []interface{}{i0, i1}, []interface{}{&o0, &o1}, opts...)
	return
}

// GroupReaderServerMethods is the interface a server writer
// implements for GroupReader.
//
// GroupReader implements methods to read or query a group's membership
// information.
type GroupReaderServerMethods interface {
	// Relate determines the relationships between the provided blessing
	// names and the members of the group.
	//
	// Given an input set of blessing names and a group defined by a set of
	// blessing patterns S, for each blessing name B in the input, Relate(B)
	// returns a set of "remainders" consisting of every blessing name B"
	// such that there exists some B' for which B = B' B" and B' is in S,
	// and "" if B is a member of S.
	//
	// For example, if a group is defined as S = {n1, n1:n2, n1:n2:n3}, then
	// Relate(n1:n2) = {n2, ""}.
	//
	// reqVersion specifies the expected version of the group's membership
	// information. If this version is set and matches the Group's current
	// version, the response will indicate that fact but will otherwise be
	// empty.
	//
	// visitedGroups is the set of groups already visited in a particular
	// chain of Relate calls, and is used to detect the presence of
	// cycles. When a cycle is detected, it is treated just like any other
	// error, and the result is approximated.
	//
	// Relate also returns information about all the errors encountered that
	// resulted in approximations, if any.
	//
	// TODO(hpucha): scrub "Approximation" for preserving privacy. Flesh
	// versioning out further. Other args we may need: option to Get() the
	// membership set when allowed (to avoid an extra RPC), options related
	// to caching this information.
	Relate(_ *context.T, _ rpc.ServerCall, blessings map[string]struct{}, hint ApproximationType, reqVersion string, visitedGroups map[string]struct{}) (remainder map[string]struct{}, approximations []Approximation, version string, _ error)
	// Get returns all entries in the group.
	// TODO(sadovsky): Flesh out this API.
	Get(_ *context.T, _ rpc.ServerCall, req GetRequest, reqVersion string) (res GetResponse, version string, _ error)
}

// GroupReaderServerStubMethods is the server interface containing
// GroupReader methods, as expected by rpc.Server.
// There is no difference between this interface and GroupReaderServerMethods
// since there are no streaming methods.
type GroupReaderServerStubMethods GroupReaderServerMethods

// GroupReaderServerStub adds universal methods to GroupReaderServerStubMethods.
type GroupReaderServerStub interface {
	GroupReaderServerStubMethods
	// Describe the GroupReader interfaces.
	Describe__() []rpc.InterfaceDesc
}

// GroupReaderServer returns a server stub for GroupReader.
// It converts an implementation of GroupReaderServerMethods into
// an object that may be used by rpc.Server.
func GroupReaderServer(impl GroupReaderServerMethods) GroupReaderServerStub {
	stub := implGroupReaderServerStub{
		impl: impl,
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

type implGroupReaderServerStub struct {
	impl GroupReaderServerMethods
	gs   *rpc.GlobState
}

func (s implGroupReaderServerStub) Relate(ctx *context.T, call rpc.ServerCall, i0 map[string]struct{}, i1 ApproximationType, i2 string, i3 map[string]struct{}) (map[string]struct{}, []Approximation, string, error) {
	return s.impl.Relate(ctx, call, i0, i1, i2, i3)
}

func (s implGroupReaderServerStub) Get(ctx *context.T, call rpc.ServerCall, i0 GetRequest, i1 string) (GetResponse, string, error) {
	return s.impl.Get(ctx, call, i0, i1)
}

func (s implGroupReaderServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implGroupReaderServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{GroupReaderDesc}
}

// GroupReaderDesc describes the GroupReader interface.
var GroupReaderDesc rpc.InterfaceDesc = descGroupReader

// descGroupReader hides the desc to keep godoc clean.
var descGroupReader = rpc.InterfaceDesc{
	Name:    "GroupReader",
	PkgPath: "v.io/v23/services/groups",
	Doc:     "// GroupReader implements methods to read or query a group's membership\n// information.",
	Methods: []rpc.MethodDesc{
		{
			Name: "Relate",
			Doc:  "// Relate determines the relationships between the provided blessing\n// names and the members of the group.\n//\n// Given an input set of blessing names and a group defined by a set of\n// blessing patterns S, for each blessing name B in the input, Relate(B)\n// returns a set of \"remainders\" consisting of every blessing name B\"\n// such that there exists some B' for which B = B' B\" and B' is in S,\n// and \"\" if B is a member of S.\n//\n// For example, if a group is defined as S = {n1, n1:n2, n1:n2:n3}, then\n// Relate(n1:n2) = {n2, \"\"}.\n//\n// reqVersion specifies the expected version of the group's membership\n// information. If this version is set and matches the Group's current\n// version, the response will indicate that fact but will otherwise be\n// empty.\n//\n// visitedGroups is the set of groups already visited in a particular\n// chain of Relate calls, and is used to detect the presence of\n// cycles. When a cycle is detected, it is treated just like any other\n// error, and the result is approximated.\n//\n// Relate also returns information about all the errors encountered that\n// resulted in approximations, if any.\n//\n// TODO(hpucha): scrub \"Approximation\" for preserving privacy. Flesh\n// versioning out further. Other args we may need: option to Get() the\n// membership set when allowed (to avoid an extra RPC), options related\n// to caching this information.",
			InArgs: []rpc.ArgDesc{
				{"blessings", ``},     // map[string]struct{}
				{"hint", ``},          // ApproximationType
				{"reqVersion", ``},    // string
				{"visitedGroups", ``}, // map[string]struct{}
			},
			OutArgs: []rpc.ArgDesc{
				{"remainder", ``},      // map[string]struct{}
				{"approximations", ``}, // []Approximation
				{"version", ``},        // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Resolve"))},
		},
		{
			Name: "Get",
			Doc:  "// Get returns all entries in the group.\n// TODO(sadovsky): Flesh out this API.",
			InArgs: []rpc.ArgDesc{
				{"req", ``},        // GetRequest
				{"reqVersion", ``}, // string
			},
			OutArgs: []rpc.ArgDesc{
				{"res", ``},     // GetResponse
				{"version", ``}, // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Read"))},
		},
	},
}

// GroupClientMethods is the client interface
// containing Group methods.
//
// A group's version covers its Permissions as well as any other data stored in
// the group. Clients should treat versions as opaque identifiers. For both Get
// and Relate, if version is set and matches the Group's current version, the
// response will indicate that fact but will otherwise be empty.
type GroupClientMethods interface {
	// GroupReader implements methods to read or query a group's membership
	// information.
	GroupReaderClientMethods
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
	// Create creates a new group if it doesn't already exist.
	// If perms is nil, a default Permissions is used, providing Admin access to
	// the caller.
	// Create requires the caller to have Write permission at the GroupServer.
	Create(_ *context.T, perms access.Permissions, entries []BlessingPatternChunk, _ ...rpc.CallOpt) error
	// Delete deletes the group.
	// Permissions for all group-related methods except Create() are checked
	// against the Group object.
	Delete(_ *context.T, version string, _ ...rpc.CallOpt) error
	// Add adds an entry to the group.
	Add(_ *context.T, entry BlessingPatternChunk, version string, _ ...rpc.CallOpt) error
	// Remove removes an entry from the group.
	Remove(_ *context.T, entry BlessingPatternChunk, version string, _ ...rpc.CallOpt) error
}

// GroupClientStub adds universal methods to GroupClientMethods.
type GroupClientStub interface {
	GroupClientMethods
	rpc.UniversalServiceMethods
}

// GroupClient returns a client stub for Group.
func GroupClient(name string) GroupClientStub {
	return implGroupClientStub{name, GroupReaderClient(name), permissions.ObjectClient(name)}
}

type implGroupClientStub struct {
	name string

	GroupReaderClientStub
	permissions.ObjectClientStub
}

func (c implGroupClientStub) Create(ctx *context.T, i0 access.Permissions, i1 []BlessingPatternChunk, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Create", []interface{}{i0, i1}, nil, opts...)
	return
}

func (c implGroupClientStub) Delete(ctx *context.T, i0 string, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Delete", []interface{}{i0}, nil, opts...)
	return
}

func (c implGroupClientStub) Add(ctx *context.T, i0 BlessingPatternChunk, i1 string, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Add", []interface{}{i0, i1}, nil, opts...)
	return
}

func (c implGroupClientStub) Remove(ctx *context.T, i0 BlessingPatternChunk, i1 string, opts ...rpc.CallOpt) (err error) {
	err = v23.GetClient(ctx).Call(ctx, c.name, "Remove", []interface{}{i0, i1}, nil, opts...)
	return
}

// GroupServerMethods is the interface a server writer
// implements for Group.
//
// A group's version covers its Permissions as well as any other data stored in
// the group. Clients should treat versions as opaque identifiers. For both Get
// and Relate, if version is set and matches the Group's current version, the
// response will indicate that fact but will otherwise be empty.
type GroupServerMethods interface {
	// GroupReader implements methods to read or query a group's membership
	// information.
	GroupReaderServerMethods
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
	// Create creates a new group if it doesn't already exist.
	// If perms is nil, a default Permissions is used, providing Admin access to
	// the caller.
	// Create requires the caller to have Write permission at the GroupServer.
	Create(_ *context.T, _ rpc.ServerCall, perms access.Permissions, entries []BlessingPatternChunk) error
	// Delete deletes the group.
	// Permissions for all group-related methods except Create() are checked
	// against the Group object.
	Delete(_ *context.T, _ rpc.ServerCall, version string) error
	// Add adds an entry to the group.
	Add(_ *context.T, _ rpc.ServerCall, entry BlessingPatternChunk, version string) error
	// Remove removes an entry from the group.
	Remove(_ *context.T, _ rpc.ServerCall, entry BlessingPatternChunk, version string) error
}

// GroupServerStubMethods is the server interface containing
// Group methods, as expected by rpc.Server.
// There is no difference between this interface and GroupServerMethods
// since there are no streaming methods.
type GroupServerStubMethods GroupServerMethods

// GroupServerStub adds universal methods to GroupServerStubMethods.
type GroupServerStub interface {
	GroupServerStubMethods
	// Describe the Group interfaces.
	Describe__() []rpc.InterfaceDesc
}

// GroupServer returns a server stub for Group.
// It converts an implementation of GroupServerMethods into
// an object that may be used by rpc.Server.
func GroupServer(impl GroupServerMethods) GroupServerStub {
	stub := implGroupServerStub{
		impl: impl,
		GroupReaderServerStub: GroupReaderServer(impl),
		ObjectServerStub:      permissions.ObjectServer(impl),
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

type implGroupServerStub struct {
	impl GroupServerMethods
	GroupReaderServerStub
	permissions.ObjectServerStub
	gs *rpc.GlobState
}

func (s implGroupServerStub) Create(ctx *context.T, call rpc.ServerCall, i0 access.Permissions, i1 []BlessingPatternChunk) error {
	return s.impl.Create(ctx, call, i0, i1)
}

func (s implGroupServerStub) Delete(ctx *context.T, call rpc.ServerCall, i0 string) error {
	return s.impl.Delete(ctx, call, i0)
}

func (s implGroupServerStub) Add(ctx *context.T, call rpc.ServerCall, i0 BlessingPatternChunk, i1 string) error {
	return s.impl.Add(ctx, call, i0, i1)
}

func (s implGroupServerStub) Remove(ctx *context.T, call rpc.ServerCall, i0 BlessingPatternChunk, i1 string) error {
	return s.impl.Remove(ctx, call, i0, i1)
}

func (s implGroupServerStub) Globber() *rpc.GlobState {
	return s.gs
}

func (s implGroupServerStub) Describe__() []rpc.InterfaceDesc {
	return []rpc.InterfaceDesc{GroupDesc, GroupReaderDesc, permissions.ObjectDesc}
}

// GroupDesc describes the Group interface.
var GroupDesc rpc.InterfaceDesc = descGroup

// descGroup hides the desc to keep godoc clean.
var descGroup = rpc.InterfaceDesc{
	Name:    "Group",
	PkgPath: "v.io/v23/services/groups",
	Doc:     "// A group's version covers its Permissions as well as any other data stored in\n// the group. Clients should treat versions as opaque identifiers. For both Get\n// and Relate, if version is set and matches the Group's current version, the\n// response will indicate that fact but will otherwise be empty.",
	Embeds: []rpc.EmbedDesc{
		{"GroupReader", "v.io/v23/services/groups", "// GroupReader implements methods to read or query a group's membership\n// information."},
		{"Object", "v.io/v23/services/permissions", "// Object provides access control for Vanadium objects.\n//\n// Vanadium services implementing dynamic access control would typically embed\n// this interface and tag additional methods defined by the service with one of\n// Admin, Read, Write, Resolve etc. For example, the VDL definition of the\n// object would be:\n//\n//   package mypackage\n//\n//   import \"v.io/v23/security/access\"\n//   import \"v.io/v23/services/permissions\"\n//\n//   type MyObject interface {\n//     permissions.Object\n//     MyRead() (string, error) {access.Read}\n//     MyWrite(string) error    {access.Write}\n//   }\n//\n// If the set of pre-defined tags is insufficient, services may define their\n// own tag type and annotate all methods with this new type.\n//\n// Instead of embedding this Object interface, define SetPermissions and\n// GetPermissions in their own interface. Authorization policies will typically\n// respect annotations of a single type. For example, the VDL definition of an\n// object would be:\n//\n//  package mypackage\n//\n//  import \"v.io/v23/security/access\"\n//\n//  type MyTag string\n//\n//  const (\n//    Blue = MyTag(\"Blue\")\n//    Red  = MyTag(\"Red\")\n//  )\n//\n//  type MyObject interface {\n//    MyMethod() (string, error) {Blue}\n//\n//    // Allow clients to change access via the access.Object interface:\n//    SetPermissions(perms access.Permissions, version string) error         {Red}\n//    GetPermissions() (perms access.Permissions, version string, err error) {Blue}\n//  }"},
	},
	Methods: []rpc.MethodDesc{
		{
			Name: "Create",
			Doc:  "// Create creates a new group if it doesn't already exist.\n// If perms is nil, a default Permissions is used, providing Admin access to\n// the caller.\n// Create requires the caller to have Write permission at the GroupServer.",
			InArgs: []rpc.ArgDesc{
				{"perms", ``},   // access.Permissions
				{"entries", ``}, // []BlessingPatternChunk
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
		{
			Name: "Delete",
			Doc:  "// Delete deletes the group.\n// Permissions for all group-related methods except Create() are checked\n// against the Group object.",
			InArgs: []rpc.ArgDesc{
				{"version", ``}, // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
		{
			Name: "Add",
			Doc:  "// Add adds an entry to the group.",
			InArgs: []rpc.ArgDesc{
				{"entry", ``},   // BlessingPatternChunk
				{"version", ``}, // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
		{
			Name: "Remove",
			Doc:  "// Remove removes an entry from the group.",
			InArgs: []rpc.ArgDesc{
				{"entry", ``},   // BlessingPatternChunk
				{"version", ``}, // string
			},
			Tags: []*vdl.Value{vdl.ValueOf(access.Tag("Write"))},
		},
	},
}
