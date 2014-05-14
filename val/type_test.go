package val

import (
	"fmt"
	"strings"
	"testing"

	testutil "veyron/lib/testutil"
)

var singletons = []struct {
	k Kind
	t *Type
	s string
}{
	{Any, AnyType, "any"},
	{Bool, BoolType, "bool"},
	{Int32, Int32Type, "int32"},
	{Int64, Int64Type, "int64"},
	{Uint32, Uint32Type, "uint32"},
	{Uint64, Uint64Type, "uint64"},
	{Float32, Float32Type, "float32"},
	{Float64, Float64Type, "float64"},
	{Complex64, Complex64Type, "complex64"},
	{Complex128, Complex128Type, "complex128"},
	{String, StringType, "string"},
	{Bytes, BytesType, "bytes"},
	{TypeVal, TypeValType, "typeval"},
}

type l []string

var enums = []struct {
	name   string
	labels l
	str    string
	errstr string
}{
	{"", l{"A"}, "enum{A}", "must be named"},
	{"FailNoLabels", l{}, "", "no Enum labels"},
	{"FailEmptyLabel", l{""}, "", "empty Enum label"},
	{"A", l{"A"}, "A enum{A}", ""},
	{"AB", l{"A", "B"}, "AB enum{A;B}", ""},
}

type f []StructField

var structs = []struct {
	name   string
	fields f
	str    string
	errstr string
}{
	{"", f{}, "struct{}", "must be named"},
	{"FailFieldName", f{{"", BoolType}}, "", "empty Struct field name"},
	{"FailDupFields", f{{"A", BoolType}, {"A", Int32Type}}, "", "duplicate field name"},
	{"FailNilFieldType", f{{"A", nil}}, "", "nil Struct field type"},
	{"Empty", f{}, "Empty struct{}", ""},
	{"A", f{{"A", BoolType}}, "A struct{A bool}", ""},
	{"AB", f{{"A", BoolType}, {"B", Int32Type}}, "AB struct{A bool;B int32}", ""},
	{"ABC", f{{"A", BoolType}, {"B", Int32Type}, {"C", Uint64Type}}, "ABC struct{A bool;B int32;C uint64}", ""},
	{"ABCD", f{{"A", BoolType}, {"B", Int32Type}, {"C", Uint64Type}, {"D", StringType}}, "ABCD struct{A bool;B int32;C uint64;D string}", ""},
}

type t []*Type

var oneofs = []struct {
	name   string
	types  t
	str    string
	errstr string
}{
	{"", t{BoolType}, "oneof{bool}", "must be named"},
	{"FailNoTypes", t{}, "", "no OneOf types"},
	{"FailAny", t{AnyType}, "", "type in OneOf must not be nil, OneOf or Any"},
	{"FailDup", t{BoolType, BoolType}, "", "duplicate OneOf type"},
	{"A", t{BoolType}, "A oneof{bool}", ""},
	{"AB", t{BoolType, Int32Type}, "AB oneof{bool;int32}", ""},
	{"ABC", t{BoolType, Int32Type, Uint64Type}, "ABC oneof{bool;int32;uint64}", ""},
	{"ABCD", t{BoolType, Int32Type, Uint64Type, StringType}, "ABCD oneof{bool;int32;uint64;string}", ""},
}

func allTypes() (types []*Type) {
	for _, test := range singletons {
		types = append(types, test.t, ListType(test.t))
		for _, test2 := range singletons {
			types = append(types, MapType(test.t, test2.t))
		}
	}
	for _, test := range enums {
		if test.errstr == "" {
			types = append(types, EnumType(test.name, []string(test.labels)))
		}
	}
	for _, test := range structs {
		if test.errstr == "" {
			types = append(types, StructType(test.name, []StructField(test.fields)))
		}
	}
	for _, test := range oneofs {
		if test.errstr == "" {
			types = append(types, OneOfType(test.name, []*Type(test.types)))
		}
	}
	return
}

func expectPanic(t *testing.T, f func(), wantstr string, format string, args ...interface{}) {
	got := testutil.CallAndRecover(f)
	gotstr := fmt.Sprint(got)
	msg := fmt.Sprintf(format, args...)
	if wantstr != "" && !strings.Contains(gotstr, wantstr) {
		t.Errorf(`%s got panic %q, want substr %q`, msg, gotstr, wantstr)
	}
	if wantstr == "" && got != nil {
		t.Errorf(`%s got panic %q, want nil`, msg, gotstr)
	}
}

