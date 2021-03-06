// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file was auto-generated by the vanadium vdl tool.
// Package: verror

package verror

import (
	"v.io/v23/context"
	"v.io/v23/i18n"
)

var _ = __VDLInit() // Must be first; see __VDLInit comments for details.

//////////////////////////////////////////////////
// Error definitions

var (

	// Unknown means the error has no known Id.  A more specific error should
	// always be used, if possible.  Unknown is typically only used when
	// automatically converting errors that do not contain an Id.
	ErrUnknown = Register("v.io/v23/verror.Unknown", NoRetry, "{1:}{2:} Error{:_}")
	// Internal means an internal error has occurred.  A more specific error
	// should always be used, if possible.
	ErrInternal = Register("v.io/v23/verror.Internal", NoRetry, "{1:}{2:} Internal error{:_}")
	// NotImplemented means that the request type is valid but that the method to
	// handle the request has not been implemented.
	ErrNotImplemented = Register("v.io/v23/verror.NotImplemented", NoRetry, "{1:}{2:} Not implemented{:_}")
	// EndOfFile means the end-of-file has been reached; more generally, no more
	// input data is available.
	ErrEndOfFile = Register("v.io/v23/verror.EndOfFile", NoRetry, "{1:}{2:} End of file{:_}")
	// BadArg means the arguments to an operation are invalid or incorrectly
	// formatted.
	ErrBadArg = Register("v.io/v23/verror.BadArg", NoRetry, "{1:}{2:} Bad argument{:_}")
	// BadState means an operation was attempted on an object while the object was
	// in an incompatible state.
	ErrBadState = Register("v.io/v23/verror.BadState", NoRetry, "{1:}{2:} Invalid state{:_}")
	// BadVersion means the version presented by the client (e.g. to a service
	// that supports content-hash-based caching or atomic read-modify-write) was
	// out of date or otherwise invalid, likely because some other request caused
	// the version at the server to change. The client should get a fresh version
	// and try again.
	ErrBadVersion = Register("v.io/v23/verror.BadVersion", NoRetry, "{1:}{2:} Version is out of date")
	// Exist means that the requested item already exists; typically returned when
	// an attempt to create an item fails because it already exists.
	ErrExist = Register("v.io/v23/verror.Exist", NoRetry, "{1:}{2:} Already exists{:_}")
	// NoExist means that the requested item does not exist; typically returned
	// when an attempt to lookup an item fails because it does not exist.
	ErrNoExist       = Register("v.io/v23/verror.NoExist", NoRetry, "{1:}{2:} Does not exist{:_}")
	ErrUnknownMethod = Register("v.io/v23/verror.UnknownMethod", NoRetry, "{1:}{2:} Method does not exist{:_}")
	ErrUnknownSuffix = Register("v.io/v23/verror.UnknownSuffix", NoRetry, "{1:}{2:} Suffix does not exist{:_}")
	// NoExistOrNoAccess means that either the requested item does not exist, or
	// is inaccessible.  Typically returned when the distinction between existence
	// and inaccessiblity should be hidden to preserve privacy.
	ErrNoExistOrNoAccess = Register("v.io/v23/verror.NoExistOrNoAccess", NoRetry, "{1:}{2:} Does not exist or access denied{:_}")
	// NoServers means a name was resolved to unusable or inaccessible servers.
	ErrNoServers = Register("v.io/v23/verror.NoServers", RetryRefetch, "{1:}{2:} No usable servers found{:_}")
	// NoAccess means the server does not authorize the client for access.
	ErrNoAccess = Register("v.io/v23/verror.NoAccess", RetryRefetch, "{1:}{2:} Access denied{:_}")
	// NotTrusted means the client does not trust the server.
	ErrNotTrusted = Register("v.io/v23/verror.NotTrusted", RetryRefetch, "{1:}{2:} Client does not trust server{:_}")
	// Aborted means that an operation was not completed because it was aborted by
	// the receiver.  A more specific error should be used if it would help the
	// caller decide how to proceed.
	ErrAborted = Register("v.io/v23/verror.Aborted", NoRetry, "{1:}{2:} Aborted{:_}")
	// BadProtocol means that an operation was not completed because of a protocol
	// or codec error.
	ErrBadProtocol = Register("v.io/v23/verror.BadProtocol", NoRetry, "{1:}{2:} Bad protocol or type{:_}")
	// Canceled means that an operation was not completed because it was
	// explicitly cancelled by the caller.
	ErrCanceled = Register("v.io/v23/verror.Canceled", NoRetry, "{1:}{2:} Canceled{:_}")
	// Timeout means that an operation was not completed before the time deadline
	// for the operation.
	ErrTimeout = Register("v.io/v23/verror.Timeout", NoRetry, "{1:}{2:} Timeout{:_}")
)

// NewErrUnknown returns an error with the ErrUnknown ID.
func NewErrUnknown(ctx *context.T) error {
	return New(ErrUnknown, ctx)
}

// NewErrInternal returns an error with the ErrInternal ID.
func NewErrInternal(ctx *context.T) error {
	return New(ErrInternal, ctx)
}

