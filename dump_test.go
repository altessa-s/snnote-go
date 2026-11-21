// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/altessa-s/snnote-go"
)

// The canonical text keeps exactly the order and precision of sndump. Only
// composition results are excluded: raster_status, gray and alpha lines.
type dumper struct{ w io.Writer }

func (d dumper) line(k string, v any)    { fmt.Fprintf(d.w, "%s\t%v\n", k, v) }
func (d dumper) num(k string, v float64) { d.line(k, fmt.Sprintf("%.17e", v)) }
func (d dumper) text(k string, v string) { d.str(k, &v) }
func (d dumper) str(k string, v *string) {
	if v == nil {
		d.line(k, "null")
		return
	}
	var b strings.Builder
	b.WriteByte('"')
	for _, ch := range []byte(*v) {
		switch {
		case ch == '"' || ch == '\\':
			b.WriteByte('\\')
			b.WriteByte(ch)
		case ch < 32:
			fmt.Fprintf(&b, "\\u%04x", ch)
		default:
			b.WriteByte(ch)
		}
	}
	b.WriteByte('"')
	d.line(k, b.String())
}
func (d dumper) bytes(k string, v []byte) {
	if v == nil {
		d.line(k, "null")
	} else {
		d.line(k, "\""+hex.EncodeToString(v)+"\"")
	}
}
func (d dumper) rect(v *snnote.Rect) {
	d.line("null", v == nil)
	if v == nil {
		v = &snnote.Rect{}
	}
	d.num("x", v.X)
	d.num("y", v.Y)
	d.num("width", v.Width)
	d.num("height", v.Height)
}
func (d dumper) point(v snnote.Point) { d.num("x", v.X); d.num("y", v.Y) }
func (d dumper) metadata(v snnote.Metadata) {
	d.line("entries", len(v.Entries))
	for _, e := range v.Entries {
		d.text("key", e.Key)
		d.text("value", e.Value)
	}
}
func (d dumper) bitmap(v snnote.Bitmap) {
	d.line("encoding", int(v.Encoding))
	d.str("name", v.Name)
	d.bytes("bytes", v.Bytes)
}
func (d dumper) optionalBitmap(k string, v *snnote.Bitmap) {
	d.line(k, v != nil)
	if v != nil {
		d.bitmap(*v)
	}
}
func (d dumper) layerInfo(v snnote.LayerInfo) {
	d.line("id", v.LayerID)
	d.text("name", v.Name)
	d.line("background", v.IsBackground)
	d.line("visible", v.IsVisible)
	d.line("deleted", v.IsDeleted)
	d.line("current", v.IsCurrent)
}
func (d dumper) layer(v snnote.Layer) {
	d.line("name", int(v.Name))
	d.optionalBitmap("bitmap", v.Bitmap)
	d.line("info", v.Info != nil)
	if v.Info != nil {
		d.layerInfo(*v.Info)
	}
	d.metadata(v.Metadata)
	d.line("stacked", v.IsStacked)
	d.line("visible", v.IsVisible)
}
func value[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}
func (d dumper) stroke(v snnote.Stroke) {
	d.line("raw_role", v.RoleRaw)
	d.line("id", v.ID)
	d.line("pen", int32(v.Pen))
	d.line("color", uint32(v.Color))
	d.line("thickness", v.Thickness)
	d.text("kind", v.KindName)
	d.line("role", int(v.Role))
	d.line("status", v.StatusRaw)
	d.line("operation", v.Operation)
	d.line("layer", v.Layer)
	d.line("page", v.PageNumber)
	d.line("drawn", v.IsDrawn)
	d.line("has_style", v.Style != nil)
	d.line("category", int(v.Category))
	d.line("style", value(v.Style))
	d.line("points", len(v.Points))
	for _, p := range v.Points {
		d.point(p.Position)
		d.line("pressure", p.Pressure)
	}
	d.line("contours", len(v.Contours))
	for _, c := range v.Contours {
		d.line("points", len(c.Points))
		for _, p := range c.Points {
			d.point(p)
		}
	}
	d.rect(v.Bounds)
	d.line("attributes", len(v.Attributes))
	for _, a := range v.Attributes {
		d.text("item", a)
	}
	d.line("pressures", len(v.Pressures))
	for _, p := range v.Pressures {
		d.line("pressure", p)
	}
}
func (d dumper) recognition(v snnote.Recognition) {
	d.text("text", v.Text)
	d.line("elements", len(v.Elements))
	for _, e := range v.Elements {
		d.text("kind", e.Kind)
		d.str("label", e.Label)
		d.line("words", len(e.Words))
		for _, w := range e.Words {
			d.text("label", w.Label)
			d.rect(w.Box)
		}
	}
}
func (d dumper) textbox(v snnote.TextBox) {
	d.line("id", v.ID)
	d.rect(v.Rect)
	d.text("text", v.Text)
	d.line("has_font", v.FontSize != nil)
	d.num("font", value(v.FontSize))
	d.str("font_path", v.FontPath)
	d.str("timestamp", v.Timestamp)
	d.line("attributes", len(v.Attributes))
	for _, a := range v.Attributes {
		d.text("item", a)
	}
}
func (d dumper) page(v snnote.Page, palette snnote.Palette) {
	d.line("raw_orientation", v.OrientationRaw)
	d.line("number", v.Number)
	d.str("id", v.ID)
	d.str("template", v.TemplateName)
	d.num("width", v.Size.Width)
	d.num("height", v.Size.Height)
	d.line("layers", len(v.Layers))
	for _, l := range v.Layers {
		d.layer(l)
	}
	d.line("layer_info", len(v.LayerInfo))
	for _, l := range v.LayerInfo {
		d.layerInfo(l)
	}
	d.line("strokes", len(v.Strokes))
	for _, s := range v.Strokes {
		d.stroke(s)
	}
	d.line("skipped_strokes", v.SkippedStrokes)
	d.line("recognition", v.Recognition != nil)
	if v.Recognition != nil {
		d.recognition(*v.Recognition)
	}
	d.bytes("archive", v.RecognitionArchive)
	d.line("disabled", len(v.DisabledAreas))
	for _, r := range v.DisabledAreas {
		d.rect(r)
	}
	d.metadata(v.Metadata)
	d.line("text_boxes", len(v.TextBoxes))
	for _, b := range v.TextBoxes {
		d.textbox(b)
	}
	d.line("drawn_strokes", len(v.DrawnStrokes))
	for _, s := range v.DrawnStrokes {
		d.line("drawn_id", s.ID)
	}
	for _, ink := range palette.Inks {
		d.line("palette_gray", ink.Gray)
		d.line("palette_alpha", ink.Alpha)
	}
	d.line("orientation", int(v.Orientation))
}
func (d dumper) notebook(v *snnote.Notebook) {
	d.line("kind", int(v.Kind))
	d.line("version", v.FormatVersion)
	d.text("equipment", v.EquipmentRaw)
	d.str("file_id", v.FileID)
	d.num("width", v.Sheet.Width)
	d.num("height", v.Sheet.Height)
	d.line("skipped_pages", v.SkippedPages)
	d.metadata(v.Header)
	d.metadata(v.Footer)
	d.line("pages", len(v.Pages))
	for _, p := range v.Pages {
		d.page(p, v.Palette)
	}
	d.line("titles", len(v.Titles))
	for _, t := range v.Titles {
		d.line("page", t.PageNumber)
		d.line("level", t.Level)
		d.line("style", int64(t.Style))
		d.rect(t.Rect)
		d.rect(t.OriginalRect)
		d.optionalBitmap("bitmap", t.Bitmap)
		d.metadata(t.Metadata)
	}
	d.line("keywords", len(v.Keywords))
	for _, k := range v.Keywords {
		d.line("page", k.PageNumber)
		d.text("text", k.Text)
		d.rect(k.Rect)
		d.rect(k.OriginalRect)
		d.metadata(k.Metadata)
	}
	d.line("links", len(v.Links))
	for _, l := range v.Links {
		d.line("page", l.PageNumber)
		d.line("kind", l.Kind)
		d.line("direction", l.Direction)
		d.line("style", int64(l.Style))
		d.rect(l.Rect)
		d.str("path", l.TargetPath)
		d.str("file_id", l.TargetFileID)
		d.str("page_id", l.TargetPageID)
		d.line("has_page", l.TargetPageNumber != nil)
		d.line("target_page", value(l.TargetPageNumber))
		d.str("timestamp", l.Timestamp)
		d.str("text", l.Text)
		d.optionalBitmap("bitmap", l.Bitmap)
		d.metadata(l.Metadata)
	}
	d.line("highlights", len(v.Highlights))
	for _, h := range v.Highlights {
		d.line("page", h.PageNumber)
		d.line("kind", h.Kind)
		d.line("color", h.ColorType)
		d.str("comment", h.Comment)
		d.line("has_date", h.Date != nil)
		d.num("date", value(h.Date))
		d.line("rects", len(h.Rects))
		for _, r := range h.Rects {
			d.rect(r)
		}
		d.line("document_rects", len(h.DocumentRects))
		for _, r := range h.DocumentRects {
			d.rect(r)
		}
	}
	d.optionalBitmap("cover", v.Cover)
	d.line("cover_type", v.CoverType)
	order := append([]snnote.Template(nil), v.Templates...)
	sort.Slice(order, func(i, j int) bool { return order[i].Name < order[j].Name })
	d.line("templates", len(order))
	for _, t := range order {
		d.text("name", t.Name)
		d.bitmap(t.Bitmap)
	}
}
func modelHash(v *snnote.Notebook) string {
	h := sha256.New()
	dumper{h}.notebook(v)
	return hex.EncodeToString(h.Sum(nil))
}

