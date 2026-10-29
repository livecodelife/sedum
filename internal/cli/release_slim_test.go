package cli

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	osexec "os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Every full archive is about 1.7 GB, so a release also publishes a slim one
// that holds the binary alone, and install.sh installs it on request
// (prov-2026-4746d9ed). The configuration is checked as text; the installer is
// run for real against a fake curl that serves a fixture release from a
// directory, so nothing here touches the network.

// --- .goreleaser.yml ----------------------------------------------------------

func goreleaserArchives(t *testing.T) []map[string]any {
	t.Helper()
	var cfg struct {
		Archives []map[string]any `yaml:"archives"`
	}
	if err := yaml.Unmarshal([]byte(read(t, ".goreleaser.yml")), &cfg); err != nil {
		t.Fatalf(".goreleaser.yml: %v", err)
	}
	return cfg.Archives
}

// fileSources lists the src (or the bare name) of each entry of an archive's
// files.
func fileSources(t *testing.T, archive map[string]any) []string {
	t.Helper()
	var out []string
	files, _ := archive["files"].([]any)
	for _, f := range files {
		switch v := f.(type) {
		case string:
			out = append(out, v)
		case map[string]any:
			src, _ := v["src"].(string)
			out = append(out, src)
		}
	}
	return out
}

func archiveByName(t *testing.T, suffix string) map[string]any {
	t.Helper()
	for _, a := range goreleaserArchives(t) {
		name, _ := a["name_template"].(string)
		if strings.HasSuffix(name, suffix) {
			return a
		}
	}
	t.Fatalf(".goreleaser.yml has no archive whose name_template ends with %q", suffix)
	return nil
}

func TestReleasePublishesASlimArchiveBesideTheFullOne(t *testing.T) {
	const base = "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
	full := archiveByName(t, "{{ .Arch }}")
	slim := archiveByName(t, "{{ .Arch }}_slim")

	if got := full["name_template"]; got != base {
		t.Errorf("the full archive name_template = %v, want %q", got, base)
	}
	if got := slim["name_template"]; got != base+"_slim" {
		t.Errorf("the slim archive name_template = %v, want %q", got, base+"_slim")
	}
	if full["id"] == slim["id"] {
		t.Errorf("both archives have the id %v", full["id"])
	}

	// The same build goes into both.
	fullIDs, _ := full["ids"].([]any)
	slimIDs, _ := slim["ids"].([]any)
	if fmt.Sprint(fullIDs) != fmt.Sprint(slimIDs) || fmt.Sprint(slimIDs) != "[sedum]" {
		t.Errorf("the archives must both pack the sedum build: full %v, slim %v", fullIDs, slimIDs)
	}

	has := func(sources []string, fragment string) bool {
		for _, s := range sources {
			if strings.Contains(s, fragment) {
				return true
			}
		}
		return false
	}
	fullFiles, slimFiles := fileSources(t, full), fileSources(t, slim)
	for _, want := range []string{"README.md", "LICENSE", "TOOL_BOUNDARIES.md", "conformance/"} {
		if !has(fullFiles, want) || !has(slimFiles, want) {
			t.Errorf("%s is not in both archives (full %v, slim %v)", want, fullFiles, slimFiles)
		}
	}
	for _, bundled := range []string{"goinfer-serve", "sedum-default-model.gguf"} {
		if !has(fullFiles, bundled) {
			t.Errorf("the full archive no longer holds %s", bundled)
		}
		if has(slimFiles, bundled) {
			t.Errorf("the slim archive holds %s", bundled)
		}
	}
}

// checksums.txt lists every archive the release publishes: the checksum block
// names that file and is not narrowed to some of the archives.
func TestChecksumsCoverEveryArchive(t *testing.T) {
	var cfg struct {
		Checksum map[string]any `yaml:"checksum"`
	}
	if err := yaml.Unmarshal([]byte(read(t, ".goreleaser.yml")), &cfg); err != nil {
		t.Fatal(err)
	}
	if got := cfg.Checksum["name_template"]; got != "checksums.txt" {
		t.Errorf("checksum name_template = %v, want checksums.txt", got)
	}
	if _, narrowed := cfg.Checksum["ids"]; narrowed {
		t.Error("the checksum block lists ids, so it covers only some of the archives")
	}
}

