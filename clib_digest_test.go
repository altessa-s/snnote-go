// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

package snnote_test

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"testing"
)

// clibDigest matches the sorted sha256-of-sha256 listing in update-clib.sh.
func clibDigest(t *testing.T, root string) string {
	t.Helper()
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			paths = append(paths, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(paths)
	digest := sha256.New()
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintf(digest, "%x ./%s\n", sha256.Sum256(data), path)
	}
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func TestCLibDigest(t *testing.T) {
	header, err := os.ReadFile("clib_digest.h")
	if err != nil {
		t.Fatal(err)
	}
	match := regexp.MustCompile(`(?m)^#define SNNOTE_GO_CLIB_DIGEST "([0-9a-f]{64})"$`).FindSubmatch(header)
	if match == nil {
		t.Fatal("missing clib digest definition")
	}
	if got := clibDigest(t, "internal/clib"); got != string(match[1]) {
		t.Fatalf("vendored tree digest %s != header %s; re-vendor with scripts/update-clib.sh", got, match[1])
	}
}

func TestCLibDigestDetectsChanges(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "source.c")
	if err := os.WriteFile(path, []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	before := clibDigest(t, root)
	if err := os.WriteFile(path, []byte("modified"), 0o600); err != nil {
		t.Fatal(err)
	}
	if clibDigest(t, root) == before {
		t.Fatal("source edit did not change digest")
	}
	before = clibDigest(t, root)
	if err := os.Rename(path, filepath.Join(root, "renamed.c")); err != nil {
		t.Fatal(err)
	}
	if clibDigest(t, root) == before {
		t.Fatal("source rename did not change digest")
	}
}
