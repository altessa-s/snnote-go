// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"errors"
	"github.com/altessa-s/snnote-go"
	"reflect"
	"testing"
)

func TestAccessors(t *testing.T) {
	for pen := -3; pen < 25; pen++ {
		if got := snnote.Pen(pen).IsMarker(); got != (pen == 5 || pen == 11) {
			t.Errorf("marker %d", pen)
		}
	}
	for color := 0; color < 300; color++ {
		gray, ok := snnote.InkColor(color).Gray()
		if ok != (color <= 255) {
			t.Errorf("gray presence %d", color)
		}
		if color == 49 && gray != 157 {
			t.Errorf("legacy marker gray %d", gray)
		}
	}
	if !snnote.ColorMarkerGray.IsMarker() || snnote.ColorWhite.IsMarker() {
		t.Fatal("marker colors")
	}
	s := snnote.TitleStyle(1255254)
	if s.Background() != 255 || s.Foreground() != 254 || s.IsHatched() || !snnote.TitleStyle(1000000).IsHatched() {
		t.Fatal("title components")
	}
	if snnote.LayerBackground.Identifier() != -1 || snnote.LayerBackground.Key() != "BGLAYER" || snnote.LayerName(99).Identifier() != -2 || snnote.LayerName(99).Key() != "" {
		t.Fatal("layer names")
	}
	if snnote.EquipmentManta.Sheet() != (snnote.Size{Width: 1920, Height: 2560}) || snnote.EquipmentOther.Sheet() != (snnote.Size{}) {
		t.Fatal("sheet")
	}
	if snnote.PaletteForFormat(20230015) != snnote.CurrentPalette() || snnote.PaletteForFormat(20220011) != snnote.LegacyPalette() {
		t.Fatal("palette selection")
	}
	for _, v := range []interface{ String() string }{snnote.EncodingRattaRLE, snnote.EquipmentNomad, snnote.KindMark, snnote.OrientationLandscape, snnote.LayerBackground, snnote.StrokeKindHeadingBox, snnote.RoleEraser, snnote.StrokeStatusDeleted, snnote.CategoryTextBox, snnote.PenInk, snnote.ColorBlack, snnote.LogError} {
		if v.String() == "" {
			t.Fatal("empty enum name")
		}
	}
	if snnote.Pen(999).String() != "pen(999)" || snnote.InkColor(999).String() != "inkcolor(999)" {
		t.Fatal("raw unknown values lost")
	}
}
func TestMetadata(t *testing.T) {
	m := snnote.Metadata{Entries: []snnote.MetadataEntry{{Key: "a", Value: ""}, {Key: "a", Value: "two\x00tail"}, {Key: "a\x00b", Value: "three"}}}
	if s, ok := m.Get("a"); !ok || s != "" {
		t.Fatal("empty value lost")
	}
	if _, ok := m.Get("absent"); ok {
		t.Fatal("missing key found")
	}
	if got := m.Values("a"); !reflect.DeepEqual(got, []string{"", "two\x00tail"}) {
		t.Fatal(got)
	}
	if s, ok := m.Get("a\x00b"); !ok || s != "three" {
		t.Fatal("embedded NUL key")
	}
}
func TestBitmaps(t *testing.T) {
	b := snnote.Bitmap{Encoding: snnote.EncodingRattaRLE, Bytes: []byte{0x61, 2, 0x62, 0}}
	codes, err := b.Codes(4, 1)
	if err != nil || !reflect.DeepEqual(codes, []byte{0x61, 0x61, 0x61, 0x62}) {
		t.Fatalf("%v %v", codes, err)
	}
	palette := snnote.CurrentPalette()
	palette.Inks[0x61] = snnote.Ink{Gray: 42, Alpha: 128}
	r, err := b.Raster(4, 1, palette)
	if err != nil || !reflect.DeepEqual(r.Gray, []byte{42, 42, 42, 255}) || !reflect.DeepEqual(r.Alpha, []byte{128, 128, 128, 0}) {
		t.Fatalf("%+v %v", r, err)
	}
	for _, tc := range []struct {
		encoding snnote.Encoding
		w, h     int
		want     error
	}{
		{snnote.EncodingPNG, 1, 1, snnote.ErrImageEncoding}, {snnote.EncodingJPEG, 1, 1, snnote.ErrImageEncoding}, {snnote.EncodingOther, 1, 1, snnote.ErrUnsupportedEncoding}, {snnote.EncodingRattaRLE, 3, 1, snnote.ErrMalformed}, {snnote.EncodingRattaRLE, 0, 1, snnote.ErrInvalidArgument}, {snnote.EncodingRattaRLE, -1, 1, snnote.ErrInvalidArgument}, {snnote.EncodingRattaRLE, snnote.MaxRasterPixels + 1, 1, snnote.ErrLimit},
	} {
		b.Encoding = tc.encoding
		if _, err := b.Codes(tc.w, tc.h); !errors.Is(err, tc.want) {
			t.Errorf("codes: %v, want %v", err, tc.want)
		}
		if _, err := b.Raster(tc.w, tc.h, palette); !errors.Is(err, tc.want) {
			t.Errorf("raster: %v, want %v", err, tc.want)
		}
	}
}
func TestCorpusBitmaps(t *testing.T) {
	var rle, png int
	for _, path := range fixtures(t) {
		n := parse(t, read(t, path))
		for _, p := range n.Pages {
			for _, l := range p.Layers {
				if l.Bitmap == nil {
					continue
				}
				b := l.Bitmap
				if b.Encoding == snnote.EncodingPNG {
					png++
					if _, err := b.Codes(int(p.Size.Width), int(p.Size.Height)); !errors.Is(err, snnote.ErrImageEncoding) {
						t.Fatal(err)
					}
					continue
				}
				if b.Encoding != snnote.EncodingRattaRLE {
					continue
				}
				codes, err := b.Codes(int(p.Size.Width), int(p.Size.Height))
				if err != nil {
					t.Fatalf("%s %s: %v", path, l.Name, err)
				}
				want := int(p.Size.Width) * int(p.Size.Height)
				if len(codes) != want {
					t.Fatalf("pixels %d != %d", len(codes), want)
				}
				raster, err := b.Raster(int(p.Size.Width), int(p.Size.Height), n.Palette)
				if err != nil {
					t.Fatal(err)
				}
				for i, c := range codes {
					ink := n.Palette.Inks[c]
					if raster.Gray[i] != ink.Gray || raster.Alpha[i] != ink.Alpha {
						t.Fatalf("pixel %d", i)
					}
				}
				rle++
			}
		}
	}
	if rle == 0 || png == 0 {
		t.Fatalf("missing encoding coverage: %d RLE, %d PNG", rle, png)
	}
}

