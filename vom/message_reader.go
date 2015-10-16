// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vom

import (
	"v.io/v23/vdl"
	"v.io/v23/verror"
)

var (
	errGotValueWantType   = verror.Register(pkgPath+".errGotValueWantType", verror.NoRetry, "{1:}{2:} vom: got value chunk, expected type chunk {:_}")
	errGotTypeWantValue   = verror.Register(pkgPath+".errGotTypeWantValue", verror.NoRetry, "{1:}{2:} vom: got type chunk, expected value chunk {:_}")
	errContChunkBeforeNew = verror.Register(pkgPath+".errContChunkBeforeNew", verror.NoRetry, "{1:}{2:} vom: received continuation chunk before new message {:_}")
	errInvalidTypeIdIndex = verror.Register(pkgPath+".errInvalidTypeIdIndex", verror.NoRetry, "{1:}{2:} vom: value referenced invalid index into type id table {:_}")
)

// Create a messageReader with the specified decbuf.
func newMessageReader(buf *decbuf) *messageReader {
	return &messageReader{
		buf: buf,
	}
}

type messageKind int

const (
	typeMessage messageKind = iota
	valueMessage
)

func (mk messageKind) String() string {
	switch mk {
	case typeMessage:
		return "TypeMessage"
	case valueMessage:
		return "ValueMessage"
	default:
		panic("unknown message kind")
	}
}

type messageData struct {
	tid      typeId
	refTypes referencedTypes
}

func (m *messageData) Reset() error {
	m.tid = 0
	return m.refTypes.Reset()
}

func (m *messageData) HasMessage() bool {
	return m.tid != 0
}

// messageReader assists in reading messages across multiple chunks.
type messageReader struct {
	// Stream data:
	buf     *decbuf
	version Version

	// Message data:
	curMsgKind messageKind
	typeMsg    messageData
	valueMsg   messageData

	lookupType     func(tid typeId) (*vdl.Type, error)
	readSingleType func() error
}

func (mr *messageReader) SetCallbacks(lookupType func(tid typeId) (*vdl.Type, error), readSingleType func() error) {
	mr.lookupType = lookupType
	mr.readSingleType = readSingleType
}

func (mr *messageReader) hasInterleavedTypesAndValues() bool {
	return mr.readSingleType != nil
}

func (mr *messageReader) startChunk() (msgConsumed bool, err error) {
	if leftover := mr.buf.RemoveLimit(); leftover > 0 {
		return false, verror.New(errLeftOverBytes, nil)
	}

	if mr.version == 0 {
		version, err := mr.buf.ReadByte()
		if err != nil {
			return false, verror.New(errEndedBeforeVersionByte, nil, err)
		}
		mr.version = Version(version)
		if !isAllowedVersion(mr.version) {
			return false, verror.New(errBadVersionByte, nil)
		}
	}

	mid, cr, err := decbufBinaryDecodeIntWithControl(mr.buf)
	if err != nil {
		return false, err
	}

	var activeMsg *messageData
	switch cr {
	case 0:
		var tid typeId
		if mid < 0 {
			mr.curMsgKind = typeMessage
			activeMsg = &mr.typeMsg
			tid = typeId(-mid)
		} else if mid > 0 {
			mr.curMsgKind = valueMessage
			activeMsg = &mr.valueMsg
			tid = typeId(mid)
		} else {
			return false, verror.New(errDecodeZeroTypeID, nil)
		}
		if activeMsg.HasMessage() {
			// Message continuation before previous ended.
			return false, verror.New(errInvalid, nil)
		}
		activeMsg.tid = tid
	case WireCtrlTypeCont:
		mr.curMsgKind = typeMessage
		activeMsg = &mr.typeMsg
		if !activeMsg.HasMessage() {
			return false, verror.New(errContChunkBeforeNew, nil)
		}
	case WireCtrlValueCont:
		mr.curMsgKind = valueMessage
		activeMsg = &mr.valueMsg
		if !activeMsg.HasMessage() {
			return false, verror.New(errContChunkBeforeNew, nil)
		}
	default:
		return false, verror.New(errBadControlCode, nil)
	}

	var hasAnyOrTypeObject bool
	var hasLength bool
	switch mr.curMsgKind {
	case typeMessage:
		hasLength = true
		hasAnyOrTypeObject = false
	case valueMessage:
		t, err := mr.lookupType(activeMsg.tid)
		if err != nil {
			return false, err
		}
		hasLength = hasChunkLen(t)
		hasAnyOrTypeObject = t.ContainsAnyOrTypeObject()
	}

	if hasAnyOrTypeObject && mr.version != Version80 {
		l, err := decbufBinaryDecodeUint(mr.buf)
		if err != nil {
			return false, err
		}
		for i := 0; i < int(l); i++ {
			refId, err := decbufBinaryDecodeUint(mr.buf)
			if err != nil {
				return false, err
			}
			activeMsg.refTypes.AddTypeID(typeId(refId))
		}
	}

	if hasLength {
		chunkLen, err := decbufBinaryDecodeUint(mr.buf)
		if err != nil {
			return true, err
		}
		mr.buf.SetLimit(int(chunkLen))
	}

	if mr.curMsgKind == typeMessage && cr == 0 {
		if mr.hasInterleavedTypesAndValues() {
			err := mr.readSingleType()
			return true, err
		}
	}

	return false, nil
}

