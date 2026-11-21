// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

import (
	"errors"
	"math"
	"testing"
)

func TestConversionBudgetAccounting(t *testing.T) {
	c := converter{budget: 16}
	if !c.charge(2, 8) || c.used != 16 {
		t.Fatal("exact budget rejected")
	}
	if c.charge(1, 1) || !errors.Is(c.err, ErrLimit) || c.used != 16 {
		t.Fatal("exhausted budget accepted")
	}
	first := c.err
	if c.charge(0, 0) || !errors.Is(c.err, first) {
		t.Fatal("conversion continued after error")
	}
	c = converter{budget: 16}
	if c.charge(math.MaxInt, 16) || !errors.Is(c.err, ErrLimit) || c.used != 0 {
		t.Fatal("overflowing allocation accepted")
	}
	c = converter{}
	if !c.charge(ParseBudgetFloor, 1) || c.charge(1, 1) {
		t.Fatal("zero-value budget is not the floor")
	}
}

func TestConversionBudget(t *testing.T) {
	if got := conversionBudget(0); got != ParseBudgetFloor {
		t.Fatalf("empty budget = %d", got)
	}
	if got := conversionBudget(1024); got != ParseBudgetFloor+1024*ParseBudgetMultiple {
		t.Fatalf("budget = %d", got)
	}
	// This architecture may not represent an input large enough to saturate.
	if uint64(math.MaxInt) > (math.MaxUint64-ParseBudgetFloor)/ParseBudgetMultiple {
		if got := conversionBudget(math.MaxInt); got != math.MaxUint64 {
			t.Fatalf("saturated budget = %d", got)
		}
	}
}