func expectMismatchedKind(t *testing.T, f func()) {
	expectPanic(t, f, "mismatched kind", "")
}

func TestTypeMismatch(t *testing.T) {
	// Make sure we panic if a method is called for a mismatched kind.
	for _, ty := range allTypes() {
		if ty.Kind() != Enum {
			expectMismatchedKind(t, func() { ty.EnumLabel(0) })
			expectMismatchedKind(t, func() { ty.EnumIndex("") })
			expectMismatchedKind(t, func() { ty.NumEnumLabel() })
		}
		if ty.Kind() != List && ty.Kind() != Map {
			expectMismatchedKind(t, func() { ty.Elem() })
		}
		if ty.Kind() != Map {
			expectMismatchedKind(t, func() { ty.Key() })
		}
		if ty.Kind() != Struct {
			expectMismatchedKind(t, func() { ty.Field(0) })
			expectMismatchedKind(t, func() { ty.FieldByName("") })
			expectMismatchedKind(t, func() { ty.NumField() })
		}
		if ty.Kind() != OneOf {
			expectMismatchedKind(t, func() { ty.OneOfType(0) })
			expectMismatchedKind(t, func() { ty.OneOfIndex(AnyType) })
			expectMismatchedKind(t, func() { ty.NumOneOfType() })
		}
	}
}

func TestSingletonTypes(t *testing.T) {
	for _, test := range singletons {
		if got, want := test.t.Kind(), test.k; got != want {
			t.Errorf(`%s got kind %q, want %q`, test.k, got, want)
		}
		if got, want := test.t.Name(), ""; got != want {
			t.Errorf(`%s got name %q, want %q`, test.k, got, want)
		}
		if got, want := test.k.String(), test.s; got != want {
			t.Errorf(`%s got kind %q, want %q`, test.k, got, want)
		}
		if got, want := test.t.String(), test.s; got != want {
			t.Errorf(`%s got string %q, want %q`, test.k, got, want)
		}
	}
}

func TestEnumTypes(t *testing.T) {
	for _, test := range enums {
		var x *Type
		create := func() { x = EnumType(test.name, []string(test.labels)) }
		expectPanic(t, create, test.errstr, "%s EnumOf", test.name)
		if x == nil {
			continue
		}
		if got, want := x.Kind(), Enum; got != want {
			t.Errorf(`Enum %s got kind %q, want %q`, test.name, got, want)
		}
		if got, want := x.Name(), test.name; got != want {
			t.Errorf(`Enum %s got name %q, want %q`, test.name, got, want)
		}
		if got, want := x.String(), test.str; got != want {
			t.Errorf(`Enum %s got string %q, want %q`, test.name, got, want)
		}
		if got, want := x.NumEnumLabel(), len(test.labels); got != want {
			t.Errorf(`Enum %s got num labels %d, want %d`, test.name, got, want)
		}
		for index, label := range test.labels {
			if got, want := x.EnumLabel(index), label; got != want {
				t.Errorf(`Enum %s got label[%d] %s, want %s`, test.name, index, got, want)
			}
			if got, want := x.EnumIndex(label), index; got != want {
				t.Errorf(`Enum %s got index[%s] %d, want %d`, test.name, label, got, want)
			}
		}
	}
}

func TestListTypes(t *testing.T) {
	for _, test := range singletons {
		x := ListType(test.t)
		if got, want := x.Kind(), List; got != want {
			t.Errorf(`List %s got kind %q, want %q`, test.k, got, want)
		}
		if got, want := x.Name(), ""; got != want {
			t.Errorf(`List %s got name %q, want %q`, test.k, got, want)
		}
		if got, want := x.String(), "[]"+test.s; got != want {
			t.Errorf(`List %s got string %q, want %q`, test.k, got, want)
		}
		if got, want := x.Elem(), test.t; got != want {
			t.Errorf(`List %s got elem %q, want %q`, test.k, got, want)
		}
	}
}

