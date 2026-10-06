package cli

import (
	"os"
	osexec "os/exec"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

// These tests guard prov-2026-4746d9ed's constraints against the release
// configuration drifting away from them. They assert about files rather than
// behavior because that is where these particular decisions live: a version
// stamped by a linker flag, an installer that stops verifying, or a module path
// that no longer resolves are all things nothing else in the suite would catch,
// and all things whose first symptom is a published artifact.

// repoRoot is where this package sits relative to the module root.
const repoRoot = "../../"

// pkgProbe exists so the test can ask the runtime what import path this package
// was compiled under, which is the module path plus its directory.
type pkgProbe struct{}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(repoRoot + path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(b)
}

// modulePath returns the path declared by the module line in go.mod.
func modulePath(t *testing.T) string {
	t.Helper()
	for _, line := range strings.Split(read(t, "go.mod"), "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "module "); ok {
			return strings.TrimSpace(rest)
		}
	}
	t.Fatal("go.mod declares no module path")
	return ""
}

// A rename that updates go.mod and misses an import, or the reverse, leaves a
// module that builds in a checkout and cannot be fetched. The compiler resolves
// imports against the module path, so comparing the two catches a partial one.
func TestModulePathAndImportPathAgree(t *testing.T) {
	declared := modulePath(t)
	compiled := reflect.TypeOf(pkgProbe{}).PkgPath()

	want := declared + "/internal/cli"
	if compiled != want {
		t.Errorf("package compiled as %q, but go.mod declares module %q (expected %q)",
			compiled, declared, want)
	}
}

// The module path is resolved as a URL by `go install`, which then checks that
// the go.mod it fetched declares the path that was asked for. A module path
// naming anything but the repository serving it is installable by no means at
// all, which is the failure this repository shipped with (prov-2026-4746d9ed).
func TestModulePathResolvesToTheRepository(t *testing.T) {
	out, err := osexec.Command("git", "-C", repoRoot, "remote", "get-url", "origin").Output()
	if err != nil {
		t.Skip("no git origin to compare against")
	}

	// git@github.com:owner/repo.git and https://github.com/owner/repo.git both
	// reduce to github.com/owner/repo.
	remote := strings.TrimSpace(string(out))
	remote = strings.TrimSuffix(remote, ".git")
	remote = strings.TrimPrefix(remote, "https://")
	remote = strings.TrimPrefix(remote, "git@")
	remote = strings.Replace(remote, ":", "/", 1)

	if declared := modulePath(t); remote != declared {
		t.Errorf("go.mod declares module %q but origin is %q; `go install %s/cmd/sedum@latest` cannot work",
			declared, remote, declared)
	}
}

// prov-2026-b5465dfa made the version a constant rather than a build-time stamp,
// because a version injected by -ldflags is absent from `go build` and `go run`
// — exactly when someone is developing against it. Release engineering's reflex
// is to add -X here, and this fails when it is added.
func TestReleaseDoesNotStampTheVersion(t *testing.T) {
	cfg := read(t, ".goreleaser.yml")

	stamp := regexp.MustCompile(`-X\s+\S*[Vv]ersion\s*=`)
	for _, line := range strings.Split(cfg, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue // the header explains why there is no -X; it may say so.
		}
		if stamp.MatchString(line) {
			t.Errorf(".goreleaser.yml stamps the version at link time (%q).\n"+
				"prov-2026-b5465dfa requires it to stay a const in internal/recording; "+
				"the tag is checked against that const in release.yml instead.", strings.TrimSpace(line))
		}
	}
}

// Because the version is not stamped, nothing about the build makes the tag and
// the constant agree. The release workflow is what does, and a release that
// published a v0.2.0 binary reporting 0.1.0 would break every caller pinning a
// version floor.
func TestReleaseWorkflowGatesTheTagAgainstTheConstant(t *testing.T) {
	wf := read(t, ".github/workflows/release.yml")

	for _, want := range []string{
		"internal/recording/version.go", // it reads the constant from source
		"GITHUB_REF_NAME",               // and compares it to the tag
		"exit 1",                        // and fails rather than warning
	} {
		if !strings.Contains(wf, want) {
			t.Errorf("release workflow does not mention %q, so nothing stops the tag and "+
				"the version constant from disagreeing", want)
		}
	}
}

