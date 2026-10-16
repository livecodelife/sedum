#!/bin/sh
# Regenerates internal/builtin/packages/go.zip, the copy of generators/go that
# the sedum binary carries, from generators/go. The archive is never edited by
# hand. The generators directory is its own Go module, so the root module
# cannot embed it directly, and go:embed refuses some of its file names (the
# literal "*" directories under files/), so the package travels as one
# deterministic zip (prov-2026-1da47b75). Run this after changing generators/go.
set -eu

root=$(cd "$(dirname "$0")/.." && pwd)
mkdir -p "$root/internal/builtin/packages"

python3 - "$root/generators/go" "$root/internal/builtin/packages/go.zip" <<'PY'
import os, sys, zipfile

src, dst = sys.argv[1], sys.argv[2]
names = []
for base, dirs, files in os.walk(src):
    for f in files:
        full = os.path.join(base, f)
        names.append(os.path.relpath(full, src).replace(os.sep, "/"))
names.sort()

with zipfile.ZipFile(dst, "w", zipfile.ZIP_DEFLATED) as z:
    for name in names:
        info = zipfile.ZipInfo(name, date_time=(1980, 1, 1, 0, 0, 0))
        info.compress_type = zipfile.ZIP_DEFLATED
        info.external_attr = 0o644 << 16
        with open(os.path.join(src, name), "rb") as fh:
            z.writestr(info, fh.read())
print("synced", dst, "from", src, "(%d files)" % len(names))
PY
