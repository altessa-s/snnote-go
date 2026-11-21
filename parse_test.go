// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/altessa-s/snnote-go"
)

func block(data []byte, payload []byte) []byte {
	data = binary.LittleEndian.AppendUint32(data, uint32(len(payload)))
	return append(data, payload...)
}
func begin() []byte {
	data := []byte("noteSN_FILE_VER_20230015")
	return block(data, []byte("<APPLY_EQUIPMENT:N5>"))
}
func finish(data []byte, footer string) []byte {
	address := len(data)
	data = block(data, []byte(footer+"<FILE_FEATURE:24>"))
	data = append(data, []byte("tail")...)
	return binary.LittleEndian.AppendUint32(data, uint32(address))
}
func damaged() []byte { return finish(begin(), "<PAGE1:999999>") }

// Each page key names a distinct length word over the same megabyte. This
// consumes work without needing a large file or exhausting the test process.
func hostile() []byte {
	data := begin()
	region := len(data)
	const length = 1048576
	const count = 20000
	for range length/4 + count + 1 {
		data = binary.LittleEndian.AppendUint32(data, length)
	}
	var footer strings.Builder
	for i := range count {
		fmt.Fprintf(&footer, "<PAGE%d:%d>", i+1, region+4*i)
	}
	return finish(data, footer.String())
}
func TestParseErrors(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		want error
	}{
		{"empty", nil, snnote.ErrUnsupportedFormat},
		{"signature", make([]byte, 64), snnote.ErrUnsupportedFormat},
		{"truncated", []byte("noteSN_FILE_VER_20230015"), snnote.ErrUnexpectedEOF},
		{"budget", hostile(), snnote.ErrLimit},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := snnote.Parse(tc.data)
			var detail *snnote.Error
			if n != nil || !errors.Is(err, tc.want) || !errors.As(err, &detail) {
				t.Fatalf("got %v, %v, want %v", n, err, tc.want)
			}
			if detail.Status == snnote.StatusLimit && (errors.Is(err, snnote.ErrMalformed) || !strings.Contains(detail.Message, "work budget")) {
				t.Fatalf("limit lost: %v", err)
			}
		})
	}
}
func TestParseConcurrently(t *testing.T) {
	paths := fixtures(t)
	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := worker; i < len(paths); i += 8 {
				data := read(t, paths[i])
				want := parse(t, data)
				got := parse(t, data)
				if !reflect.DeepEqual(want, got) {
					t.Errorf("concurrent model differs: %s", paths[i])
				}
			}
		}()
	}
	wg.Wait()
}
func TestInputAndResultOwnership(t *testing.T) {
	data := read(t, "testdata/samples/layers-dense-n6-20260016.note")
	n := parse(t, data)
	before := modelHash(n)
	clear(data)
	if after := modelHash(n); after != before {
		t.Fatal("Go model retains input")
	}
	for _, p := range n.Pages {
		for _, s := range p.DrawnStrokes {
			found := false
			for i := range p.Strokes {
				if s == &p.Strokes[i] {
					found = true
					break
				}
			}
			if !found {
				t.Fatal("drawn stroke does not reference Go stroke slice")
			}
		}
	}
}
func TestErrorsMatchByStatus(t *testing.T) {
	for status := snnote.StatusInvalidArgument; status <= snnote.StatusLimit; status++ {
		e := &snnote.Error{Status: status, Value: 17, Message: "detail"}
		if !errors.Is(e, &snnote.Error{Status: status}) || e.Error() != "snnote: "+status.String()+": detail" {
			t.Fatal(e)
		}
		var nilError *snnote.Error
		if e.Is(nilError) || nilError.Is(e) || nilError.Error() != "snnote: <nil>" || e.Is(errors.New(e.Error())) {
			t.Fatal("incorrect nil or foreign error match")
		}
	}
	if snnote.Status(999).String() != "status(999)" {
		t.Fatal("unknown status")
	}
}
