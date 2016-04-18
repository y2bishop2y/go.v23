// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engine_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"v.io/v23"
	"v.io/v23/context"
	"v.io/v23/query/engine"
	ds "v.io/v23/query/engine/datasource"
	td "v.io/v23/query/engine/internal/testdata"
	"v.io/v23/query/engine/public"
	"v.io/v23/query/syncql"
	"v.io/v23/vdl"
	"v.io/v23/verror"
	"v.io/v23/vom"
	_ "v.io/x/ref/runtime/factories/roaming"
	"v.io/x/ref/test"
)

type execSelectTest struct {
	query   string
	headers []string
	r       [][]*vom.RawBytes
}

type execDeleteTest struct {
	delQuery   string
	delHeaders []string
	delResults [][]*vom.RawBytes
	selQuery   string
	selHeaders []string
	selResults [][]*vom.RawBytes
}

type execSelectErrorTest struct {
	query string
	err   error
}

type prepareSelectTest struct {
	query        string
	paramValues1 []*vom.RawBytes
	paramValues2 []*vom.RawBytes
	headers      []string
	r1           [][]*vom.RawBytes
	r2           [][]*vom.RawBytes
}

type prepareDeleteTest struct {
	delQuery    string
	paramValues []*vom.RawBytes
	delHeaders  []string
	delResults  [][]*vom.RawBytes
	selQuery    string
	selHeaders  []string
	selResults  [][]*vom.RawBytes
}

// Expect errors when calling QueryEngine.PrepareStatement()
type prepareSelectErrorOnPrepareTest struct {
	query string
	err   error
}

// Expect errors when calling PreparedStatement.Exec()
type prepareSelectErrorOnExecTest struct {
	query       string
	paramValues []*vom.RawBytes
	err         error
}

type mockDB struct {
	ctx    *context.T
	tables []*table
}

type table struct {
	name string
	rows []kv
}

type keyValueStreamImpl struct {
	table           *table
	cursor          int
	keyIndexRanges  ds.IndexRanges
	keyRangesCursor int
}

func compareKeyToLimit(key, limit string) int {
	if limit == "" || key < limit {
		return -1
	} else if key == limit {
		return 0
	} else {
		return 1
	}
}

func copyTable(src *table) *table {
	var tgt table
	tgt.name = src.name
	for _, row := range src.rows {
		tgt.rows = append(tgt.rows, row)
	}
	return &tgt
}

func (kvs *keyValueStreamImpl) Advance() bool {
	for true {
		kvs.cursor++ // initialized to -1
		if kvs.cursor >= len(kvs.table.rows) {
			return false
		}
		// does it match any keyRange
		for kvs.keyRangesCursor < len(*kvs.keyIndexRanges.StringRanges) {
			if kvs.table.rows[kvs.cursor].key >= (*kvs.keyIndexRanges.StringRanges)[kvs.keyRangesCursor].Start && compareKeyToLimit(kvs.table.rows[kvs.cursor].key, (*kvs.keyIndexRanges.StringRanges)[kvs.keyRangesCursor].Limit) < 0 {
				return true
			}
			// Keys and keyIndexRanges.StringRanges are both sorted low to high,
			// so we can increment keyRangesCursor if the keyRange.Limit is < the key.
			if compareKeyToLimit(kvs.table.rows[kvs.cursor].key, (*kvs.keyIndexRanges.StringRanges)[kvs.keyRangesCursor].Limit) > 0 {
				kvs.keyRangesCursor++
				if kvs.keyRangesCursor >= len(*kvs.keyIndexRanges.StringRanges) {
					return false
				}
			} else {
				break
			}
		}
	}
	return false
}

func (kvs *keyValueStreamImpl) KeyValue() (string, *vom.RawBytes) {
	return kvs.table.rows[kvs.cursor].key, kvs.table.rows[kvs.cursor].value
}

func (kvs *keyValueStreamImpl) Err() error {
	return nil
}

func (kvs *keyValueStreamImpl) Cancel() {
}

func (t *table) GetIndexFields() []ds.Index {
	return []ds.Index{}
}

func (t *table) Scan(indexRanges ...ds.IndexRanges) (ds.KeyValueStream, error) {
	var keyValueStreamImpl keyValueStreamImpl
	keyValueStreamImpl.table = copyTable(t)
	keyValueStreamImpl.cursor = -1
	keyValueStreamImpl.keyIndexRanges = indexRanges[0]
	return &keyValueStreamImpl, nil
}

func (t *table) Delete(k string) (bool, error) {
	for i, kv := range t.rows {
		if kv.key == k {
			t.rows = append(t.rows[:i], t.rows[i+1:]...)
			return true, nil
		}
	}
	// nothing to delete
	return false, nil
}

func (db mockDB) GetContext() *context.T {
	return db.ctx
}

