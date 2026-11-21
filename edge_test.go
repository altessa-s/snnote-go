// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/altessa-s/snnote-go"
)

func TestBitmapErrorsAreIndependent(t *testing.T) {
	b := snnote.Bitmap{}
	for _, decode := range []func() error{
		func() error { _, err := b.Codes(0, 1); return err },
		func() error { _, err := b.Raster(0, 1, snnote.Palette{}); return err },
	} {
		err := decode()
		var detail *snnote.Error
		if !errors.As(err, &detail) || !errors.Is(err, snnote.ErrInvalidArgument) || detail == snnote.ErrInvalidArgument || detail.Message == "" {
			t.Fatalf("bad error: %v", err)
		}
		detail.Status = snnote.StatusOK
		if snnote.ErrInvalidArgument.Status != snnote.StatusInvalidArgument {
			t.Fatal("sentinel modified")
		}
	}
}

func TestWideEnumsDoNotWrap(t *testing.T) {
	if strconv.IntSize < 64 {
		t.Skip("requires 64-bit int")
	}
	wide := int64(1) << 32
	b := snnote.Bitmap{Encoding: snnote.Encoding(wide), Bytes: []byte{0x61, 0}}
	if _, err := b.Codes(1, 1); !errors.Is(err, snnote.ErrInvalidArgument) {
		t.Fatal(err)
	}
	if _, err := b.Raster(1, 1, snnote.Palette{}); !errors.Is(err, snnote.ErrInvalidArgument) {
		t.Fatal(err)
	}
	if snnote.Status(wide).String() != "status(4294967296)" || snnote.LogLevel(wide).String() != "log_level(4294967296)" {
		t.Fatal("unknown number lost")
	}
	if snnote.LayerName(wide).Key() != "" || snnote.LayerName(wide).Identifier() != -2 || snnote.Equipment(wide+1).Sheet() != (snnote.Size{}) {
		t.Fatal("accessor truncated enum")
	}
}
