//go:build snnote_test

// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#cgo CFLAGS: -DSNNOTE_GO_TEST
#include <stdlib.h>
int sngo_log_installs(void);

#include "snnote/snnote.h"
static ptrdiff_t sngo_drawn_index(const sn_page *p, size_t i) {
 return p->drawn_strokes[i] - p->strokes;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"unsafe"
)

// cLookupMatches keeps the Go lookup implementations checked against the
// library without making C ownership part of the public model.
func cLookupMatches(data []byte, n *Notebook) bool {
	var raw *C.sn_notebook
	var cerr C.sn_error
	if C.sn_parse_notebook(bytesPointer(data), C.size_t(len(data)), &raw, &cerr) != C.SN_OK {
		return false
	}
	defer C.sn_notebook_free(raw)
	c := converter{budget: conversionBudget(len(data))}
	count := c.arrayLength(unsafe.Pointer(raw.pages), raw.page_count, "lookup pages")
	if c.err != nil || count != len(n.Pages) {
		return false
	}
	if Equipment(raw.equipment) != n.Equipment {
		return false
	}
	for i, ink := range raw.palette.inks {
		if n.Palette.Inks[i] != (Ink{Gray: uint8(ink.gray), Alpha: uint8(ink.alpha)}) {
			return false
		}
	}
	if int(raw.link_count) != len(n.Links) || int(raw.template_count) != len(n.Templates) {
		return false
	}
	for i, link := range unsafe.Slice(raw.links, int(raw.link_count)) {
		if !reflect.DeepEqual(n.Links[i].URL, c.string(link.url)) {
			return false
		}
	}
	for i, tmpl := range unsafe.Slice(raw.templates, int(raw.template_count)) {
		if n.Templates[i].Name != c.requiredString(tmpl.name, "template") || !reflect.DeepEqual(n.Templates[i].Bitmap, c.bitmap(tmpl.bitmap)) {
			return false
		}
	}
	for i, p := range unsafe.Slice(raw.pages, count) {
		page := &n.Pages[i]
		if int(p.stroke_count) != len(page.Strokes) || int(p.drawn_stroke_count) != len(page.DrawnStrokes) {
			return false
		}
		for j, stroke := range unsafe.Slice(p.strokes, int(p.stroke_count)) {
			if page.Strokes[j].Kind != StrokeKind(stroke.kind) || page.Strokes[j].Status != StrokeStatus(stroke.status) {
				return false
			}
		}
		for j, stroke := range page.DrawnStrokes {
			index := int(C.sngo_drawn_index(&p, C.size_t(j)))
			if index < 0 || index >= len(page.Strokes) || stroke != &page.Strokes[index] {
				return false
			}
		}
		for name := LayerName(-1); name <= LayerName(6); name++ {
			got := n.Pages[i].Layer(name)
			want := C.sn_page_layer(&p, C.sn_layer_name(name))
			if (got == nil) != (want == nil) {
				return false
			}
			if want != nil && !reflect.DeepEqual(*got, c.layer(*want)) {
				return false
			}
		}
		if !metadataMatches(p.metadata, n.Pages[i].Metadata, conversionBudget(len(data))) {
			return false
		}
	}
	return c.err == nil && metadataMatches(raw.header, n.Header, conversionBudget(len(data))) && metadataMatches(raw.footer, n.Footer, conversionBudget(len(data)))
}
func metadataMatches(raw C.sn_metadata, m Metadata, budget uint64) bool {
	keys := []string{"missing-key"}
	for _, e := range m.Entries {
		keys = append(keys, e.Key)
	}
	for _, key := range keys {
		// The C lookup takes a NUL-terminated key. Go deliberately permits the
		// complete UTF-8 key, so only its representable domain is compared here.
		if strings.ContainsRune(key, 0) {
			continue
		}
		ck := C.CString(key)
		want := C.sn_metadata_get(&raw, ck)
		got, ok := m.Get(key)
		c := converter{budget: budget}
		match := ok == (want != nil)
		if want != nil {
			value := c.string(*want)
			match = match && value != nil && got == *value
		}
		var index C.size_t
		var values []string
		for {
			v := C.sn_metadata_next(&raw, ck, &index)
			if v == nil {
				break
			}
			value := c.string(*v)
			if value == nil {
				match = false
				break
			}
			values = append(values, *value)
		}
		C.free(unsafe.Pointer(ck))
		if !match || !reflect.DeepEqual(m.Values(key), values) {
			return false
		}
	}
	return true
}

// requiredStringNulls audits the C values before required strings are normalized.
func requiredStringNulls(data []byte) ([]string, error) {
	var raw *C.sn_notebook
	var cerr C.sn_error
	if err := errorOf(C.sn_parse_notebook(bytesPointer(data), C.size_t(len(data)), &raw, &cerr), &cerr); err != nil {
		return nil, err
	}
	defer C.sn_notebook_free(raw)
	var missing []string
	c := converter{budget: conversionBudget(len(data)), requiredNulls: &missing}
	_ = c.notebook(*raw)
	return missing, c.err
}

// conversionBoundaryChecks exercises C views that the parser does not normally emit.
func conversionBoundaryChecks() error {
	var storage C.char
	for _, tc := range []struct {
		name    string
		pointer unsafe.Pointer
		count   C.size_t
		want    error
	}{
		{"empty null", nil, 0, nil},
		{"empty allocated", unsafe.Pointer(&storage), 0, nil},
		{"missing storage", nil, 1, ErrMalformed},
		{"oversized null", nil, C.size_t(math.MaxInt) + 1, ErrOutOfMemory},
		{"oversized allocated", unsafe.Pointer(&storage), C.size_t(math.MaxInt) + 1, ErrOutOfMemory},
	} {
		var c converter
		if n := c.arrayLength(tc.pointer, tc.count, tc.name); n != 0 || !errors.Is(c.err, tc.want) {
			return fmt.Errorf("%s: count %d, error %w", tc.name, n, c.err)
		}
	}
	var c converter
	if c.requiredString(C.sn_string{}, "test") != "" || c.string(C.sn_string{}) != nil || c.bytes(C.sn_bytes{}) != nil {
		return fmt.Errorf("absent value conversion")
	}
	empty := C.sn_string{data: &storage}
	if s := c.string(empty); s == nil || *s != "" {
		return fmt.Errorf("present empty string")
	}
	if b := c.bytes(C.sn_bytes{data: (*C.uint8_t)(unsafe.Pointer(&storage))}); b == nil || len(b) != 0 {
		return fmt.Errorf("present empty bytes")
	}
	if c.length(C.size_t(math.MaxInt), "maximum") != math.MaxInt || c.err != nil {
		return fmt.Errorf("maximum int rejected")
	}
	var stroke C.sn_stroke
	drawn := &stroke
	c = converter{}
	_ = c.page(C.sn_page{drawn_strokes: &drawn, drawn_stroke_count: 1})
	if !errors.Is(c.err, ErrMalformed) {
		return fmt.Errorf("out-of-range drawn stroke accepted")
	}
	return nil
}

func logInstallCount() int32 { return int32(C.sngo_log_installs()) }
