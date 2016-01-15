// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	ds "v.io/v23/query/engine/datasource"
	"v.io/v23/query/engine/internal/query_parser"
	"v.io/v23/query/syncql"
	"v.io/v23/vdl"
)

// Select result stream
type selectResultStreamImpl struct {
	db              ds.Database
	selectStatement *query_parser.SelectStatement
	resultCount     int64 // results served so far (needed for limit clause)
	skippedCount    int64 // skipped so far (needed for offset clause)
	keyValueStream  ds.KeyValueStream
	k               string
	v               *vdl.Value
	err             error
}

func (rs *selectResultStreamImpl) Advance() bool {
	if rs.selectStatement.Limit != nil && rs.resultCount >= rs.selectStatement.Limit.Limit.Value {
		return false
	}
	for rs.keyValueStream.Advance() {
		k, v := rs.keyValueStream.KeyValue()
		// EvalWhereUsingOnlyKey
		// INCLUDE: the row should be included in the results
		// EXCLUDE: the row should NOT be included
		// FETCH_VALUE: the value and/or type of the value are required to make determination.
		rv := EvalWhereUsingOnlyKey(rs.db, rs.selectStatement.Where, k)
		var match bool
		switch rv {
		case INCLUDE:
			match = true
		case EXCLUDE:
			match = false
		case FETCH_VALUE:
			match = Eval(rs.db, k, v, rs.selectStatement.Where.Expr)
		}
		if match {
			if rs.selectStatement.ResultsOffset == nil || rs.selectStatement.ResultsOffset.ResultsOffset.Value <= rs.skippedCount {
				rs.k = k
				rs.v = v
				rs.resultCount++
				return true
			} else {
				rs.skippedCount++
			}
		}
	}
	if err := rs.keyValueStream.Err(); err != nil {
		rs.err = syncql.NewErrKeyValueStreamError(rs.db.GetContext(), rs.selectStatement.Off, err)
	}
	return false
}

func (rs *selectResultStreamImpl) Result() []*vdl.Value {
	return ComposeProjection(rs.db, rs.k, rs.v, rs.selectStatement.Select)
}

func (rs *selectResultStreamImpl) Err() error {
	return rs.err
}

func (rs *selectResultStreamImpl) Cancel() {
	rs.keyValueStream.Cancel()
}

// Delete result stream
type deleteResultStreamImpl struct {
	db              ds.Database
	deleteStatement *query_parser.DeleteStatement
	deleteCursor    int64 // zero or one
	deleteCount     int64
	err             error
}

func (rs *deleteResultStreamImpl) Advance() bool {
	if rs.deleteCursor == 0 {
		rs.deleteCursor++
		return true
	}
	return false
}

func (rs *deleteResultStreamImpl) Result() []*vdl.Value {
	return []*vdl.Value{vdl.ValueOf(rs.deleteCount)}
}

func (rs *deleteResultStreamImpl) Err() error {
	return rs.err
}

func (rs *deleteResultStreamImpl) Cancel() {
	rs.deleteCursor++
}