// An installer run by piping a URL into a shell that then skips the checksum it
// publishes is offering a guarantee it does not keep.
func TestInstallerVerifiesWhatItDownloads(t *testing.T) {
	sh := read(t, "install.sh")

	if !strings.Contains(sh, "checksums.txt") {
		t.Error("install.sh never fetches checksums.txt")
	}
	if !strings.Contains(sh, "checksum mismatch") {
		t.Error("install.sh has no checksum mismatch failure path")
	}

	// A verification that runs after the binary is in place verifies nothing.
	verify := strings.Index(sh, "\tverify ")
	install := strings.Index(sh, `install_file "${tmp}/${BIN}"`)
	if verify < 0 || install < 0 {
		t.Fatal("install.sh no longer has the verify and install steps this test recognises")
	}
	if verify > install {
		t.Error("install.sh installs the binary before verifying the archive it came from")
	}

	for _, insecure := range []string{"--insecure", "curl -k", "-fsSLk"} {
		if strings.Contains(sh, insecure) {
			t.Errorf("install.sh disables TLS verification (%q)", insecure)
		}
	}
}

// The installer builds the archive name itself, so it and the release config
// have to agree on the template. They are edited in different files, months
// apart, and a disagreement is invisible until a release nobody can install.
func TestInstallerAndReleaseAgreeOnArchiveNames(t *testing.T) {
	if cfg := read(t, ".goreleaser.yml"); !strings.Contains(cfg,
		`name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"`) {
		t.Error(".goreleaser.yml's archive name_template changed; install.sh derives the same name by hand")
	}
	if sh := read(t, "install.sh"); !strings.Contains(sh, `name="${BIN}_${version}_${os}_${arch}.tar.gz"`) {
		t.Error("install.sh's archive name changed; it must match .goreleaser.yml's name_template")
	}
}

// prov-2026-d2f1bdf2: the bundled default model is whatever the release asset
// URL serves, and the only thing standing between a release build and an
// unverified model is the checksum this script pins. The maintainer uploads the
// asset; until they do, the pin does not match and the build fails, which is
// intended. So the pin is asserted to the digit, and the URL and filename that
// everything else (goreleaser, install.sh, bundled.go) depends on are asserted
// to have stayed where they were.

const newDefaultModelSHA = "5185a0b25d5a13262fa6dfcdb1b2091b7bb64f130d65cdb5322c99c992dd052e"

func TestFetchScriptPinsTheNewDefaultModel(t *testing.T) {
	sh := read(t, "hack/fetch-release-assets.sh")

	m := regexp.MustCompile(`(?m)^MODEL_SHA=["']?([0-9a-fA-F]+)["']?\s*$`).FindStringSubmatch(sh)
	if m == nil {
		t.Fatal("hack/fetch-release-assets.sh has no MODEL_SHA=<hex> line this test recognises")
	}
	if m[1] != newDefaultModelSHA {
		t.Errorf("MODEL_SHA = %s, want %s (sedum-binder-final-2b-le12k-e5, Qwen3.5-2B.Q8_0.gguf)", m[1], newDefaultModelSHA)
	}
}

func TestFetchScriptKeepsTheModelURLAndBundledFilename(t *testing.T) {
	sh := read(t, "hack/fetch-release-assets.sh")

	const url = `MODEL_URL="https://github.com/livecodelife/sedum/releases/download/model-assets/sedum-default-model.gguf"`
	if !strings.Contains(sh, url) {
		t.Errorf("MODEL_URL changed; the record keeps it and the bundled filename as they were:\n  want %s", url)
	}
	if !strings.Contains(sh, `"$DIST/sedum-default-model.gguf"`) {
		t.Error("the staged model filename is no longer sedum-default-model.gguf, which bundled.go looks up")
	}
	if !strings.Contains(sh, `fetch_and_verify "$MODEL_URL"`) {
		t.Error("the model is no longer fetched through fetch_and_verify, so nothing checks its checksum")
	}
}

// The pin changing without the reason being said beside it would leave the next
// person staring at a release build that fails its checksum with no way to know
// that is the intended order of operations (upload first, then tag).
func TestFetchScriptSaysTheUploadComesBeforeTheTag(t *testing.T) {
	sh := strings.ToLower(read(t, "hack/fetch-release-assets.sh"))

	for _, want := range []string{"upload", "checksum"} {
		if !strings.Contains(sh, want) {
			t.Errorf("hack/fetch-release-assets.sh does not mention %q beside the pin; the record requires the "+
				"upload-then-tag requirement to be stated wherever the pin is changed", want)
		}
	}
	if !strings.Contains(sh, "sedum-binder-final-2b-le12k-e5") {
		t.Error("hack/fetch-release-assets.sh does not name the checkpoint the pin is for")
	}
}

func TestReadmeNamesTheNewDefaultCheckpoint(t *testing.T) {
	if !strings.Contains(read(t, "README.md"), "sedum-binder-final-2b-le12k-e5") {
		t.Error("README.md does not name sedum-binder-final-2b-le12k-e5 as the default model")
	}
}
