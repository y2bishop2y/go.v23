// This file was auto-generated by the veyron vdl tool.
// Source: vomtype.vdl

package testdata

import (
	// VDL system imports
	"fmt"
	"v.io/core/veyron2/vdl"
)

// vomdata config types
type ConvertGroup struct {
	Name        string
	PrimaryType *vdl.Type
	Values      []vdl.AnyRep
}

func (ConvertGroup) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.ConvertGroup"
}) {
}

type VomdataStruct struct {
	EncodeDecodeData []vdl.AnyRep
	CompatData       map[string][]*vdl.Type
	ConvertData      map[string][]ConvertGroup
}

func (VomdataStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.VomdataStruct"
}) {
}

// Named Types
type NBool bool

func (NBool) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NBool"
}) {
}

type NString string

func (NString) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NString"
}) {
}

type NByteSlice []byte

func (NByteSlice) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NByteSlice"
}) {
}

type NByteArray [4]byte

func (NByteArray) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NByteArray"
}) {
}

type NByte byte

func (NByte) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NByte"
}) {
}

type NUint16 uint16

func (NUint16) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NUint16"
}) {
}

type NUint32 uint32

func (NUint32) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NUint32"
}) {
}

type NUint64 uint64

func (NUint64) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NUint64"
}) {
}

type NInt16 int16

func (NInt16) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NInt16"
}) {
}

type NInt32 int32

func (NInt32) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NInt32"
}) {
}

type NInt64 int64

func (NInt64) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NInt64"
}) {
}

type NFloat32 float32

func (NFloat32) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NFloat32"
}) {
}

type NFloat64 float64

func (NFloat64) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NFloat64"
}) {
}

type NComplex64 complex64

func (NComplex64) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NComplex64"
}) {
}

type NComplex128 complex128

func (NComplex128) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NComplex128"
}) {
}

type NArray2Uint64 [2]uint64

func (NArray2Uint64) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NArray2Uint64"
}) {
}

type NListUint64 []uint64

func (NListUint64) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NListUint64"
}) {
}

type NSetUint64 map[uint64]struct{}

func (NSetUint64) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NSetUint64"
}) {
}

type NMapUint64String map[uint64]string

func (NMapUint64String) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NMapUint64String"
}) {
}

type NStruct struct {
	A bool
	B string
	C int64
}

func (NStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NStruct"
}) {
}

type NEnum int

const (
	NEnumA NEnum = iota
	NEnumB
	NEnumC
)

// NEnumAll holds all labels for NEnum.
var NEnumAll = []NEnum{NEnumA, NEnumB, NEnumC}

// NEnumFromString creates a NEnum from a string label.
func NEnumFromString(label string) (x NEnum, err error) {
	err = x.Set(label)
	return
}

// Set assigns label to x.
func (x *NEnum) Set(label string) error {
	switch label {
	case "A", "a":
		*x = NEnumA
		return nil
	case "B", "b":
		*x = NEnumB
		return nil
	case "C", "c":
		*x = NEnumC
		return nil
	}
	*x = -1
	return fmt.Errorf("unknown label %q in testdata.NEnum", label)
}

// String returns the string label of x.
func (x NEnum) String() string {
	switch x {
	case NEnumA:
		return "A"
	case NEnumB:
		return "B"
	case NEnumC:
		return "C"
	}
	return ""
}

func (NEnum) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.NEnum"
	Enum struct{ A, B, C string }
}) {
}