func (db mockDB) GetTable(table string, writeAccessReq bool) (ds.Table, error) {
	for _, t := range db.tables {
		if t.name == table {
			return t, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No such table: %s.", table))

}

var db mockDB
var custTable table
var numTable table
var fooTable table
var funWithMapsTable table
var ratingsArrayTable table
var tdhApprovalsTable table
var previousRatingsTable table
var previousAddressesTable table
var manyMapsTable table
var manySetsTable table
var bigTable table
var funWithTypesTable table

type kv struct {
	key   string
	value *vom.RawBytes
}

var t2015 time.Time

var t2015_04 time.Time
var t2015_04_12 time.Time
var t2015_04_12_22 time.Time
var t2015_04_12_22_16 time.Time
var t2015_04_12_22_16_06 time.Time

var t2015_07 time.Time
var t2015_07_01 time.Time
var t2015_07_01_01_23_45 time.Time

func init() {
	var shutdown v23.Shutdown
	db.ctx, shutdown = test.V23Init()
	defer shutdown()
}

func initTables() {
	db.tables = db.tables[:0]
	custTable.rows = custTable.rows[:0]
	numTable.rows = custTable.rows[:0]
	fooTable.rows = custTable.rows[:0]
	funWithMapsTable.rows = funWithMapsTable.rows[:0]
	ratingsArrayTable.rows = ratingsArrayTable.rows[:0]
	tdhApprovalsTable.rows = tdhApprovalsTable.rows[:0]
	previousRatingsTable.rows = previousRatingsTable.rows[:0]
	previousAddressesTable.rows = previousAddressesTable.rows[:0]
	manyMapsTable.rows = manyMapsTable.rows[:0]
	manySetsTable.rows = manySetsTable.rows[:0]
	bigTable.rows = bigTable.rows[:0]
	funWithTypesTable.rows = funWithTypesTable.rows[:0]

	t20150122131101, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Jan 22 2015 13:11:01 -0800 PST")
	t20150210161202, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Feb 10 2015 16:12:02 -0800 PST")
	t20150311101303, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Mar 11 2015 10:13:03 -0700 PDT")
	t20150317111404, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Mar 17 2015 11:14:04 -0700 PDT")
	t20150317131505, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Mar 17 2015 13:15:05 -0700 PDT")
	t20150412221606, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Apr 12 2015 22:16:06 -0700 PDT")
	t20150413141707, _ := time.Parse("Jan 2 2006 15:04:05 -0700 MST", "Apr 13 2015 14:17:07 -0700 PDT")

	t2015, _ = time.Parse("2006 MST", "2015 PST")

	t2015_04, _ = time.Parse("Jan 2006 MST", "Apr 2015 PDT")
	t2015_07, _ = time.Parse("Jan 2006 MST", "Jul 2015 PDT")

	t2015_04_12, _ = time.Parse("Jan 2 2006 MST", "Apr 12 2015 PDT")
	t2015_07_01, _ = time.Parse("Jan 2 2006 MST", "Jul 01 2015 PDT")

	t2015_04_12_22, _ = time.Parse("Jan 2 2006 15 MST", "Apr 12 2015 22 PDT")
	t2015_04_12_22_16, _ = time.Parse("Jan 2 2006 15:04 MST", "Apr 12 2015 22:16 PDT")
	t2015_04_12_22_16_06, _ = time.Parse("Jan 2 2006 15:04:05 MST", "Apr 12 2015 22:16:06 PDT")
	t2015_07_01_01_23_45, _ = time.Parse("Jan 2 2006 15:04:05 MST", "Jul 01 2015 01:23:45 PDT")

	custTable.name = "Customer"
	custTable.rows = []kv{
		kv{
			"001",
			vom.RawBytesOf(td.Customer{"John Smith", 1, true, td.AddressInfo{"1 Main St.", "Palo Alto", "CA", "94303"}, []td.AddressInfo{td.AddressInfo{"10 Brown St.", "Mountain View", "CA", "94043"}}, td.CreditReport{Agency: td.CreditAgencyEquifax, Report: td.AgencyReportEquifaxReport{td.EquifaxCreditReport{'A', [4]int16{87, 81, 42, 2}}}}}),
		},
		kv{
			"001001",
			vom.RawBytesOf(td.Invoice{1, 1000, t20150122131101, 42, td.AddressInfo{"1 Main St.", "Palo Alto", "CA", "94303"}}),
		},
		kv{
			"001002",
			vom.RawBytesOf(td.Invoice{1, 1003, t20150210161202, 7, td.AddressInfo{"2 Main St.", "Palo Alto", "CA", "94303"}}),
		},
		kv{
			"001003",
			vom.RawBytesOf(td.Invoice{1, 1005, t20150311101303, 88, td.AddressInfo{"3 Main St.", "Palo Alto", "CA", "94303"}}),
		},
		kv{
			"002",
			vom.RawBytesOf(td.Customer{"Bat Masterson", 2, true, td.AddressInfo{"777 Any St.", "Collins", "IA", "50055"}, []td.AddressInfo{td.AddressInfo{"19 Green St.", "Boulder", "CO", "80301"}, td.AddressInfo{"558 W. Orange St.", "Lancaster", "PA", "17603"}}, td.CreditReport{Agency: td.CreditAgencyTransUnion, Report: td.AgencyReportTransUnionReport{td.TransUnionCreditReport{80, map[string]int16{"2015Q2": 40, "2015Q1": 60}}}}}),
		},
		kv{
			"002001",
			vom.RawBytesOf(td.Invoice{2, 1001, t20150317111404, 166, td.AddressInfo{"777 Any St.", "collins", "IA", "50055"}}),
		},
		kv{
			"002002",
			vom.RawBytesOf(td.Invoice{2, 1002, t20150317131505, 243, td.AddressInfo{"888 Any St.", "collins", "IA", "50055"}}),
		},
		kv{
			"002003",
			vom.RawBytesOf(td.Invoice{2, 1004, t20150412221606, 787, td.AddressInfo{"999 Any St.", "collins", "IA", "50055"}}),
		},
		kv{
			"002004",
			vom.RawBytesOf(td.Invoice{2, 1006, t20150413141707, 88, td.AddressInfo{"101010 Any St.", "collins", "IA", "50055"}}),
		},
		kv{
			"003",
			vom.RawBytesOf(td.Customer{"John Steed", 3, true, td.AddressInfo{"100 Queen St.", "New%London", "CT", "06320"}, []td.AddressInfo{}, td.CreditReport{Agency: td.CreditAgencyExperian, Report: td.AgencyReportExperianReport{td.ExperianCreditReport{td.ExperianRatingGood, map[td.Tdh]struct{}{td.TdhTom: {}, td.TdhHarry: {}}, td.TdhTom}}}}),
		},
	}
	db.tables = append(db.tables, &custTable)

	numTable.name = "Numbers"
	numTable.rows = []kv{
		kv{
			"001",
			vom.RawBytesOf(td.Numbers{byte(12), uint16(1234), uint32(5678), uint64(999888777666), int16(9876), int32(876543), int64(128), float32(3.14159), float64(2.71828182846)}),
		},
		kv{
			"002",
			vom.RawBytesOf(td.Numbers{byte(9), uint16(99), uint32(999), uint64(9999999), int16(9), int32(99), int64(88), float32(1.41421356237), float64(1.73205080757)}),
		},
		kv{
			"003",
			vom.RawBytesOf(td.Numbers{byte(210), uint16(210), uint32(210), uint64(210), int16(210), int32(210), int64(210), float32(210.0), float64(210.0)}),
		},
	}
	db.tables = append(db.tables, &numTable)

	fooTable.name = "Foo"
	fooTable.rows = []kv{
		kv{
			"001",
			vom.RawBytesOf(td.FooType{td.BarType{td.BazType{"FooBarBaz", td.TitleOrValueTypeTitle{"Vice President"}}}}),
		},
		kv{
			"002",
			vom.RawBytesOf(td.FooType{td.BarType{td.BazType{"BazBarFoo", td.TitleOrValueTypeValue{42}}}}),
		},
	}
	db.tables = append(db.tables, &fooTable)

	funWithMapsTable.name = "FunWithMaps"
	funWithMapsTable.rows = []kv{
		kv{
			"AAA",
			vom.RawBytesOf(td.FunWithMaps{td.K{'a', "bbb"}, map[td.K]td.V{td.K{'a', "aaa"}: td.V{"bbb", 23.0}, td.K{'a', "bbb"}: td.V{"ccc", 14.7}},
				map[int16][]map[string]struct{}{
					23: []map[string]struct{}{
						map[string]struct{}{"foo": {}, "bar": {}},
					},
				},
			}),
		},
		kv{
			"BBB",
			vom.RawBytesOf(td.FunWithMaps{td.K{'x', "zzz"}, map[td.K]td.V{td.K{'x', "zzz"}: td.V{"yyy", 17.1}, td.K{'r', "sss"}: td.V{"qqq", 7.8}},
				map[int16][]map[string]struct{}{
					42: []map[string]struct{}{
						map[string]struct{}{"great": {}, "dane": {}},
						map[string]struct{}{"german": {}, "shepard": {}},
					},
				},
			}),
		},
	}
	db.tables = append(db.tables, &funWithMapsTable)

	ratingsArrayTable.name = "RatingsArray"
	ratingsArrayTable.rows = []kv{
		kv{
			"000",
			vom.RawBytesOf(td.RatingsArray{40, 20, 10, 0}),
		},
		kv{
			"111",
			vom.RawBytesOf(td.RatingsArray{17, 18, 19, 20}),
		},
	}
	db.tables = append(db.tables, &ratingsArrayTable)

	tdhApprovalsTable.name = "TdhApprovals"
	tdhApprovalsTable.rows = []kv{
		kv{
			"yyy",
			vom.RawBytesOf(map[td.Tdh]struct{}{td.TdhTom: {}}),
		},
		kv{
			"zzz",
			vom.RawBytesOf(map[td.Tdh]struct{}{td.TdhDick: {}, td.TdhHarry: {}}),
		},
	}
	db.tables = append(db.tables, &tdhApprovalsTable)

	previousRatingsTable.name = "PreviousRatings"
	previousRatingsTable.rows = []kv{
		kv{
			"x1",
			vom.RawBytesOf(map[string]int16{"1Q2015": 1, "2Q2015": 2}),
		},
		kv{
			"x2",
			vom.RawBytesOf(map[string]int16{"2Q2015": 3}),
		},
	}
	db.tables = append(db.tables, &previousRatingsTable)

	previousAddressesTable.name = "PreviousAddresses"
	previousAddressesTable.rows = []kv{
		kv{
			"a1",
			vom.RawBytesOf([]td.AddressInfo{
				td.AddressInfo{"100 Main St.", "Anytown", "CA", "94303"},
				td.AddressInfo{"200 Main St.", "Othertown", "IA", "51050"},
			}),
		},
		kv{
			"a2",
			vom.RawBytesOf([]td.AddressInfo{
				td.AddressInfo{"500 Orange St", "Uptown", "ID", "83209"},
				td.AddressInfo{"200 Fulton St", "Downtown", "MT", "59001"},
			}),
		},
	}
	db.tables = append(db.tables, &previousAddressesTable)

	manyMapsTable.name = "ManyMaps"
	manyMapsTable.rows = []kv{
		kv{
			"0",
			vom.RawBytesOf(td.ManyMaps{
				map[bool]string{true: "It was the best of times,"},
				map[byte]string{10: "it was the worst of times,"},
				map[uint16]string{16: "it was the age of wisdom,"},
				map[uint32]string{32: "it was the age of foolishness,"},
				map[uint64]string{64: "it was the epoch of belief,"},
				map[int16]string{17: "it was the epoch of incredulity,"},
				map[int32]string{33: "it was the season of Light,"},
				map[int64]string{65: "it was the season of Darkness,"},
				map[float32]string{32.1: "it was the spring of hope,"},
				map[float64]string{64.2: "it was the winter of despair,"}, map[string]string{"Dickens": "we are all going direct to Heaven,"},
				map[string]map[string]string{
					"Charles": map[string]string{"Dickens": "we are all going direct to Heaven,"},
				},
				map[time.Time]string{t2015_07_01_01_23_45: "we are all going direct the other way"},
			}),
		},
	}
	db.tables = append(db.tables, &manyMapsTable)

	manySetsTable.name = "ManySets"
	manySetsTable.rows = []kv{
		kv{
			"0",
			vom.RawBytesOf(td.ManySets{
				map[bool]struct{}{true: {}},
				map[byte]struct{}{10: {}},
				map[uint16]struct{}{16: {}},
				map[uint32]struct{}{32: {}},
				map[uint64]struct{}{64: {}},
				map[int16]struct{}{17: {}},
				map[int32]struct{}{33: {}},
				map[int64]struct{}{65: {}},
				map[float32]struct{}{32.1: {}},
				map[float64]struct{}{64.2: {}},
				map[string]struct{}{"Dickens": {}},
				map[time.Time]struct{}{t2015_07_01_01_23_45: {}},
			}),
		},
	}
	db.tables = append(db.tables, &manySetsTable)

	bigTable.name = "BigTable"

	for i := 100; i < 301; i++ {
		k := fmt.Sprintf("%d", i)
		b := vom.RawBytesOf(td.BigData{k})
		bigTable.rows = append(bigTable.rows, kv{k, b})
	}
	db.tables = append(db.tables, &bigTable)

	custType := vdl.TypeOf(td.Customer{"John Steed", 3, true, td.AddressInfo{"100 Queen St.", "New London", "CT", "06320"}, []td.AddressInfo{}, td.CreditReport{Agency: td.CreditAgencyExperian, Report: td.AgencyReportExperianReport{td.ExperianCreditReport{td.ExperianRatingGood, map[td.Tdh]struct{}{td.TdhTom: {}, td.TdhHarry: {}}, td.TdhTom}}}})
	invType := vdl.TypeOf(td.Invoice{2, 1006, t20150413141707, 88, td.AddressInfo{"101010 Any St.", "collins", "IA", "50055"}})

	funWithTypesTable.name = "FunWithTypes"
	funWithTypesTable.rows = []kv{
		kv{
			"1",
			vom.RawBytesOf(td.FunWithTypes{T1: custType, T2: invType}),
		},
		kv{
			"2",
			vom.RawBytesOf(td.FunWithTypes{T1: custType, T2: custType}),
		},
	}
	db.tables = append(db.tables, &funWithTypesTable)
}

func TestSelect(t *testing.T) {
	initTables()
	basic := []execSelectTest{
		{
			// Select values for all customer records.
			"select v from Customer where Type(v) like \"%.Customer\"",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
				{custTable.rows[4].value},
				{custTable.rows[9].value},
			},
		},
		{
			// Find a city that contains the literal '%' character in it.
			"select v.Address.City from Customer where Type(v) like \"%.Customer\" and v.Address.City like \"%^%%\" escape '^'",
			[]string{"v.Address.City"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("New%London")},
			},
		},
		{
			// Select values for all customer records.
			"select v from Customer where Type(v) not like \"%.Customer\"",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[1].value},
				{custTable.rows[2].value},
				{custTable.rows[3].value},
				{custTable.rows[5].value},
				{custTable.rows[6].value},
				{custTable.rows[7].value},
				{custTable.rows[8].value},
			},
		},
		{
			// Select values for all customer records.
			"select Type(v) from Customer where Type(v) not like \"%Customer\"",
			[]string{"Type"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
			},
		},
		{
			// All customers have a v.Credit with type CreditReport.
			"select v from Customer where Type(v.Credit) like \"%.CreditReport\"",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
				{custTable.rows[4].value},
				{custTable.rows[9].value},
			},
		},
		{
			// Only customer "001" has an equifax report.
			"select v from Customer where Type(v.Credit.Report.EquifaxReport) like \"%.EquifaxCreditReport\"",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
			},
		},
		{
			// Print the types of every record
			"select Type(v) from Customer",
			[]string{"Type"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer")},
			},
		},
		{
			// Print the types of every credit report
			"select Type(v.Credit.Report) from Customer",
			[]string{"Type"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.AgencyReport")},
				{vom.RawBytesOf(nil)},
				{vom.RawBytesOf(nil)},
				{vom.RawBytesOf(nil)},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.AgencyReport")},
				{vom.RawBytesOf(nil)},
				{vom.RawBytesOf(nil)},
				{vom.RawBytesOf(nil)},
				{vom.RawBytesOf(nil)},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.AgencyReport")},
			},
		},
		{
			// Print the types of every cusomer's v.Credit.Report.EquifaxReport
			"select Type(v.Credit.Report.EquifaxReport) from Customer where Type(v) like \"%.Customer\"",
			[]string{"Type"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.EquifaxCreditReport")},
				{vom.RawBytesOf(vom.RawBytesOf(nil))},
				{vom.RawBytesOf(vom.RawBytesOf(nil))},
			},
		},
		{
			// Select values where v.InvoiceNum is nil
			// Since InvoiceNum does not exist for Invoice,
			// this will return just customers.
			"select v from Customer where v.InvoiceNum is nil",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
				{custTable.rows[4].value},
				{custTable.rows[9].value},
			},
		},
		{
			// Select values where v.InvoiceNum is nil
			// or v.Name is nil This will select all customers
			// with the former and all invoices with the latter.
			// Hence, all records are returned.
			"select v from Customer where v.InvoiceNum is nil or v.Name is nil",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
				{custTable.rows[1].value},
				{custTable.rows[2].value},
				{custTable.rows[3].value},
				{custTable.rows[4].value},
				{custTable.rows[5].value},
				{custTable.rows[6].value},
				{custTable.rows[7].value},
				{custTable.rows[8].value},
				{custTable.rows[9].value},
			},
		},
		{
			// Select values where v.InvoiceNum is nil
			// and v.Name is nil.  Expect nothing returned.
			"select v from Customer where v.InvoiceNum is nil and v.Name is nil",
			[]string{"v"},
			[][]*vom.RawBytes{},
		},
		{
			// Select values where v.InvoiceNum is not nil
			// This will return just invoices.
			"select v from Customer where v.InvoiceNum is not nil",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[1].value},
				{custTable.rows[2].value},
				{custTable.rows[3].value},
				{custTable.rows[5].value},
				{custTable.rows[6].value},
				{custTable.rows[7].value},
				{custTable.rows[8].value},
			},
		},
		{
			// Select values where v.InvoiceNum is not nil
			// or v.Name is not nil. All records are returned.
			"select v from Customer where v.InvoiceNum is not nil or v.Name is not nil",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
				{custTable.rows[1].value},
				{custTable.rows[2].value},
				{custTable.rows[3].value},
				{custTable.rows[4].value},
				{custTable.rows[5].value},
				{custTable.rows[6].value},
				{custTable.rows[7].value},
				{custTable.rows[8].value},
				{custTable.rows[9].value},
			},
		},
		{
			// Select values where v.InvoiceNum is nil and v.Name is not nil.
			// All customers are returned.
			"select v from Customer where v.InvoiceNum is nil and v.Name is not nil",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
				{custTable.rows[4].value},
				{custTable.rows[9].value},
			},
		},
		{
			// Select values where v.InvoiceNum is not nil
			// and v.Name is not nil.  Expect nothing returned.
			"select v from Customer where v.InvoiceNum is not nil and v.Name is not nil",
			[]string{"v"},
			[][]*vom.RawBytes{},
		},
		{
			// Select keys & values for all customer records.
			"select k, v from Customer where \"v.io/v23/query/engine/internal/testdata.Customer\" = Type(v)",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), custTable.rows[0].value},
				{vom.RawBytesOf(custTable.rows[4].key), custTable.rows[4].value},
				{vom.RawBytesOf(custTable.rows[9].key), custTable.rows[9].value},
			},
		},
		{
			// Select keys & names for all customer records.
			"select k, v.Name from Customer where Type(v) like \"%.Customer\"",
			[]string{"k", "v.Name"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), vom.RawBytesOf("John Smith")},
				{vom.RawBytesOf(custTable.rows[4].key), vom.RawBytesOf("Bat Masterson")},
				{vom.RawBytesOf(custTable.rows[9].key), vom.RawBytesOf("John Steed")},
			},
		},
		{
			// Select both customer and invoice records.
			// Customer records have Id.
			// Invoice records have CustId.
			"select v.Id, v.CustId from Customer",
			[]string{"v.Id", "v.CustId"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(1)), vom.RawBytesOf(nil)},
				{vom.RawBytesOf(nil), vom.RawBytesOf(int64(1))},
				{vom.RawBytesOf(nil), vom.RawBytesOf(int64(1))},
				{vom.RawBytesOf(nil), vom.RawBytesOf(int64(1))},
				{vom.RawBytesOf(int64(2)), vom.RawBytesOf(nil)},
				{vom.RawBytesOf(nil), vom.RawBytesOf(int64(2))},
				{vom.RawBytesOf(nil), vom.RawBytesOf(int64(2))},
				{vom.RawBytesOf(nil), vom.RawBytesOf(int64(2))},
				{vom.RawBytesOf(nil), vom.RawBytesOf(int64(2))},
				{vom.RawBytesOf(int64(3)), vom.RawBytesOf(nil)},
			},
		},
		{
			// Select keys & values fo all invoice records.
			"select k, v from Customer where Type(v) like \"%.Invoice\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[1].key), custTable.rows[1].value},
				{vom.RawBytesOf(custTable.rows[2].key), custTable.rows[2].value},
				{vom.RawBytesOf(custTable.rows[3].key), custTable.rows[3].value},
				{vom.RawBytesOf(custTable.rows[5].key), custTable.rows[5].value},
				{vom.RawBytesOf(custTable.rows[6].key), custTable.rows[6].value},
				{vom.RawBytesOf(custTable.rows[7].key), custTable.rows[7].value},
				{vom.RawBytesOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Select key, cust id, invoice number and amount for $88 invoices.
			"select k, v.CustId as ID, v.InvoiceNum as InvoiceNumber, v.Amount as Amt from Customer where Type(v) like \"%.Invoice\" and v.Amount = 88",
			[]string{"k", "ID", "InvoiceNumber", "Amt"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[3].key), vom.RawBytesOf(int64(1)), vom.RawBytesOf(int64(1005)), vom.RawBytesOf(int64(88))},
				{vom.RawBytesOf(custTable.rows[8].key), vom.RawBytesOf(int64(2)), vom.RawBytesOf(int64(1006)), vom.RawBytesOf(int64(88))},
			},
		},
		{
			// Select keys & values for all records with a key prefix of "001".
			"select k, v from Customer where k like \"001%\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), custTable.rows[0].value},
				{vom.RawBytesOf(custTable.rows[1].key), custTable.rows[1].value},
				{vom.RawBytesOf(custTable.rows[2].key), custTable.rows[2].value},
				{vom.RawBytesOf(custTable.rows[3].key), custTable.rows[3].value},
			},
		},
		{
			// Select keys & values for all records with a key prefix of "002".
			"select k, v from Customer where k like \"002%\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[4].key), custTable.rows[4].value},
				{vom.RawBytesOf(custTable.rows[5].key), custTable.rows[5].value},
				{vom.RawBytesOf(custTable.rows[6].key), custTable.rows[6].value},
				{vom.RawBytesOf(custTable.rows[7].key), custTable.rows[7].value},
				{vom.RawBytesOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Select keys & values for all records with a key prefix NOT LIKE "002%".
			"select k, v from Customer where k not like \"002%\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), custTable.rows[0].value},
				{vom.RawBytesOf(custTable.rows[1].key), custTable.rows[1].value},
				{vom.RawBytesOf(custTable.rows[2].key), custTable.rows[2].value},
				{vom.RawBytesOf(custTable.rows[3].key), custTable.rows[3].value},
				{vom.RawBytesOf(custTable.rows[9].key), custTable.rows[9].value},
			},
		},
		{
			// Select keys & values for all records with a key prefix NOT LIKE "002".
			// will be optimized to k <> "002"
			"select k, v from Customer where k not like \"002\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), custTable.rows[0].value},
				{vom.RawBytesOf(custTable.rows[1].key), custTable.rows[1].value},
				{vom.RawBytesOf(custTable.rows[2].key), custTable.rows[2].value},
				{vom.RawBytesOf(custTable.rows[3].key), custTable.rows[3].value},
				{vom.RawBytesOf(custTable.rows[5].key), custTable.rows[5].value},
				{vom.RawBytesOf(custTable.rows[6].key), custTable.rows[6].value},
				{vom.RawBytesOf(custTable.rows[7].key), custTable.rows[7].value},
				{vom.RawBytesOf(custTable.rows[8].key), custTable.rows[8].value},
				{vom.RawBytesOf(custTable.rows[9].key), custTable.rows[9].value},
			},
		},
		{
			// Select keys & values for all records with a key prefix of "001".
			// or a key prefix of "002".
			"select k, v from Customer where k like \"001%\" or k like \"002%\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), custTable.rows[0].value},
				{vom.RawBytesOf(custTable.rows[1].key), custTable.rows[1].value},
				{vom.RawBytesOf(custTable.rows[2].key), custTable.rows[2].value},
				{vom.RawBytesOf(custTable.rows[3].key), custTable.rows[3].value},
				{vom.RawBytesOf(custTable.rows[4].key), custTable.rows[4].value},
				{vom.RawBytesOf(custTable.rows[5].key), custTable.rows[5].value},
				{vom.RawBytesOf(custTable.rows[6].key), custTable.rows[6].value},
				{vom.RawBytesOf(custTable.rows[7].key), custTable.rows[7].value},
				{vom.RawBytesOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Select keys & values for all records with a key prefix of "001".
			// or a key prefix of "002".
			"select k, v from Customer where k like \"002%\" or k like \"001%\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), custTable.rows[0].value},
				{vom.RawBytesOf(custTable.rows[1].key), custTable.rows[1].value},
				{vom.RawBytesOf(custTable.rows[2].key), custTable.rows[2].value},
				{vom.RawBytesOf(custTable.rows[3].key), custTable.rows[3].value},
				{vom.RawBytesOf(custTable.rows[4].key), custTable.rows[4].value},
				{vom.RawBytesOf(custTable.rows[5].key), custTable.rows[5].value},
				{vom.RawBytesOf(custTable.rows[6].key), custTable.rows[6].value},
				{vom.RawBytesOf(custTable.rows[7].key), custTable.rows[7].value},
				{vom.RawBytesOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Let's play with whitespace and mixed case.
			"   sElEcT  k,  v from \n  Customer WhErE k lIkE \"002%\" oR k LiKe \"001%\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), custTable.rows[0].value},
				{vom.RawBytesOf(custTable.rows[1].key), custTable.rows[1].value},
				{vom.RawBytesOf(custTable.rows[2].key), custTable.rows[2].value},
				{vom.RawBytesOf(custTable.rows[3].key), custTable.rows[3].value},
				{vom.RawBytesOf(custTable.rows[4].key), custTable.rows[4].value},
				{vom.RawBytesOf(custTable.rows[5].key), custTable.rows[5].value},
				{vom.RawBytesOf(custTable.rows[6].key), custTable.rows[6].value},
				{vom.RawBytesOf(custTable.rows[7].key), custTable.rows[7].value},
				{vom.RawBytesOf(custTable.rows[8].key), custTable.rows[8].value},
			},
		},
		{
			// Add in a like clause that accepts all strings.
			"   sElEcT  k,  v from \n  Customer WhErE k lIkE \"002%\" oR k LiKe \"001%\" or k lIkE \"%\"",
			[]string{"k", "v"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key), custTable.rows[0].value},
				{vom.RawBytesOf(custTable.rows[1].key), custTable.rows[1].value},
				{vom.RawBytesOf(custTable.rows[2].key), custTable.rows[2].value},
				{vom.RawBytesOf(custTable.rows[3].key), custTable.rows[3].value},
				{vom.RawBytesOf(custTable.rows[4].key), custTable.rows[4].value},
				{vom.RawBytesOf(custTable.rows[5].key), custTable.rows[5].value},
				{vom.RawBytesOf(custTable.rows[6].key), custTable.rows[6].value},
				{vom.RawBytesOf(custTable.rows[7].key), custTable.rows[7].value},
				{vom.RawBytesOf(custTable.rows[8].key), custTable.rows[8].value},
				{vom.RawBytesOf(custTable.rows[9].key), custTable.rows[9].value},
			},
		},
		{
			// Select id, name for customers whose last name is Masterson.
			"select v.Id as ID, v.Name as Name from Customer where Type(v) like \"%.Customer\" and v.Name like \"%Masterson\"",
			[]string{"ID", "Name"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(2)), vom.RawBytesOf("Bat Masterson")},
			},
		},
		{
			// Select records where v.Address.City is "Collins" or type is Invoice.
			"select v from Customer where v.Address.City = \"Collins\" or Type(v) like \"%.Invoice\"",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[1].value},
				{custTable.rows[2].value},
				{custTable.rows[3].value},
				{custTable.rows[4].value},
				{custTable.rows[5].value},
				{custTable.rows[6].value},
				{custTable.rows[7].value},
				{custTable.rows[8].value},
			},
		},
		{
			"select k from Customer where k >= \"002001\" and k <= \"002002\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
			},
		},
		{
			"select k from Customer where k > \"002001\" and k <= \"002002\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002002")},
			},
		},
		{
			"select k from Customer where k > \"002001\" and k < \"002002\"",
			[]string{"k"},
			[][]*vom.RawBytes{},
		},
		{
			"select k from Customer where k > \"002001\" or k < \"002002\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			"select k from Customer where k <> \"002\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			"select k from Customer where k <> \"002\" or k like \"002\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			"select k from Customer where k <> \"002\" and k like \"002\"",
			[]string{"k"},
			[][]*vom.RawBytes{},
		},
		{
			"select k from Customer where k <> \"002\" and k like \"002%\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			"select k from Customer where k <> \"002\" and k like \"%002\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("002002")},
			},
		},
		{
			// Select records where v.Address.City is "Collins" and v.InvoiceNum is not nil.
			"select v from Customer where v.Address.City = \"Collins\" and v.InvoiceNum is not nil",
			[]string{"v"},
			[][]*vom.RawBytes{},
		},
		{
			// Select records where v.Address.City is "Collins" and v.InvoiceNum is nil.
			"select v from Customer where v.Address.City = \"Collins\" and v.InvoiceNum is nil",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[4].value},
			},
		},
		{
			// Select customer name for customer Id (i.e., key) "001".
			"select v.Name as Name from Customer where Type(v) like \"%.Customer\" and k = \"001\"",
			[]string{"Name"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("John Smith")},
			},
		},
		{
			// Select v where v.Credit.Report.EquifaxReport.Rating = 'A'
			"select v from Customer where v.Credit.Report.EquifaxReport.Rating = 'A'",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
			},
		},
		{
			// Select v where v.AgencyRating = "Bad"
			"select v from Customer where v.Credit.Report.EquifaxReport.Rating < 'A' or v.Credit.Report.ExperianReport.Rating = \"Bad\" or v.Credit.Report.TransUnionReport.Rating < 90",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[4].value},
			},
		},
		{
			// Select records where v.Bar.Baz.Name = "FooBarBaz"
			"select v from Foo where v.Bar.Baz.Name = \"FooBarBaz\"",
			[]string{"v"},
			[][]*vom.RawBytes{
				[]*vom.RawBytes{fooTable.rows[0].value},
			},
		},
		{
			// Select records where v.Bar.Baz.TitleOrValue.Value = 42
			"select v from Foo where v.Bar.Baz.TitleOrValue.Value = 42",
			[]string{"v"},
			[][]*vom.RawBytes{
				[]*vom.RawBytes{fooTable.rows[1].value},
			},
		},
		{
			// Select records where v.Bar.Baz.TitleOrValue.Title = "Vice President"
			"select v from Foo where v.Bar.Baz.TitleOrValue.Title = \"Vice President\"",
			[]string{"v"},
			[][]*vom.RawBytes{
				[]*vom.RawBytes{fooTable.rows[0].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" or type is Invoice.
			// Limit 3
			"select v from Customer where v.Address.City = \"Collins\" or Type(v) like \"%.Invoice\" limit 3",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[1].value},
				{custTable.rows[2].value},
				{custTable.rows[3].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" or type is Invoice.
			// Offset 5
			"select v from Customer where v.Address.City = \"Collins\" or Type(v) like \"%.Invoice\" offset 5",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[6].value},
				{custTable.rows[7].value},
				{custTable.rows[8].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" is "Mountain View".
			"select v from Customer where v.Address.City = \"Mountain View\"",
			[]string{"v"},
			[][]*vom.RawBytes{},
		},
		{
			// Select records where v.Address.City = "Collins" or type is Invoice.
			// Offset 8
			"select v from Customer where v.Address.City = \"Collins\" or Type(v) like \"%.Invoice\" offset 8",
			[]string{"v"},
			[][]*vom.RawBytes{},
		},
		{
			// Select records where v.Address.City = "Collins" or type is Invoice.
			// Offset 23
			"select v from Customer where v.Address.City = \"Collins\" or Type(v) like \"%.Invoice\" offset 23",
			[]string{"v"},
			[][]*vom.RawBytes{},
		},
		{
			// Select records where v.Address.City = "Collins" is 84 or type is Invoice.
			// Limit 3 Offset 2
			"select v from Customer where v.Address.City = \"Collins\" or Type(v) like \"%.Invoice\" limit 3 offset 2",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[3].value},
				{custTable.rows[4].value},
				{custTable.rows[5].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" or (type is Invoice and v.InvoiceNum is not nil).
			// Limit 3 Offset 2
			"select v from Customer where v.Address.City = \"Collins\" or (Type(v) like \"%.Invoice\" and v.InvoiceNum is not nil) limit 3 offset 2",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[3].value},
				{custTable.rows[4].value},
				{custTable.rows[5].value},
			},
		},
		{
			// Select records where v.Address.City = "Collins" or (type is Invoice and v.InvoiceNum is nil).
			// Limit 3 Offset 2
			"select v from Customer where v.Address.City = \"Collins\" or (Type(v) like \"%.Invoice\" and v.InvoiceNum is nil) limit 3 offset 2",
			[]string{"v"},
			[][]*vom.RawBytes{},
		},
		// Test functions.
		{
			// Select invoice records where date is 2015-03-17
			"select v from Customer where Type(v) like \"%.Invoice\" and Year(v.InvoiceDate, \"America/Los_Angeles\") = 2015 and Month(v.InvoiceDate, \"America/Los_Angeles\") = 3 and Day(v.InvoiceDate, \"America/Los_Angeles\") = 17",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[5].value},
				{custTable.rows[6].value},
			},
		},
		{
			// Now will always be > 2012, so all customer records will be returned.
			"select v from Customer where Now() > Time(\"2006-01-02 MST\", \"2012-03-17 PDT\")",
			[]string{"v"},
			[][]*vom.RawBytes{
				{custTable.rows[0].value},
				{custTable.rows[1].value},
				{custTable.rows[2].value},
				{custTable.rows[3].value},
				{custTable.rows[4].value},
				{custTable.rows[5].value},
				{custTable.rows[6].value},
				{custTable.rows[7].value},
				{custTable.rows[8].value},
				{custTable.rows[9].value},
			},
		},
		{
			// Select April 2015 PT invoices.
			// Note: this wouldn't work for March as daylight saving occurs March 8
			// and causes comparisons for those days to be off 1 hour.
			// It would work to use UTC -- see next test.
			"select k from Customer where Year(v.InvoiceDate, \"America/Los_Angeles\") = 2015 and Month(v.InvoiceDate, \"America/Los_Angeles\") = 4",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[7].key)},
				{vom.RawBytesOf(custTable.rows[8].key)},
			},
		},
		{
			// Select March 2015 UTC invoices.
			"select k from Customer where Year(v.InvoiceDate, \"UTC\") = 2015 and Month(v.InvoiceDate, \"UTC\") = 3",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[3].key)},
				{vom.RawBytesOf(custTable.rows[5].key)},
				{vom.RawBytesOf(custTable.rows[6].key)},
			},
		},
		{
			// Select 2015 UTC invoices.
			"select k from Customer where Year(v.InvoiceDate, \"UTC\") = 2015",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[1].key)},
				{vom.RawBytesOf(custTable.rows[2].key)},
				{vom.RawBytesOf(custTable.rows[3].key)},
				{vom.RawBytesOf(custTable.rows[5].key)},
				{vom.RawBytesOf(custTable.rows[6].key)},
				{vom.RawBytesOf(custTable.rows[7].key)},
				{vom.RawBytesOf(custTable.rows[8].key)},
			},
		},
		{
			// Select the Mar 17 2015 11:14:04 America/Los_Angeles invoice.
			"select k from Customer where v.InvoiceDate = Time(\"2006-01-02 15:04:05 MST\", \"2015-03-17 11:14:04 PDT\")",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[5].key)},
			},
		},
		{
			// Select invoices in the minute Mar 17 2015 11:14 America/Los_Angeles invoice.
			"select k from Customer where Year(v.InvoiceDate, \"America/Los_Angeles\") = 2015 and Month(v.InvoiceDate, \"America/Los_Angeles\") = 3 and Day(v.InvoiceDate, \"America/Los_Angeles\") = 17 and Hour(v.InvoiceDate, \"America/Los_Angeles\") = 11 and Minute(v.InvoiceDate, \"America/Los_Angeles\") = 14",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[5].key)},
			},
		},
		{
			// Select invoices in the hour Mar 17 2015 11 hundred America/Los_Angeles invoice.
			"select k from Customer where Year(v.InvoiceDate, \"America/Los_Angeles\") = 2015 and Month(v.InvoiceDate, \"America/Los_Angeles\") = 3 and Day(v.InvoiceDate, \"America/Los_Angeles\") = 17 and Hour(v.InvoiceDate, \"America/Los_Angeles\") = 11",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[5].key)},
			},
		},
		{
			// Select invoices on the day Mar 17 2015 America/Los_Angeles invoice.
			"select k from Customer where Year(v.InvoiceDate, \"America/Los_Angeles\") = 2015 and Month(v.InvoiceDate, \"America/Los_Angeles\") = 3 and Day(v.InvoiceDate, \"America/Los_Angeles\") = 17",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[5].key)},
				{vom.RawBytesOf(custTable.rows[6].key)},
			},
		},
		{
			"select Str(v.Address) from Customer where Str(v.Id) = \"1\"",
			[]string{"Str"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.AddressInfo struct{Street string;City string;State string;Zip string}{Street: \"1 Main St.\", City: \"Palo Alto\", State: \"CA\", Zip: \"94303\"}")},
			},
		},
		// Test string functions in where clause.
		{
			"select k from Customer where Str(v.Id) = \"1\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[0].key)},
			},
		},
		{
			// Select invoices shipped to Any street -- using Lowercase.
			"select k from Customer where Type(v) like \"%.Invoice\" and Lowercase(v.ShipTo.Street) like \"%any%\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[5].key)},
				{vom.RawBytesOf(custTable.rows[6].key)},
				{vom.RawBytesOf(custTable.rows[7].key)},
				{vom.RawBytesOf(custTable.rows[8].key)},
			},
		},
		{
			"select HtmlEscape(\"<a img='foo'>Foo Image</a>\") from Customer where k = \"001\"",
			[]string{"HtmlEscape"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("&lt;a img=&#39;foo&#39;&gt;Foo Image&lt;/a&gt;")},
			},
		},
		{
			"select HtmlUnescape(\"&lt;a img=&#39;foo&#39;&gt;Foo Image&lt;/a&gt;\") from Customer where k = \"001\"",
			[]string{"HtmlUnescape"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("<a img='foo'>Foo Image</a>")},
			},
		},
		{
			// Select invoices shipped to Any street -- using Uppercase.
			"select k from Customer where Type(v) like \"%.Invoice\" and Uppercase(v.ShipTo.Street) like \"%ANY%\"",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(custTable.rows[5].key)},
				{vom.RawBytesOf(custTable.rows[6].key)},
				{vom.RawBytesOf(custTable.rows[7].key)},
				{vom.RawBytesOf(custTable.rows[8].key)},
			},
		},
		// Select clause functions.
		// Time function
		{
			"select Time(\"2006-01-02 MST\", \"2015-07-01 PDT\") from Customer",
			[]string{"Time"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
				{vom.RawBytesOf(t2015_07_01)},
			},
		},
		// Time function
		{
			"select Time(\"2006-01-02 15:04:05 MST\", \"2015-07-01 01:23:45 PDT\") from Customer",
			[]string{"Time"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
				{vom.RawBytesOf(t2015_07_01_01_23_45)},
			},
		},
		// Lowercase function
		{
			"select Lowercase(v.Name) as name from Customer where Type(v) like \"%.Customer\"",
			[]string{"name"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("john smith")},
				{vom.RawBytesOf("bat masterson")},
				{vom.RawBytesOf("john steed")},
			},
		},
		// Uppercase function
		{
			"select Uppercase(v.Name) as NAME from Customer where Type(v) like \"%.Customer\"",
			[]string{"NAME"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("JOHN SMITH")},
				{vom.RawBytesOf("BAT MASTERSON")},
				{vom.RawBytesOf("JOHN STEED")},
			},
		},
		// Second function
		{
			"select k, Second(v.InvoiceDate, \"America/Los_Angeles\") from Customer where Type(v) like \"%.Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"Second",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002003"), vom.RawBytesOf(int64(6))},
			},
		},
		// Minute function
		{
			"select k, Minute(v.InvoiceDate, \"America/Los_Angeles\") from Customer where Type(v) like \"%.Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"Minute",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002003"), vom.RawBytesOf(int64(16))},
			},
		},
		// Hour function
		{
			"select k, Hour(v.InvoiceDate, \"America/Los_Angeles\") from Customer where Type(v) like \"%.Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"Hour",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002003"), vom.RawBytesOf(int64(22))},
			},
		},
		// Day function
		{
			"select k, Day(v.InvoiceDate, \"America/Los_Angeles\") from Customer where Type(v) like \"%.Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"Day",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002003"), vom.RawBytesOf(int64(12))},
			},
		},
		// Month function
		{
			"select k, Month(v.InvoiceDate, \"America/Los_Angeles\") from Customer where Type(v) like \"%.Invoice\" and k = \"002003\"",
			[]string{
				"k",
				"Month",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002003"), vom.RawBytesOf(int64(4))},
			},
		},
		// Year function
		{
			"select k, Year(v.InvoiceDate, \"America/Los_Angeles\") from Customer where Type(v) like \"%.Invoice\" and k = \"001001\"",
			[]string{
				"k",
				"Year",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001001"), vom.RawBytesOf(int64(2015))},
			},
		},
		// Nested functions
		{
			"select Year(Time(\"2006-01-02 15:04:05 MST\", \"2015-07-01 01:23:45 PDT\"), \"America/Los_Angeles\")  from Customer where Type(v) like \"%.Invoice\" and k = \"001001\"",
			[]string{"Year"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(2015))},
			},
		},
		// Bad arg to function.  Expression is false.
		{
			"select v from Customer where Type(v) like \"%.Invoice\" and Day(v.InvoiceDate, v.Foo) = v.InvoiceDate",
			[]string{"v"},
			[][]*vom.RawBytes{},
		},
		// Map in selection
		{
			"select v.Credit.Report.TransUnionReport.PreviousRatings[\"2015Q2\"] from Customer where v.Name = \"Bat Masterson\"",
			[]string{"v.Credit.Report.TransUnionReport.PreviousRatings[2015Q2]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int16(40))},
			},
		},
		// Map in selection using function as key.
		{
			"select v.Credit.Report.TransUnionReport.PreviousRatings[Uppercase(\"2015q2\")] from Customer where v.Name = \"Bat Masterson\"",
			[]string{"v.Credit.Report.TransUnionReport.PreviousRatings[Uppercase]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int16(40))},
			},
		},
		// Map in selection using struct as key.
		{
			"select v.Map[v.Key] from FunWithMaps",
			[]string{"v.Map[v.Key]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(td.V{"ccc", 14.7})},
				{vom.RawBytesOf(td.V{"yyy", 17.1})},
			},
		},
		// map of int16 to array of sets of strings
		{
			"select v.Confusing[23][0][\"foo\"] from FunWithMaps",
			[]string{"v.Confusing[23][0][foo]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(true)},
				{vom.RawBytesOf(nil)},
			},
		},
		// FunWithTypes
		{
			"select Type(v.T1), Type(v.T2) from FunWithTypes",
			[]string{"Type", "Type"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("typeobject"), vom.RawBytesOf("typeobject")},
				{vom.RawBytesOf("typeobject"), vom.RawBytesOf("typeobject")},
			},
		},
		// FunWithTypes
		{
			"select Str(v.T1), Str(v.T2) from FunWithTypes",
			[]string{"Str", "Str"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer"), vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer"), vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer")},
			},
		},
		// FunWithTypes
		{
			"select Str(v.T1), Str(v.T2) from FunWithTypes where v.T1 = v.T2",
			[]string{"Str", "Str"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer"), vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer")},
			},
		},
		// FunWithTypes
		{
			"select Str(v.T1), Str(v.T2) from FunWithTypes where Str(v.T1) = Str(v.T2)",
			[]string{"Str", "Str"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer"), vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer")},
			},
		},
		// FunWithTypes
		{
			"select Str(v.T1), Str(v.T2) from FunWithTypes where Type(v.T1) = Type(v.T2)",
			[]string{"Str", "Str"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer"), vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Invoice")},
				{vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer"), vom.RawBytesOf("v.io/v23/query/engine/internal/testdata.Customer")},
			},
		},
		// Function using a map lookup as arg
		{
			"select Uppercase(v.B[true]) from ManyMaps",
			[]string{"Uppercase"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("IT WAS THE BEST OF TIMES,")},
			},
		},
		// Set in selection
		{
			"select v.Credit.Report.ExperianReport.TdhApprovals[\"Tom\"], v.Credit.Report.ExperianReport.TdhApprovals[\"Dick\"], v.Credit.Report.ExperianReport.TdhApprovals[\"Harry\"] from Customer where v.Name = \"John Steed\"",
			[]string{
				"v.Credit.Report.ExperianReport.TdhApprovals[Tom]",
				"v.Credit.Report.ExperianReport.TdhApprovals[Dick]",
				"v.Credit.Report.ExperianReport.TdhApprovals[Harry]",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(true), vom.RawBytesOf(false), vom.RawBytesOf(true)},
			},
		},
		// List in selection
		{
			"select v.PreviousAddresses[0].Street from Customer where v.Name = \"Bat Masterson\"",
			[]string{"v.PreviousAddresses[0].Street"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("19 Green St.")},
			},
		},
		// List in selection (index out of bounds)
		{
			"select v.PreviousAddresses[2].Street from Customer where v.Name = \"Bat Masterson\"",
			[]string{"v.PreviousAddresses[2].Street"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(nil)},
			},
		},
		// Array in selection
		{
			"select v.Credit.Report.EquifaxReport.FourScoreRatings[2] from Customer where v.Name = \"John Smith\"",
			[]string{"v.Credit.Report.EquifaxReport.FourScoreRatings[2]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int16(42))},
			},
		},
		// Array in selection (using an array as the index)
		// Note: v.Credit.Report.EquifaxReport.FourScoreRatings[3] is 2
		// and v.Credit.Report.EquifaxReport.FourScoreRatings[2] is 42
		{
			"select v.Credit.Report.EquifaxReport.FourScoreRatings[v.Credit.Report.EquifaxReport.FourScoreRatings[3]] from Customer where v.Name = \"John Smith\"",
			[]string{"v.Credit.Report.EquifaxReport.FourScoreRatings[v.Credit.Report.EquifaxReport.FourScoreRatings[3]]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int16(42))},
			},
		},
		// Array in selection (index out of bounds)
		{
			"select v.Credit.Report.EquifaxReport.FourScoreRatings[4] from Customer where v.Name = \"John Smith\"",
			[]string{"v.Credit.Report.EquifaxReport.FourScoreRatings[4]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(nil)},
			},
		},
		// Map in where expression
		{
			"select v.Name from Customer where v.Credit.Report.TransUnionReport.PreviousRatings[\"2015Q2\"] = 40",
			[]string{"v.Name"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Bat Masterson")},
			},
		},
		// Set in where expression (convert string to enum to do lookup)
		{
			"select v.Name from Customer where v.Credit.Report.ExperianReport.TdhApprovals[\"Tom\"] = true",
			[]string{
				"v.Name",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("John Steed")},
			},
		},
		// Negative case: Set in where expression (convert string to enum to do lookup)
		{
			"select v.Name from Customer where v.Credit.Report.ExperianReport.TdhApprovals[\"Dick\"] = true",
			[]string{
				"v.Name",
			},
			[][]*vom.RawBytes{},
		},
		// Set in where expression (use another field as lookup key)
		// Find all customers where experian auditor was also an approver.
		{
			"select v.Name from Customer where v.Credit.Report.ExperianReport.TdhApprovals[v.Credit.Report.ExperianReport.Auditor] = true",
			[]string{
				"v.Name",
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("John Steed")},
			},
		},
		// List in where expression
		{
			"select v.Name from Customer where v.PreviousAddresses[0].Street = \"19 Green St.\"",
			[]string{"v.Name"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Bat Masterson")},
			},
		},
		// List in where expression (index out of bounds)
		{
			"select v.Name from Customer where v.PreviousAddresses[10].Street = \"19 Green St.\"",
			[]string{"v.Name"},
			[][]*vom.RawBytes{},
		},
		// Array in where expression
		{
			"select v.Name from Customer where v.Credit.Report.EquifaxReport.FourScoreRatings[2] = 42",
			[]string{"v.Name"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("John Smith")},
			},
		},
		// Array in where expression (using another field as index)
		// Note: v.Credit.Report.EquifaxReport.FourScoreRatings[3] is 2
		// and v.Credit.Report.EquifaxReport.FourScoreRatings[2] is 42
		{
			"select v.Name from Customer where v.Credit.Report.EquifaxReport.FourScoreRatings[v.Credit.Report.EquifaxReport.FourScoreRatings[3]] = 42",
			[]string{"v.Name"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("John Smith")},
			},
		},
		// Array in where expression (index out of bounds, using another field as index)
		{
			"select v.Name from Customer where v.Credit.Report.EquifaxReport.FourScoreRatings[v.Credit.Report.EquifaxReport.FourScoreRatings[2]] = 42",
			[]string{"v.Name"},
			[][]*vom.RawBytes{},
		},
		// Array in select and where expressions (top level value is the array)
		{
			"select v[0], v[1], v[2], v[3] from RatingsArray where v[0] = 40",
			[]string{"v[0]", "v[1]", "v[2]", "v[3]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int16(40)), vom.RawBytesOf(int16(20)), vom.RawBytesOf(int16(10)), vom.RawBytesOf(int16(0))},
			},
		},
		// List in select and where expressions (top level value is the list)
		{
			"select v[-1].City, v[0].City, v[1].City, v[2].City from PreviousAddresses where v[1].Street = \"200 Main St.\"",
			[]string{"v[-1].City", "v[0].City", "v[1].City", "v[2].City"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(nil), vom.RawBytesOf("Anytown"), vom.RawBytesOf("Othertown"), vom.RawBytesOf(nil)},
			},
		},
		// Set in select and where expressions (top level value is the set)
		{
			"select v[\"Tom\"], v[\"Dick\"], v[\"Harry\"] from TdhApprovals where v[\"Dick\"] = true",
			[]string{"v[Tom]", "v[Dick]", "v[Harry]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(false), vom.RawBytesOf(true), vom.RawBytesOf(true)},
			},
		},
		// Map in select and where expressions (top level value is the map)
		{
			"select v[\"1Q2015\"], v[\"2Q2015\"] from PreviousRatings where v[\"2Q2015\"] = 3",
			[]string{"v[1Q2015]", "v[2Q2015]"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(nil), vom.RawBytesOf(int16(3))},
			},
		},
		// Test lots of types as map keys
		{
			"select v.B[true], v.By[10], v.U16[16], v.U32[32], v.U64[64], v.I16[17], v.I32[33], v.I64[65], v.F32[32.1], v.F64[64.2], v.S[\"Dickens\"], v.Ms[\"Charles\"][\"Dickens\"], v.T[Time(\"2006-01-02 15:04:05 MST\", \"2015-07-01 01:23:45 PDT\")] from ManyMaps",
			[]string{"v.B[true]", "v.By[10]", "v.U16[16]", "v.U32[32]", "v.U64[64]", "v.I16[17]", "v.I32[33]", "v.I64[65]", "v.F32[32.1]", "v.F64[64.2]", "v.S[Dickens]", "v.Ms[Charles][Dickens]", "v.T[Time]"},
			[][]*vom.RawBytes{
				[]*vom.RawBytes{
					vom.RawBytesOf("It was the best of times,"),
					vom.RawBytesOf("it was the worst of times,"),
					vom.RawBytesOf("it was the age of wisdom,"),
					vom.RawBytesOf("it was the age of foolishness,"),
					vom.RawBytesOf("it was the epoch of belief,"),
					vom.RawBytesOf("it was the epoch of incredulity,"),
					vom.RawBytesOf("it was the season of Light,"),
					vom.RawBytesOf("it was the season of Darkness,"),
					vom.RawBytesOf("it was the spring of hope,"),
					vom.RawBytesOf("it was the winter of despair,"),
					vom.RawBytesOf("we are all going direct to Heaven,"),
					vom.RawBytesOf("we are all going direct to Heaven,"),
					vom.RawBytesOf("we are all going direct the other way"),
				},
			},
		},
		// Test lots of types as set keys
		{
			"select v.B[true], v.By[10], v.U16[16], v.U32[32], v.U64[64], v.I16[17], v.I32[33], v.I64[65], v.F32[32.1], v.F64[64.2], v.S[\"Dickens\"], v.T[Time(\"2006-01-02 15:04:05 MST\", \"2015-07-01 01:23:45 PDT\")] from ManySets",
			[]string{"v.B[true]", "v.By[10]", "v.U16[16]", "v.U32[32]", "v.U64[64]", "v.I16[17]", "v.I32[33]", "v.I64[65]", "v.F32[32.1]", "v.F64[64.2]", "v.S[Dickens]", "v.T[Time]"},
			[][]*vom.RawBytes{
				[]*vom.RawBytes{
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
					vom.RawBytesOf(true),
				},
			},
		},
		{
			"select k, v.Key from BigTable where k < \"101\" or k = \"200\" or k like \"300%\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{svPair("100"), svPair("200"), svPair("300")},
		},
		{
			"select k, v.Key from BigTable where \"101\" > k or \"200\" = k or k like \"300%\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{svPair("100"), svPair("200"), svPair("300")},
		},
		{
			"select k, v.Key from BigTable where k is nil",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where k is not nil and \"103\" = v.Key",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{svPair("103")},
		},
		{
			"select k, v.Key from BigTable where k like \"10_\" or k like \"20_\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("100"),
				svPair("101"),
				svPair("102"),
				svPair("103"),
				svPair("104"),
				svPair("105"),
				svPair("106"),
				svPair("107"),
				svPair("108"),
				svPair("109"),
				svPair("200"),
				svPair("201"),
				svPair("202"),
				svPair("203"),
				svPair("204"),
				svPair("205"),
				svPair("206"),
				svPair("207"),
				svPair("208"),
				svPair("209"),
			},
		},
		{
			"select k, v.Key from BigTable where k like \"_%9\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("109"),
				svPair("119"),
				svPair("129"),
				svPair("139"),
				svPair("149"),
				svPair("159"),
				svPair("169"),
				svPair("179"),
				svPair("189"),
				svPair("199"),
				svPair("209"),
				svPair("219"),
				svPair("229"),
				svPair("239"),
				svPair("249"),
				svPair("259"),
				svPair("269"),
				svPair("279"),
				svPair("289"),
				svPair("299"),
			},
		},
		{
			"select k, v.Key from BigTable where k like \"__0\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("100"),
				svPair("110"),
				svPair("120"),
				svPair("130"),
				svPair("140"),
				svPair("150"),
				svPair("160"),
				svPair("170"),
				svPair("180"),
				svPair("190"),
				svPair("200"),
				svPair("210"),
				svPair("220"),
				svPair("230"),
				svPair("240"),
				svPair("250"),
				svPair("260"),
				svPair("270"),
				svPair("280"),
				svPair("290"),
				svPair("300"),
			},
		},
		{
			"select k, v.Key from BigTable where k like \"10%\" or  k like \"20%\" or  k like \"30%\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("100"),
				svPair("101"),
				svPair("102"),
				svPair("103"),
				svPair("104"),
				svPair("105"),
				svPair("106"),
				svPair("107"),
				svPair("108"),
				svPair("109"),
				svPair("200"),
				svPair("201"),
				svPair("202"),
				svPair("203"),
				svPair("204"),
				svPair("205"),
				svPair("206"),
				svPair("207"),
				svPair("208"),
				svPair("209"),
				svPair("300"),
			},
		},
		{
			"select k, v.Key from BigTable where k like \"1__\" and  k like \"_2_\" and  k like \"__3\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{svPair("123")},
		},
		{
			"select k, v.Key from BigTable where (k >  \"100\" and k < \"103\") or (k > \"205\" and k < \"208\")",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("101"),
				svPair("102"),
				svPair("206"),
				svPair("207"),
			},
		},
		{
			"select k, v.Key from BigTable where ( \"100\" < k and \"103\" > k) or (\"205\" < k and \"208\" > k)",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("101"),
				svPair("102"),
				svPair("206"),
				svPair("207"),
			},
		},
		{
			"select k, v.Key from BigTable where k <=  \"100\" or k = \"101\" or k >= \"300\" or (k <> \"299\" and k not like \"300\" and k >= \"298\")",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("100"),
				svPair("101"),
				svPair("298"),
				svPair("300"),
			},
		},
		{
			"select k, v.Key from BigTable where \"100\" >= k or \"101\" = k or \"300\" <= k or (\"299\" <> k and k not like \"300\" and \"298\" <= k)",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("100"),
				svPair("101"),
				svPair("298"),
				svPair("300"),
			},
		},
		{
			"select k, v.Key from BigTable where k like  \"1%\" and k like \"%9\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{
				svPair("109"),
				svPair("119"),
				svPair("129"),
				svPair("139"),
				svPair("149"),
				svPair("159"),
				svPair("169"),
				svPair("179"),
				svPair("189"),
				svPair("199"),
			},
		},
		{
			"select k, v.Key from BigTable where k like  \"3%\" and k like \"30%\" and k like \"300%\"",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{svPair("300")},
		},
		{
			"select k, v.Key from BigTable where \"110\" > k",
			[]string{"k", "v.Key"},
			svPairs(100, 109),
		},
		{
			"select k, v.Key from BigTable where \"110\" < k and \"205\" > k",
			[]string{"k", "v.Key"},
			svPairs(111, 204),
		},
		{
			"select k, v.Key from BigTable where \"110\" <= k and \"205\" >= k",
			[]string{"k", "v.Key"},
			svPairs(110, 205),
		},
		{
			"select k, v.Key from BigTable where k is nil",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where k is not nil",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where k <> k",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where k < k",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where k > k",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where k = k",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where k <= k",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where k >= k",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where k = v.Key",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where v.Key = k",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where v.Key = k",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where k <> v.key",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where v.key <> k",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where k < v.key",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where v.key < k",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where k > v.key",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where v.key > k",
			[]string{"k", "v.Key"},
			[][]*vom.RawBytes{},
		},
		{
			"select k, v.Key from BigTable where k >= v.Key",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where v.Key >= k",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where k <= v.Key",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			"select k, v.Key from BigTable where v.Key <= k",
			[]string{"k", "v.Key"},
			svPairs(100, 300),
		},
		{
			// Split on .
			"select Split(Type(v), \".\") from Customer where k = \"001\"",
			[]string{"Split"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf([]string{"v", "io/v23/query/engine/internal/testdata", "Customer"})},
			},
		},
		{
			// Split on /
			"select Split(Type(v), \"/\") from Customer where k = \"001\"",
			[]string{"Split"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf([]string{"v.io", "v23", "query", "engine", "internal", "testdata.Customer"})},
			},
		},
		{
			// Split on /, Len of array
			"select Len(Split(Type(v), \"/\")) from Customer where k = \"001\"",
			[]string{"Len"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(6))},
			},
		},
		{
			// Split on empty string, Len of array
			// Split with sep == empty string splits on chars.
			"select Len(Split(Type(v), \"\")) from Customer where k = \"001\"",
			[]string{"Len"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(len("v.io/v23/query/engine/internal/testdata.Customer")))},
			},
		},
		{
			// Len of string, list and struct.
			// returns len of string, elements in list and nil for struct.
			"select Len(v.Name), Len(v.PreviousAddresses), Len(v.Credit) from Customer where k = \"002\"",
			[]string{"Len", "Len", "Len"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(13)), vom.RawBytesOf(int64(2)), vom.RawBytesOf(nil)},
			},
		},
		{
			// Len of set
			"select Len(v.Credit.Report.ExperianReport.TdhApprovals) from Customer where Type(v) like \"%.Customer\" and k = \"003\"",
			[]string{"Len"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(2))},
			},
		},
		{
			// Len of map
			"select Len(v.Map) from FunWithMaps",
			[]string{"Len"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(2))},
				{vom.RawBytesOf(int64(2))},
			},
		},
		{
			// Ceiling
			"select k, Ceiling(v.F64) from Numbers where k = \"001\"",
			[]string{"k", "Ceiling"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(3.0))},
			},
		},
		{
			// Floor
			"select k, Floor(v.F64) from Numbers where k = \"001\"",
			[]string{"k", "Floor"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(2.0))},
			},
		},
		{
			// Truncate
			"select k, Truncate(v.F64) from Numbers where k = \"001\"",
			[]string{"k", "Truncate"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(2.0))},
			},
		},
		{
			// Log
			"select k, Log(v.F64) from Numbers where k = \"001\"",
			[]string{"k", "Log"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(1.0000000000003513))},
			},
		},
		{
			// Log10
			"select k, Log10(v.F64) from Numbers where k = \"001\"",
			[]string{"k", "Log10"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(0.43429448190340436))},
			},
		},
		{
			// Pow
			"select k, Pow(10.0, 4.0) from Numbers where k = \"001\"",
			[]string{"k", "Pow"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(10000))},
			},
		},
		{
			// Pow10
			"select k, Pow10(5) from Numbers where k = \"001\"",
			[]string{"k", "Pow10"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(100000))},
			},
		},
		{
			// Mod
			"select k, Mod(v.F32, v.F64) from Numbers where k = \"001\"",
			[]string{"k", "Mod"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(0.42330828994820324))},
			},
		},
		{
			// Remainder
			"select k, Remainder(v.F32, v.F64) from Numbers where k = \"001\"",
			[]string{"k", "Remainder"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001"), vom.RawBytesOf(float64(0.42330828994820324))},
			},
		},
		{
			// Sprintf
			"select Sprintf(\"%d, %d, %d, %d, %d, %d, %d, %g, %g\", v.B, v.Ui16, v.Ui32, v.Ui64, v.I16, v.I32, v.I64, v.F32, v.F64) from Numbers where k = \"001\"",
			[]string{"Sprintf"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("12, 1234, 5678, 999888777666, 9876, 876543, 128, 3.141590118408203, 2.71828182846")},
			},
		},
		{
			// Sprintf
			"select Sprintf(\"%s, %s %s\", v.Address.City, v.Address.State, v.Address.Zip) from Customer where v.Name = \"John Smith\"",
			[]string{"Sprintf"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Palo Alto, CA 94303")},
			},
		},
		{
			// StrCat
			"select StrCat(v.Address.City, v.Address.State) from Customer where v.Name = \"John Smith\"",
			[]string{"StrCat"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Palo AltoCA")},
			},
		},
		{
			// StrCat
			"select StrCat(v.Address.City, \", \", v.Address.State) from Customer where v.Name = \"John Smith\"",
			[]string{"StrCat"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Palo Alto, CA")},
			},
		},
		{
			// StrCat
			"select StrCat(v.Address.City, \"42\") from Customer where v.Name = \"John Smith\"",
			[]string{"StrCat"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Palo Alto42")},
			},
		},
		{
			// StrIndex
			"select StrIndex(v.Address.City, \"lo\") from Customer where v.Name = \"John Smith\"",
			[]string{"StrIndex"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(2))},
			},
		},
		{
			// StrLastIndex
			"select StrLastIndex(v.Address.City, \"l\") from Customer where v.Name = \"John Smith\"",
			[]string{"StrLastIndex"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(int64(6))},
			},
		},
		{
			// StrRepeat
			"select StrRepeat(v.Address.City, 3) from Customer where v.Name = \"John Smith\"",
			[]string{"StrRepeat"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Palo AltoPalo AltoPalo Alto")},
			},
		},
		{
			// StrReplace
			"select StrReplace(v.Address.City, \"Palo\", \"Shallow\") from Customer where v.Name = \"John Smith\"",
			[]string{"StrReplace"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Shallow Alto")},
			},
		},
		{
			// Trim
			"select Trim(\"   Foo     \") from Customer where v.Name = \"John Smith\"",
			[]string{"Trim"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Foo")},
			},
		},
		{
			// TrimLeft
			"select TrimLeft(\"   Foo     \") from Customer where v.Name = \"John Smith\"",
			[]string{"TrimLeft"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("Foo     ")},
			},
		},
		{
			// TrimRight
			"select TrimRight(\"   Foo     \") from Customer where v.Name = \"John Smith\"",
			[]string{"TrimRight"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("   Foo")},
			},
		},
		{
			// RuneCount
			"select Len(\"Hello, 世界\"), RuneCount(\"Hello, 世界\") from Customer where v.Name = \"John Smith\"",
			[]string{"Len", "RuneCount"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(13), vom.RawBytesOf(9)},
			},
		},
	}

	qe := engine.Create(db)

	for _, test := range basic {
		headers, rs, err := qe.Exec(test.query)
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, err)
		} else {
			// Collect results.
			r := [][]*vom.RawBytes{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if got, want := vdl.ValueOf(r), vdl.ValueOf(test.r); !vdl.EqualValue(got, want) {
				t.Errorf("query: %s; got %v, want %v", test.query, got, want)
			}
			if !reflect.DeepEqual(test.headers, headers) {
				t.Errorf("query: %s; got %v, want %v", test.query, headers, test.headers)
			}
		}
	}

	// Do the same thing with a prepared statement.  (Even though there are no parameters, it's
	// good to know that all will still work.  Supplying parameters is a separate test.)
	// For good measure, exersize the Handle() and GetPreparedStatement functions.
	for _, test := range basic {
		p, err := qe.PrepareStatement(test.query)
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, err)
		}
		h := p.Handle()
		p2, err := qe.GetPreparedStatement(h)

		headers, rs, err := p2.Exec()
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, err)
		} else {
			// Collect results.
			r := [][]*vom.RawBytes{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if got, want := vdl.ValueOf(r), vdl.ValueOf(test.r); !vdl.EqualValue(got, want) {
				t.Errorf("query: %s; got %v, want %v", test.query, r, test.r)
			}
			if !reflect.DeepEqual(test.headers, headers) {
				t.Errorf("query: %s; got %v, want %v", test.query, headers, test.headers)
			}
		}
		p2.Close()
	}
}

