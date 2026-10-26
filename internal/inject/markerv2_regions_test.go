package inject

import (
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
)

// prov-2026-f6bd09e5: FindRegions reports, for each region, the hash of the text
// actually present between its marker lines.

func TestFindRegionsBodyHash(t *testing.T) {
	const (
		open  = `# sedum:a {"v":2,"tier":"owned"}`
		close = `# /sedum:a`
	)
	tests := []struct {
		name    string
		content string
		body    string // the text between the marker lines, as it is in content
	}{
		{"one line", "head\n" + open + "\nx := 1\n" + close + "\ntail\n", "x := 1\n"},
		{"several lines with their own indentation", "head\n" + open + "\n  a\n\tb\n    c\n" + close + "\n", "  a\n\tb\n    c\n"},
		{"the marker lines' indentation is not part of the body",
			"class X\n  " + open + "\n  def a\n  end\n  " + close + "\nend\n", "  def a\n  end\n"},
		{"markers indented deeper than the body",
			"    " + open + "\nx\n    " + close + "\n", "x\n"},
		{"a closing line directly after the opening line", open + "\n" + close + "\n", ""},
		{"blank lines inside are part of the body", open + "\n\nx\n\n" + close + "\n", "\nx\n\n"},
		{"trailing spaces are part of the body", open + "\nx \n" + close + "\n", "x \n"},
		{"CRLF line endings are taken as they are", open + "\r\na\r\nb\r\n" + close + "\r\n", "a\r\nb\r\n"},
		{"a lone carriage return is taken as it is", open + "\na\rb\n" + close + "\n", "a\rb\n"},
		{"a body line that looks like a marker for another label is body text",
			open + "\n# sedum:b\n" + close + "\n", "# sedum:b\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			regions, err := FindRegions(lineComment("#"), tc.content)
			if err != nil {
				t.Fatalf("FindRegions: %v", err)
			}
			if len(regions) != 1 {
				t.Fatalf("found %d regions, want 1", len(regions))
			}
			if want := BodyHash(tc.body); regions[0].BodyHash != want {
				t.Errorf("BodyHash = %s, want the hash of %q (%s)", regions[0].BodyHash, tc.body, want)
			}
			if want := sha256Of(strings.ReplaceAll(tc.body, "\r\n", "\n")); regions[0].BodyHash != want {
				t.Errorf("BodyHash = %s, want the hash of the LF-normalized body (%s)", regions[0].BodyHash, want)
			}
		})
	}
}

func TestFindRegionsOffsetsAndHashesPerRegion(t *testing.T) {
	first := "# sedum:a {\"v\":2,\"tier\":\"owned\"}\none\n# /sedum:a\n"
	second := "  # sedum:b:x {\"tier\":\"owned\",\"kwargs\":{\"n\":1}}\n  two\n  # /sedum:b:x\n"
	content := "head\n" + first + "mid\n" + second + "tail\n"

	regions, err := FindRegions(lineComment("#"), content)
	if err != nil {
		t.Fatalf("FindRegions: %v", err)
	}
	if len(regions) != 2 {
		t.Fatalf("found %d regions, want 2", len(regions))
	}

	if got := content[regions[0].Start:regions[0].End]; got != first {
		t.Errorf("first region extent = %q, want %q", got, first)
	}
	if got := content[regions[1].Start:regions[1].End]; got != second {
		t.Errorf("second region extent = %q, want %q", got, second)
	}
	if regions[0].BodyHash != BodyHash("one\n") {
		t.Errorf("first BodyHash = %s, want the hash of %q", regions[0].BodyHash, "one\n")
	}
	if regions[1].BodyHash != BodyHash("  two\n") {
		t.Errorf("second BodyHash = %s, want the hash of %q", regions[1].BodyHash, "  two\n")
	}
	if regions[0].Marker.Version != 2 || regions[1].Marker.Version != 0 {
		t.Errorf("versions = %d, %d; want 2 then 0", regions[0].Marker.Version, regions[1].Marker.Version)
	}
}

// BodyHash is set whatever the marker's version, and for a foreign marker too.
func TestFindRegionsHashesEveryRegionWhateverItsVersion(t *testing.T) {
	content := `# sedum:a {"v":2,"tier":"owned","body":"sha256:stale"}` + "\nv2 text\n# /sedum:a\n" +
		`# sedum:b {"tier":"owned","kwargs":{"n":1}}` + "\nv1 text\n# /sedum:b\n" +
		`# sedum:c {"v":3,"tier":"owned","anything":true}` + "\nforeign text\n# /sedum:c\n" +
		`# sedum:d` + "\nbare text\n# /sedum:d\n"

	regions, err := FindRegions(lineComment("#"), content)
	if err != nil {
		t.Fatalf("FindRegions: %v", err)
	}
	wantBodies := []string{"v2 text\n", "v1 text\n", "foreign text\n", "bare text\n"}
	if len(regions) != len(wantBodies) {
		t.Fatalf("found %d regions, want %d", len(regions), len(wantBodies))
	}
	for i, body := range wantBodies {
		if want := BodyHash(body); regions[i].BodyHash != want {
			t.Errorf("region %d (%s) BodyHash = %s, want %s", i, regions[i].Marker.Label(), regions[i].BodyHash, want)
		}
	}
	if !regions[2].Marker.Foreign {
		t.Error("the v3 region was not reported as foreign: a foreign region is still found")
	}
	// The hash on the marker is what was recorded; BodyHash is what is there.
	if regions[0].Marker.Body != "sha256:stale" || regions[0].BodyHash == regions[0].Marker.Body {
		t.Errorf("Marker.Body = %q, BodyHash = %q; want the recorded hash kept apart from the actual one",
			regions[0].Marker.Body, regions[0].BodyHash)
	}
}

func TestFindRegionsBodyHashThroughADelimitedComment(t *testing.T) {
	comment := genpkg.Comment{Prefix: "<!--", Suffix: "-->"}
	content := "<div>\n" +
		`  <!-- sedum:addCol:a {"v":2,"tier":"owned"} -->` + "\n" +
		"  <td>a</td>\n" +
		"  <!-- /sedum:addCol:a -->\n" +
		"</div>\n"
	regions, err := FindRegions(comment, content)
	if err != nil {
		t.Fatalf("FindRegions: %v", err)
	}
	if len(regions) != 1 {
		t.Fatalf("found %d regions, want 1", len(regions))
	}
	if want := BodyHash("  <td>a</td>\n"); regions[0].BodyHash != want {
		t.Errorf("BodyHash = %s, want %s", regions[0].BodyHash, want)
	}
}

// A v2 marker with an unreadable kwargs is a read error, which FindRegions
// returns rather than skipping the region.
func TestFindRegionsReturnsAParseError(t *testing.T) {
	content := `# sedum:a {"v":2,"kwargs":{"x":1}}` + "\nbody\n# /sedum:a\n"
	_, err := FindRegions(lineComment("#"), content)
	if err == nil || !strings.Contains(err.Error(), "kwargs must be a string") {
		t.Errorf("FindRegions error = %v, want one naming that kwargs must be a string", err)
	}
}
