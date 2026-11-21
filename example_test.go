// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"errors"
	"fmt"
	"github.com/altessa-s/snnote-go"
	"log"
	"os"
)

func ExampleParse() {
	data, err := os.ReadFile("testdata/supernote-typescript/blank-a6x-3.15.27-two-pages.note")
	if err != nil {
		log.Fatal(err)
	}
	notebook, err := snnote.Parse(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(notebook.Pages), "pages")
	// Output: 2 pages
}

func ExampleBitmap_Raster() {
	bitmap := snnote.Bitmap{Encoding: snnote.EncodingRattaRLE, Bytes: []byte{0x61, 2, 0x62, 0}}
	raster, err := bitmap.Raster(4, 1, snnote.CurrentPalette())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d %d %v\n", raster.Width, raster.Height, raster.Alpha)
	// Output: 4 1 [255 255 255 0]
}

func ExampleSetLogger() {
	snnote.SetLogger(func(level snnote.LogLevel, message string) {
		log.Printf("snnote %v: %s", level, message)
	})
	defer snnote.SetLogger(nil)
	// Diagnostics from subsequent calls go to this handler.
	// Output:
}

func ExampleError() {
	// A tiny encoded bitmap can request too many output pixels or fail to fill them.
	bitmap := snnote.Bitmap{Encoding: snnote.EncodingRattaRLE, Bytes: []byte{0x61, 0}}
	for _, width := range []int{snnote.MaxRasterPixels + 1, 2} {
		_, err := bitmap.Codes(width, 1)
		switch {
		case errors.Is(err, snnote.ErrLimit):
			fmt.Println("resource budget refused the bitmap")
		case errors.Is(err, snnote.ErrMalformed):
			fmt.Println(err)
		}
	}
	// Output:
	// resource budget refused the bitmap
	// snnote: malformed data: RATTA_RLE runs do not fill bitmap dimensions
}
