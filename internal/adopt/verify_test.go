package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/verify_test.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/verify_test.go","path":"testing"}}
import "testing"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/verify_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestStripMarkers(t *testing.T) {\n\tin := \"package p\\n\\n// sedum:anchor:decls\\n// sedum:setFunction {\\\"a\\\":1}\\nfunc f() {}\\n\\n// /sedum:setFunction\\n\"\n\twant := \"package p\\n\\nfunc f() {}\\n\\n\"\n\tgot := StripMarkers(genpkg.Comment{Prefix: \"//\"}, in)\n\tif got != want {\n\t\tt.Errorf(\"StripMarkers() = %q, want %q\", got, want)\n\t}\n}","file":"internal/adopt/verify_test.go","name":"TestStripMarkers"}}
func TestStripMarkers(t *testing.T) {
	in := "package p\n\n// sedum:anchor:decls\n// sedum:setFunction {\"a\":1}\nfunc f() {}\n\n// /sedum:setFunction\n"
	want := "package p\n\nfunc f() {}\n\n"
	got := StripMarkers(genpkg.Comment{Prefix: "//"}, in)
	if got != want {
		t.Errorf("StripMarkers() = %q, want %q", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestStripMarkersKeepsIndentedMarkerLookalikes(t *testing.T) {\n\tin := \"\\tx := \\\"// sedum:\\\"\\n\\t// sedum:setFunction {}\\n\"\n\twant := \"\\tx := \\\"// sedum:\\\"\\n\"\n\tgot := StripMarkers(genpkg.Comment{Prefix: \"//\"}, in)\n\tif got != want {\n\t\tt.Errorf(\"StripMarkers() = %q, want %q\", got, want)\n\t}\n}","file":"internal/adopt/verify_test.go","name":"TestStripMarkersKeepsIndentedMarkerLookalikes"}}
func TestStripMarkersKeepsIndentedMarkerLookalikes(t *testing.T) {
	in := "\tx := \"// sedum:\"\n\t// sedum:setFunction {}\n"
	want := "\tx := \"// sedum:\"\n"
	got := StripMarkers(genpkg.Comment{Prefix: "//"}, in)
	if got != want {
		t.Errorf("StripMarkers() = %q, want %q", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestVerifyIgnoresBlankLines(t *testing.T) {\n\toriginal := \"a\\nb\\n\"\n\tadopted := \"a\\n\\n// sedum:x {}\\nb\\n// /sedum:x\\n\\n\"\n\tif err := Verify(genpkg.Comment{Prefix: \"//\"}, original, adopted); err != nil {\n\t\tt.Fatalf(\"unexpected error: %v\", err)\n\t}\n}","file":"internal/adopt/verify_test.go","name":"TestVerifyIgnoresBlankLines"}}
func TestVerifyIgnoresBlankLines(t *testing.T) {
	original := "a\nb\n"
	adopted := "a\n\n// sedum:x {}\nb\n// /sedum:x\n\n"
	if err := Verify(genpkg.Comment{Prefix: "//"}, original, adopted); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestVerifyReportsFirstDifference(t *testing.T) {\n\terr := Verify(genpkg.Comment{Prefix: \"//\"}, \"a\\nb\\nc\\n\", \"a\\nB\\nc\\n\")\n\tvar verr *VerifyError\n\tif !errors.As(err, &verr) {\n\t\tt.Fatalf(\"want *VerifyError, got %v\", err)\n\t}\n\tif verr.Line != 2 || verr.Want != \"b\" || verr.Got != \"B\" {\n\t\tt.Errorf(\"got %+v\", *verr)\n\t}\n}","file":"internal/adopt/verify_test.go","name":"TestVerifyReportsFirstDifference"}}
func TestVerifyReportsFirstDifference(t *testing.T) {
	err := Verify(genpkg.Comment{Prefix: "//"}, "a\nb\nc\n", "a\nB\nc\n")
	var verr *VerifyError
	if !errors.As(err, &verr) {
		t.Fatalf("want *VerifyError, got %v", err)
	}
	if verr.Line != 2 || verr.Want != "b" || verr.Got != "B" {
		t.Errorf("got %+v", *verr)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestVerifyLengthMismatch(t *testing.T) {\n\terr := Verify(genpkg.Comment{Prefix: \"//\"}, \"a\\nb\\n\", \"a\\n\")\n\tvar verr *VerifyError\n\tif !errors.As(err, &verr) {\n\t\tt.Fatalf(\"want *VerifyError, got %v\", err)\n\t}\n\tif verr.Line != 2 || verr.Want != \"b\" || verr.Got != \"<end of file>\" {\n\t\tt.Errorf(\"got %+v\", *verr)\n\t}\n\terr = Verify(genpkg.Comment{Prefix: \"//\"}, \"a\\n\", \"a\\nb\\n\")\n\tif !errors.As(err, &verr) {\n\t\tt.Fatalf(\"want *VerifyError, got %v\", err)\n\t}\n\tif verr.Want != \"<end of file>\" {\n\t\tt.Errorf(\"got %+v\", *verr)\n\t}\n}","file":"internal/adopt/verify_test.go","name":"TestVerifyLengthMismatch"}}
func TestVerifyLengthMismatch(t *testing.T) {
	err := Verify(genpkg.Comment{Prefix: "//"}, "a\nb\n", "a\n")
	var verr *VerifyError
	if !errors.As(err, &verr) {
		t.Fatalf("want *VerifyError, got %v", err)
	}
	if verr.Line != 2 || verr.Want != "b" || verr.Got != "<end of file>" {
		t.Errorf("got %+v", *verr)
	}
	err = Verify(genpkg.Comment{Prefix: "//"}, "a\n", "a\nb\n")
	if !errors.As(err, &verr) {
		t.Fatalf("want *VerifyError, got %v", err)
	}
	if verr.Want != "<end of file>" {
		t.Errorf("got %+v", *verr)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestVerifyLeadingWhitespaceMatters(t *testing.T) {\n\tif err := Verify(genpkg.Comment{Prefix: \"//\"}, \"a\\n\\tb\\n\", \"a\\n    b\\n\"); err == nil {\n\t\tt.Error(\"expected a difference for changed leading whitespace\")\n\t}\n}","file":"internal/adopt/verify_test.go","name":"TestVerifyLeadingWhitespaceMatters"}}
func TestVerifyLeadingWhitespaceMatters(t *testing.T) {
	if err := Verify(genpkg.Comment{Prefix: "//"}, "a\n\tb\n", "a\n    b\n"); err == nil {
		t.Error("expected a difference for changed leading whitespace")
	}
}

// /sedum:setFunction