func TestDelete(t *testing.T) {
	basic := []execDeleteTest{
		{
			// Multi-delete based on key (using like).
			"delete from Customer where k like \"001%\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(4)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Single delete based on key.
			"delete from Customer where k = \"001\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(1)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete based on value satisfying where clause.
			"delete from Customer where Type(v) like \"%Customer\" and v.Name like \"%Masterson\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(1)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete based on value satisfying where clause WITH LIMIT.
			"delete from Customer where Type(v) like \"%Invoice\" and v.CustId = 1 limit 2",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(2)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete all type Customer k/v pairs.
			"delete from Customer where Type(v) like \"%Customer\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(3)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			// Delete Customers that have an equifax credit report.
			"delete from Customer where Type(v) like \"%Customer\" and Type(v.Credit.Report.EquifaxReport) like \"%.EquifaxCreditReport\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(1)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete k/v pairs where InvoiceNum is nil (i.e., delete Customer types)
			"delete from Customer where v.InvoiceNum is nil",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(3)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			// Delete k/v pairs where InvoiceNum is NOT nil (i.e., delete Invoice types)
			"delete from Customer where v.InvoiceNum is not nil",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(7)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete k/v pairs where InvoiceNum is nil and Name is not nil (i.e., delete Customer types)
			"delete from Customer where v.InvoiceNum is nil and v.Name is not nil",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(3)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			// Delete invoices where amount is $88.
			"delete from Customer where Type(v) like \"%.Invoice\" and v.Amount = 88",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(2)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete everything except those associated with Customer 002.
			"delete from Customer where k not like \"002%\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(5)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			// Play with mixed case keywords
			// Delete everything except those associated with Customer 002.
			"dElEtE fRoM Customer wHeRe k NoT lIkE \"002%\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(5)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			// Play with white space.
			// Delete everything except those associated with Customer 002.
			"      delete    from    Customer    where \n k   not \n like \"002%\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(5)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			// Delete a key range
			"delete from Customer where k >= \"002002\" and k <= \"002003\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(2)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete three keys
			"delete from Customer where k = \"001002\" or k = \"002003\" or k = \"003\"",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(3)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			// Delete customers where with a bad credit rating (leaves the invoices alone).
			"delete from Customer where v.Credit.Report.EquifaxReport.Rating < 'A' or v.Credit.Report.ExperianReport.Rating = \"Bad\" or v.Credit.Report.TransUnionReport.Rating < 90",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(1)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete invoices where date is 2015-03-17
			"delete from Customer where Type(v) like \"%.Invoice\" and Year(v.InvoiceDate, \"America/Los_Angeles\") = 2015 and Month(v.InvoiceDate, \"America/Los_Angeles\") = 3 and Day(v.InvoiceDate, \"America/Los_Angeles\") = 17",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(2)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete all April 2015 invoices.
			"delete from Customer where Type(v) like \"%.Invoice\" and Year(v.InvoiceDate, \"America/Los_Angeles\") = 2015 and Month(v.InvoiceDate, \"America/Los_Angeles\") = 4",
			[]string{"Count"},
			[][]*vom.RawBytes{{vom.RawBytesOf(2)}},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("003")},
			},
		},
	}

	for _, test := range basic {
		initTables()

		qe := engine.Create(db)
		// Delete
		headers, rs, err := qe.Exec(test.delQuery)
		if err != nil {
			t.Errorf("delQuery: %s; got %v, want nil", test.delQuery, err)
		} else {
			// Collect results.
			r := [][]*vom.RawBytes{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if !reflect.DeepEqual(test.delResults, r) {
				t.Errorf("delQuery: %s; got %v, want %v", test.delQuery, r, test.delResults)
			}
			if !reflect.DeepEqual(test.delHeaders, headers) {
				t.Errorf("delQuery: %s; got %v, want %v", test.delQuery, headers, test.delHeaders)
			}
		}
		// Select
		headers, rs, err = qe.Exec(test.selQuery)
		if err != nil {
			t.Errorf("selQuery: %s; got %v, want nil", test.selQuery, err)
		} else {
			// Collect results.
			r := [][]*vom.RawBytes{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if !reflect.DeepEqual(test.selResults, r) {
				t.Errorf("delQuery: %s; got %v, want %v", test.delQuery, r, test.selResults)
			}
			if !reflect.DeepEqual(test.selHeaders, headers) {
				t.Errorf("delQuery: %s; got %v, want %v", test.delQuery, headers, test.selHeaders)
			}
		}
	}
}

func svPair(s string) []*vom.RawBytes {
	v := vom.RawBytesOf(s)
	return []*vom.RawBytes{v, v}
}

// Genearate k,v pairs for start to finish (*INCLUSIVE*)
func svPairs(start, finish int64) [][]*vom.RawBytes {
	retVal := [][]*vom.RawBytes{}
	for i := start; i <= finish; i++ {
		v := vom.RawBytesOf(fmt.Sprintf("%d", i))
		retVal = append(retVal, []*vom.RawBytes{v, v})
	}
	return retVal
}

func TestQueryPrepare(t *testing.T) {
	initTables()
	basic := []prepareSelectTest{
		{
			// Select based on parameterized key
			"select k from Customer where k = ?",
			[]*vom.RawBytes{vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Select based on parameterized type like expression
			"select k from Customer where Type(v) like ?",
			[]*vom.RawBytes{vom.RawBytesOf("%Customer")},
			[]*vom.RawBytes{vom.RawBytesOf("%Invoice")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("003")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
			},
		},
		{
			// Select based on parameterized type and key like expressions
			"select k from Customer where Type(v) like ? and k like ?",
			[]*vom.RawBytes{vom.RawBytesOf("%Customer"), vom.RawBytesOf("%2")},
			[]*vom.RawBytes{vom.RawBytesOf("%Invoice"), vom.RawBytesOf("00_002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("002002")},
			},
		},
		{
			// Select with parameters embedded in function args
			"select k from Customer where Type(v) like \"%.Invoice\" and Year(v.InvoiceDate, ?) = 2015 and Month(v.InvoiceDate, ?) = 3 and Day(v.InvoiceDate, ?) = 17",
			[]*vom.RawBytes{vom.RawBytesOf("America/Los_Angeles"), vom.RawBytesOf("America/Los_Angeles"), vom.RawBytesOf("America/Los_Angeles")},
			[]*vom.RawBytes{vom.RawBytesOf("Australia/Sydney"), vom.RawBytesOf("Australia/Sydney"), vom.RawBytesOf("Australia/Sydney")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
			},
			[][]*vom.RawBytes{},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where 1 = ? and (2 = ? and (3 = ? and 4 = ?) and 5 = ?) and k like ?",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where (1 = ? and 2 = ?) and (3 = ? and 4 = ? and 5 = ?) and k like ?",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where ((1 = ? and 2 = ? and 3 = ? and 4 = ?) and 5 = ?) and k like ?",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where ? = 1 and (2 = ? and (? = 3 and 4 = ?) and ? = 5) and k like ?",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where (1 = ? and ? = 2) and (3 = ? and ? = 4 and 5 = ?) and ? = k",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where ((? = 1 and ? = 2 and ? = 3 and ? = 4) and ? = 5) and ? = k",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where (? = 1 and ? = 2) and ? = 3 and (? = 4 and ? = 5) and ? = k",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where (1 = ? and 2 = ?) and 3 = ? and (4 = ? and 5 = ?) and k = ?",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where (? = 1 and 2 = ?) and ? = 3 and (4 = ? and ? = 5) and k = ?",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
		{
			// Nested expresions, do we get the values in the right order?
			"select k from Customer where (1 = ? and ? = 2) and 3 = ? and (? = 4 and 5 = ?) and k = ?",
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("001")},
			[]*vom.RawBytes{vom.RawBytesOf(1), vom.RawBytesOf(2), vom.RawBytesOf(3), vom.RawBytesOf(4), vom.RawBytesOf(5), vom.RawBytesOf("002")},
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
			},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("002")},
			},
		},
	}

	qe := engine.Create(db)

	// Prepare all of the statements ahead of time.
	preparedStatements := []public.PreparedStatement{}
	for _, test := range basic {
		p, err := qe.PrepareStatement(test.query)
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, err)
		}
		preparedStatements = append(preparedStatements, p)
	}

	// First execution
	for i, p := range preparedStatements {
		headers, rs, err := p.Exec(basic[i].paramValues1...)
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", basic[i].query, err)
		} else {
			// Collect results.
			r := [][]*vom.RawBytes{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if !reflect.DeepEqual(basic[i].r1, r) {
				t.Errorf("query: %s; got %v, want %v", basic[i].query, r, basic[i].r1)
			}
			if !reflect.DeepEqual(basic[i].headers, headers) {
				t.Errorf("query: %s; got %v, want %v", basic[i].query, headers, basic[i].headers)
			}
		}
	}

	// Second execution
	for i, p := range preparedStatements {
		headers, rs, err := p.Exec(basic[i].paramValues2...)
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", basic[i].query, err)
		} else {
			// Collect results.
			r := [][]*vom.RawBytes{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if !reflect.DeepEqual(basic[i].r2, r) {
				t.Errorf("query: %s; got %v, want %v", basic[i].query, r, basic[i].r2)
			}
			if !reflect.DeepEqual(basic[i].headers, headers) {
				t.Errorf("query: %s; got %v, want %v", basic[i].query, headers, basic[i].headers)
			}
		}
	}
}

