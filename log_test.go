// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/altessa-s/snnote-go"
)

type collected struct {
	mu     sync.Mutex
	levels []snnote.LogLevel
	lines  []string
}

func (c *collected) log(level snnote.LogLevel, message string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.levels = append(c.levels, level)
	c.lines = append(c.lines, message)
}

func (c *collected) saw(level snnote.LogLevel, needle string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, line := range c.lines {
		if c.levels[i] == level && strings.Contains(line, needle) {
			return true
		}
	}
	return false
}

func (c *collected) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.lines)
}

func TestALostPageReachesTheLogger(t *testing.T) {
	var seen collected
	snnote.SetLogger(seen.log)
	defer snnote.SetLogger(nil)

	page := parse(t, damaged())
	if page.SkippedPages != 1 {
		t.Errorf("skipped = %d, want 1", page.SkippedPages)
	}
	if !seen.saw(snnote.LogError, "page") {
		t.Errorf("the loss was not said out loud: %q", seen.lines)
	}
}

func TestNilSilencesTheLogger(t *testing.T) {
	var seen collected
	snnote.SetLogger(seen.log)
	snnote.SetLogger(nil)
	parse(t, damaged())
	if seen.count() != 0 {
		t.Errorf("a silenced logger heard %q", seen.lines)
	}
}

func TestTheLatestLoggerStands(t *testing.T) {
	var first, second collected
	snnote.SetLogger(first.log)
	snnote.SetLogger(second.log)
	defer snnote.SetLogger(nil)
	parse(t, damaged())
	if first.count() != 0 || second.count() == 0 {
		t.Errorf("first heard %d lines, second %d", first.count(), second.count())
	}
}

// Concurrent parses use separate handler goroutines; replacement is safe.
func TestTheLoggerUnderConcurrentParsing(t *testing.T) {
	var seen collected
	snnote.SetLogger(seen.log)
	defer snnote.SetLogger(nil)
	damaged := damaged()

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 20 {
				if _, err := snnote.Parse(damaged); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 20 {
			if i%2 == 0 {
				snnote.SetLogger(nil)
			} else {
				snnote.SetLogger(seen.log)
			}
		}
	}()
	wg.Wait()
}

// A handler that panics loses its line, not the parse.
func TestAPanickingLoggerDoesNotStopTheParse(t *testing.T) {
	calls := 0
	snnote.SetLogger(func(snnote.LogLevel, string) {
		calls++
		panic("handler failure")
	})
	defer snnote.SetLogger(nil)
	page, err := snnote.Parse(damaged())
	if err != nil || page.SkippedPages != 1 {
		t.Fatalf("parse = %+v, %v; want one skipped block", page, err)
	}
	if calls == 0 {
		t.Errorf("the handler was never called")
	}
}

// Reentrant library calls run outside the enclosing parse's work meter, so
// neither normal returns nor recovered panics can strand its recovery frame.
func TestLoggerCanCallLibrary(t *testing.T) {
	calls := 0
	snnote.SetLogger(func(snnote.LogLevel, string) {
		calls++
		if snnote.CurrentPalette() != snnote.PaletteForFormat(20230015) {
			t.Error("palette changed in callback")
		}
		if _, err := snnote.Parse(finish(begin(), "")); err != nil {
			t.Error(err)
		}
	})
	defer snnote.SetLogger(nil)
	parse(t, damaged())
	if calls == 0 {
		t.Fatal("callback never ran")
	}
	if _, err := snnote.Parse(hostile()); !errors.Is(err, snnote.ErrLimit) {
		t.Fatalf("enclosing work meter was not restored: %v", err)
	}
}

func TestInFlightLineSurvivesSilencing(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	delivered := make(chan struct{})
	snnote.SetLogger(func(snnote.LogLevel, string) {
		close(started)
		<-release
		close(delivered)
	})
	defer snnote.SetLogger(nil)
	done := make(chan error, 1)
	go func() { _, err := snnote.Parse(damaged()); done <- err }()
	<-started
	snnote.SetLogger(nil)
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	<-delivered
}
