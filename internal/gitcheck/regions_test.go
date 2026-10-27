package gitcheck

import (
	"reflect"
	"strings"
	"testing"
)

// The region reader splits a file into the text outside regions and the
// regions, and compares two such splits (the region reader record of
// sedum check).

func open(leader, action, tier, record, gen, pkg string) string {
	return leader + " sedum:" + action + ` {"tier":"` + tier + `","record":"` + record + `","gen":"` + gen + `","pkg":"` + pkg + `"}`
}

func closeLine(leader, action string) string {
	return leader + " /sedum:" + action
}

func region(action, record, gen, pkg, body string) string {
	return open("//", action, "owned", record, gen, pkg) + "\n" + body + closeLine("//", action) + "\n"
}

func TestReadRegionsReadsAnOpeningLineWithEitherLeaderAndIndent(t *testing.T) {
	for _, leader := range []string{"//", "#"} {
		src := "head\n  \t" + open(leader, "addField", "owned", "prov-1", "g1", "p1") + "\nbody\n" + closeLine(leader, "addField") + "\ntail\n"

		outside, regions, ok := ReadRegions([]byte(src))

		if !ok || len(regions) != 1 {
			t.Fatalf("%s: ok=%v regions=%+v", leader, ok, regions)
		}
		want := Region{Tier: "owned", Record: "prov-1", Gen: "g1", Pkg: "p1",
			Prefix: "  \t" + leader + " sedum:addField ", Body: "body\n", Close: closeLine(leader, "addField")}
		if regions[0] != want {
			t.Errorf("%s: region = %+v, want %+v", leader, regions[0], want)
		}
		if !strings.HasPrefix(outside, "head\n") || !strings.HasSuffix(outside, "\ntail\n") || strings.Contains(outside, "body") {
			t.Errorf("%s: outside = %q", leader, outside)
		}
	}
}

func TestTheClosingLineMayBeIndentedDifferentlyAndTrailingWhitespaceIsIgnored(t *testing.T) {
	src := open("//", "a", "owned", "r", "g", "p") + " \t\r\nx\n\t" + closeLine("//", "a") + " \t\r\n"

	_, regions, ok := ReadRegions([]byte(src))

	if !ok || len(regions) != 1 || regions[0].Body != "x\n" {
		t.Fatalf("ok=%v regions=%+v", ok, regions)
	}
	if regions[0].Close != "\t// /sedum:a \t\r" {
		t.Errorf("Close = %q, want the line as it is in the file", regions[0].Close)
	}
}

func TestAnUnclosedRegionMakesTheFileUnreadable(t *testing.T) {
	if _, _, ok := ReadRegions([]byte(open("//", "a", "owned", "r", "g", "p") + "\nx\n")); ok {
		t.Error("an unclosed region was read")
	}
}

func TestARegionOpenedInsideAnotherMakesTheFileUnreadable(t *testing.T) {
	src := open("//", "a", "owned", "r", "g", "p") + "\n" + open("//", "b", "owned", "r2", "g", "p") + "\n" + closeLine("//", "b") + "\n" + closeLine("//", "a") + "\n"
	if _, _, ok := ReadRegions([]byte(src)); ok {
		t.Error("a nested region was read")
	}
}

func TestAClosingLineNeedsTheSameLeaderAndActionAsTheOpeningLine(t *testing.T) {
	src := open("//", "a", "owned", "r", "g", "p") + "\n" + closeLine("#", "a") + "\n" + closeLine("//", "b") + "\n"
	if _, _, ok := ReadRegions([]byte(src)); ok {
		t.Error("a closing line with another leader or action closed the region")
	}
}

func TestAnOpeningLineThatDoesNotParseIsTextOutsideRegions(t *testing.T) {
	for name, line := range map[string]string{
		"no json":        "// sedum:addField",
		"bad json":       `// sedum:addField {"tier":"owned",`,
		"missing key":    `// sedum:addField {"tier":"owned","record":"r","gen":"g"}`,
		"non string key": `// sedum:addField {"tier":"owned","record":"r","gen":1,"pkg":"p"}`,
		"empty action":   `// sedum: {"tier":"owned","record":"r","gen":"g","pkg":"p"}`,
		"no space":       `//sedum:a {"tier":"owned","record":"r","gen":"g","pkg":"p"}`,
	} {
		outside, regions, ok := ReadRegions([]byte(line + "\n"))
		if !ok || len(regions) != 0 || outside != line+"\n" {
			t.Errorf("%s: ok=%v regions=%+v outside=%q", name, ok, regions, outside)
		}
	}
}

