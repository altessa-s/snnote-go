//go:build darwin || linux

// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/altessa-s/snnote-go"
)

func residentPeak(t *testing.T) int64 {
	t.Helper()
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		t.Fatal(err)
	}
	rss := int64(usage.Maxrss)
	if runtime.GOOS == "linux" {
		rss *= 1024
	}
	return rss
}

// C storage is not counted by the Go heap, so use the process resident peak.
func TestLoggerGoexitReturnsAndReleasesStorage(t *testing.T) {
	input := damaged()
	snnote.SetLogger(func(snnote.LogLevel, string) { runtime.Goexit() })
	defer snnote.SetLogger(nil)
	run := func(count int) {
		t.Helper()
		done := make(chan error, 1)
		go func() {
			for range count {
				n, err := snnote.Parse(input)
				if err != nil {
					done <- err
					return
				}
				if n == nil || n.SkippedPages != 1 {
					done <- &snnote.Error{Status: snnote.StatusMalformed, Message: "unexpected parse result"}
					return
				}
			}
			done <- nil
		}()
		select {
		case err := <-done:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(30 * time.Second):
			t.Fatal("parse did not return after handler Goexit")
		}
	}
	run(1000)
	runtime.GC()
	before := residentPeak(t)
	run(12000)
	runtime.GC()
	growth := residentPeak(t) - before
	t.Logf("resident peak growth over 12000 parses: %d bytes", growth)
	// The leaked C arena would exceed 190 MiB. Allow runtime/allocator noise.
	if growth > 64<<20 {
		t.Fatalf("resident growth %d exceeds 64 MiB", growth)
	}
}
