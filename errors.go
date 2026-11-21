// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#include "snnote/snnote.h"
*/
import "C"

import (
	"strconv"
	"strings"
)

// Status identifies the reason a call failed.
type Status int

// Status values agree with the C library.
const (
	StatusOK Status = iota
	StatusInvalidArgument
	StatusOutOfMemory
	StatusUnsupportedFormat
	StatusUnexpectedEOF
	StatusMalformed
	StatusImageEncoding
	StatusUnsupportedEncoding
	StatusLimit
)

// String returns the library's English description or status(number) for unknown values.
func (s Status) String() string {
	if s < StatusOK || s > StatusLimit {
		return "status(" + strconv.Itoa(int(s)) + ")"
	}
	return C.GoString(C.sn_status_string(C.sn_status(s)))
}

// Error records a failure, its format location or count, and a content-free diagnostic.
type Error struct {
	// Status identifies the failure category used by errors.Is.
	Status Status
	// Value is the format offset or count named by Message; zero when no numeric
	// location applies.
	Value int64
	// Message is a content-free diagnostic without the trailing status phrase.
	Message string
}

// Sentinels match errors by status through errors.Is.
var (
	ErrInvalidArgument     = &Error{Status: StatusInvalidArgument}
	ErrOutOfMemory         = &Error{Status: StatusOutOfMemory}
	ErrUnsupportedFormat   = &Error{Status: StatusUnsupportedFormat}
	ErrUnexpectedEOF       = &Error{Status: StatusUnexpectedEOF}
	ErrMalformed           = &Error{Status: StatusMalformed}
	ErrImageEncoding       = &Error{Status: StatusImageEncoding}
	ErrUnsupportedEncoding = &Error{Status: StatusUnsupportedEncoding}
	ErrLimit               = &Error{Status: StatusLimit}
)

// Error spells the status followed by its optional diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "snnote: <nil>"
	}
	text := "snnote: " + e.Status.String()
	if e.Message != "" {
		text += ": " + e.Message
	}
	return text
}

// Is compares statuses, so callers can distinguish limits from malformed files.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t != nil && e != nil && t.Status == e.Status
}
func errorOf(status C.sn_status, e *C.sn_error) error {
	if status == C.SN_OK {
		return nil
	}
	return &Error{Status: Status(status), Value: int64(e.value), Message: strings.TrimSuffix(C.GoString(&e.message[0]), ": "+Status(status).String())}
}
