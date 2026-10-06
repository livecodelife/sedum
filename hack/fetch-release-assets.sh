#!/bin/sh
#
# Fetches the binary assets a Sedum release bundles alongside the sedum
# binary itself: goinfer-serve per platform, and the fine-tuned gguf Sedum
# defaults to when neither --model nor --local-model is given. Downloaded
# once, verified against pinned checksums, and staged where
# .goreleaser.yml's archive step expects them.
#
# Run by goreleaser's own before.hooks. Safe to run by hand too, to
# reproduce what a release build fetches, from the repo root.

set -eu

# Outside dist/: goreleaser requires dist/ empty before it runs, and this is
# populated by a before hook that runs ahead of that check.
DIST=.release-assets
mkdir -p "$DIST"

info() { printf '%s\n' "$*" >&2; }
fail() { printf 'error: %s\n' "$*" >&2; exit 1; }

sha256_of() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | cut -d' ' -f1
	else
		shasum -a 256 "$1" | cut -d' ' -f1
	fi
}

fetch_and_verify() {
	url=$1
	dest=$2
	expected=$3

	if [ -f "$dest" ] && [ "$(sha256_of "$dest")" = "$expected" ]; then
		info "already staged: $dest"
		return
	fi

	info "fetching $dest"
	curl -fsSL -o "$dest" "$url" || fail "could not download $url"

	actual=$(sha256_of "$dest")
	[ "$actual" = "$expected" ] ||
		fail "checksum mismatch for $dest: expected $expected, got $actual"
}

# --- goinfer-serve, pinned to the release this script was verified against -
#
# Pinned rather than "latest" for the reason install.sh verifies what it
# downloads: a release rebuilt next year should fetch the exact binary this
# script was written and checked against, not whatever goinfer ships by
# then.
GOINFER_VERSION=v0.19.0
GOINFER_BASE="https://github.com/townsendmerino/goinfer/releases/download/${GOINFER_VERSION}"

# One line per Sedum build target: os, arch, and goinfer-serve's own
# published sha256 for that platform's binary at $GOINFER_VERSION.
GOINFER_TARGETS='
darwin amd64 eccb4f1afc430f766a70543e1bc4baf466500cfca1195d1c60bee995fec00eb2
darwin arm64 9922d2d153b3585089bf56c9237b5a4a147d270ee9ee30d514e76a0bcc674b2c
linux  amd64 898acb6e1e8b9da9e5f6d21e86b5e17e2832195fa90c05a29bf6b9d2c4c5cabb
linux  arm64 d169fd568ac7006062ade7b4166e6662f5b8a69593c9152d521ef3d105e76785
'

echo "$GOINFER_TARGETS" | while read -r os arch sha; do
	[ -n "$os" ] || continue
	mkdir -p "$DIST/${os}_${arch}"
	fetch_and_verify \
		"${GOINFER_BASE}/goinfer-serve-${os}-${arch}" \
		"$DIST/${os}_${arch}/goinfer-serve" \
		"$sha"
	chmod +x "$DIST/${os}_${arch}/goinfer-serve"
done

# --- the bundled default model ---------------------------------------------
#
# Hosted as a release asset in this repo, not goinfer's: it is Sedum's own
# fine-tuned checkpoint, unrelated to the goinfer-serve version above.
# "model-assets" is a standing, non-version tag - never referenced by a v*
# release tag, so it never triggers .github/workflows/release.yml on its own.
MODEL_URL="https://github.com/livecodelife/sedum/releases/download/model-assets/sedum-default-model.gguf"
#
# The pin below is the sha256 of Qwen3.5-2B.Q8_0.gguf from the
# sedum-binder-final-2b-le12k-e5 checkpoint. The maintainer must UPLOAD that
# file to the model-assets release, replacing sedum-default-model.gguf, BEFORE
# tagging a release. Until then the URL serves the old model, the checksum does
# not match, and the release build fails its checksum check. That failure is
# intended: a build must never ship an unverified model. Order: upload, then tag.
MODEL_SHA=5185a0b25d5a13262fa6dfcdb1b2091b7bb64f130d65cdb5322c99c992dd052e

fetch_and_verify "$MODEL_URL" "$DIST/sedum-default-model.gguf" "$MODEL_SHA"

info "release assets staged in $DIST"
