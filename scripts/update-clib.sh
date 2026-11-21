#!/usr/bin/env bash
# Copyright 2026 ALTESSA SOLUTIONS INC.
# Use of this source code is governed by an MIT-style license that can be
# found in the LICENSE file.

# Re-vendor snnote-lib into internal/clib: the library sources and headers,
# its licence, tag and commit. Refresh the shared fixtures and their filtered
# reference hashes from the same checkout.
#
#   scripts/update-clib.sh ../snnote-lib
#
# The new tree is assembled and checked beside the old one and swapped in
# whole: on any failure the previous tree and digest stay as they were.
# clib_digest.h in the package root is rewritten with a digest of the vendored
# files, because Go's build cache tracks the package directory and not
# internal/clib; the digest is what makes a change down there rebuild the
# package.

set -euo pipefail

allow_dirty=false
if [[ ${1:-} == --allow-dirty ]]
then
    allow_dirty=true
    shift
fi
if [[ $# -ne 1 ]]
then
    echo "usage: $0 [--allow-dirty] <path-to-snnote-lib>" >&2
    exit 2
fi
if ! source_dir=$(cd "$1" && pwd -P)
then
    echo "library directory does not exist or is inaccessible: $1" >&2
    exit 1
fi
if ! module_dir=$(cd "$(dirname "$0")/.." && pwd -P)
then
    echo "cannot resolve module directory" >&2
    exit 1
fi
if [[ ! -f "$source_dir/include/snnote/snnote.h" ]]
then
    echo "not a snnote-lib directory: $source_dir (missing include/snnote/snnote.h)" >&2
    exit 1
fi
internal="$module_dir/internal"
clib="$internal/clib"

# Never copy from a tree that would be moved away during replacement.
case "$source_dir/" in "$clib/"*)
    echo "refusing source inside vendored destination: $source_dir" >&2
    exit 1
    ;;
esac
case "$module_dir/" in "$source_dir/"*)
    echo "refusing source containing module: $source_dir" >&2
    exit 1
    ;;
esac
if ! state=$(GIT_OPTIONAL_LOCKS=0 git --no-optional-locks -C "$source_dir" status --porcelain)
then
    echo "not a readable library checkout: $source_dir" >&2
    exit 1
fi
if [[ -n "$state" && $allow_dirty == false ]]
then
    echo "refusing dirty library checkout (use --allow-dirty explicitly): $source_dir" >&2
    exit 1
fi
version=$(GIT_OPTIONAL_LOCKS=0 git --no-optional-locks -C "$source_dir" describe --tags --always)
commit=$(git -C "$source_dir" rev-parse HEAD)
if [[ -n "$state" ]]
then
    version="$version-dirty"
fi

if command -v sha256sum > /dev/null
then
    sha256() { sha256sum | cut -d' ' -f1
    }
elif command -v shasum > /dev/null
then
    sha256() { shasum -a 256 | cut -d' ' -f1
    }
else
    echo "neither sha256sum nor shasum is installed" >&2
    exit 1
fi

sources=(sn_annotations sn_base64 sn_bytes sn_decimal sn_json sn_log sn_lz4 sn_metadata sn_model sn_parser sn_raster sn_rle sn_strokes sn_support)
headers=(export log parser raster snnote types)

mkdir -p "$internal" "$module_dir/build"
staging=$(mktemp -d "$module_dir/build/clib-staging.XXXXXX")
outputs=(internal/clib clib_digest.h testdata/samples testdata/supernote-typescript testdata/dumps.sha256)
installed=()
finished=false

cleanup() {
    local status=$?
    if [[ $finished == false ]]
    then
        for path in "${installed[@]+"${installed[@]}"}"
        do
            rm -rf "$module_dir/$path"
        done
        for path in "${outputs[@]}"
        do
            if [[ -e "$staging/old/$path" ]]
            then
                mv "$staging/old/$path" "$module_dir/$path"
            fi
        done
    fi
    rm -rf "$staging"
    exit "$status"
}
trap cleanup EXIT

tree="$staging/new/internal/clib"
mkdir -p "$tree/src" "$tree/include/snnote"
for name in "${sources[@]}"
do
    cp "$source_dir/src/$name.c" "$tree/src/"