func TestMapTypes(t *testing.T) {
	for _, key := range singletons {
		for _, elem := range singletons {
			x := MapType(key.t, elem.t)
			if got, want := x.Kind(), Map; got != want {
				t.Errorf(`Map[%s]%s got kind %q, want %q`, key.k, elem.k, got, want)
			}
			if got, want := x.Name(), ""; got != want {
				t.Errorf(`Map[%s]%s got name %q, want %q`, key.k, elem.k, got, want)
			}
			if got, want := x.String(), fmt.Sprintf("map[%s]%s", key.s, elem.s); got != want {
				t.Errorf(`Map[%s]%s got name %q, want %q`, key.k, elem.k, got, want)
			}
			if got, want := x.Key(), key.t; got != want {
				t.Errorf(`Map[%s]%s got key %q, want %q`, key.k, elem.k, got, want)
			}
			if got, want := x.Elem(), elem.t; got != want {
				t.Errorf(`Map[%s]%s got elem %q, want %q`, key.k, elem.k, got, want)
			}
		}
	}
}

func TestStructTypes(t *testing.T) {
	for _, test := range structs {
		var x *Type
		create := func() { x = StructType(test.name, []StructField(test.fields)) }
		expectPanic(t, create, test.errstr, "%s StructOf", test.name)
		if x == nil {
			continue
		}
		if got, want := x.Kind(), Struct; got != want {
			t.Errorf(`Struct %s got kind %q, want %q`, test.name, got, want)
		}
		if got, want := x.Name(), test.name; got != want {
			t.Errorf(`Struct %s got name %q, want %q`, test.name, got, want)
		}
		if got, want := x.String(), test.str; got != want {
			t.Errorf(`Struct %s got string %q, want %q`, test.name, got, want)
		}
		if got, want := x.NumField(), len(test.fields); got != want {
			t.Errorf(`Struct %s got num fields %d, want %d`, test.name, got, want)
		}
		for index, field := range test.fields {
			if got, want := x.Field(index), field; got != want {
				t.Errorf(`Struct %s got field[%d] %v, want %v`, test.name, index, got, want)
			}
			gotf, goti := x.FieldByName(field.Name)
			if wantf := field; gotf != wantf {
				t.Errorf(`Struct %s got field[%s] %v, want %v`, test.name, field.Name, gotf, wantf)
			}
			if wanti := index; goti != wanti {
				t.Errorf(`Struct %s got field[%s] index %d, want %d`, test.name, field.Name, goti, wanti)
			}
		}
	}
	// Make sure hash consing of struct types respects the ordering of the fields.
	A, B, C := BoolType, Int32Type, Uint64Type
	x := StructType("X", []StructField{{"A", A}, {"B", B}, {"C", C}})
	for iter := 0; iter < 10; iter++ {
		abc := StructType("X", []StructField{{"A", A}, {"B", B}, {"C", C}})
		acb := StructType("X", []StructField{{"A", A}, {"C", C}, {"B", B}})
		bac := StructType("X", []StructField{{"B", B}, {"A", A}, {"C", C}})
		bca := StructType("X", []StructField{{"B", B}, {"C", C}, {"A", A}})
		cab := StructType("X", []StructField{{"C", C}, {"A", A}, {"B", B}})
		cba := StructType("X", []StructField{{"C", C}, {"B", B}, {"A", A}})
		if x != abc || x == acb || x == bac || x == bca || x == cab || x == cba {
			t.Errorf(`Struct ABC hash consing broken: %v, %v, %v, %v, %v, %v, %v`, x, abc, acb, bac, bca, cab, cba)
		}
		ac := StructType("X", []StructField{{"A", A}, {"C", C}})
		ca := StructType("X", []StructField{{"C", C}, {"A", A}})
		if x == ac || x == ca {
			t.Errorf(`Struct ABC / AC hash consing broken: %v, %v, %v`, x, ac, ca)
		}
	}
}

