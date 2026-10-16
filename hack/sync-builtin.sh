#!/bin/sh
# Regenerates internal/builtin/_go, the copy of generators/go that the sedum
# binary carries, from generators/go. The copy is never edited by hand: the
# generators directory is its own Go module, so the root module cannot embed it
# directly (prov-2026-1da47b75). Run this after changing generators/go.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
src="$root/generators/go"
dst="$root/internal/builtin/_go"

rm -rf "$dst"
mkdir -p "$dst"
cp -R "$src/." "$dst/"
echo "synced $dst from $src"