type (
	// NUnion represents any single field of the NUnion union type.
	NUnion interface {
		// Index returns the field index.
		Index() int
		// Interface returns the field value as an interface.
		Interface() interface{}
		// Name returns the field name.
		Name() string
		// __VDLReflect describes the NUnion union type.
		__VDLReflect(__NUnionReflect)
	}
	// NUnionA represents field A of the NUnion union type.
	NUnionA struct{ Value bool }
	// NUnionB represents field B of the NUnion union type.
	NUnionB struct{ Value string }
	// NUnionC represents field C of the NUnion union type.
	NUnionC struct{ Value int64 }
	// __NUnionReflect describes the NUnion union type.
	__NUnionReflect struct {
		Name  string "v.io/core/veyron2/vom/testdata.NUnion"
		Type  NUnion
		Union struct {
			A NUnionA
			B NUnionB
			C NUnionC
		}
	}
)

func (x NUnionA) Index() int                   { return 0 }
func (x NUnionA) Interface() interface{}       { return x.Value }
func (x NUnionA) Name() string                 { return "A" }
func (x NUnionA) __VDLReflect(__NUnionReflect) {}

func (x NUnionB) Index() int                   { return 1 }
func (x NUnionB) Interface() interface{}       { return x.Value }
func (x NUnionB) Name() string                 { return "B" }
func (x NUnionB) __VDLReflect(__NUnionReflect) {}

func (x NUnionC) Index() int                   { return 2 }
func (x NUnionC) Interface() interface{}       { return x.Value }
func (x NUnionC) Name() string                 { return "C" }
func (x NUnionC) __VDLReflect(__NUnionReflect) {}

// Nested Custom Types
type MBool NBool

func (MBool) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.MBool"
}) {
}

type MStruct struct {
	A bool
	B NBool
	C MBool
	D *NStruct
	E *vdl.Type
	F vdl.AnyRep
}

func (MStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.MStruct"
}) {
}

type MList []NListUint64

func (MList) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.MList"
}) {
}

type MMap map[NFloat32]NListUint64

func (MMap) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.MMap"
}) {
}

// Recursive Type Definitions
type RecA []RecA

func (RecA) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.RecA"
}) {
}

type RecX []RecY

func (RecX) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.RecX"
}) {
}

type RecY []RecX

func (RecY) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.RecY"
}) {
}

// Additional types for compatibility and conversion checks
type ListString []string

func (ListString) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.ListString"
}) {
}

type Array3String [3]string

func (Array3String) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.Array3String"
}) {
}

type ABCStruct struct {
	A bool
	B string
	C int64
}

func (ABCStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.ABCStruct"
}) {
}

type ADEStruct struct {
	A bool
	D vdl.AnyRep
	E *vdl.Type
}

func (ADEStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.ADEStruct"
}) {
}

type XYZStruct struct {
	X bool
	Y vdl.AnyRep
	Z string
}

func (XYZStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.XYZStruct"
}) {
}

type YZStruct struct {
	Y NBool
	Z NString
}

func (YZStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.YZStruct"
}) {
}

type ZStruct struct {
	Z string
}

func (ZStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.ZStruct"
}) {
}

type MapOnlyStruct struct {
	Key1 int64
	Key2 uint32
	Key3 complex128
}

func (MapOnlyStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.MapOnlyStruct"
}) {
}

type StructOnlyMap map[string]uint64

func (StructOnlyMap) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.StructOnlyMap"
}) {
}

type MapSetStruct struct {
	Key bool
}

func (MapSetStruct) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.MapSetStruct"
}) {
}

type SetStructMap map[string]bool

func (SetStructMap) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.SetStructMap"
}) {
}

type MapStructSet map[string]struct{}

func (MapStructSet) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.MapStructSet"
}) {
}

type SetOnlyMap map[int64]bool

func (SetOnlyMap) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.SetOnlyMap"
}) {
}

type MapOnlySet map[uint16]struct{}

func (MapOnlySet) __VDLReflect(struct {
	Name string "v.io/core/veyron2/vom/testdata.MapOnlySet"
}) {
}