func fixtures(t testing.TB) []string {
	t.Helper()
	var paths []string
	err := filepath.WalkDir("testdata", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && (strings.HasSuffix(path, ".note") || strings.HasSuffix(path, ".mark")) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 32 {
		t.Fatalf("fixtures = %d, want 32", len(paths))
	}
	return paths
}
func read(t testing.TB, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func parse(t testing.TB, b []byte) *snnote.Notebook {
	t.Helper()
	n, err := snnote.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// Reference hashes travel with the fixtures so CI needs no sibling checkout.
// SNNOTE_REFERENCE asks this test to rebuild sndump from that read-only checkout
// and compare its actual filtered output as well as the checked-in hashes.
func TestCanonicalDumps(t *testing.T) {
	expected := map[string]string{}
	for _, line := range strings.Split(string(read(t, "testdata/dumps.sha256")), "\n") {
		if fields := strings.Fields(line); len(fields) == 2 {
			expected[fields[1]] = fields[0]
		}
	}
	source := os.Getenv("SNNOTE_REFERENCE")
	tool := ""
	if source != "" {
		if err := os.MkdirAll("build/reference", 0o755); err != nil {
			t.Fatal(err)
		}
		tool = "build/reference/sndump"
		cc := os.Getenv("CC")
		if cc == "" {
			cc = "cc"
		}
		args := []string{"-std=c2x", "-O2", "-I" + filepath.Join(source, "include"), "-I" + filepath.Join(source, "src"), "-I" + filepath.Join(source, "tests"), filepath.Join(source, "tools/sndump.c"), filepath.Join(source, "tests/sn_composition.c")}
		sources, err := filepath.Glob(filepath.Join(source, "src/*.c"))
		if err != nil {
			t.Fatal(err)
		}
		args = append(args, sources...)
		args = append(args, "-lm", "-o", tool)
		if output, err := exec.Command(cc, args...).CombinedOutput(); err != nil {
			t.Fatalf("build sndump: %v\n%s", err, output)
		}
	}
	for _, path := range fixtures(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			got := modelHash(parse(t, read(t, path)))
			if want := expected[path]; got != want {
				t.Errorf("Go dump SHA-256 %s, expected %s", got, want)
			}
			if tool != "" {
				cmd := exec.Command(tool, path)
				pipe, err := cmd.StdoutPipe()
				if err != nil {
					t.Fatal(err)
				}
				if err = cmd.Start(); err != nil {
					t.Fatal(err)
				}
				h := sha256.New()
				err = filterDump(h, pipe)
				waitErr := cmd.Wait()
				if err != nil || waitErr != nil {
					t.Fatalf("sndump: %v, %v", err, waitErr)
				}
				want := hex.EncodeToString(h.Sum(nil))
				if want != got {
					t.Errorf("C dump SHA-256 %s, Go %s", want, got)
				}
			}
			t.Logf("%s %s", got, path)
		})
	}
}
func filterDump(h hash.Hash, r io.Reader) error {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 && !strings.HasPrefix(line, "raster_status\t") && !strings.HasPrefix(line, "gray\t") && !strings.HasPrefix(line, "alpha\t") {
			if _, writeErr := io.WriteString(h, line); writeErr != nil {
				return writeErr
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