func TestOneOfTypes(t *testing.T) {
	for _, test := range oneofs {
		var x *Type
		create := func() { x = OneOfType(test.name, []*Type(test.types)) }
		expectPanic(t, create, test.errstr, "%s OneOfOf", test.name)
		if x == nil {
			continue
		}
		if got, want := x.Kind(), OneOf; got != want {
			t.Errorf(`OneOf %s got kind %q, want %q`, test.name, got, want)
		}
		if got, want := x.Name(), test.name; got != want {
			t.Errorf(`OneOf %s got name %q, want %q`, test.name, got, want)
		}
		if got, want := x.String(), test.str; got != want {
			t.Errorf(`OneOf %s got string %q, want %q`, test.name, got, want)
		}
		if got, want := x.NumOneOfType(), len(test.types); got != want {
			t.Errorf(`OneOf %s got num types %d, want %d`, test.name, got, want)
		}
		for index, one := range test.types {
			if got, want := x.OneOfType(index), one; got != want {
				t.Errorf(`OneOf %s got type[%d] %s, want %s`, test.name, index, got, want)
			}
			if got, want := x.OneOfIndex(one), index; got != want {
				t.Errorf(`OneOf %s got index[%s] %d, want %d`, test.name, one, got, want)
			}
		}
	}
	// Make sure hash consing of oneof types respects ordering.
	A, B, C := BoolType, Int32Type, Uint64Type
	x := OneOfType("X", []*Type{A, B, C})
	for iter := 0; iter < 10; iter++ {
		abc := OneOfType("X", []*Type{A, B, C})
		acb := OneOfType("X", []*Type{A, C, B})
		bac := OneOfType("X", []*Type{B, A, C})
		bca := OneOfType("X", []*Type{B, C, A})
		cab := OneOfType("X", []*Type{C, A, B})
		cba := OneOfType("X", []*Type{C, B, A})
		if x != abc || x == acb || x == bac || x == bca || x == cab || x == cba {
			t.Errorf(`OneOf ABC hash consing broken: %v, %v, %v, %v, %v, %v, %v`, x, abc, acb, bac, bca, cab, cba)
		}
		ac := OneOfType("X", []*Type{A, C})
		ca := OneOfType("X", []*Type{C, A})
		if x == ac || x == ca {
			t.Errorf(`OneOf ABC / AC hash consing broken: %v, %v, %v`, x, ac, ca)
		}
	}
}

func TestNamedTypes(t *testing.T) {
	for _, test := range singletons {
		switch test.k {
		case Any, TypeVal:
			continue // can't name Any or TypeVal
		}
		name := "Named" + test.s
		x := NamedType(name, test.t)
		if got, want := x.Kind(), test.k; got != want {
			t.Errorf(`Named %s got kind %q, want %q`, test.k, got, want)
		}
		if got, want := x.Name(), name; got != want {
			t.Errorf(`Named %s got name %q, want %q`, test.k, got, want)
		}
		if got, want := x.String(), name+" "+test.s; got != want {
			t.Errorf(`Named %s got string %q, want %q`, test.k, got, want)
		}
	}
	// Try a chain of named types:
	// type A B
	// type B C
	// type C D
	// type D struct{X []C}
	var builder TypeBuilder
	a, b, c := builder.Named("A"), builder.Named("B"), builder.Named("C")
	d := builder.Struct("D")
	a.SetBase(b)
	b.SetBase(c)
	c.SetBase(d)
	d.AppendField("X", builder.List().SetElem(c))
	builder.Build()
	bA, errA := a.Built()
	bB, errB := b.Built()
	bC, errC := c.Built()
	bD, errD := d.Built()
	if errA != nil || errB != nil || errC != nil || errD != nil {
		t.Errorf(`Named chain got (%q,%q,%q,%q), want nil`, errA, errB, errC, errD)
	}
	if got, want := bA.Kind(), Struct; got != want {
		t.Errorf(`Named chain got kind %q, want %q`, got, want)
	}
	if got, want := bB.Kind(), Struct; got != want {
		t.Errorf(`Named chain got kind %q, want %q`, got, want)
	}
	if got, want := bC.Kind(), Struct; got != want {
		t.Errorf(`Named chain got kind %q, want %q`, got, want)
	}
	if got, want := bD.Kind(), Struct; got != want {
		t.Errorf(`Named chain got kind %q, want %q`, got, want)
	}
	if got, want := bA.Name(), "A"; got != want {
		t.Errorf(`Named chain got name %q, want %q`, got, want)
	}
	if got, want := bB.Name(), "B"; got != want {
		t.Errorf(`Named chain got name %q, want %q`, got, want)
	}
	if got, want := bC.Name(), "C"; got != want {
		t.Errorf(`Named chain got name %q, want %q`, got, want)
	}
	if got, want := bD.Name(), "D"; got != want {
		t.Errorf(`Named chain got name %q, want %q`, got, want)
	}
	if got, want := bA.String(), "A struct{X []C struct{X []C}}"; got != want {
		t.Errorf(`Named chain got name %q, want %q`, got, want)
	}
	if got, want := bB.String(), "B struct{X []C struct{X []C}}"; got != want {
		t.Errorf(`Named chain got name %q, want %q`, got, want)
	}
	if got, want := bC.String(), "C struct{X []C}"; got != want {
		t.Errorf(`Named chain got name %q, want %q`, got, want)
	}
	if got, want := bD.String(), "D struct{X []C struct{X []C}}"; got != want {
		t.Errorf(`Named chain got name %q, want %q`, got, want)
	}
	if got, want := bA.NumField(), 1; got != want {
		t.Errorf(`Named chain got NumField %q, want %q`, got, want)
	}
	if got, want := bB.NumField(), 1; got != want {
		t.Errorf(`Named chain got NumField %q, want %q`, got, want)
	}
	if got, want := bC.NumField(), 1; got != want {
		t.Errorf(`Named chain got NumField %q, want %q`, got, want)
	}
	if got, want := bD.NumField(), 1; got != want {
		t.Errorf(`Named chain got NumField %q, want %q`, got, want)
	}
	if got, want := bA.Field(0).Name, "X"; got != want {
		t.Errorf(`Named chain got Field(0).Name %q, want %q`, got, want)
	}
	if got, want := bB.Field(0).Name, "X"; got != want {
		t.Errorf(`Named chain got Field(0).Name %q, want %q`, got, want)
	}
	if got, want := bC.Field(0).Name, "X"; got != want {
		t.Errorf(`Named chain got Field(0).Name %q, want %q`, got, want)
	}
	if got, want := bD.Field(0).Name, "X"; got != want {
		t.Errorf(`Named chain got Field(0).Name %q, want %q`, got, want)
	}
	listC := ListType(bC)
	if got, want := bA.Field(0).Type, listC; got != want {
		t.Errorf(`Named chain got Field(0).Type %q, want %q`, got, want)
	}
	if got, want := bB.Field(0).Type, listC; got != want {
		t.Errorf(`Named chain got Field(0).Type %q, want %q`, got, want)
	}
	if got, want := bC.Field(0).Type, listC; got != want {
		t.Errorf(`Named chain got Field(0).Type %q, want %q`, got, want)
	}
	if got, want := bD.Field(0).Type, listC; got != want {
		t.Errorf(`Named chain got Field(0).Type %q, want %q`, got, want)
	}
}

