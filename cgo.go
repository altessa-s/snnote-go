// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#cgo CFLAGS: -std=c2x -I${SRCDIR}/internal/clib/include -I${SRCDIR}/internal/clib/src
#cgo linux LDFLAGS: -lm
#include "snnote/snnote.h"
*/
import "C"

// Version is the vendored library release declared by its public header.
const Version = "1.0.0"

const (
	// ParseBudgetFloor is the library limit for parse budget floor.
	ParseBudgetFloor = 1048576
	// ParseBudgetMultiple is the library limit for parse budget multiple.
	ParseBudgetMultiple = 128
	// ParseWorkFloor is the library limit for parse work floor.
	ParseWorkFloor = 1048576
	// ParseWorkMultiple is the library limit for parse work multiple.
	ParseWorkMultiple = 256
	// ProtocolMaxLength is the library limit for protocol max length.
	ProtocolMaxLength = 32
	// JSONMaxDepth is the library limit for json max depth.
	JSONMaxDepth = 64
	// MaxRasterPixels is the library limit for max raster pixels.
	MaxRasterPixels = 16777216
)
