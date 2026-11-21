// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/altessa-s/snnote-go"
)

func TestSharedPageCopyBudget(t *testing.T) {
	for _, pages := range []int{2, 256} {
		t.Run(fmt.Sprint(pages), func(t *testing.T) {
			data := begin()
			bitmap := len(data)
			data = block(data, make([]byte, 32<<10))
			layer := len(data)
			data = block(data, []byte(fmt.Sprintf("<LAYERBITMAP:%d><LAYERPROTOCOL:RATTA_RLE>", bitmap)))
			page := len(data)
			data = block(data, []byte(fmt.Sprintf("<MAINLAYER:%d>", layer)))
			var footer strings.Builder
			for i := 1; i <= pages; i++ {
				fmt.Fprintf(&footer, "<PAGE%d:%d>", i, page)
			}
			data = finish(data, footer.String())
			n, err := snnote.Parse(data)
			if pages == 256 {
				var status *snnote.Error
				if n != nil || !errors.Is(err, snnote.ErrLimit) || !errors.As(err, &status) || !strings.Contains(status.Message, "Go model copy") {
					t.Fatalf("shared page amplification: notebook present=%v, error=%v", n != nil, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(n.Pages) != pages {
				t.Fatalf("pages = %d, want %d", len(n.Pages), pages)
			}
			first := n.Pages[0].Layers[0].Bitmap.Bytes
			second := n.Pages[1].Layers[0].Bitmap.Bytes
			if len(first) != 32<<10 || len(second) != len(first) {
				t.Fatal("missing copied bitmap")
			}
			first[0] = 1
			if second[0] != 0 {
				t.Fatal("shared C bitmap became aliased Go storage")
			}
		})
	}
}

func TestBitmapCopyAboveParseFloor(t *testing.T) {
	// A compact RLE stream expands beyond the minimum notebook-copy budget.
	const pixels = snnote.ParseBudgetFloor + 16384
	var encoded []byte
	for i := 0; i < pixels/16384; i++ {
		encoded = append(encoded, 0x61, 255)
	}
	b := snnote.Bitmap{Encoding: snnote.EncodingRattaRLE, Bytes: encoded}
	codes, err := b.Codes(pixels, 1)
	if err != nil || len(codes) != pixels {
		t.Fatalf("codes length=%d, error=%v", len(codes), err)
	}
	raster, err := b.Raster(pixels, 1, snnote.CurrentPalette())
	if err != nil {
		t.Fatal(err)
	}
	if len(raster.Gray) != pixels || len(raster.Alpha) != pixels {
		t.Fatal("incomplete raster planes")
	}
}