func TestPrepareDelete(t *testing.T) {
	initTables()
	basic := []prepareDeleteTest{
		{
			// Delete based on parameterized key
			"delete from Customer where k = ?",
			[]*vom.RawBytes{vom.RawBytesOf("002")},
			[]string{"Count"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(1)},
			},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete based on parameterized type like expression
			"delete from Customer where Type(v) like ?",
			[]*vom.RawBytes{vom.RawBytesOf("%Invoice")},
			[]string{"Count"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(7)},
			},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete based on parameterized type and key like expressions
			"delete from Customer where Type(v) like ? and k like ?",
			[]*vom.RawBytes{vom.RawBytesOf("%Invoice"), vom.RawBytesOf("00_002")},
			[]string{"Count"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(2)},
			},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002001")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
		{
			// Delete with parameters embedded in function args
			"delete from Customer where Type(v) like \"%.Invoice\" and Year(v.InvoiceDate, ?) = 2015 and Month(v.InvoiceDate, ?) = 3 and Day(v.InvoiceDate, ?) = 17",
			[]*vom.RawBytes{vom.RawBytesOf("America/Los_Angeles"), vom.RawBytesOf("America/Los_Angeles"), vom.RawBytesOf("America/Los_Angeles")},
			[]string{"Count"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf(2)},
			},
			"select k from Customer",
			[]string{"k"},
			[][]*vom.RawBytes{
				{vom.RawBytesOf("001")},
				{vom.RawBytesOf("001001")},
				{vom.RawBytesOf("001002")},
				{vom.RawBytesOf("001003")},
				{vom.RawBytesOf("002")},
				{vom.RawBytesOf("002003")},
				{vom.RawBytesOf("002004")},
				{vom.RawBytesOf("003")},
			},
		},
	}
	qe := engine.Create(db)

	// Prepare all of the statements ahead of time.
	preparedStatements := []public.PreparedStatement{}
	for _, test := range basic {
		p, err := qe.PrepareStatement(test.delQuery)
		if err != nil {
			t.Errorf("delQuery: %s; got %v, want nil", test.delQuery, err)
		}
		preparedStatements = append(preparedStatements, p)
	}

	// Execute, then exec select to check resulting db.
	for i, p := range preparedStatements {
		initTables()
		headers, rs, err := p.Exec(basic[i].paramValues...)
		if err != nil {
			t.Errorf("delQuery: %s; got %v, want nil", basic[i].delQuery, err)
		} else {
			// Collect results.
			r := [][]*vom.RawBytes{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if !reflect.DeepEqual(basic[i].delResults, r) {
				t.Errorf("delQuery: %s; got %v, want %v", basic[i].delQuery, r, basic[i].delResults)
			}
			if !reflect.DeepEqual(basic[i].delHeaders, headers) {
				t.Errorf("delQuery: %s; got %v, want %v", basic[i].delQuery, headers, basic[i].delHeaders)
			}
		}
		headers, rs, err = qe.Exec(basic[i].selQuery)
		if err != nil {
			t.Errorf("delQuery: %s; got %v, want nil", basic[i].delQuery, err)
		} else {
			// Collect results.
			r := [][]*vom.RawBytes{}
			for rs.Advance() {
				r = append(r, rs.Result())
			}
			if !reflect.DeepEqual(basic[i].selResults, r) {
				t.Errorf("delQuery: %s; got %v, want %v", basic[i].delQuery, r, basic[i].selResults)
			}
			if !reflect.DeepEqual(basic[i].selHeaders, headers) {
				t.Errorf("delQuery: %s; got %v, want %v", basic[i].delQuery, headers, basic[i].selHeaders)
			}
		}
	}
}