func (mr *messageReader) nextReadableChunk() (err error) {
	for !mr.buf.HasDataAvailable() {
		if err := mr.nextChunk(); err != nil {
			return err
		}
	}
	return nil
}

func (mr *messageReader) nextChunk() (err error) {
	consumed := true
	for consumed {
		// Types will be read within startChunk if this is a mixed stream.
		consumed, err = mr.startChunk()
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr *messageReader) StartTypeMessage() (typeId, error) {
	if !mr.hasInterleavedTypesAndValues() {
		// Two stream case.
		consumed, err := mr.startChunk()
		if err != nil {
			return 0, err
		}
		if consumed {
			return 0, verror.New(verror.ErrInternal, nil, "chunk unexpectedly consumed while reading single type message")
		}
	}

	if mr.curMsgKind == valueMessage {
		return 0, verror.New(errGotValueWantType, nil)
	}
	return mr.typeMsg.tid, nil
}

func (mr *messageReader) StartValueMessage() (typeId, error) {
	if err := mr.nextChunk(); err != nil {
		return 0, err
	}

	if mr.curMsgKind == typeMessage {
		return 0, verror.New(errGotTypeWantValue, nil)
	}
	return mr.valueMsg.tid, nil
}

func (mr *messageReader) EndMessage() error {
	if leftover := mr.buf.RemoveLimit(); leftover > 0 {
		return verror.New(errLeftOverBytes, nil)
	}
	switch mr.curMsgKind {
	case typeMessage:
		mr.curMsgKind = valueMessage // Set kind to value message to return to reading value if we temporarily left to decode a type.
		return mr.typeMsg.Reset()
	case valueMessage:
		return mr.valueMsg.Reset()
	default:
		panic("unknown msg kind")
	}
}

// ReadSmall returns a buffer with the next n bytes, and increments the read
// position past those bytes.  Returns an error if fewer than n bytes are
// available.
//
// The returned slice points directly at our internal buffer, and is only valid
// until the next messageReader call.
//
// REQUIRES: n >= 0 && n <= 9
func (mr *messageReader) ReadSmall(n int) ([]byte, error) {
	if err := mr.nextReadableChunk(); err != nil {
		return nil, err
	}
	return mr.buf.ReadSmall(n)
}

// PeekSmall returns a buffer with at least the next n bytes, but possibly
// more.  The read position isn't incremented.  Returns an error if fewer than
// min bytes are available.
//
// The returned slice points directly at our internal buffer, and is only valid
// until the next messageReader call.
//
// REQUIRES: min >= 0 && n <= 9
func (mr *messageReader) PeekSmall(n int) ([]byte, error) {
	if err := mr.nextReadableChunk(); err != nil {
		return nil, err
	}
	return mr.buf.PeekSmall(n)
}

// Skip increments the read position past the next n bytes in the message,
// which may be across multiple chunks.
//
// REQUIRES: n >= 0
func (mr *messageReader) Skip(n int) error {
	for n > 0 {
		if err := mr.nextReadableChunk(); err != nil {
			return err
		}

		nSkipped, err := mr.buf.Skip(n)
		if err != nil {
			return err
		}

		n -= nSkipped
	}
	return nil
}

// ReadByte returns the next byte, and increments the read position.
func (mr *messageReader) ReadByte() (byte, error) {
	if err := mr.nextReadableChunk(); err != nil {
		return 0, err
	}
	return mr.buf.ReadByte()
}

// PeekByte returns the next byte, without changing the read position.
func (mr *messageReader) PeekByte() (byte, error) {
	if err := mr.nextReadableChunk(); err != nil {
		return 0, err
	}
	return mr.buf.PeekByte()
}

// ReadIntoBuf reads bytes from multiple chunks into the specified buffer.
// The entire buffer should have been written at the end of a successful call.
func (mr *messageReader) ReadIntoBuf(p []byte) error {
	for len(p) > 0 {
		if err := mr.nextReadableChunk(); err != nil {
			return err
		}

		n, err := mr.buf.ReadIntoBuf(p)
		if err != nil {
			return err
		}

		p = p[n:]
	}
	return nil
}

func (mr *messageReader) ReferencedType(index uint64) (t *vdl.Type, err error) {
	if mr.version == Version80 {
		panic("type references shouldn't be used for version 0x80")
	}

	var tid typeId
	switch mr.curMsgKind {
	case typeMessage:
		tid, err = mr.typeMsg.refTypes.ReferencedTypeId(index)
	case valueMessage:
		tid, err = mr.valueMsg.refTypes.ReferencedTypeId(index)
	default:
		panic("unknown message kind")
	}
	if err != nil {
		return
	}
	t, err = mr.lookupType(tid)
	return
}

type referencedTypes struct {
	tids []typeId
}

func (refTypes *referencedTypes) Reset() (err error) {
	refTypes.tids = refTypes.tids[:0]
	return
}

func (refTypes *referencedTypes) AddTypeID(tid typeId) {
	refTypes.tids = append(refTypes.tids, tid)
}

func (refTypes *referencedTypes) ReferencedTypeId(index uint64) (typeId, error) {
	if index >= uint64(len(refTypes.tids)) {
		return 0, verror.New(errInvalidTypeIdIndex, nil)
	}
	return refTypes.tids[index], nil
}
