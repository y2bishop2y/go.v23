// This file was auto-generated by the veyron vdl tool.
// Source: native.vdl

// Package native tests a package with native type conversions.
package native

import (
	// VDL system imports
	"v.io/core/veyron2/vdl"

	// VDL user imports
	"time"
)

type WireString int32

func (WireString) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vdl/testdata/native.WireString"
}) {
}

type WireMapStringInt int32

func (WireMapStringInt) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vdl/testdata/native.WireMapStringInt"
}) {
}

type WireTime int32

func (WireTime) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vdl/testdata/native.WireTime"
}) {
}

type WireSamePkg int32

func (WireSamePkg) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vdl/testdata/native.WireSamePkg"
}) {
}

type WireMultiImport int32

func (WireMultiImport) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vdl/testdata/native.WireMultiImport"
}) {
}

type WireAll struct {
	A string
	B map[string]int
	C time.Time
	D NativeSamePkg
	E map[NativeSamePkg]time.Time
}

func (WireAll) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vdl/testdata/native.WireAll"
}) {
}

// WireMapStringInt must implement native type conversions.
var _ interface {
	VDLToNative(*map[string]int) error
	VDLFromNative(map[string]int) error
} = (*WireMapStringInt)(nil)

// WireMultiImport must implement native type conversions.
var _ interface {
	VDLToNative(*map[NativeSamePkg]time.Time) error
	VDLFromNative(map[NativeSamePkg]time.Time) error
} = (*WireMultiImport)(nil)

// WireSamePkg must implement native type conversions.
var _ interface {
	VDLToNative(*NativeSamePkg) error
	VDLFromNative(NativeSamePkg) error
} = (*WireSamePkg)(nil)

// WireString must implement native type conversions.
var _ interface {
	VDLToNative(*string) error
	VDLFromNative(string) error
} = (*WireString)(nil)

// WireTime must implement native type conversions.
var _ interface {
	VDLToNative(*time.Time) error
	VDLFromNative(time.Time) error
} = (*WireTime)(nil)

func init() {
	vdl.Register((*WireString)(nil))
	vdl.Register((*WireMapStringInt)(nil))
	vdl.Register((*WireTime)(nil))
	vdl.Register((*WireSamePkg)(nil))
	vdl.Register((*WireMultiImport)(nil))
	vdl.Register((*WireAll)(nil))
}
