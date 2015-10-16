// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vom

import (
	"fmt"
	"io"
)

const minBufFree = 1024 // buffers always have at least 1K free after growth

// encbuf manages the write buffer for encoders.  The approach is similar to
// bytes.Buffer, but the implementation is simplified to only deal with many
// writes followed by a read of the whole buffer.
type encbuf struct {
	buf []byte
	end int // [0, end) is data that's already written

	// It's faster to hold end than to use both the len and cap properties of buf,
	// since end is cheaper to update than buf.
	maxSize int
}

func newEncbuf(maxSize int) *encbuf {
	return &encbuf{
		buf:     make([]byte, minBufFree),
		maxSize: maxSize,
	}
}

// Bytes returns a slice of the bytes written so far.
func (b *encbuf) Bytes() []byte { return b.buf[:b.end] }

// Len returns the number of bytes written so far.
func (b *encbuf) Len() int { return b.end }

// Reset the length to 0 to start a new round of writes.
func (b *encbuf) Reset() { b.end = 0 }

// SpaceAvailable returns the amount of space available on the buffer for writing.
func (b *encbuf) SpaceAvailable() int {
	return b.maxSize - b.end
}

// assertEnoughSpace if there is not enough space to write the value.
// TODO(bprosnitz) Consider removing this after things stabilize
func (b *encbuf) assertEnoughSpace(size int) {
	if size > b.SpaceAvailable() {
		panic(fmt.Sprintf("insufficient space to write %d bytes", size))
	}
}

// reserve at least min free bytes in the buffer.
func (b *encbuf) tryReserve(min int) {
	if len(b.buf) == b.maxSize {
		return // can't expand further
	}
	if len(b.buf)-b.end < min {
		newlen := len(b.buf) * 2
		if newlen-b.end < min {
			newlen = b.end + min + minBufFree
		}
		if newlen > b.maxSize {
			newlen = b.maxSize
		}
		newbuf := make([]byte, newlen)
		copy(newbuf, b.buf[:b.end])
		b.buf = newbuf
	}
}

// Grow the buffer by n bytes, and returns those bytes.
//
// Different from bytes.Buffer.Grow, which doesn't return the bytes.  Although
// this makes expandingEncbuf slightly easier to misuse, it helps to improve performance
// by avoiding unnecessary copying.
func (b *encbuf) Grow(n int) []byte {
	b.assertEnoughSpace(n)
	b.tryReserve(n)
	oldend := b.end
	b.end += n
	return b.buf[oldend:b.end]
}

// WriteOneByte writes byte c into the buffer.
func (b *encbuf) WriteOneByte(c byte) {
	b.assertEnoughSpace(1)
	b.tryReserve(1)
	b.buf[b.end] = c
	b.end++
}

// WriteMaximumPossible writes the maximum possible length of slice p, given
// limit constraints, into the buffer.
func (b *encbuf) WriteMaximumPossible(p []byte) (amtWritten int) {
	b.tryReserve(len(p))
	amtWritten = copy(b.buf[b.end:], p)
	b.end += amtWritten
	return
}

// decbuf manages the read buffer for decoders.  The approach is similar to
// bufio.Reader, but the API is better suited for fast decoding.
type decbuf struct {
	buf      []byte
	beg, end int // [beg, end) is data read from reader but unread by the user
	lim      int // number of bytes left in limit, or -1 for no limit
	reader   io.Reader

	// It's faster to hold end than to use the len and cap properties of buf,
	// since end is cheaper to update than buf.
}

// newDecbuf returns a new decbuf that fills its internal buffer by reading r.
func newDecbuf(r io.Reader) *decbuf {
	return &decbuf{
		buf:    make([]byte, minBufFree),
		lim:    -1,
		reader: r,
	}
}

// newDecbufFromBytes returns a new decbuf that reads directly from b.
func newDecbufFromBytes(b []byte) *decbuf {
	return &decbuf{
		buf:    b,
		end:    len(b),
		lim:    -1,
		reader: alwaysEOFReader{},
	}
}

type alwaysEOFReader struct{}

func (alwaysEOFReader) Read([]byte) (int, error) { return 0, io.EOF }

// Reset resets the buffer so it has no data.
func (b *decbuf) Reset() {
	b.beg = 0
	b.end = 0
	b.lim = -1
}

// SetLimit sets a limit to the bytes that are returned by decbuf; after a limit
// is set, subsequent reads cannot read past the limit, even if more bytes are
// available.  Attempts to read past the limit return io.EOF.  Call RemoveLimit
// to remove the limit.
//
// REQUIRES: limit >=0,
func (b *decbuf) SetLimit(limit int) {
	b.lim = limit
}

func (b *decbuf) HasDataAvailable() bool {
	return b.lim != 0 // -1 or positive
}

// RemoveLimit removes the limit, and returns the number of leftover bytes.
// Returns -1 if no limit was set.
func (b *decbuf) RemoveLimit() int {
	leftover := b.lim
	b.lim = -1
	return leftover
}