func TestHashConsTypes(t *testing.T) {
	// Create a bunch of distinct types multiple times.
	var types [3][]*Type
	for iter := 0; iter < 3; iter++ {
		for _, a := range singletons {
			types[iter] = append(types[iter], ListType(a.t))
			for _, b := range singletons {
				lA, lB := "A"+a.s, "B"+b.s
				types[iter] = append(types[iter], EnumType("Enum"+lA+lB, []string{lA, lB}))
				types[iter] = append(types[iter], MapType(a.t, b.t))
				types[iter] = append(types[iter], StructType("Struct"+lA+lB, []StructField{{lA, a.t}, {lB, b.t}}))
				if a.t != b.t && a.k != Any && b.k != Any {
					types[iter] = append(types[iter], OneOfType("OneOf"+lA+lB, []*Type{a.t, b.t}))
				}
			}
		}
	}

	// Make sure the pointers are the same across iterations, and different within
	// an iteration.
	seen := map[*Type]bool{}
	for ix := 0; ix < len(types[0]); ix++ {
		a, b, c := types[0][ix], types[1][ix], types[2][ix]
		if a != b || a != c {
			t.Errorf(`HashCons mismatched pointer[%d]: %p %p %p`, ix, a, b, c)
		}
		if seen[a] {
			t.Errorf(`HashCons dup pointer[%d]: %v %v`, ix, a, seen)
		}
		seen[a] = true
	}
}