done
cp "$source_dir/src/sn_internal.h" "$tree/src/"
for name in "${headers[@]}"
do
    cp "$source_dir/include/snnote/$name.h" "$tree/include/snnote/"
done
cp "$source_dir/LICENSE" "$source_dir/COPYRIGHT" "$tree/"
printf '%s\n%s\n' "$version" "$commit" > "$tree/VERSION"

# Every file this module compiles or ships must be there and not empty; a
# source the library gains or loses means the wrappers change first.
expected=(VERSION LICENSE COPYRIGHT src/sn_internal.h)
for name in "${sources[@]}"
do
    expected+=("src/$name.c")
done
for name in "${headers[@]}"
do
    expected+=("include/snnote/$name.h")
done
for file in "${expected[@]}"
do
    if [[ ! -s "$tree/$file" ]]
    then
        echo "missing or empty: $file" >&2
        exit 1
    fi
done
for file in "$source_dir"/src/*.c
do
    name=$(basename "$file" .c)
    if [[ ! -f "$tree/src/$name.c" ]]
    then
        echo "$name.c is not listed in sources" >&2
        exit 1
    fi
    if [[ ! -f "$module_dir/clib_${name#sn_}.c" ]]
    then
        echo "$name.c has no wrapper; add clib_${name#sn_}.c and list it here" >&2
        exit 1
    fi
done

for file in "$source_dir"/include/snnote/*.h
do
    name=$(basename "$file")
    if [[ ! -f "$tree/include/snnote/$name" ]]
    then
        echo "$name is not listed in headers" >&2
        exit 1
    fi
done

digest=$(
cd "$tree"
find . -type f | LC_ALL=C sort | while read -r file
do
    printf '%s %s\n' "$(sha256 < "$file")" "$file"
done | sha256)

cat > "$staging/new/clib_digest.h" <<EOF
// Copyright 2026 ALTESSA SOLUTIONS INC.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file.

// Generated by scripts/update-clib.sh. Do not edit.
//
// The vendored snnote-lib, named by its tag and by a digest of every file
// under internal/clib. Each clib_*.c includes this file, so a change to the
// vendored sources changes a file in the package directory, which is what
// Go's build cache tracks.

#ifndef SNNOTE_GO_CLIB_DIGEST_H
#define SNNOTE_GO_CLIB_DIGEST_H

#define SNNOTE_GO_CLIB_VERSION "$version"
#define SNNOTE_GO_CLIB_DIGEST "$digest"

#endif
EOF

fixtures="$staging/new/testdata"
mkdir -p "$fixtures"
cp -R "$source_dir/tests/fixtures/samples" "$source_dir/tests/fixtures/supernote-typescript" "$fixtures/"
# The samples are one directory closer to the module licence here.
sed 's|(../../../LICENSE)|(../../LICENSE)|g' "$source_dir/tests/fixtures/samples/README.md" > "$fixtures/samples/README.md"

"${CC:-cc}" -std=c2x -O2 -I"$source_dir/include" -I"$source_dir/src" -I"$source_dir/tests" \
    "$source_dir/tools/sndump.c" "$source_dir/tests/sn_composition.c" "$source_dir"/src/*.c -lm -o "$staging/sndump"

(
    cd "$staging/new"
    export LC_ALL=C
    for path in testdata/{samples,supernote-typescript}/*
    do
        case "$path" in
            *.note|*.mark)
                hash=$("$staging/sndump" "$path" | awk -F '\t' '$1 != "raster_status" && $1 != "gray" && $1 != "alpha"' | sha256)
                printf '%s  %s\n' "$hash" "$path"
                ;;
        esac
    done
) > "$fixtures/dumps.sha256"
test -s "$fixtures/dumps.sha256"

mkdir -p "$staging/old/internal" "$staging/old/testdata"
for path in "${outputs[@]}"
do
    if [[ -e "$module_dir/$path" ]]
    then
        mv "$module_dir/$path" "$staging/old/$path"
    fi
    installed+=("$path")
    mv "$staging/new/$path" "$module_dir/$path"
done
finished=true

echo "vendored snnote-lib $version at $commit ($digest)"
