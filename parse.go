// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#include "snnote/snnote.h"
*/
import "C"

import (
	"math"
	"unsafe"
)

// Parse reads a complete notebook and copies its model before releasing all C memory.
// Optional damaged content is counted in SkippedPages or SkippedStrokes.
// The Go copy has its own ParseBudgetFloor + ParseBudgetMultiple*len(data)
// allocation budget; exceeding it returns ErrLimit and no partial notebook.
func Parse(data []byte) (*Notebook, error) {
	var notebook *C.sn_notebook
	var cerr C.sn_error
	status := C.sn_parse_notebook(bytesPointer(data), C.size_t(len(data)), &notebook, &cerr)
	if err := errorOf(status, &cerr); err != nil {
		return nil, err
	}
	defer C.sn_notebook_free(notebook)
	c := converter{budget: conversionBudget(len(data))}
	if !c.charge(1, unsafe.Sizeof(Notebook{})) {
		return nil, c.err
	}
	out := c.notebook(*notebook)
	if c.err != nil {
		return nil, c.err
	}
	return &out, nil
}

// An empty input is still a valid buffer argument, rather than a null caller error.
var emptyInput [1]byte

func bytesPointer(data []byte) *C.uint8_t {
	if len(data) == 0 {
		return (*C.uint8_t)(unsafe.Pointer(&emptyInput[0]))
	}
	return (*C.uint8_t)(unsafe.Pointer(&data[0]))
}

type converter struct {
	budget        uint64
	used          uint64
	err           error
	requiredNulls *[]string
}

// conversionBudget saturates rather than wrapping for very large inputs.
func conversionBudget(inputLength int) uint64 {
	if uint64(inputLength) > (math.MaxUint64-ParseBudgetFloor)/ParseBudgetMultiple {
		return math.MaxUint64
	}
	return ParseBudgetFloor + ParseBudgetMultiple*uint64(inputLength)
}

// charge accounts for Go payload storage, not allocator rounding or GC overhead.
// A zero-value converter has the minimum budget, including in test helpers.
func (c *converter) charge(count int, size uintptr) bool {
	if c.err != nil {
		return false
	}
	budget := c.budget
	if budget == 0 {
		budget = ParseBudgetFloor
	}
	if size != 0 && uint64(count) > (budget-c.used)/uint64(size) {
		c.err = &Error{Status: StatusLimit, Message: "Go model copy allocation budget exceeded"}
		return false
	}
	c.used += uint64(count) * uint64(size)
	return true
}

func (c *converter) allocationLength(p unsafe.Pointer, count C.size_t, what string, size uintptr) int {
	n := c.arrayLength(p, count, what)
	if !c.charge(n, size) {
		return 0
	}
	return n
}

func (c *converter) length(n C.size_t, what string) int {
	if c.err != nil {
		return 0
	}
	if uint64(n) > math.MaxInt {
		c.err = &Error{Status: StatusOutOfMemory, Message: what + " too large for a Go int"}
		return 0
	}
	return int(n)
}

// arrayLength validates a C pointer/count pair without inspecting its memory.
// A non-null zero-length allocation is valid; callers must skip the slice.
func (c *converter) arrayLength(p unsafe.Pointer, count C.size_t, what string) int {
	n := c.length(count, what)
	if c.err != nil || n == 0 {
		return 0
	}
	if p == nil {
		c.err = &Error{Status: StatusMalformed, Message: what + " has null storage"}
		return 0
	}
	return n
}

// Lengths preserve embedded NULs and avoid the narrower C int used by GoStringN.
func (c *converter) string(v C.sn_string) *string {
	n := c.arrayLength(unsafe.Pointer(v.data), v.length, "string length")
	if c.err != nil || v.data == nil || !c.charge(n, 1) || !c.charge(1, unsafe.Sizeof("")) {
		return nil
	}
	s := ""
	if n > 0 {
		s = string(unsafe.Slice((*byte)(unsafe.Pointer(v.data)), n))
	}
	return &s
}

func (c *converter) requiredString(v C.sn_string, field string) string {
	if v.data == nil && c.requiredNulls != nil {
		*c.requiredNulls = append(*c.requiredNulls, field)
	}
	if s := c.string(v); s != nil {
		return *s
	}
	return ""
}

func (c *converter) bytes(v C.sn_bytes) []byte {
	n := c.arrayLength(unsafe.Pointer(v.data), v.size, "byte length")
	if c.err != nil || v.data == nil || !c.charge(n, 1) {
		return nil
	}
	out := make([]byte, n)
	if n > 0 {
		copy(out, unsafe.Slice((*byte)(unsafe.Pointer(v.data)), n))
	}
	return out
}
