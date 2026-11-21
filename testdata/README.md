# Parser fixtures

The [five device samples](samples/README.md) exercise the N6 notebook and annotation formats. The [27 third-party fixtures](supernote-typescript/README.md) come from supernote-typescript under Apache-2.0.

`dumps.sha256` records the filtered canonical model dump for each fixture. The ten small fuzz seeds in `fuzz/FuzzParse/` come from `seeds_test.go`. Regenerate them with `SNNOTE_WRITE_SEEDS=1 go test -tags snnote_test -run '^TestSmallSeeds$' .`.
