// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/altessa-s/snnote-go"
)

// seedNotebook builds complete small containers with reachable footer offsets.
func seedNotebook(kind int) []byte {
	data := begin()
	add := func(payload []byte) int { at := len(data); data = block(data, payload); return at }
	tags := func(text string) int { return add([]byte(text)) }
	encoded := func(text string) int { return add([]byte(base64.StdEncoding.EncodeToString([]byte(text)))) }
	footer, pageText := "", "<PAGEID:p>"
	switch kind {
	case 0, 1:
		// Two records deliberately share an ID, making drawn-pointer identity observable.
		record := make([]byte, 208)
		binary.LittleEndian.PutUint32(record, 16)
		binary.LittleEndian.PutUint32(record[40:], 5000)
		copy(record[48:], "others")
		for _, v := range []uint32{0, 2, 100, 100, 200, 200} {
			record = binary.LittleEndian.AppendUint32(record, v)
		}
		record = append(record, make([]byte, 20)...)
		for _, v := range []uint32{0, 0, 1} {
			record = binary.LittleEndian.AppendUint32(record, v)
		}
		record = append(record, make([]byte, 40)...)
		record = binary.LittleEndian.AppendUint32(record, 0)
		record = append(record, make([]byte, 10)...)
		record = binary.LittleEndian.AppendUint32(record, 0)
		records := block(block(nil, record), record)
		var log []byte
		if kind == 0 {
			log = append(binary.LittleEndian.AppendUint32(nil, 2), records...)
		} else {
			// Literal-only raw LZ4 needs no compressor dependency.
			packed := []byte{0xf0}
			length := len(records) - 15
			for length >= 255 {
				packed = append(packed, 255)
				length -= 255
			}
			packed = append(packed, byte(length))
			packed = append(packed, records...)
			for _, v := range []uint32{^uint32(0), 2, uint32(len(records)), uint32(len(packed))} {
				log = binary.LittleEndian.AppendUint32(log, v)
			}
			log = append(log, packed...)
		}
		pageText += fmt.Sprintf("<TOTALPATH:%d>", add(log))
	case 2, 3, 4:
		texts := []string{"<TITLERECT:1,2,3,4><TITLELEVEL:2>", "<KEYWORD: keyword ><KEYWORDRECT:1,2,3,4>", "<LINKTYPE:4><LINKFILE:aHR0cHM6Ly9leGFtcGxlLmNvbQ==><SHOWTEXT:77u/bGluaw==>"}
		prefixes := []string{"TITLE_", "KEYWORD_", "LINK_"}
		footer += fmt.Sprintf("<%s10001:%d>", prefixes[kind-2], tags(texts[kind-2]))
	case 5:
		copy(data, "mark")
		at := encoded(`{"0":[{"annotationContent":"mark","time":1000,"rnRectList":[{"left":1,"top":2,"right":3,"bottom":4}]}]}`)
		header := tags(fmt.Sprintf("<APPLY_EQUIPMENT:N5><HIGHLIGHTINFO:%d>", at))
		footer += fmt.Sprintf("<FILE_FEATURE:%d>", header)
	case 6:
		at := encoded(`{"elements":[{"type":"Text","label":"text","words":[{"label":"text","bounding-box":{"x":1,"y":2,"width":3,"height":4}}]}]}`)
		pageText += fmt.Sprintf("<RECOGNTEXT:%d>", at)
	case 7:
		// Full 16 Ki-pixel blocks and a short tail cover a complete N5 sheet.
		var rle []byte
		for left := 1404 * 1872; left > 0; {
			count := min(left, 128)
			length := byte(count - 1)
			if left >= 16384 {
				count, length = 16384, 255
			}
			rle = append(rle, 0x61, length)
			left -= count
		}
		bitmap := add(rle)
		layer := tags(fmt.Sprintf("<LAYERBITMAP:%d><LAYERPROTOCOL:RATTA_RLE>", bitmap))
		pageText += fmt.Sprintf("<MAINLAYER:%d>", layer)
	case 8:
		// First insertion order differs from lexical order.
		z := add([]byte{0x61, 0})
		a := add([]byte{0x62, 0})
		footer += fmt.Sprintf("<STYLE_z:%d><STYLE_a:%d>", z, a)
	case 9:
		return finish(data, "")
	}
	footer += fmt.Sprintf("<PAGE1:%d>", tags(pageText))
	return finish(data, footer)
}

func smallSeeds() map[string][]byte {
	seeds := make(map[string][]byte)
	for i := range 10 {
		seeds[fmt.Sprintf("small-%02d", i)] = seedNotebook(i)
	}
	return seeds
}

// Set SNNOTE_WRITE_SEEDS=1 to regenerate the checked-in Go fuzz corpus.
func TestSmallSeeds(t *testing.T) {
	for name, data := range smallSeeds() {
		t.Run(name, func(t *testing.T) {
			n := parse(t, data)
			index, err := strconv.Atoi(name[len("small-"):])
			if err != nil {
				t.Fatal(err)
			}
			if index == 9 {
				if len(n.Pages) != 0 {
					t.Fatal("expected empty notebook")
				}
			} else if len(n.Pages) != 1 {
				t.Fatalf("pages: %d", len(n.Pages))
			}
			switch index {
			case 0, 1:
				if len(n.Pages[0].Strokes) != 2 || len(n.Pages[0].DrawnStrokes) != 2 {
					t.Fatalf("stroke seed: %+v", n.Pages[0])
				}
			case 2:
				if len(n.Titles) != 1 {
					t.Fatal("title missing")
				}
			case 3:
				if len(n.Keywords) != 1 {
					t.Fatal("keyword missing")
				}
			case 4:
				if len(n.Links) != 1 || n.Links[0].URL == nil {
					t.Fatal("link URL missing")
				}
			case 5:
				if n.Kind != snnote.KindMark || len(n.Highlights) != 1 {
					t.Fatal("highlight missing")
				}
			case 6:
				if n.Pages[0].Recognition == nil {
					t.Fatal("recognition missing")
				}
			case 7:
				if len(n.Pages[0].Layers) != 1 || n.Pages[0].Layers[0].Bitmap == nil {
					t.Fatal("bitmap missing")
				}
				b := n.Pages[0].Layers[0].Bitmap
				if b.Encoding != snnote.EncodingRattaRLE {
					t.Fatal("RLE missing")
				}
				if _, err := b.Codes(1404, 1872); err != nil {
					t.Fatal(err)
				}
			case 8:
				if len(n.Templates) != 2 || n.Templates[0].Name != "z" || n.Templates[1].Name != "a" {
					t.Fatalf("templates: %+v", n.Templates)
				}
			}
			content := []byte("go test fuzz v1\n[]byte(" + strconv.Quote(string(data)) + ")\n")
			path := filepath.Join("testdata/fuzz/FuzzParse", name)
			if os.Getenv("SNNOTE_WRITE_SEEDS") != "" {
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, content, 0o644); err != nil {
					t.Fatal(err)
				}
			} else if string(read(t, path)) != string(content) {
				t.Fatal("stored seed differs from builder")
			}
		})
	}
}
