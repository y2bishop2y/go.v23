// This file was auto-generated by the veyron vdl tool.
// Source: time.vdl

// Package time provides a standard representation of time.
//
// The representations described below are required to provide wire
// compatibility between different programming environments.  Generated code for
// different environments typically provide automatic conversions into native
// representations, for simpler idiomatic usage.
package time

import (
	// VDL system imports
	"v.io/core/veyron2/vdl"

	// VDL user imports
	"time"
)

// Time represents an absolute point in time with nanosecond precision.
//
// Time is represented as the duration before or after a fixed epoch.  The zero
// Time represents the epoch 0001-01-01T00:00:00.000000000Z.  This uses the
// proleptic Gregorian calendar; the calendar runs on an exact 400 year cycle.
// Leap seconds are "smeared", ensuring that no leap second table is necessary
// for interpretation.
//
// This is similar to Go time.Time, but always in the UTC location.
// http://golang.org/pkg/time/#Time
//
// This is similar to conventional "unix time", but with the epoch defined at
// year 1 rather than year 1970.  This allows the zero Time to be used as a
// natural sentry, since it isn't a valid time for many practical applications.
// http://en.wikipedia.org/wiki/Unix_time
type Time struct {
	Seconds int64
	Nano    int32
}

func (Time) __VDLReflect(struct {
	Name string "time.Time"
}) {
}

// Duration represents the elapsed duration between two points in time, with
// nanosecond precision.
type Duration struct {
	// Seconds represents the seconds in the duration.  The range is roughly
	// +/-290 billion years, larger than the estimated age of the universe.
	Seconds int64
	// Nano represents the fractions of a second at nanosecond resolution.  Must
	// be in the inclusive range between +/-999,999,999.
	//
	// In normalized form, durations less than one second are represented with 0
	// Seconds and +/-Nanos.  For durations one second or more, the sign of Nanos
	// must match Seconds, or be 0.
	Nano int32
}

func (Duration) __VDLReflect(struct {
	Name string "time.Duration"
}) {
}

// Duration must implement native type conversions.
var _ interface {
	VDLToNative(*time.Duration) error
	VDLFromNative(time.Duration) error
} = (*Duration)(nil)

// Time must implement native type conversions.
var _ interface {
	VDLToNative(*time.Time) error
	VDLFromNative(time.Time) error
} = (*Time)(nil)

func init() {
	vdl.Register((*Time)(nil))
	vdl.Register((*Duration)(nil))
}
