#!/bin/sh
#
# Sedum installer.
#
#   curl -fsSL https://raw.githubusercontent.com/livecodelife/sedum/main/install.sh | sh
#
# Resolves your OS and architecture, downloads the matching release archive,
# verifies it against the release's published checksums.txt, and installs the
# binary. The verification step is not optional: an installer that is run by
# piping a URL into a shell and then skips the checksum it publishes is offering
# a guarantee it does not keep (prov-2026-4746d9ed).
#
# Environment:
#   SEDUM_VERSION       version to install, with or without a leading "v".
#                       Default: the latest published release.
#   SEDUM_INSTALL_DIR   directory to install into. Default: $HOME/.local/bin.
#
# Exits nonzero on any failure, leaving nothing behind.

set -eu

REPO="livecodelife/sedum"
BIN="sedum"

info()  { printf '%s\n' "$*" >&2; }
warn()  { printf 'warning: %s\n' "$*" >&2; }
fail()  { printf 'error: %s\n' "$*" >&2; exit 1; }

need() {
	command -v "$1" >/dev/null 2>&1 || fail "$1 is required but was not found on PATH"
}

# --- platform -----------------------------------------------------------------
# These must agree with .goreleaser.yml's archive name_template. If that template
# changes, every installer already in someone's shell history breaks.

detect_os() {
	case "$(uname -s)" in
		Linux)  echo linux ;;
		Darwin) echo darwin ;;
		*)      fail "unsupported operating system: $(uname -s). Build from source: go install github.com/${REPO}/cmd/${BIN}@latest" ;;
	esac
}

detect_arch() {
	case "$(uname -m)" in
		x86_64|amd64)  echo amd64 ;;
		arm64|aarch64) echo arm64 ;;
		*)             fail "unsupported architecture: $(uname -m). Build from source: go install github.com/${REPO}/cmd/${BIN}@latest" ;;
	esac
}

# --- version ------------------------------------------------------------------

latest_version() {
	# The redirect target of /releases/latest names the tag, which avoids the
	# API's unauthenticated rate limit that a JSON call would hit.
	url=$(curl -fsSLo /dev/null -w '%{url_effective}' "https://github.com/${REPO}/releases/latest") ||
		fail "could not reach GitHub to resolve the latest release"
	tag=${url##*/}
	[ -n "$tag" ] && [ "$tag" != "releases" ] ||
		fail "could not resolve the latest release. Set SEDUM_VERSION to install a specific one."
	echo "$tag"
}

# --- checksum -----------------------------------------------------------------

sha256_of() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | cut -d' ' -f1
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$1" | cut -d' ' -f1
	else
		fail "neither sha256sum nor shasum is available, so the download cannot be verified"
	fi
}

verify() {
	archive=$1 sums=$2 name=$3

	expected=$(awk -v n="$name" '$2 == n || $2 == "*" n { print $1; exit }' "$sums")
	[ -n "$expected" ] || fail "$name is not listed in checksums.txt for this release"

	actual=$(sha256_of "$archive")
	[ "$actual" = "$expected" ] ||
		fail "checksum mismatch for ${name}: expected ${expected}, got ${actual}. Refusing to install."
}

# --- main ---------------------------------------------------------------------

main() {
	need curl
	need tar
	need uname

	os=$(detect_os)
	arch=$(detect_arch)

	tag=${SEDUM_VERSION:-}
	if [ -z "$tag" ]; then
		tag=$(latest_version)
	fi
	case "$tag" in v*) ;; *) tag="v${tag}" ;; esac
	version=${tag#v}

	dest=${SEDUM_INSTALL_DIR:-$HOME/.local/bin}

	name="${BIN}_${version}_${os}_${arch}.tar.gz"
	base="https://github.com/${REPO}/releases/download/${tag}"

	tmp=$(mktemp -d) || fail "could not create a temporary directory"
	trap 'rm -rf "$tmp"' EXIT INT TERM

	info "Installing ${BIN} ${tag} (${os}/${arch})"

	curl -fsSL -o "${tmp}/${name}" "${base}/${name}" ||
		fail "could not download ${name}. Check that ${tag} has a ${os}/${arch} build at ${base}"
	curl -fsSL -o "${tmp}/checksums.txt" "${base}/checksums.txt" ||
		fail "could not download checksums.txt for ${tag}, so the archive cannot be verified"

	verify "${tmp}/${name}" "${tmp}/checksums.txt" "$name"
	info "Checksum verified."

	tar -xzf "${tmp}/${name}" -C "$tmp" ||
		fail "could not extract ${name}"
	[ -f "${tmp}/${BIN}" ] || fail "${name} did not contain a ${BIN} binary"

	mkdir -p "$dest" || fail "could not create ${dest}"
	# Staged next to the target and moved into place, so an interrupted install
	# cannot leave a half-written binary where a working one used to be.
	cp "${tmp}/${BIN}" "${dest}/.${BIN}.tmp" || fail "could not write to ${dest}"
	chmod +x "${dest}/.${BIN}.tmp"
	mv "${dest}/.${BIN}.tmp" "${dest}/${BIN}" || fail "could not install into ${dest}"

	installed=$("${dest}/${BIN}" --version 2>/dev/null || true)
	if [ "$installed" != "$version" ]; then
		warn "installed binary reports version '${installed}', expected '${version}'"
	fi

	info "Installed ${dest}/${BIN} (${installed:-$version})"

	case ":${PATH}:" in
		*":${dest}:"*) ;;
		*)
			info ""
			info "${dest} is not on your PATH. Add it:"
			info "    export PATH=\"${dest}:\$PATH\""
			;;
	esac
}

main "$@"