func TestKnownSheetPixelCount(t *testing.T) {
	n := parse(t, read(t, "testdata/supernote-typescript/blank-a6x-3.15.27-two-pages.note"))
	for _, p := range n.Pages {
		if p.Size != (snnote.Size{Width: 1404, Height: 1872}) {
			t.Fatalf("unexpected sheet %+v", p.Size)
		}
		for _, l := range p.Layers {
			if l.Bitmap == nil || l.Bitmap.Encoding != snnote.EncodingRattaRLE {
				continue
			}
			codes, err := l.Bitmap.Codes(1404, 1872)
			if err != nil || len(codes) != 2628288 {
				t.Fatalf("known sheet pixels %d: %v", len(codes), err)
			}
		}
	}
}

func TestDrawnStrokeBackingArray(t *testing.T) {
	p := snnote.Page{Strokes: make([]snnote.Stroke, 1)}
	p.Strokes[0].ID = 42
	p.DrawnStrokes = []*snnote.Stroke{&p.Strokes[0]}
	p.Strokes = append(p.Strokes, snnote.Stroke{ID: 43})
	p.Strokes[0].ID = 44
	if p.DrawnStrokes[0].ID != 42 {
		t.Fatal("append retargeted original pointer")
	}
	p.Strokes = p.Strokes[:0]
	if p.DrawnStrokes[0].ID != 42 {
		t.Fatal("reslicing lost original record")
	}
}
