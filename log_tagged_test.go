//go:build snnote_test

// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"github.com/altessa-s/snnote-go"
	"testing"
)

// The C library keeps every sink it is given, so the sink is installed once
// however often the logger changes.
func TestTheSinkIsInstalledOnce(t *testing.T) {
	var seen collected
	for range 1000 {
		snnote.SetLogger(seen.log)
		snnote.SetLogger(nil)
	}
	if got := snnote.LogInstalls(); got != 1 {
		t.Errorf("installed %d times, want 1", got)
	}
	snnote.SetLogger(seen.log)
	defer snnote.SetLogger(nil)
	parse(t, damaged())
	if seen.count() == 0 {
		t.Errorf("the logger set after the cycles hears nothing")
	}
}
