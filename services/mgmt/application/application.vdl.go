// This file was auto-generated by the veyron vdl tool.
// Source: application.vdl

// Package application defines the type for describing an application.
package application

import (
	// VDL system imports
	"v.io/core/veyron2/vdl"

	// VDL user imports
	"v.io/core/veyron2/security"
)

// Envelope is a collection of metadata that describes an application.
type Envelope struct {
	// Title is the publisher-assigned application title.  Application
	// installations with the same title are considered as belonging to the
	// same application by the application management system.
	//
	// A change in the title signals a new application.
	Title string
	// Args is an array of command-line arguments to be used when executing
	// the binary.
	Args []string
	// Binary is an object name that identifies the application binary.
	Binary string
	// Signature represents a signature on the sha256 hash of the application
	// binary by the publisher principal.
	Signature security.Signature
	// Publisher represents the set of blessings that have been bound to
	// the principal who published this binary.
	Publisher security.WireBlessings
	// Env is an array that stores the environment variable values to be
	// used when executing the binary.
	Env []string
	// Packages is a map of packages to install on the local filesystem
	// before executing the binary. The map key is the local file/directory
	// name, relative to the instance's packages directory, where the
	// package should be installed. For archives, this name represents a
	// directory into which the archive is to be extracted, and for regular
	// files it represents the name for the file.
	// The map value is the object name of the package.
	// Each object's media type determines how to install it.
	//
	// For example, with key=pkg1,value=binaryrepo/configfiles (an archive),
	// the "configfiles" package will be installed under the "pkg1"
	// directory. With key=pkg2,value=binaryrepo/binfile (a binary), the
	// "binfile" file will be installed as the "pkg2" file.
	//
	// The keys must be valid file/directory names, without path separators.
	//
	// Any number of packages may be specified.
	Packages map[string]string
}

func (Envelope) __VDLReflect(struct {
	Name string "v.io/core/veyron2/services/mgmt/application.Envelope"
}) {
}

func init() {
	vdl.Register((*Envelope)(nil))
}

// Device manager application envelopes must present this title.
const DeviceManagerTitle = "device manager"