// NewErrNotImplemented returns an error with the ErrNotImplemented ID.
func NewErrNotImplemented(ctx *context.T) error {
	return New(ErrNotImplemented, ctx)
}

// NewErrEndOfFile returns an error with the ErrEndOfFile ID.
func NewErrEndOfFile(ctx *context.T) error {
	return New(ErrEndOfFile, ctx)
}

// NewErrBadArg returns an error with the ErrBadArg ID.
func NewErrBadArg(ctx *context.T) error {
	return New(ErrBadArg, ctx)
}

// NewErrBadState returns an error with the ErrBadState ID.
func NewErrBadState(ctx *context.T) error {
	return New(ErrBadState, ctx)
}

// NewErrBadVersion returns an error with the ErrBadVersion ID.
func NewErrBadVersion(ctx *context.T) error {
	return New(ErrBadVersion, ctx)
}

// NewErrExist returns an error with the ErrExist ID.
func NewErrExist(ctx *context.T) error {
	return New(ErrExist, ctx)
}

// NewErrNoExist returns an error with the ErrNoExist ID.
func NewErrNoExist(ctx *context.T) error {
	return New(ErrNoExist, ctx)
}

// NewErrUnknownMethod returns an error with the ErrUnknownMethod ID.
func NewErrUnknownMethod(ctx *context.T) error {
	return New(ErrUnknownMethod, ctx)
}

// NewErrUnknownSuffix returns an error with the ErrUnknownSuffix ID.
func NewErrUnknownSuffix(ctx *context.T) error {
	return New(ErrUnknownSuffix, ctx)
}

// NewErrNoExistOrNoAccess returns an error with the ErrNoExistOrNoAccess ID.
func NewErrNoExistOrNoAccess(ctx *context.T) error {
	return New(ErrNoExistOrNoAccess, ctx)
}

// NewErrNoServers returns an error with the ErrNoServers ID.
func NewErrNoServers(ctx *context.T) error {
	return New(ErrNoServers, ctx)
}

// NewErrNoAccess returns an error with the ErrNoAccess ID.
func NewErrNoAccess(ctx *context.T) error {
	return New(ErrNoAccess, ctx)
}

// NewErrNotTrusted returns an error with the ErrNotTrusted ID.
func NewErrNotTrusted(ctx *context.T) error {
	return New(ErrNotTrusted, ctx)
}

// NewErrAborted returns an error with the ErrAborted ID.
func NewErrAborted(ctx *context.T) error {
	return New(ErrAborted, ctx)
}

// NewErrBadProtocol returns an error with the ErrBadProtocol ID.
func NewErrBadProtocol(ctx *context.T) error {
	return New(ErrBadProtocol, ctx)
}

// NewErrCanceled returns an error with the ErrCanceled ID.
func NewErrCanceled(ctx *context.T) error {
	return New(ErrCanceled, ctx)
}

// NewErrTimeout returns an error with the ErrTimeout ID.
func NewErrTimeout(ctx *context.T) error {
	return New(ErrTimeout, ctx)
}

var __VDLInitCalled bool

// __VDLInit performs vdl initialization.  It is safe to call multiple times.
// If you have an init ordering issue, just insert the following line verbatim
// into your source files in this package, right after the "package foo" clause:
//
//    var _ = __VDLInit()
//
// The purpose of this function is to ensure that vdl initialization occurs in
// the right order, and very early in the init sequence.  In particular, vdl
// registration and package variable initialization needs to occur before
// functions like vdl.TypeOf will work properly.
//
// This function returns a dummy value, so that it can be used to initialize the
// first var in the file, to take advantage of Go's defined init order.
func __VDLInit() struct{} {
	if __VDLInitCalled {
		return struct{}{}
	}
	__VDLInitCalled = true

	// Set error format strings.
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrUnknown.ID), "{1:}{2:} Error{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrInternal.ID), "{1:}{2:} Internal error{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNotImplemented.ID), "{1:}{2:} Not implemented{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrEndOfFile.ID), "{1:}{2:} End of file{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrBadArg.ID), "{1:}{2:} Bad argument{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrBadState.ID), "{1:}{2:} Invalid state{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrBadVersion.ID), "{1:}{2:} Version is out of date")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrExist.ID), "{1:}{2:} Already exists{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNoExist.ID), "{1:}{2:} Does not exist{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrUnknownMethod.ID), "{1:}{2:} Method does not exist{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrUnknownSuffix.ID), "{1:}{2:} Suffix does not exist{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNoExistOrNoAccess.ID), "{1:}{2:} Does not exist or access denied{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNoServers.ID), "{1:}{2:} No usable servers found{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNoAccess.ID), "{1:}{2:} Access denied{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrNotTrusted.ID), "{1:}{2:} Client does not trust server{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrAborted.ID), "{1:}{2:} Aborted{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrBadProtocol.ID), "{1:}{2:} Bad protocol or type{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrCanceled.ID), "{1:}{2:} Canceled{:_}")
	i18n.Cat().SetWithBase(i18n.LangID("en"), i18n.MsgID(ErrTimeout.ID), "{1:}{2:} Timeout{:_}")

	return struct{}{}
}
