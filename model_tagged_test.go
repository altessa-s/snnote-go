//go:build snnote_test

// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"github.com/altessa-s/snnote-go"
	"strings"
	"testing"
)

func TestConstants(t *testing.T) {
	version, constants := snnote.CConstants()
	if version != snnote.Version {
		t.Fatalf("version %s != %s", version, snnote.Version)
	}
	for name, pair := range constants {
		if pair[0] != pair[1] {
			t.Errorf("%s: Go %d, C %d", name, pair[0], pair[1])
		}
	}
}

func TestGoLookupsMatchC(t *testing.T) {
	for _, path := range fixtures(t) {
		data := read(t, path)
		if !snnote.CLookupMatches(data, parse(t, data)) {
			t.Errorf("lookup mismatch: %s", path)
		}
	}
	var p *snnote.Page
	if p.Layer(snnote.LayerMain) != nil {
		t.Fatal("nil page has a layer")
	}
}

func TestRequiredStringsPresentInC(t *testing.T) {
	for _, path := range fixtures(t) {
		t.Run(path, func(t *testing.T) {
			missing, err := snnote.RequiredStringNulls(read(t, path))
			if err != nil {
				t.Fatal(err)
			}
			if len(missing) != 0 {
				t.Fatalf("null required C strings: %v", missing)
			}
		})
	}
}

func TestConversionBoundaries(t *testing.T) {
	if err := snnote.ConversionBoundaryChecks(); err != nil {
		t.Fatal(err)
	}
}

func TestSmallSeedValuesMatchC(t *testing.T) {
	for name, data := range smallSeeds() {
		n := parse(t, data)
		if !snnote.CLookupMatches(data, n) {
			t.Errorf("C comparison: %s", name)
		}
		missing, err := snnote.RequiredStringNulls(data)
		if err != nil || len(missing) != 0 {
			t.Errorf("required strings: %s: %v %v", name, missing, err)
		}
	}
}

// Equal identifiers must not hide a wrong DrawnStrokes pointer.
func TestCComparisonDetectsOmittedFieldChanges(t *testing.T) {
	cases := []struct {
		seed   int
		change func(*snnote.Notebook)
	}{
		{0, func(n *snnote.Notebook) { n.Pages[0].Strokes[0].Kind++ }},
		{0, func(n *snnote.Notebook) { n.Pages[0].Strokes[0].Status++ }},
		{0, func(n *snnote.Notebook) { n.Equipment++ }},
		{4, func(n *snnote.Notebook) { *n.Links[0].URL += "changed" }},
		{9, func(n *snnote.Notebook) { n.Palette.Inks[0].Gray++ }},
		{8, func(n *snnote.Notebook) { n.Templates[0], n.Templates[1] = n.Templates[1], n.Templates[0] }},
		{0, func(n *snnote.Notebook) { n.Pages[0].DrawnStrokes[0] = &n.Pages[0].Strokes[1] }},
	}
	for i, tc := range cases {
		data := seedNotebook(tc.seed)
		n := parse(t, data)
		if !snnote.CLookupMatches(data, n) {
			t.Fatalf("case %d baseline mismatch", i)
		}
		tc.change(n)
		if snnote.CLookupMatches(data, n) {
			t.Errorf("case %d mutation went undetected", i)
		}
	}
}

func TestCLookupLargeMetadata(t *testing.T) {
	for _, tc := range []struct {
		name        string
		count, size int
	}{
		{"single", 1, snnote.ParseBudgetFloor + 1},
		{"repeated", 3, snnote.ParseBudgetFloor / 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			value := strings.Repeat("x", tc.size)
			data := finish(begin(), strings.Repeat("<CUSTOM:"+value+">", tc.count))
			n := parse(t, data)
			if !snnote.CLookupMatches(data, n) {
				t.Fatal("large metadata lookup mismatch")
			}
			n.Footer.Entries[len(n.Footer.Entries)-1].Value = "changed"
			if snnote.CLookupMatches(data, n) {
				t.Fatal("metadata mutation went undetected")
			}
		})
	}
}