// --- install.sh against a fake release ---------------------------------------

const (
	releaseTag     = "v9.9.9"
	releaseVersion = "9.9.9"
	fullName       = "sedum_9.9.9_linux_amd64.tar.gz"
	slimName       = "sedum_9.9.9_linux_amd64_slim.tar.gz"
)

// tarGz builds an archive of the given files, executable ones included.
func tarGz(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		mode := int64(0o644)
		if n == "sedum" || n == "goinfer-serve" {
			mode = 0o755
		}
		if err := tw.WriteHeader(&tar.Header{Name: n, Mode: mode, Size: int64(len(files[n]))}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(files[n])); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// release is a directory of the files a release publishes. The names are the
// ones in the URLs the installer asks for.
type release struct {
	dir     string
	listed  []string
	archive map[string][]byte
}

// newRelease publishes the full and the slim archive and lists both in
// checksums.txt.
func newRelease(t *testing.T) *release {
	t.Helper()
	r := &release{dir: t.TempDir(), archive: map[string][]byte{}}
	sedum := "#!/bin/sh\necho " + releaseVersion + "\n"
	r.add(t, fullName, tarGz(t, map[string]string{
		"sedum": sedum, "goinfer-serve": "server", "sedum-default-model.gguf": "weights",
		"README.md": "readme", "LICENSE": "license"}), true)
	r.add(t, slimName, tarGz(t, map[string]string{
		"sedum": sedum, "README.md": "readme", "LICENSE": "license"}), true)
	return r
}

// add publishes an archive; listed says whether checksums.txt names it.
func (r *release) add(t *testing.T, name string, data []byte, listed bool) {
	t.Helper()
	if data != nil {
		if err := os.WriteFile(filepath.Join(r.dir, name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if listed {
		sum := sha256.Sum256(data)
		r.listed = append(r.listed, hex.EncodeToString(sum[:])+"  "+name)
	}
	r.writeChecksums(t)
}

func (r *release) writeChecksums(t *testing.T) {
	t.Helper()
	text := strings.Join(r.listed, "\n") + "\n"
	if err := os.WriteFile(filepath.Join(r.dir, "checksums.txt"), []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// withoutSlim removes the slim archive and its checksums line, as in a release
// made before slim archives existed.
func (r *release) withoutSlim(t *testing.T) {
	t.Helper()
	if err := os.Remove(filepath.Join(r.dir, slimName)); err != nil {
		t.Fatal(err)
	}
	r.unlist(t)
}

// unlist drops the slim archive's line from checksums.txt and keeps the file.
func (r *release) unlist(t *testing.T) {
	t.Helper()
	var kept []string
	for _, l := range r.listed {
		if !strings.HasSuffix(l, slimName) {
			kept = append(kept, l)
		}
	}
	r.listed = kept
	r.writeChecksums(t)
}

// installResult is what one run of install.sh did.
type installResult struct {
	stdout, stderr string
	err            error
	dest           string
}

// install runs install.sh from the repository with a fake curl serving r and a
// fake uname reporting linux/amd64. env is added to the environment; args go to
// the script. dest is the install directory it is told to use.
func install(t *testing.T, r *release, dest string, env []string, args ...string) installResult {
	t.Helper()
	return runInstaller(t, r, dest, env, false, args)
}

// installPiped is install() the documented way: the script comes in on standard
// input and the arguments follow `sh -s --`.
func installPiped(t *testing.T, r *release, dest string, args ...string) installResult {
	t.Helper()
	return runInstaller(t, r, dest, nil, true, args)
}

func runInstaller(t *testing.T, r *release, dest string, env []string, piped bool, args []string) installResult {
	t.Helper()
	bin := t.TempDir()
	curl := `#!/bin/sh
# Serves files from $FAKE_RELEASE_DIR by the last path element of the URL.
out=""
url=""
while [ $# -gt 0 ]; do
	case "$1" in
		-o|-fsSLo) out=$2; shift 2 ;;
		-w) shift 2 ;;
		-*) shift ;;
		*) url=$1; shift ;;
	esac
done
name=${url##*/}
[ -f "$FAKE_RELEASE_DIR/$name" ] || exit 22
cp "$FAKE_RELEASE_DIR/$name" "$out"
`
	uname := "#!/bin/sh\ncase \"$1\" in -s) echo Linux ;; -m) echo x86_64 ;; *) echo Linux ;; esac\n"
	for name, body := range map[string]string{"curl": curl, "uname": uname} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	script, err := filepath.Abs(repoRoot + "install.sh")
	if err != nil {
		t.Fatal(err)
	}
	var cmd *osexec.Cmd
	if piped {
		cmd = osexec.Command("sh", append([]string{"-s", "--"}, args...)...)
		in, err := os.Open(script)
		if err != nil {
			t.Fatal(err)
		}
		defer in.Close()
		cmd.Stdin = in
	} else {
		cmd = osexec.Command("sh", append([]string{script}, args...)...)
	}
	cmd.Env = append(os.Environ(),
		"PATH="+bin+":"+os.Getenv("PATH"),
		"FAKE_RELEASE_DIR="+r.dir,
		"SEDUM_VERSION="+releaseTag,
		"SEDUM_INSTALL_DIR="+dest,
		"HOME="+t.TempDir(),
	)
	cmd.Env = append(cmd.Env, env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()
	return installResult{stdout: stdout.String(), stderr: stderr.String(), err: err, dest: dest}
}

// files lists the names in a directory; a directory that does not exist has none.
func files(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out
}

func wantFiles(t *testing.T, dir string, want ...string) {
	t.Helper()
	sort.Strings(want)
	if got := files(t, dir); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("%s holds %v, want %v", dir, got, want)
	}
}

func (i installResult) mustSucceed(t *testing.T) {
	t.Helper()
	if i.err != nil {
		t.Fatalf("install.sh failed: %v\nstderr:\n%s", i.err, i.stderr)
	}
}

func (i installResult) mustFail(t *testing.T) {
	t.Helper()
	if i.err == nil {
		t.Fatalf("install.sh succeeded; stderr:\n%s", i.stderr)
	}
}

func wantIn(t *testing.T, what, text string, needles ...string) {
	t.Helper()
	for _, n := range needles {
		if !strings.Contains(text, n) {
			t.Errorf("%s does not contain %q:\n%s", what, n, text)
		}
	}
}

func TestInstallerInstallsTheFullArchiveByDefault(t *testing.T) {
	r := newRelease(t)
	dest := filepath.Join(t.TempDir(), "bin")

	got := install(t, r, dest, nil)

	got.mustSucceed(t)
	wantFiles(t, dest, "sedum", "goinfer-serve", "sedum-default-model.gguf")
	wantIn(t, "stderr", got.stderr, "bundled local model")
}

func TestNoModelInstallsOnlyTheSlimBinary(t *testing.T) {
	r := newRelease(t)
	dest := filepath.Join(t.TempDir(), "bin")

	got := install(t, r, dest, nil, "--no-model")

	got.mustSucceed(t)
	wantFiles(t, dest, "sedum")
	wantIn(t, "stderr", got.stderr, "no local model",
		"install the full archive", "without --no-model", "SEDUM_NO_MODEL")
	if strings.Contains(got.stderr, "bundled local model") {
		t.Errorf("a slim install reported a bundled local model:\n%s", got.stderr)
	}
}

func TestOnlyTheExactValueOneOfSedumNoModelSelectsSlim(t *testing.T) {
	cases := []struct {
		value string
		slim  bool
	}{{"1", true}, {"0", false}, {"true", false}, {"", false}, {"11", false}}
	for _, c := range cases {
		t.Run("SEDUM_NO_MODEL="+c.value, func(t *testing.T) {
			r := newRelease(t)
			dest := filepath.Join(t.TempDir(), "bin")

			got := install(t, r, dest, []string{"SEDUM_NO_MODEL=" + c.value})

			got.mustSucceed(t)
			if c.slim {
				wantFiles(t, dest, "sedum")
			} else {
				wantFiles(t, dest, "sedum", "goinfer-serve", "sedum-default-model.gguf")
			}
		})
	}
}

// The documented way to run the installer is piped into sh, where the argument
// comes after -s --.
func TestNoModelIsReadWhenTheInstallerIsPipedIntoSh(t *testing.T) {
	r := newRelease(t)
	dest := filepath.Join(t.TempDir(), "bin")

	got := installPiped(t, r, dest, "--no-model")

	got.mustSucceed(t)
	wantFiles(t, dest, "sedum")
}

// A slim install writes the binary and leaves what an earlier full install put
// there alone.
func TestSlimInstallLeavesAnEarlierFullInstallAlone(t *testing.T) {
	r := newRelease(t)
	dest := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{"goinfer-serve": "old server", "sedum-default-model.gguf": "old weights"} {
		if err := os.WriteFile(filepath.Join(dest, name), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	got := install(t, r, dest, nil, "--no-model")

	got.mustSucceed(t)
	wantFiles(t, dest, "sedum", "goinfer-serve", "sedum-default-model.gguf")
	for name, body := range map[string]string{"goinfer-serve": "old server", "sedum-default-model.gguf": "old weights"} {
		data, err := os.ReadFile(filepath.Join(dest, name))
		if err != nil || string(data) != body {
			t.Errorf("%s was changed: %q, %v", name, data, err)
		}
	}
}

func TestAnUnknownArgumentFailsNamingItAndInstallsNothing(t *testing.T) {
	r := newRelease(t)
	dest := filepath.Join(t.TempDir(), "bin")

	got := install(t, r, dest, nil, "--no-modle")

	got.mustFail(t)
	wantIn(t, "stderr", got.stderr, "--no-modle")
	if _, err := os.Stat(dest); err == nil {
		t.Errorf("an install directory was created for a refused run: %v", files(t, dest))
	}
}

// A release made before slim archives existed lists none in checksums.txt, and
// may not hold the file at all. A slim request then fails: nonzero exit, a
// message naming the tag and the slim archive, nothing installed, and no quiet
// fallback to the full archive.
func TestNoModelOnAReleaseWithoutTheSlimArchiveFailsAndInstallsNothing(t *testing.T) {
	cases := []struct {
		name  string
		setup func(*release, *testing.T)
	}{
		{"neither the file nor a checksum line", (*release).withoutSlim},
		{"the file but no checksum line", (*release).unlist},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newRelease(t)
			c.setup(r, t)
			dest := filepath.Join(t.TempDir(), "bin")
			if err := os.MkdirAll(dest, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dest, "keep"), []byte("mine"), 0o644); err != nil {
				t.Fatal(err)
			}

			got := install(t, r, dest, nil, "--no-model")

			got.mustFail(t)
			wantIn(t, "stderr", got.stderr, releaseTag, slimName)
			wantFiles(t, dest, "keep")
			if data, err := os.ReadFile(filepath.Join(dest, "keep")); err != nil || string(data) != "mine" {
				t.Errorf("the install directory was touched: %q, %v", data, err)
			}
		})
	}
}

// A slim archive is checked against checksums.txt like the full one.
func TestSlimArchiveWithAWrongChecksumIsRefused(t *testing.T) {
	r := newRelease(t)
	r.unlist(t)
	r.listed = append(r.listed, strings.Repeat("0", 64)+"  "+slimName)
	r.writeChecksums(t)
	dest := filepath.Join(t.TempDir(), "bin")

	got := install(t, r, dest, nil, "--no-model")

	got.mustFail(t)
	wantIn(t, "stderr", got.stderr, "checksum mismatch")
	wantFiles(t, dest)
}
