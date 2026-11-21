// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#include "snnote/snnote.h"
*/
import "C"
import "strconv"

// TitleStyle retains the decimal background and foreground components without clamping.
type TitleStyle int64

// String preserves the stored style as a decimal number.
func (s TitleStyle) String() string { return strconv.FormatInt(int64(s), 10) }

// Background returns the stored BBB component of 1BBBFFF.
func (s TitleStyle) Background() int64 { return int64(C.sn_title_style_background(C.int64_t(s))) }

// Foreground returns the stored FFF component of 1BBBFFF.
func (s TitleStyle) Foreground() int64 { return int64(C.sn_title_style_foreground(C.int64_t(s))) }

// IsHatched identifies black-on-black title boxes.
func (s TitleStyle) IsHatched() bool { return bool(C.sn_title_style_is_hatched(C.int64_t(s))) }

// IsMarker recognizes both firmware spellings of the marker pen.
func (p Pen) IsMarker() bool { return bool(C.sn_pen_is_marker(C.int32_t(p))) }

// Gray maps a known byte color to intensity, returning false above the byte range.
func (v InkColor) Gray() (uint8, bool) {
	var gray C.uint8_t
	ok := C.sn_ink_color_gray(C.uint32_t(v), &gray)
	return uint8(gray), bool(ok)
}

// IsMarker distinguishes marker colors from ordinary ink of the same intensity.
func (v InkColor) IsMarker() bool { return bool(C.sn_ink_color_is_marker(C.uint32_t(v))) }

// Key returns the metadata key, or an empty string for an invalid layer name.
func (v LayerName) Key() string {
	if v < LayerMain || v > LayerBackground {
		return ""
	}
	return C.GoString(C.sn_layer_name_key(C.sn_layer_name(v)))
}

// Identifier returns -1 for the background and -2 for an invalid name.
func (v LayerName) Identifier() int64 {
	if v < LayerMain || v > LayerBackground {
		return -2
	}
	return int64(C.sn_layer_name_identifier(C.sn_layer_name(v)))
}

// Sheet returns portrait pixel dimensions, or zero dimensions for unknown equipment.
func (v Equipment) Sheet() Size {
	if v < EquipmentOther || v > EquipmentManta {
		return Size{}
	}
	var c converter
	return c.size(C.sn_equipment_sheet(C.sn_equipment(v)))
}

// LegacyPalette returns the palette used before format suffix 15.
func LegacyPalette() Palette { var c converter; return c.palette(C.sn_palette_legacy()) }

// CurrentPalette returns the palette used from format suffix 15 onward.
func CurrentPalette() Palette { var c converter; return c.palette(C.sn_palette_current()) }

// PaletteForFormat selects by the format version modulo 10000.
func PaletteForFormat(version int64) Palette {
	var c converter
	return c.palette(C.sn_palette_for_format(C.int64_t(version)))
}

// Layer finds the first layer of the given identity; a nil page has none.
func (p *Page) Layer(name LayerName) *Layer {
	if p != nil {
		for i := range p.Layers {
			if p.Layers[i].Name == name {
				return &p.Layers[i]
			}
		}
	}
	return nil
}

// Get returns the first value for a key, distinguishing absence from an empty value.
func (m Metadata) Get(key string) (string, bool) {
	for _, entry := range m.Entries {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return "", false
}

// Values returns all values for a key in stored order, including repetitions.
func (m Metadata) Values(key string) []string {
	var out []string
	for _, entry := range m.Entries {
		if entry.Key == key {
			out = append(out, entry.Value)
		}
	}
	return out
}
