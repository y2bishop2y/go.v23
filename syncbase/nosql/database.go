// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nosql

import (
	wire "v.io/syncbase/v23/services/syncbase/nosql"
	"v.io/syncbase/v23/syncbase/util"
	"v.io/v23/context"
	"v.io/v23/naming"
	"v.io/v23/security/access"
)

func NewDatabase(parentFullName, relativeName string) *database {
	fullName := naming.Join(parentFullName, relativeName)
	return &database{
		c:              wire.DatabaseClient(fullName),
		parentFullName: parentFullName,
		fullName:       fullName,
		name:           relativeName,
	}
}

type database struct {
	c              wire.DatabaseClientMethods
	parentFullName string
	fullName       string
	name           string
}

var _ Database = (*database)(nil)

// TODO(sadovsky): Validate names before sending RPCs.

// Name implements Database.Name.
func (d *database) Name() string {
	return d.name
}

// FullName implements Database.FullName.
func (d *database) FullName() string {
	return d.fullName
}

// Exists implements Database.Exists.
func (d *database) Exists(ctx *context.T) (bool, error) {
	return d.c.Exists(ctx)
}

// Table implements Database.Table.
func (d *database) Table(relativeName string) Table {
	return newTable(d.fullName, relativeName)
}

// ListTables implements Database.ListTables.
func (d *database) ListTables(ctx *context.T) ([]string, error) {
	return util.List(ctx, d.fullName)
}

// Create implements Database.Create.
func (d *database) Create(ctx *context.T, perms access.Permissions) error {
	return d.c.Create(ctx, perms)
}

// Delete implements Database.Delete.
func (d *database) Delete(ctx *context.T) error {
	return d.c.Delete(ctx)
}

// CreateTable implements Database.CreateTable.
func (d *database) CreateTable(ctx *context.T, relativeName string, perms access.Permissions) error {
	return wire.TableClient(naming.Join(d.fullName, relativeName)).Create(ctx, perms)
}

// DeleteTable implements Database.DeleteTable.
func (d *database) DeleteTable(ctx *context.T, relativeName string) error {
	return wire.TableClient(naming.Join(d.fullName, relativeName)).Delete(ctx)
}

// Exec implements Database.Exec.
func (d *database) Exec(ctx *context.T, query string) ([]string, ResultStream, error) {
	ctx, cancel := context.WithCancel(ctx)
	call, err := d.c.Exec(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	resultStream := newResultStream(cancel, call)
	// The first row contains headers, pull them off the stream
	// and return them separately.
	var headers []string
	if !resultStream.Advance() {
		if err = resultStream.Err(); err != nil {
			// Since there was an error, can't get headers.
			// Just return the error.
			return nil, nil, err
		}
	}
	for _, header := range resultStream.Result() {
		headers = append(headers, header.RawString())
	}
	return headers, resultStream, nil
}

// BeginBatch implements Database.BeginBatch.
func (d *database) BeginBatch(ctx *context.T, opts wire.BatchOptions) (BatchDatabase, error) {
	relativeName, err := d.c.BeginBatch(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &batch{database: *NewDatabase(d.parentFullName, relativeName)}, nil
}

// SetPermissions implements Database.SetPermissions.
func (d *database) SetPermissions(ctx *context.T, perms access.Permissions, version string) error {
	return d.c.SetPermissions(ctx, perms, version)
}

// GetPermissions implements Database.GetPermissions.
func (d *database) GetPermissions(ctx *context.T) (perms access.Permissions, version string, err error) {
	return d.c.GetPermissions(ctx)
}

// SyncGroup implements Database.SyncGroup.
func (d *database) SyncGroup(sgName string) SyncGroup {
	return newSyncGroup(d.fullName, sgName)
}

// GetSyncGroupNames implements Database.GetSyncGroupNames.
func (d *database) GetSyncGroupNames(ctx *context.T) ([]string, error) {
	return d.c.GetSyncGroupNames(ctx)
}
