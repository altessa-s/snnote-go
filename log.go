// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote

/*
#include "snnote/log.h"

void sngo_log_install(void);
*/
import "C"

import (
	"strconv"
	"sync"
	"sync/atomic"
)

// LogLevel is how much a diagnostic line matters: a piece of the page was
// lost, or the format held something the parser does not know and read past.
type LogLevel int

// The log levels.
const (
	LogDebug LogLevel = iota
	LogError
)

// String names the level and retains unknown numbers in its fallback.
func (l LogLevel) String() string {
	if l == LogError {
		return "error"
	}
	if l == LogDebug {
		return "debug"
	}
	return "log_level(" + strconv.Itoa(int(l)) + ")"
}

// logHandler is what the C sink calls through. Read on whichever thread is
// parsing and swapped by SetLogger, hence atomic; nil drops every line.
var logHandler atomic.Pointer[func(LogLevel, string)]

// logInstall installs the C sink once per process. The C library keeps every
// sink it is ever given, so installing on every SetLogger would grow without
// bound; the sink stays installed and the Go handler decides. If the first
// installation fails to allocate, logging stays disabled for the process: the
// C API reports no installation result and sync.Once does not retry.
var logInstall sync.Once

// SetLogger names the function the parser's diagnostics go to, for the whole
// process; nil goes back to silence, which is the default. It may be called
// at any time, from any goroutine.
//
// Each handler runs on a separate goroutine, synchronously with the parse that
// produced the line. Concurrent parses may call it concurrently, so it must
// synchronize shared state. A handler that panics or exits with runtime.Goexit
// loses only that line; the parse continues normally. The parse waits for the
// handler to return, so a handler must not wait for that parse. A line already
// in flight may still be delivered after SetLogger(nil) returns.
//
// Once a logger has been set, the C library formats its lines even while the
// handler is nil; a line is a few dozen bytes, built only where something was
// lost or unknown.
//
// Nothing in a message is a secret: every value is a count, a block type or a
// raw style number, never the page's writing, text or file names.
func SetLogger(handler func(level LogLevel, message string)) {
	if handler == nil {
		logHandler.Store(nil)
		return
	}
	logHandler.Store(&handler)
	logInstall.Do(func() {
		C.sngo_log_install()
	})
}

//export sngoLog
func sngoLog(level C.int, message *C.char) {
	handler := logHandler.Load()
	if handler == nil {
		return
	}
	line := C.GoString(message)
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() { _ = recover() }()
		(*handler)(LogLevel(level), line)
	}()
	<-done
}