func TestExecErrors(t *testing.T) {
	initTables()
	basic := []execSelectErrorTest{
		{
			"select a from Customer",
			syncql.NewErrInvalidSelectField(db.GetContext(), 7),
		},
		{
			"select v from Unknown",
			// The following error text is dependent on implementation of Database.
			syncql.NewErrTableCantAccess(db.GetContext(), 14, "Unknown", errors.New("No such table: Unknown.")),
		},
		{
			"select v from Customer offset -1",
			// The following error text is dependent on implementation of Database.
			syncql.NewErrExpected(db.GetContext(), 30, "positive integer literal"),
		},
		{
			"select k, v.Key from BigTable where 110 <= k and 205 >= k",
			syncql.NewErrKeyExpressionLiteral(db.GetContext(), 36),
		},
		{
			"select k, v.Key from BigTable where Type(k) = \"BigData\"",
			syncql.NewErrArgMustBeField(db.GetContext(), 41),
		},
		{
			"select v from Customer where Len(10) = 3",
			syncql.NewErrFunctionLenInvalidArg(db.GetContext(), 33),
		},
		{
			"select K from Customer where Type(v) = \"Invoice\"",
			syncql.NewErrDidYouMeanLowercaseK(db.GetContext(), 7),
		},
		{
			"select V from Customer where Type(v) = \"Invoice\"",
			syncql.NewErrDidYouMeanLowercaseV(db.GetContext(), 7),
		},
		{
			"select k from Customer where K = \"001\"",
			syncql.NewErrDidYouMeanLowercaseK(db.GetContext(), 29),
		},
		{
			"select v from Customer where Type(V) = \"Invoice\"",
			syncql.NewErrDidYouMeanLowercaseV(db.GetContext(), 34),
		},
		{
			"select K, V from Customer where Type(V) = \"Invoice\" and K = \"001\"",
			syncql.NewErrDidYouMeanLowercaseK(db.GetContext(), 7),
		},
		{
			"select type(v) from Customer where Type(v) = \"Invoice\"",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Type"),
		},
		{
			"select Type(v) from Customer where type(v) = \"Invoice\"",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 35, "Type"),
		},
		{
			"select type(v) from Customer where type(v) = \"Invoice\"",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Type"),
		},
		{
			"select time(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Time"),
		},
		{
			"select TimE(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Time"),
		},
		{
			"select year(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Year"),
		},
		{
			"select month(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Month"),
		},
		{
			"select day(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Day"),
		},
		{
			"select hour(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Hour"),
		},
		{
			"select minute(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Minute"),
		},
		{
			"select second(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Second"),
		},
		{
			"select now() from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Now"),
		},
		{
			"select LowerCase(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Lowercase"),
		},
		{
			"select UPPERCASE(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Uppercase"),
		},
		{
			"select spliT(\"foo:bar\", \":\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Split"),
		},
		{
			"select len(\"foo\") from Customer",
			syncql.NewErrDidYouMeanFunction(db.GetContext(), 7, "Len"),
		},
		{
			"select StrRepeat(\"foo\", \"x\") from Customer",
			syncql.NewErrIntConversionError(db.GetContext(), 24, errors.New("Cannot convert operand to int64.")),
		},
		{
			"select StrCat(v.Address.City, 42) from Customer",
			syncql.NewErrStringConversionError(db.GetContext(), 30, errors.New("Cannot convert operand to string.")),
		},
		{
			"select v from Customer where k like \"abc %\" escape ' '",
			syncql.NewErrInvalidEscapeChar(db.GetContext(), 51),
		},
		{
			"select v.Foo from",
			syncql.NewErrUnexpectedEndOfStatement(db.GetContext(), 17),
		},
	}

	qe := engine.Create(db)

	for _, test := range basic {
		_, _, err := qe.Exec(test.query)
		// Test both that the IDs compare and the text compares (since the offset needs to match).
		if verror.ErrorID(err) != verror.ErrorID(test.err) || err.Error() != test.err.Error() {
			t.Errorf("query: %s; got %v, want %v", test.query, err, test.err)
		}
	}
}

