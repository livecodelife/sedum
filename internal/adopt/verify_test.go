package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/verify_test.go","path":"errors"},"kwargs":"sha256:41c176716dde17bc02724b76d52fd57d893ab281cc2af302780cb7083192e8ac","body":"sha256:2b718237ba0a56467c8fd494cc131e942ef410e6228c03b644d70476c71862dc"}
import "errors"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/verify_test.go","path":"testing"},"kwargs":"sha256:443d7e68b116949f353e2e6ed4ebd5949dad32787874102b51ddb742c624c73c","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/verify_test.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:9e2a5e4be15a8546157fc26059d0c464b99519bd06567d8049b908d105806435","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/verify_test.go","name":"TestStripMarkers"},"kwargs":"sha256:8028407c3116a3c6afc1ca8d2e83f46bd31901d19cf2f1cb6a029c0245ec796e","body":"sha256:0f3e40d48bb353867d02afce91b802eb140e05f9780048c85b980781400383f9"}
func TestStripMarkers(t *testing.T) {
	in := "package p\n\n// sedum:anchor:decls\n// sedum:setFunction {\"a\":1}\nfunc f() {}\n\n// /sedum:setFunction\n"
	want := "package p\n\nfunc f() {}\n\n"
	got := StripMarkers(genpkg.Comment{Prefix: "//"}, in)
	if got != want {
		t.Errorf("StripMarkers() = %q, want %q", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/verify_test.go","name":"TestStripMarkersKeepsIndentedMarkerLookalikes"},"kwargs":"sha256:2489ce53c5383953c13c21e77472eddb6f82f3729344ae3002fa8b00cc63e229","body":"sha256:9529d09aa73b3b5696460a25a91598f4cf16cedf2dec5ba7edb6693b458a8301"}
func TestStripMarkersKeepsIndentedMarkerLookalikes(t *testing.T) {
	in := "\tx := \"// sedum:\"\n\t// sedum:setFunction {}\n"
	want := "\tx := \"// sedum:\"\n"
	got := StripMarkers(genpkg.Comment{Prefix: "//"}, in)
	if got != want {
		t.Errorf("StripMarkers() = %q, want %q", got, want)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/verify_test.go","name":"TestVerifyIgnoresBlankLines"},"kwargs":"sha256:6ff275fe09cc78a06647f33f8b25ea355b2fe6101de041c245617a0881b90d8f","body":"sha256:ead2fd7a78ca27f3abb515a00f634aa4be4e2ce9bb9affbb2e44c68e169ace14"}
func TestVerifyIgnoresBlankLines(t *testing.T) {
	original := "a\nb\n"
	adopted := "a\n\n// sedum:x {}\nb\n// /sedum:x\n\n"
	if err := Verify(genpkg.Comment{Prefix: "//"}, original, adopted); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/verify_test.go","name":"TestVerifyReportsFirstDifference"},"kwargs":"sha256:a5fc1233943acde0e80024a528400b4f16bcf1ffe64d50e9a4735060b1691279","body":"sha256:6c04be1b2fb992802b4d424c23191501fa8072f801b0b4d8b62b3b611123e525"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/verify_test.go","name":"TestVerifyLengthMismatch"},"kwargs":"sha256:d774484abc91c2215f3b464ce499af00dac828935908e8bff8b0a24391d2cb8c","body":"sha256:75ebc327278ba3ae0d744009978e6209de2f11b223959c2eaead903881258c69"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/verify_test.go","name":"TestVerifyLeadingWhitespaceMatters"},"kwargs":"sha256:f6b05e23a0906bc8f8a9c8bb025ed2dd18ff51a47e6be192860aec9618b636f5","body":"sha256:4f37541324baeb2372c813563cf58b4b3b617e3133495a250695383ec653c914"}
func TestVerifyLeadingWhitespaceMatters(t *testing.T) {
	if err := Verify(genpkg.Comment{Prefix: "//"}, "a\n\tb\n", "a\n    b\n"); err == nil {
		t.Error("expected a difference for changed leading whitespace")
	}
}

// /sedum:setFunction