func TestAssignableFrom(t *testing.T) {
	// Systematic testing of AssignableFrom over allTypes will just duplicate the
	// actual logic, so we just spot-check some results manually.
	tests := []struct {
		t, f   *Type
		expect bool
	}{
		{BoolType, BoolType, true},
		{AnyType, BoolType, true},
		{OneOfType("U", []*Type{BoolType}), BoolType, true},
		{OneOfType("U", []*Type{BoolType, Int32Type}), Int32Type, true},

		{BoolType, Int32Type, false},
		{BoolType, AnyType, false},
		{BoolType, OneOfType("U", []*Type{BoolType}), false},
		{BoolType, OneOfType("U", []*Type{BoolType, Int32Type}), false},
		{OneOfType("U", []*Type{BoolType}), StringType, false},
		{OneOfType("U", []*Type{BoolType, Int32Type}), StringType, false},
	}
	for _, test := range tests {
		if test.t.AssignableFrom(test.f) != test.expect {
			t.Errorf(`%v.AssignableFrom(%v) expect %v`, test.t, test.f, test.expect)
		}
	}
}

// TODO(toddw): Add tests for ValidMapKey, ValidOneOfType

func TestSelfRecursiveType(t *testing.T) {
	buildTree := func() (*Type, error, *Type, error) {
		// type Node struct {
		//   Val      string
		//   Children []Node
		// }
		var builder TypeBuilder
		pendN := builder.Struct("Node")
		pendC := builder.List().SetElem(pendN)
		pendN.AppendField("Val", StringType)
		pendN.AppendField("Children", pendC)
		builder.Build()
		c, cerr := pendC.Built()
		n, nerr := pendN.Built()
		return c, cerr, n, nerr
	}
	c, cerr, n, nerr := buildTree()
	if cerr != nil || nerr != nil {
		t.Errorf(`build got cerr %q nerr %q, want nil`, cerr, nerr)
	}
	// Check node
	if got, want := n.Kind(), Struct; got != want {
		t.Errorf(`node Kind got %s, want %s`, got, want)
	}
	if got, want := n.Name(), "Node"; got != want {
		t.Errorf(`node Name got %q, want %q`, got, want)
	}
	if got, want := n.String(), "Node struct{Val string;Children []Node}"; got != want {
		t.Errorf(`node String got %q, want %q`, got, want)
	}
	if got, want := n.NumField(), 2; got != want {
		t.Errorf(`node NumField got %q, want %q`, got, want)
	}
	if got, want := n.Field(0).Name, "Val"; got != want {
		t.Errorf(`node Field(0).Name got %q, want %q`, got, want)
	}
	if got, want := n.Field(0).Type, StringType; got != want {
		t.Errorf(`node Field(0).Type got %q, want %q`, got, want)
	}
	if got, want := n.Field(1).Name, "Children"; got != want {
		t.Errorf(`node Field(1).Name got %q, want %q`, got, want)
	}
	if got, want := n.Field(1).Type, c; got != want {
		t.Errorf(`node Field(1).Type got %q, want %q`, got, want)
	}
	// Check children
	if got, want := c.Kind(), List; got != want {
		t.Errorf(`children Kind got %s, want %s`, got, want)
	}
	if got, want := c.Name(), ""; got != want {
		t.Errorf(`children Name got %q, want %q`, got, want)
	}
	if got, want := c.String(), "[]Node struct{Val string;Children []Node}"; got != want {
		t.Errorf(`children String got %q, want %q`, got, want)
	}
	if got, want := c.Elem(), n; got != want {
		t.Errorf(`children Elem got %q, want %q`, got, want)
	}
	// Check hash-consing
	for iter := 0; iter < 5; iter++ {
		c2, cerr2, n2, nerr2 := buildTree()
		if cerr2 != nil || nerr2 != nil {
			t.Errorf(`build got cerr %q nerr %q, want nil`, cerr, nerr)
		}
		if got, want := c2, c; got != want {
			t.Errorf(`cons children got %q, want %q`, got, want)
		}
		if got, want := n2, n; got != want {
			t.Errorf(`cons node got %q, want %q`, got, want)
		}
	}
}

