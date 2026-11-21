// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#include "snnote/snnote.h"

// Construct the pointer-bearing view on the C stack; Go passes only byte storage.
static sn_status sngo_codes(sn_encoding encoding, const uint8_t *data, size_t length,
 size_t width, size_t height, uint8_t **out, size_t *size, sn_error *err) {
 sn_bitmap bitmap = {.encoding = encoding, .bytes = {data, length}};
 return sn_bitmap_codes(&bitmap, width, height, out, size, err);
}
static sn_status sngo_raster(sn_encoding encoding, const uint8_t *data, size_t length,
 size_t width, size_t height, const sn_palette *palette, sn_raster **out, sn_error *err) {
 sn_bitmap bitmap = {.encoding = encoding, .bytes = {data, length}};
 return sn_bitmap_raster(&bitmap, width, height, palette, out, err);
}
*/
import "C"

import "math"

// Raster contains independent row-major intensity and coverage planes, top row first.
type Raster struct {
	Width, Height int
	Gray, Alpha   []byte
}

// Codes decodes RATTA_RLE into one code per pixel. PNG and JPEG return ErrImageEncoding.
func (b Bitmap) Codes(width, height int) ([]byte, error) {
	if width <= 0 || height <= 0 {
		return nil, &Error{Status: StatusInvalidArgument, Message: "bitmap dimensions must be positive"}
	}
	if b.Encoding < math.MinInt32 || b.Encoding > math.MaxInt32 {
		return nil, &Error{Status: StatusInvalidArgument, Message: "encoding does not fit C int"}
	}
	var out *C.uint8_t
	var size C.size_t
	var cerr C.sn_error
	status := C.sngo_codes(C.sn_encoding(b.Encoding), bytesPointer(b.Bytes), C.size_t(len(b.Bytes)), C.size_t(width), C.size_t(height), &out, &size, &cerr)
	if err := errorOf(status, &cerr); err != nil {
		return nil, err
	}
	defer C.sn_codes_free(out)
	// Bitmap decoding retains its existing pixel ceiling, independently of Parse.
	c := converter{budget: MaxRasterPixels}
	data := c.bytes(C.sn_bytes{data: out, size: size})
	return data, c.err
}

// Raster decodes RATTA_RLE through the supplied palette, copying both pixel planes.
func (b Bitmap) Raster(width, height int, palette Palette) (*Raster, error) {
	if width <= 0 || height <= 0 {
		return nil, &Error{Status: StatusInvalidArgument, Message: "bitmap dimensions must be positive"}
	}
	if b.Encoding < math.MinInt32 || b.Encoding > math.MaxInt32 {
		return nil, &Error{Status: StatusInvalidArgument, Message: "encoding does not fit C int"}
	}
	var cp C.sn_palette
	for i, ink := range palette.Inks {
		cp.inks[i] = C.sn_ink{gray: C.uint8_t(ink.Gray), alpha: C.uint8_t(ink.Alpha)}
	}
	var out *C.sn_raster
	var cerr C.sn_error
	status := C.sngo_raster(C.sn_encoding(b.Encoding), bytesPointer(b.Bytes), C.size_t(len(b.Bytes)), C.size_t(width), C.size_t(height), &cp, &out, &cerr)
	if err := errorOf(status, &cerr); err != nil {
		return nil, err
	}
	defer C.sn_raster_free(out)
	// Both decoded planes may reach the existing pixel ceiling.
	c := converter{budget: 2 * MaxRasterPixels}
	w, h := c.length(out.width, "width"), c.length(out.height, "height")
	if c.err != nil {
		return nil, c.err
	}
	if h != 0 && w > math.MaxInt/h {
		return nil, &Error{Status: StatusOutOfMemory, Message: "raster pixel count exceeds Go int"}
	}
	pixels := C.size_t(w * h)
	raster := &Raster{Width: w, Height: h, Gray: c.bytes(C.sn_bytes{data: out.gray, size: pixels}), Alpha: c.bytes(C.sn_bytes{data: out.alpha, size: pixels})}
	if c.err != nil {
		return nil, c.err
	}
	return raster, nil
}
