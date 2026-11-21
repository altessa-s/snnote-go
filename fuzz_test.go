// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"bytes"
	"errors"
	"github.com/altessa-s/snnote-go"
	"os"
	"reflect"
	"sort"
	"testing"
)

// walk touches every copied field and slice, including values omitted by sndump.
func walk(v reflect.Value) uint64 {
	if !v.IsValid() {
		return 0
	}
	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			return walk(v.Elem())
		}
	case reflect.Struct:
		var sum uint64
		for i := range v.NumField() {
			sum += walk(v.Field(i))
		}
		return sum
	case reflect.Slice, reflect.Array:
		var sum uint64
		for i := range v.Len() {
			sum += walk(v.Index(i))
		}
		return sum
	case reflect.String:
		var sum uint64
		for _, b := range []byte(v.String()) {
			sum += uint64(b)
		}
		return sum
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Int, reflect.Int32, reflect.Int64:
		return uint64(v.Int())
	case reflect.Float64:
		return uint64(v.Float())
	case reflect.Bool:
		if v.Bool() {
			return 1
		}
	}
	return 0
}
func FuzzParse(f *testing.F) {
	paths := fixtures(f)
	sort.Slice(paths, func(i, j int) bool {
		a, err := os.Stat(paths[i])
		if err != nil {
			f.Fatal(err)
		}
		b, err := os.Stat(paths[j])
		if err != nil {
			f.Fatal(err)
		}
		return a.Size() < b.Size()
	})
	for _, path := range paths[:2] {
		f.Add(read(f, path))
	}
	f.Add([]byte{})
	f.Add(damaged())
	f.Fuzz(func(t *testing.T, data []byte) {
		owned := bytes.Clone(data)
		n, err := snnote.Parse(owned)
		if err != nil {
			if n != nil {
				t.Fatal("result with error")
			}
			return
		}
		if n == nil {
			t.Fatal("nil successful model")
		}
		_ = walk(reflect.ValueOf(n))
		clear(owned)
		again, err := snnote.Parse(data)
		if err != nil || !reflect.DeepEqual(n, again) {
			t.Fatal("strings or slices depend on input storage or parse state")
		}
		for _, p := range n.Pages {
			indices := make(map[*snnote.Stroke]bool, len(p.Strokes))
			for i := range p.Strokes {
				indices[&p.Strokes[i]] = true
			}
			for _, stroke := range p.DrawnStrokes {
				if !indices[stroke] {
					t.Fatal("drawn stroke outside Strokes")
				}
			}
			// N5/N6 pages fit this cap; larger pages are traversed but not decoded.
			const maxPixels = 3_000_000
			if p.Size.Width <= 0 || p.Size.Height <= 0 || p.Size.Width > maxPixels || p.Size.Height > maxPixels {
				continue
			}
			w, h := int(p.Size.Width), int(p.Size.Height)
			if w == 0 || h == 0 || w > maxPixels/h {
				continue
			}
			check := func(b *snnote.Bitmap) {
				if b == nil || b.Encoding != snnote.EncodingRattaRLE {
					return
				}
				codes, codeErr := b.Codes(w, h)
				raster, rasterErr := b.Raster(w, h, n.Palette)
				if (codeErr == nil) != (rasterErr == nil) {
					t.Fatalf("decode errors differ: %v, %v", codeErr, rasterErr)
				}
				if codeErr != nil {
					var detail *snnote.Error
					if !errors.As(codeErr, &detail) || !errors.Is(rasterErr, detail) {
						t.Fatal("decode statuses differ")
					}
					return
				}
				if len(codes) != w*h || raster.Width != w || raster.Height != h || len(raster.Gray) != w*h || len(raster.Alpha) != w*h {
					t.Fatal("decode dimensions differ")
				}
				for i, code := range codes {
					if raster.Gray[i] != n.Palette.Inks[code].Gray || raster.Alpha[i] != n.Palette.Inks[code].Alpha {
						t.Fatal("palette mapping differs")
					}
				}
			}
			for _, l := range p.Layers {
				check(l.Bitmap)
			}
			for _, title := range n.Titles {
				if title.PageNumber == p.Number {
					check(title.Bitmap)
				}
			}
			for _, link := range n.Links {
				if link.PageNumber == p.Number {
					check(link.Bitmap)
				}
			}
			for _, tmpl := range n.Templates {
				check(&tmpl.Bitmap)
			}
			check(n.Cover)
		}
	})
}