// fill the buffer with at least min bytes of data.  Returns an error if fewer
// than min bytes could be filled.  Doesn't advance the read position.
func (b *decbuf) fill(min int) error {
	switch avail := b.end - b.beg; {
	case avail >= min:
		// Fastpath - enough bytes are available.
		return nil
	case len(b.buf) < min:
		// The buffer isn't big enough.  Make a new buffer that's big enough and
		// copy existing data to the front.
		newlen := len(b.buf) * 2
		if newlen < min+minBufFree {
			newlen = min + minBufFree
		}
		newbuf := make([]byte, newlen)
		b.end = copy(newbuf, b.buf[b.beg:b.end])
		b.beg = 0
		b.buf = newbuf
	default:
		// The buffer is big enough.  Move existing data to the front.
		b.moveDataToFront()
	}
	// INVARIANT: len(b.buf)-b.beg >= min
	//
	// Fill [b.end:] until min bytes are available.  We must loop since Read may
	// return success with fewer bytes than requested.
	for b.end-b.beg < min {
		switch nread, err := b.reader.Read(b.buf[b.end:]); {
		case nread > 0:
			b.end += nread
		case err != nil:
			return err
		}
	}
	return nil
}

// moveDataToFront moves existing data in buf to the front, so that b.beg is 0.
func (b *decbuf) moveDataToFront() {
	b.end = copy(b.buf, b.buf[b.beg:b.end])
	b.beg = 0
}

// ReadSmall returns a buffer with the next n bytes, and increments the read
// position past those bytes.  Returns an error if fewer than n bytes are
// available.
//
// The returned slice points directly at our internal buffer, and is only valid
// until the next decbuf call.
//
// REQUIRES: n >= 0
func (b *decbuf) ReadSmall(n int) ([]byte, error) {
	if b.lim > -1 {
		if b.lim < n {
			b.lim = 0
			return nil, io.EOF
		}
		b.lim -= n
	}
	if err := b.fill(n); err != nil {
		return nil, err
	}
	buf := b.buf[b.beg : b.beg+n]
	b.beg += n
	return buf, nil
}

// PeekSmall returns a buffer with at least the next n bytes, but possibly
// more.  The read position isn't incremented.  Returns an error if fewer than
// min bytes are available.
//
// The returned slice points directly at our internal buffer, and is only valid
// until the next decbuf call.
//
// REQUIRES: min >= 0
func (b *decbuf) PeekSmall(min int) ([]byte, error) {
	if b.lim > -1 && b.lim < min {
		return nil, io.EOF
	}
	if err := b.fill(min); err != nil {
		return nil, err
	}
	return b.buf[b.beg:b.end], nil
}

// Skip increments the read position past the next n bytes.  Returns an error if
// fewer than n bytes are available.
//
// REQUIRES: n >= 0
func (b *decbuf) Skip(n int) (int, error) {
	readAmt := n
	if b.lim > -1 {
		if b.lim < n {
			n = b.lim
			readAmt = n
		}
		b.lim -= n
	}
	// If enough bytes are available, just update indices.
	avail := b.end - b.beg
	if avail >= n {
		b.beg += n
		return readAmt, nil
	}
	n -= avail
	// Keep reading into buf until we've read enough bytes.
	for {
		switch nread, err := b.reader.Read(b.buf); {
		case nread > 0:
			if nread >= n {
				b.beg = n
				b.end = nread
				return readAmt, nil
			}
			n -= nread
		case err != nil:
			return 0, err
		}
	}
}

// ReadByte returns the next byte, and increments the read position.
func (b *decbuf) ReadByte() (byte, error) {
	if b.lim > -1 {
		if b.lim == 0 {
			return 0, io.EOF
		}
		b.lim--
	}
	if err := b.fill(1); err != nil {
		return 0, err
	}
	ret := b.buf[b.beg]
	b.beg++
	return ret, nil
}

// PeekByte returns the next byte, without changing the read position.
func (b *decbuf) PeekByte() (byte, error) {
	if b.lim == 0 {
		return 0, io.EOF
	}
	if err := b.fill(1); err != nil {
		return 0, err
	}
	return b.buf[b.beg], nil
}

// ReadIntoBuf reads the next len(p) bytes into p, and increments the read position
// past those bytes.  Returns an error if fewer than len(p) bytes are available.
func (b *decbuf) ReadIntoBuf(p []byte) (int, error) {
	if b.lim > -1 {
		if b.lim < len(p) {
			p = p[:b.lim]
		}
		b.lim -= len(p)
	}
	amtRead := len(p)
	// Copy bytes from the buffer.
	ncopy := copy(p, b.buf[b.beg:b.end])
	b.beg += ncopy
	p = p[ncopy:]
	// Keep reading into p until we've read enough bytes.
	for len(p) > 0 {
		switch nread, err := b.reader.Read(p); {
		case nread > 0:
			p = p[nread:]
		case err != nil:
			return 0, err
		}
	}
	return amtRead, nil
}
