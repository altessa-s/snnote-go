// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

// Package snnote reads Supernote .note notebooks and .mark annotations through
// the vendored snnote-lib C23 library. It requires cgo and a C23 compiler.
//
// Parse returns Go values and releases the C notebook before returning. All
// strings are copied by byte length, preserving embedded NULs. Optional strings,
// numbers, rectangles and nested records are pointers: nil means absent, while
// a pointer to an empty string or zero means present. Required text and metadata
// keys and values use plain strings. Byte payloads are copied slices; nil
// means absent. Array counts are represented by slice lengths. DrawnStrokes
// refers into the page's Go Strokes slice. Metadata preserves repeated keys
// and their order; Get returns the first match and Values returns all matches.
//
// Bitmap.Codes and Bitmap.Raster decode RATTA_RLE only. PNG and JPEG require
// an external image decoder; page composition belongs to a rendering library.
// Failures are *Error values, matched by status through errors.Is. ErrLimit
// distinguishes refused resource use from malformed content. Counts larger
// than a Go int can hold return ErrOutOfMemory.
//
// # Concurrency
//
// Calls keep no state between calls, so any number may run at once. A goroutine
// stays on the same OS thread for the duration of a C call. The library's
// thread-local work meter and its internal setjmp recovery live entirely
// inside that call; no thread binding is needed between calls. Go memory is
// never retained by C. The logger is the one process-wide setting and may be
// replaced while calls run. Each handler runs on a separate goroutine,
// synchronously with the parse that produced its line. A handler that panics or
// exits loses only that line. A line already in flight may still be delivered
// after SetLogger(nil) returns.
// Returned values may be shared for reading; callers synchronize mutations.
package snnote