func TestErrorOnPrepare(t *testing.T) {
	initTables()
	basic := []prepareSelectErrorOnPrepareTest{
		{
			"select v.Foo from",
			syncql.NewErrUnexpectedEndOfStatement(db.GetContext(), 17),
		},
		{
			"delete v from Customer",
			syncql.NewErrExpectedFrom(db.GetContext(), 7, "v"),
		},
	}

	qe := engine.Create(db)

	for _, test := range basic {
		_, err := qe.PrepareStatement(test.query)
		// Test both that the IDs compare and the text compares (since the offset needs to match).
		if verror.ErrorID(err) != verror.ErrorID(test.err) || err.Error() != test.err.Error() {
			t.Errorf("query: %s; got %v, want %v", test.query, err, test.err)
		}
	}
}

func TestErrorOnPrepareExec(t *testing.T) {
	initTables()
	basic := []prepareSelectErrorOnExecTest{
		{
			"select k, v from Customer where v.Foo = ?",
			[]*vom.RawBytes{},
			syncql.NewErrNotEnoughParamValuesSpecified(db.GetContext(), 40),
		},
		{
			"select k, v from Customer where v.Foo = 42",
			[]*vom.RawBytes{vom.RawBytesOf(77)},
			syncql.NewErrTooManyParamValuesSpecified(db.GetContext(), 26),
		},
		{
			"select k, v from FooTable where v.Foo = ?",
			[]*vom.RawBytes{vom.RawBytesOf("abc")},
			syncql.NewErrTableCantAccess(db.GetContext(), 17, "FooTable", errors.New("No such table: FooTable.")),
		},
		{
			"select k, v from Customer where v.Foo like ?",
			[]*vom.RawBytes{vom.RawBytesOf(42)},
			syncql.NewErrLikeExpressionsRequireRhsString(db.GetContext(), 43),
		},
		{
			"select k, v from Customer where k = ?",
			[]*vom.RawBytes{vom.RawBytesOf(42)},
			syncql.NewErrKeyExpressionLiteral(db.GetContext(), 36),
		},
		{
			"delete from Customer where v.Foo = ?",
			[]*vom.RawBytes{},
			syncql.NewErrNotEnoughParamValuesSpecified(db.GetContext(), 35),
		},
		{
			"delete from Customer where v.Foo = 42",
			[]*vom.RawBytes{vom.RawBytesOf(77)},
			syncql.NewErrTooManyParamValuesSpecified(db.GetContext(), 21),
		},
		{
			"delete from FooTable where v.Foo = ?",
			[]*vom.RawBytes{vom.RawBytesOf("abc")},
			syncql.NewErrTableCantAccess(db.GetContext(), 12, "FooTable", errors.New("No such table: FooTable.")),
		},
		{
			"delete from Customer where v.Foo like ?",
			[]*vom.RawBytes{vom.RawBytesOf(42)},
			syncql.NewErrLikeExpressionsRequireRhsString(db.GetContext(), 38),
		},
		{
			"delete from Customer where k = ?",
			[]*vom.RawBytes{vom.RawBytesOf(42)},
			syncql.NewErrKeyExpressionLiteral(db.GetContext(), 31),
		},
	}

	qe := engine.Create(db)

	for _, test := range basic {
		ps, err := qe.PrepareStatement(test.query)
		if err != nil {
			t.Errorf("query: %s; got %v, want nil", test.query, err)
		}
		_, _, err = ps.Exec(test.paramValues...)
		// Test both that the IDs compare and the text compares (since the offset needs to match).
		if verror.ErrorID(err) != verror.ErrorID(test.err) || err.Error() != test.err.Error() {
			t.Errorf("query: %s; got %v, want %v", test.query, err, test.err)
		}
	}
}