type (
	// BDEunion represents any single field of the BDEunion union type.
	BDEunion interface {
		// Index returns the field index.
		Index() int
		// Interface returns the field value as an interface.
		Interface() interface{}
		// Name returns the field name.
		Name() string
		// __VDLReflect describes the BDEunion union type.
		__VDLReflect(__BDEunionReflect)
	}
	// BDEunionB represents field B of the BDEunion union type.
	BDEunionB struct{ Value string }
	// BDEunionD represents field D of the BDEunion union type.
	BDEunionD struct{ Value vdl.AnyRep }
	// BDEunionE represents field E of the BDEunion union type.
	BDEunionE struct{ Value *vdl.Type }
	// __BDEunionReflect describes the BDEunion union type.
	__BDEunionReflect struct {
		Name  string "v.io/core/veyron2/vom/testdata.BDEunion"
		Type  BDEunion
		Union struct {
			B BDEunionB
			D BDEunionD
			E BDEunionE
		}
	}
)

func (x BDEunionB) Index() int                     { return 0 }
func (x BDEunionB) Interface() interface{}         { return x.Value }
func (x BDEunionB) Name() string                   { return "B" }
func (x BDEunionB) __VDLReflect(__BDEunionReflect) {}

func (x BDEunionD) Index() int                     { return 1 }
func (x BDEunionD) Interface() interface{}         { return x.Value }
func (x BDEunionD) Name() string                   { return "D" }
func (x BDEunionD) __VDLReflect(__BDEunionReflect) {}

func (x BDEunionE) Index() int                     { return 2 }
func (x BDEunionE) Interface() interface{}         { return x.Value }
func (x BDEunionE) Name() string                   { return "E" }
func (x BDEunionE) __VDLReflect(__BDEunionReflect) {}

func init() {
	vdl.Register((*ConvertGroup)(nil))
	vdl.Register((*VomdataStruct)(nil))
	vdl.Register((*NBool)(nil))
	vdl.Register((*NString)(nil))
	vdl.Register((*NByteSlice)(nil))
	vdl.Register((*NByteArray)(nil))
	vdl.Register((*NByte)(nil))
	vdl.Register((*NUint16)(nil))
	vdl.Register((*NUint32)(nil))
	vdl.Register((*NUint64)(nil))
	vdl.Register((*NInt16)(nil))
	vdl.Register((*NInt32)(nil))
	vdl.Register((*NInt64)(nil))
	vdl.Register((*NFloat32)(nil))
	vdl.Register((*NFloat64)(nil))
	vdl.Register((*NComplex64)(nil))
	vdl.Register((*NComplex128)(nil))
	vdl.Register((*NArray2Uint64)(nil))
	vdl.Register((*NListUint64)(nil))
	vdl.Register((*NSetUint64)(nil))
	vdl.Register((*NMapUint64String)(nil))
	vdl.Register((*NStruct)(nil))
	vdl.Register((*NEnum)(nil))
	vdl.Register((*NUnion)(nil))
	vdl.Register((*MBool)(nil))
	vdl.Register((*MStruct)(nil))
	vdl.Register((*MList)(nil))
	vdl.Register((*MMap)(nil))
	vdl.Register((*RecA)(nil))
	vdl.Register((*RecX)(nil))
	vdl.Register((*RecY)(nil))
	vdl.Register((*ListString)(nil))
	vdl.Register((*Array3String)(nil))
	vdl.Register((*ABCStruct)(nil))
	vdl.Register((*ADEStruct)(nil))
	vdl.Register((*XYZStruct)(nil))
	vdl.Register((*YZStruct)(nil))
	vdl.Register((*ZStruct)(nil))
	vdl.Register((*MapOnlyStruct)(nil))
	vdl.Register((*StructOnlyMap)(nil))
	vdl.Register((*MapSetStruct)(nil))
	vdl.Register((*SetStructMap)(nil))
	vdl.Register((*MapStructSet)(nil))
	vdl.Register((*SetOnlyMap)(nil))
	vdl.Register((*MapOnlySet)(nil))
	vdl.Register((*BDEunion)(nil))
}