func TestMutuallyRecursiveType(t *testing.T) {
	build := func() (*Type, error, *Type, error, *Type, error, *Type, error) {
		// type D A
		// type A struct{X int32;B B;C C}
		// type B struct{Y int32;A A;C C}
		// type C struct{Z string}
		var builder TypeBuilder
		d := builder.Named("D")
		a, b, c := builder.Struct("A"), builder.Struct("B"), builder.Struct("C")
		d.SetBase(a)
		a.AppendField("X", Int32Type).AppendField("B", b).AppendField("C", c)
		b.AppendField("Y", Int32Type).AppendField("A", a).AppendField("C", c)
		c.AppendField("Z", StringType)
		builder.Build()
		builtD, derr := d.Built()
		builtA, aerr := a.Built()
		builtB, berr := b.Built()
		builtC, cerr := c.Built()
		return builtD, derr, builtA, aerr, builtB, berr, builtC, cerr
	}
	d, derr, a, aerr, b, berr, c, cerr := build()
	if derr != nil || aerr != nil || berr != nil || cerr != nil {
		t.Errorf(`build got (%q,%q,%q,%q), want nil`, derr, aerr, berr, cerr)
	}
	// Check D
	if got, want := d.Kind(), Struct; got != want {
		t.Errorf(`D Kind got %s, want %s`, got, want)
	}
	if got, want := d.Name(), "D"; got != want {
		t.Errorf(`D Name got %q, want %q`, got, want)
	}
	if got, want := d.String(), "D struct{X int32;B B struct{Y int32;A A struct{X int32;B B;C C struct{Z string}};C C};C C}"; got != want {
		t.Errorf(`D String got %q, want %q`, got, want)
	}
	if got, want := d.NumField(), 3; got != want {
		t.Errorf(`D NumField got %q, want %q`, got, want)
	}
	if got, want := d.Field(0).Name, "X"; got != want {
		t.Errorf(`D Field(0).Name got %q, want %q`, got, want)
	}
	if got, want := d.Field(0).Type, Int32Type; got != want {
		t.Errorf(`D Field(0).Type got %q, want %q`, got, want)
	}
	if got, want := d.Field(1).Name, "B"; got != want {
		t.Errorf(`D Field(1).Name got %q, want %q`, got, want)
	}
	if got, want := d.Field(1).Type, b; got != want {
		t.Errorf(`D Field(1).Type got %q, want %q`, got, want)
	}
	if got, want := d.Field(2).Name, "C"; got != want {
		t.Errorf(`D Field(2).Name got %q, want %q`, got, want)
	}
	if got, want := d.Field(2).Type, c; got != want {
		t.Errorf(`D Field(2).Type got %q, want %q`, got, want)
	}
	// Check A
	if got, want := a.Kind(), Struct; got != want {
		t.Errorf(`A Kind got %s, want %s`, got, want)
	}
	if got, want := a.Name(), "A"; got != want {
		t.Errorf(`A Name got %q, want %q`, got, want)
	}
	if got, want := a.String(), "A struct{X int32;B B struct{Y int32;A A;C C struct{Z string}};C C}"; got != want {
		t.Errorf(`A String got %q, want %q`, got, want)
	}
	if got, want := a.NumField(), 3; got != want {
		t.Errorf(`A NumField got %q, want %q`, got, want)
	}
	if got, want := a.Field(0).Name, "X"; got != want {
		t.Errorf(`A Field(0).Name got %q, want %q`, got, want)
	}
	if got, want := a.Field(0).Type, Int32Type; got != want {
		t.Errorf(`A Field(0).Type got %q, want %q`, got, want)
	}
	if got, want := a.Field(1).Name, "B"; got != want {
		t.Errorf(`A Field(1).Name got %q, want %q`, got, want)
	}
	if got, want := a.Field(1).Type, b; got != want {
		t.Errorf(`A Field(1).Type got %q, want %q`, got, want)
	}
	if got, want := a.Field(2).Name, "C"; got != want {
		t.Errorf(`A Field(2).Name got %q, want %q`, got, want)
	}
	if got, want := a.Field(2).Type, c; got != want {
		t.Errorf(`A Field(2).Type got %q, want %q`, got, want)
	}
	// Check B
	if got, want := b.Kind(), Struct; got != want {
		t.Errorf(`B Kind got %s, want %s`, got, want)
	}
	if got, want := b.Name(), "B"; got != want {
		t.Errorf(`B Name got %q, want %q`, got, want)
	}
	if got, want := b.String(), "B struct{Y int32;A A struct{X int32;B B;C C struct{Z string}};C C}"; got != want {
		t.Errorf(`B String got %q, want %q`, got, want)
	}
	if got, want := b.NumField(), 3; got != want {
		t.Errorf(`B NumField got %q, want %q`, got, want)
	}
	if got, want := b.Field(0).Name, "Y"; got != want {
		t.Errorf(`B Field(0).Name got %q, want %q`, got, want)
	}
	if got, want := b.Field(0).Type, Int32Type; got != want {
		t.Errorf(`B Field(0).Type got %q, want %q`, got, want)
	}
	if got, want := b.Field(1).Name, "A"; got != want {
		t.Errorf(`B Field(1).Name got %q, want %q`, got, want)
	}
	if got, want := b.Field(1).Type, a; got != want {
		t.Errorf(`B Field(1).Type got %q, want %q`, got, want)
	}
	if got, want := b.Field(2).Name, "C"; got != want {
		t.Errorf(`B Field(2).Name got %q, want %q`, got, want)
	}
	if got, want := b.Field(2).Type, c; got != want {
		t.Errorf(`B Field(2).Type got %q, want %q`, got, want)
	}
	// Check C
	if got, want := c.Kind(), Struct; got != want {
		t.Errorf(`C Kind got %s, want %s`, got, want)
	}
	if got, want := c.Name(), "C"; got != want {
		t.Errorf(`C Name got %q, want %q`, got, want)
	}
	if got, want := c.String(), "C struct{Z string}"; got != want {
		t.Errorf(`C String got %q, want %q`, got, want)
	}
	if got, want := c.NumField(), 1; got != want {
		t.Errorf(`C NumField got %q, want %q`, got, want)
	}
	if got, want := c.Field(0).Name, "Z"; got != want {
		t.Errorf(`C Field(0).Name got %q, want %q`, got, want)
	}
	if got, want := c.Field(0).Type, StringType; got != want {
		t.Errorf(`C Field(0).Type got %q, want %q`, got, want)
	}
	// Check hash-consing
	for iter := 0; iter < 5; iter++ {
		d2, derr, a2, aerr, b2, berr, c2, cerr := build()
		if derr != nil || aerr != nil || berr != nil || cerr != nil {
			t.Errorf(`build got (%q,%q,%q,%q), want nil`, derr, aerr, berr, cerr)
		}
		if got, want := d2, d; got != want {
			t.Errorf(`build got %q, want %q`, got, want)
		}
		if got, want := a2, a; got != want {
			t.Errorf(`build got %q, want %q`, got, want)
		}
		if got, want := b2, b; got != want {
			t.Errorf(`build got %q, want %q`, got, want)
		}
		if got, want := c2, c; got != want {
			t.Errorf(`build got %q, want %q`, got, want)
		}
	}
}

