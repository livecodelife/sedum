package genpkg

import (
	"strings"
	"testing"
)

// Loading a package whose target closes its comments (prov-2026-a6f6bb81).

// htmlPackage is a package targeting a delimited-comment language, shaped the
// way the erp-grid package that motivated comment_suffix is shaped: one anchor
// in a file template, one action injecting a per-column region at it.
func htmlPackage() map[string]string {
	return map[string]string{
		"web/sedum.yaml": `name: web
extensions: [".html"]
comment_prefix: "<!--"
comment_suffix: "-->"
transforms:
  tsProp: [camel]
`,
		"web/files/app/features/{feature}/list.html": "<app-grid>\n" +
			"  <!-- sedum:anchor:columns -->\n" +
			"</app-grid>\n",
		"web/actions/actions.yaml": `actions:
  addTextColumn:
    kwargs:
      feature: { type: string, required: true }
      column: { type: string, required: true }
    identity: [column]
    injects_into: "app/features/{{feature}}/list.html"
    anchor: columns
`,
		"web/actions/addTextColumn.html": `  <app-grid-col [name]="'{{column|tsProp}}'" />` + "\n",
	}
}

func TestADelimitedPackageLoadsClean(t *testing.T) {
	set, findings := loadTree(t, htmlPackage())

	for _, f := range findings {
		t.Errorf("a package declaring comment_suffix reported %s: %s", f.Rule, f.Message)
	}

	pkg, ok := set.Lookup("web")
	if !ok {
		t.Fatal("the web package did not load")
	}
	if pkg.CommentSuffix != "-->" {
		t.Errorf("comment_suffix = %q, want %q", pkg.CommentSuffix, "-->")
	}

	comment := pkg.Comment()
	if comment.Prefix != "<!--" || comment.Suffix != "-->" {
		t.Errorf("Comment() = %+v", comment)
	}
	if !comment.Delimited() {
		t.Error("a package with a suffix does not report itself delimited")
	}
}

// The anchor check runs through the same shape the file template plants. Before
// the suffix was handled, a delimited package's every anchor read as unplanted
// and the package was rejected with an action targeting a marker "no file
// template plants" - which is exactly the diagnostic an author cannot act on.
func TestADelimitedPackagesAnchorsAreSeenAsPlanted(t *testing.T) {
	_, findings := loadTree(t, htmlPackage())

	for _, f := range findings {
		if f.Rule == RuleAnchorUnplanted {
			t.Fatalf("the columns anchor was not recognized as planted: %s", f.Message)
		}
	}
}

// A closing delimiter with nothing to open is a package error, not a default.
func TestACommentSuffixWithoutAPrefixIsRejected(t *testing.T) {
	files := htmlPackage()
	files["web/sedum.yaml"] = `name: web
extensions: [".html"]
comment_suffix: "-->"
`
	_, findings := loadTree(t, files)

	f := findingFor(t, findings, RuleCommentSuffixOrphan)
	if !strings.Contains(f.Message, "-->") {
		t.Errorf("the diagnostic does not quote the orphaned suffix: %s", f.Message)
	}
}

// An unknown key is still an error. comment_suffix is modelled now; a
// misspelling of it must not load silently and leave every marker unterminated.
func TestAMisspelledCommentSuffixStillFailsLoad(t *testing.T) {
	files := htmlPackage()
	files["web/sedum.yaml"] = `name: web
extensions: [".html"]
comment_prefix: "<!--"
comment_sufix: "-->"
`
	_, findings := loadTree(t, files)
	if len(findings) == 0 {
		t.Fatal("a misspelled comment_suffix loaded clean")
	}
}
