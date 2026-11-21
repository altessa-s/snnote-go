# snnote-go

Read Supernote `.note` notebooks and `.mark` annotations as Go values. The C23
[snnote-lib](https://github.com/altessa-s/snnote-lib) sources are vendored here;
there is no system library to install and no Go dependency beyond the standard library.

## Requirements

Go 1.24 or newer with cgo enabled, and a C23 compiler: clang 19 or newer,
GCC 13 or newer, or a current Xcode clang. The compiler must support `constexpr`
and typed enums. The build uses `-std=c2x`; Linux also links `libm`.

```sh
CC=clang-19 go build ./...
```

## Installing

```sh
go get github.com/altessa-s/snnote-go@latest
```

Import `github.com/altessa-s/snnote-go`; its package name is `snnote`.

## Reading

```go
data, err := os.ReadFile("testdata/supernote-typescript/blank-a6x-3.15.27-two-pages.note")
if err != nil {
	log.Fatal(err)
}
notebook, err := snnote.Parse(data)
if err != nil {
	log.Fatal(err)
}
fmt.Println(len(notebook.Pages), "pages")
```

`Parse` copies every public field of the C notebook, then frees it. There is
nothing to close. Pages include layers, layer-table rows, every readable stroke,
drawn strokes, text boxes, recognition, the raw recognition archive, disabled
areas and metadata. The notebook also carries titles, keywords, links,
highlights, templates and cover data. Unknown raw values remain available.

Optional strings, numbers, rectangles and nested records are pointers; nil means
absent. Required text and metadata keys and values use plain strings. All strings
preserve embedded NULs. Byte payloads are copied slices. `Metadata.Entries` preserves
order and repeated keys; `Get` returns the first value and `Values` returns all.
`Page.DrawnStrokes` points into the original Go `Strokes` backing array;
appending or re-slicing does not retarget those pointers. Optional content that could
not be read is reported through skipped counts and diagnostics.

## Bitmaps

```go
bitmap := snnote.Bitmap{Encoding: snnote.EncodingRattaRLE, Bytes: []byte{0x61, 2, 0x62, 0}}
raster, err := bitmap.Raster(4, 1, snnote.CurrentPalette())
if err != nil {
	log.Fatal(err)
}
fmt.Printf("%d %d %v\n", raster.Width, raster.Height, raster.Alpha)
```

`Bitmap.Codes` returns RATTA_RLE codes. `Bitmap.Raster` maps them through a
`Palette` into independent row-major `Gray` and `Alpha` planes. Dimensions are
positive pixel counts; the pixel ceiling is `MaxRasterPixels` (16,777,216).
Use `LegacyPalette`, `CurrentPalette`, or `PaletteForFormat`; parsed notebooks
already carry their palette. PNG and JPEG return `ErrImageEncoding`; unknown
protocols return `ErrUnsupportedEncoding`. Composition and image decoding belong
to a separate rendering library. This module does neither.

## Diagnostics

```go
snnote.SetLogger(func(level snnote.LogLevel, message string) {
	log.Printf("snnote %v: %s", level, message)
})
defer snnote.SetLogger(nil)
// Diagnostics from subsequent calls go to this handler.
```

Logging is silent by default. `SetLogger(nil)` silences it again. Handlers may
run on separate goroutines, synchronously with the parse that produced each
line. Concurrent parses may run handlers concurrently, so handlers must
synchronize shared state. A handler that panics or exits loses only that line.
The parse waits for its handler to return, so a handler must not wait for that
parse.
A line already in flight may still be delivered after `SetLogger(nil)` returns. The C sink is installed
exactly once because each installation retains a binding for process lifetime.
Subsequent changes replace an atomic Go pointer. Messages contain format facts,
never notebook text, paths or names.

## Limits

C parse allocation storage and the Go model copy each have a separate budget of
1 MiB plus 128 times the input length. The Go budget counts copied strings, byte
payloads, slice backing arrays and pointed-to values, including every copy of a
shared C block. It excludes allocator rounding and garbage collector overhead;
it is not a process memory limit. Exceeding either budget returns `ErrLimit`
without a partial notebook.
Parse work is separately limited to 1 MiB plus 256 units per input byte.
The corresponding constants are `ParseBudgetFloor`, `ParseBudgetMultiple`,
`ParseWorkFloor` and `ParseWorkMultiple`. JSON nesting is limited to
`JSONMaxDepth` (64); protocol names longer than `ProtocolMaxLength` (32) stay
unknown. The constants come from the C library; the Go copy enforces the same
allocation ceiling independently.

`ErrLimit` means a resource ceiling refused the operation, allowing a caller to
distinguish a hostile input from `ErrMalformed`. A count beyond Go's `int` range
or a failed C allocation returns `ErrOutOfMemory`. Every failure is an
`*Error{Status, Value, Message}`, matched through `errors.Is` by status.
Other sentinels are `ErrInvalidArgument`, `ErrUnsupportedFormat`,
`ErrUnexpectedEOF`, `ErrImageEncoding` and `ErrUnsupportedEncoding`.

```go
// A tiny encoded bitmap can request too many output pixels or fail to fill them.
bitmap := snnote.Bitmap{Encoding: snnote.EncodingRattaRLE, Bytes: []byte{0x61, 0}}
for _, width := range []int{snnote.MaxRasterPixels + 1, 2} {
	_, err := bitmap.Codes(width, 1)
	switch {
	case errors.Is(err, snnote.ErrLimit):
		fmt.Println("resource budget refused the bitmap")
	case errors.Is(err, snnote.ErrMalformed):
		fmt.Println(err)
	}
}
```

## Design

Values own their data. No finalizers, handles or non-cgo implementation are
needed. Calls keep no state between calls and may run concurrently. Each C call
stays on its OS thread until it returns, keeping thread-local work accounting
and internal recovery within one call. Returned data can be shared for reading;
callers synchronize mutations.

Derived pen, color, title, equipment and palette accessors call the C library.
Page-layer and metadata lookup operate on Go values. `Version` is the C header's
release version; `internal/clib/VERSION` records the upstream tag and full commit.

## Layout

| Path | Contents |
|---|---|
| `model.go`, `parse.go`, `convert.go` | Values and complete model copying |
| `enums.go`, `accessors.go`, `bitmap.go` | Codes, derived values and bitmap decoding |
| `errors.go`, `log.go`, `log_bridge.c` | Status errors and diagnostic bridge |
| `clib_*.c`, `clib_digest.h` | cgo compilation units and cache invalidation digest |
| `internal/clib/` | Vendored sources, public headers, license and revision |
| `scripts/update-clib.sh` | Validated replacement of the vendored tree |
| `testdata/` | All 32 fixtures and filtered canonical dump hashes |

```sh
scripts/update-clib.sh ../snnote-lib
go test -tags snnote_test -race ./...
SNNOTE_REFERENCE=../snnote-lib go test -tags snnote_test -run TestCanonicalDumps -v .
```

Contributor checks use `-tags snnote_test` to include C comparison and audit
helpers, which are excluded from consumer builds. Run `go vet -tags snnote_test ./...` and `golangci-lint run --build-tags snnote_test` as well. Vendoring rejects
dirty checkouts unless `--allow-dirty` is explicitly supplied.
`TestCLibDigest` recomputes the vendored tree digest during ordinary tests and
rejects a stale `clib_digest.h`; re-vendor to update the sources and digest together.

The optional reference test builds the library's own `sndump` under `build/`,
reading the sibling checkout without changing it. The Go dump is made entirely
from Go values. Only per-page `raster_status`, `gray` and `alpha` lines are
excluded, because they use a composition helper outside the library. CI compares
against the checked-in hashes generated by the vendoring script.

## See also

- [snnote-lib](https://github.com/altessa-s/snnote-lib) — the C23 notebook library.
- [rmlines-go](https://github.com/altessa-s/rmlines-go) — reMarkable page data for Go.
- [rmrender-go](https://github.com/altessa-s/rmrender-go) — reMarkable rendering.

## Acknowledgements

The bindings use 27 notebooks and annotations from Brandon Philips’s supernote-typescript project under Apache-2.0. The [fixture inventory](testdata/supernote-typescript/README.md) lists their source revision, earlier filenames and checksums.

## License

MIT — see [LICENSE](LICENSE) and [COPYRIGHT](COPYRIGHT). Vendored C sources retain
their own license under `internal/clib/`. The third-party corpus license is kept
beside its fixtures at `testdata/supernote-typescript/LICENSE`.

## Trademark

Supernote is a trademark of Ratta Software. This project is not affiliated with, endorsed by or sponsored by Ratta Software.