func TestUniqueTypeNames(t *testing.T) {
	var builder TypeBuilder
	var pending [2][]PendingType
	pending[0] = makeAllPending(&builder)
	pending[1] = makeAllPending(&builder)
	builder.Build()
	// The first pending types have no errors, but have nil types since the other
	// pending types fail to build.
	for _, p := range pending[0] {
		ty, err := p.Built()
		if ty != nil {
			t.Errorf(`built[0] got type %q, want nil`, ty)
		}
		if err != nil {
			t.Errorf(`built[0] got error %q, want nil`, err)
		}
	}
	// The second built types have non-unique name errors, and also nil types.
	for _, p := range pending[1] {
		ty, err := p.Built()
		if ty != nil {
			t.Errorf(`built[0] got type %q, want nil`, ty)
		}
		if got, want := fmt.Sprint(err), "duplicate type names"; !strings.Contains(got, want) {
			t.Errorf(`built[0] got error %s, want %s`, got, want)
		}
	}
}

func makeAllPending(builder *TypeBuilder) []PendingType {
	var ret []PendingType
	for _, test := range singletons {
		switch test.k {
		case Any, TypeVal:
			continue // can't name Any or TypeVal
		}
		ret = append(ret, builder.Named("Named"+test.s).SetBase(test.t))
	}
	for _, test := range enums {
		if test.errstr == "" {
			ret = append(ret, builder.Enum("Enum"+test.name).SetLabels([]string(test.labels)))
		}
	}
	for _, test := range structs {
		if test.errstr == "" {
			p := builder.Struct("Struct" + test.name)
			for _, f := range test.fields {
				p.AppendField(f.Name, f.Type)
			}
			ret = append(ret, p)
		}
	}
	for _, test := range oneofs {
		if test.errstr == "" {
			p := builder.OneOf("OneOf" + test.name)
			for _, t := range test.types {
				p.AppendType(t)
			}
			ret = append(ret, p)
		}
	}
	return ret
}
