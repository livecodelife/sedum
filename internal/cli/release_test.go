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
	install := strings.Index(sh, "mv \"${dest}")
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