func TestOtherJSONKeysAreAcceptedOnAnOpeningLine(t *testing.T) {
	line := `// sedum:setConst {"v":2,"tier":"owned","record":"r","gen":"g","pkg":"p","id":{"file":"a.go","name":"X"},"kwargs":"sha256:1","body":"sha256:2"}`
	_, regions, ok := ReadRegions([]byte(line + "\nx\n// /sedum:setConst\n"))
	if !ok || len(regions) != 1 || regions[0].Record != "r" {
		t.Fatalf("ok=%v regions=%+v", ok, regions)
	}
}

func TestTheTextOutsideRegionsReplacesEachRegionByOneFixedLine(t *testing.T) {
	a := "a\n" + region("x", "r", "g", "p", "one\ntwo\n") + "b\n"
	b := "a\n" + region("y", "r2", "g2", "p2", "other\n") + "b\n"
	outsideA, _, _ := ReadRegions([]byte(a))
	outsideB, _, _ := ReadRegions([]byte(b))
	if outsideA != outsideB {
		t.Errorf("outside differs: %q vs %q", outsideA, outsideB)
	}
}

func TestSameShape(t *testing.T) {
	base := "top\n" + region("a", "r1", "g1", "p1", "old\n") + "mid\n" + region("b", "r2", "g2", "p2", "old2\n") + "end\n"
	cases := []struct {
		name          string
		commit        string
		loose, strict bool
	}{
		{"identical", base, true, true},
		{"body changes", strings.Replace(base, "old\n", "new\n", 1), true, false},
		{"pkg changes", strings.Replace(base, `"pkg":"p1"`, `"pkg":"p9"`, 1), true, false},
		{"another JSON key changes", strings.Replace(base, `"pkg":"p1"}`, `"pkg":"p1","kwargs":"x"}`, 1), true, true},
		{"prefix changes", strings.Replace(base, "// sedum:a ", "  // sedum:a ", 1), true, false},
		{"closing line changes", strings.Replace(base, "// /sedum:a\n", "// /sedum:a  \n", 1), true, false},
		{"gen changes", strings.Replace(base, `"gen":"g1"`, `"gen":"g9"`, 1), false, false},
		{"record changes", strings.Replace(base, `"record":"r1"`, `"record":"r9"`, 1), false, false},
		{"tier changes", strings.Replace(base, `"tier":"owned"`, `"tier":"seeded"`, 1), false, false},
		{"outside text changes by one byte", strings.Replace(base, "mid\n", "mid!\n", 1), false, false},
		{"a region is dropped", "top\n" + region("a", "r1", "g1", "p1", "old\n") + "mid\nend\n", false, false},
		{"the order changes", "top\n" + region("b", "r2", "g2", "p2", "old2\n") + "mid\n" + region("a", "r1", "g1", "p1", "old\n") + "end\n", false, false},
		{"unreadable commit file", base + open("//", "c", "owned", "r", "g", "p") + "\n", false, false},
	}
	for _, c := range cases {
		if got := SameShape([]byte(base), []byte(c.commit), false); got != c.loose {
			t.Errorf("%s: SameShape(strict=false) = %v, want %v", c.name, got, c.loose)
		}
		if got := SameShape([]byte(base), []byte(c.commit), true); got != c.strict {
			t.Errorf("%s: SameShape(strict=true) = %v, want %v", c.name, got, c.strict)
		}
	}
}

func TestSameShapeIsFalseWhenTheBaseHasNoRegionOrIsUnreadable(t *testing.T) {
	if SameShape([]byte("plain\n"), []byte("plain\n"), false) {
		t.Error("a base with no region agreed")
	}
	unreadable := open("//", "a", "owned", "r", "g", "p") + "\n"
	if SameShape([]byte(unreadable), []byte(unreadable), false) {
		t.Error("an unreadable base agreed")
	}
}

func TestReadRegionsOfAFileWithoutRegionsReturnsNoRegions(t *testing.T) {
	outside, regions, ok := ReadRegions([]byte("a\nb"))
	if !ok || outside != "a\nb" || !reflect.DeepEqual(regions, []Region{}) {
		t.Errorf("outside=%q regions=%#v ok=%v", outside, regions, ok)
	}
}
