// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"v.io/syncbase/v23/syncbase/nosql/internal/query"
	"v.io/syncbase/v23/syncbase/nosql/internal/query/query_db"
)

type demoDB struct {
	tables []table
}

type table struct {
	name string
	rows []kv
}

type keyValueStreamImpl struct {
	table        table
	cursor       int
	prefixes     []string
	prefixCursor int
}

func (kvs *keyValueStreamImpl) Advance() bool {
	for true {
		kvs.cursor++ // initialized to -1
		if kvs.cursor >= len(kvs.table.rows) {
			return false
		}
		for kvs.prefixCursor < len(kvs.prefixes) {
			// does it match any prefix
			if kvs.prefixes[kvs.prefixCursor] == "" || strings.HasPrefix(kvs.table.rows[kvs.cursor].key, kvs.prefixes[kvs.prefixCursor]) {
				return true
			}
			// Keys and prefixes are both sorted low to high, so we can increment
			// prefixCursor if the prefix is < the key.
			if kvs.prefixes[kvs.prefixCursor] < kvs.table.rows[kvs.cursor].key {
				kvs.prefixCursor++
				if kvs.prefixCursor >= len(kvs.prefixes) {
					return false
				}
			} else {
				break
			}
		}
	}
	return false
}

func (kvs *keyValueStreamImpl) KeyValue() (string, interface{}) {
	return kvs.table.rows[kvs.cursor].key, kvs.table.rows[kvs.cursor].value
}

func (kvs *keyValueStreamImpl) Err() error {
	return nil
}

func (kvs *keyValueStreamImpl) Cancel() {
}

func (t table) Scan(prefixes []string) (query_db.KeyValueStream, error) {
	var keyValueStreamImpl keyValueStreamImpl
	keyValueStreamImpl.table = t
	keyValueStreamImpl.cursor = -1
	keyValueStreamImpl.prefixes = prefixes
	return &keyValueStreamImpl, nil
}

func (db demoDB) GetTable(table string) (query_db.Table, error) {
	for _, t := range db.tables {
		if t.name == table {
			return t, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No such table: %s.", table))

}

type AddressInfo struct {
	Street string
	City   string
	State  string
	Zip    string
}

type Customer struct {
	Name    string
	ID      int64
	Active  bool
	Rating  rune
	Address AddressInfo
}

type Invoice struct {
	CustID     int64
	InvoiceNum int64
	Amount     int64
	ShipTo     AddressInfo
}

type Numbers struct {
	B    byte
	UI16 uint16
	UI32 uint32
	UI64 uint64
	I16  int16
	I32  int32
	I64  int64
	BI   big.Int
	F32  float32
	F64  float64
	BR   big.Rat
}

type kv struct {
	key   string
	value interface{}
}

func createDB() query_db.Database {
	var db demoDB
	var custTable table
	custTable.name = "Customer"
	custTable.rows = []kv{
		kv{
			"001",
			Customer{"John Smith", 1, true, 'A', AddressInfo{"1 Main St.", "Palo Alto", "CA", "94303"}},
		},
		kv{
			"001001",
			Invoice{1, 1000, 42, AddressInfo{"1 Main St.", "Palo Alto", "CA", "94303"}},
		},
		kv{
			"001002",
			Invoice{1, 1003, 7, AddressInfo{"2 Main St.", "Palo Alto", "CA", "94303"}},
		},
		kv{
			"001003",
			Invoice{1, 1005, 88, AddressInfo{"3 Main St.", "Palo Alto", "CA", "94303"}},
		},
		kv{
			"002",
			Customer{"Bat Masterson", 2, true, 'B', AddressInfo{"777 Any St.", "Collins", "IA", "50055"}},
		},
		kv{
			"002001",
			Invoice{2, 1001, 166, AddressInfo{"777 Any St.", "collins", "IA", "50055"}},
		},
		kv{
			"002002",
			Invoice{2, 1002, 243, AddressInfo{"888 Any St.", "collins", "IA", "50055"}},
		},
		kv{
			"002003",
			Invoice{2, 1004, 787, AddressInfo{"999 Any St.", "collins", "IA", "50055"}},
		},
		kv{
			"002004",
			Invoice{2, 1006, 88, AddressInfo{"101010 Any St.", "collins", "IA", "50055"}},
		},
	}
	db.tables = append(db.tables, custTable)

	var numTable table
	numTable.name = "Numbers"
	numTable.rows = []kv{
		kv{
			"001",
			Numbers{byte(12), uint16(1234), uint32(5678), uint64(999888777666), int16(9876), int32(876543), int64(128), *big.NewInt(1234567890), float32(3.14159), float64(2.71828182846), *big.NewRat(123, 1)},
		},
		kv{
			"002",
			Numbers{byte(9), uint16(99), uint32(999), uint64(9999999), int16(9), int32(99), int64(88), *big.NewInt(9999), float32(1.41421356237), float64(1.73205080757), *big.NewRat(999999, 1)},
		},
	}
	db.tables = append(db.tables, numTable)
	return db
}

func dumpDB(db query_db.Database) {
	switch db := db.(type) {
	case demoDB:
		for _, t := range db.tables {
			fmt.Printf("table: %s\n", t.name)
			for _, row := range t.rows {
				fmt.Printf("key: %s, value: ", row.key)
				switch v := row.value.(type) {
				case Customer:
					fmt.Printf("type:Customer Name:%v,ID:%v,Active:%v,Rating:%v,Address{%v}\n", v.Name, v.ID, v.Active, v.Rating, dumpAddress(&v.Address))
				case Invoice:
					fmt.Printf("type:Invoice CustID:%v,InvoiceNum:%v,Amount:%v,ShipTo:{%v}\n", v.CustID, v.InvoiceNum, v.Amount, dumpAddress(&v.ShipTo))
				case Numbers:
					fmt.Printf("type:Numbers B:%v,UI16:%v,UI32:%v,UI64:%v,I16:%v,I32:%v,I64:%v,BI:%v,F32:%v,F64:%v,BR:%v\n", v.B, v.UI16, v.UI32, v.UI64, v.I16, v.I32, v.I64, v.BI, v.F32, v.F64, v.BR)
				default:
					fmt.Printf("%v\n", row.value)
				}
			}
		}
	}
}

func dumpAddress(a *AddressInfo) string {
	return fmt.Sprintf("Street:%v,City:%v,State:%v,Zip:%v", a.Street, a.City, a.State, a.Zip)
}

func queryExec(db query_db.Database, q string) {
	if rs, err := query.Exec(db, q); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", q)
		fmt.Fprintf(os.Stderr, "%s^\n", strings.Repeat(" ", int(err.Off)))
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	} else {
		for rs.Advance() {
			for i, column := range rs.Result() {
				if i != 0 {
					fmt.Printf(",")
				}
				fmt.Printf("%v", column)
			}
			fmt.Printf("\n")
		}
		if err := rs.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
	}
}

func main() {
	db := createDB()
	reader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Printf("Enter query (or 'dump' or 'exit')? ")
		q, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Input error: %s\n", err.Error()))
		} else {
			// kill the newline
			q = q[0 : len(q)-1]
			if q == "" || strings.ToLower(q) == "exit" {
				os.Exit(0)
			} else if strings.ToLower(q) == "dump" {
				dumpDB(db)
			} else {
				queryExec(db, q)
			}
		}
	}
}
