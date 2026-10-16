package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:a782092fec4bcd6da1f8a0bbb80ea5d0f1e05bfedcdd28c25bffe608fa122e91","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_test.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:a782092fec4bcd6da1f8a0bbb80ea5d0f1e05bfedcdd28c25bffe608fa122e91","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_test.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:a782092fec4bcd6da1f8a0bbb80ea5d0f1e05bfedcdd28c25bffe608fa122e91","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/goadopter_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:a782092fec4bcd6da1f8a0bbb80ea5d0f1e05bfedcdd28c25bffe608fa122e91","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// goSummary renders spans as one line each: kind, then action and name for units.\nfunc goSummary(spans []Span) string {\n\tvar lines []string\n\tfor _, s := range spans {\n\t\tline := string(s.Kind)\n\t\tif s.Kind == KindUnit {\n\t\t\tline += \" \" + s.Action\n\t\t\tif n, ok := s.Kwargs[\"name\"]; ok {\n\t\t\t\tline += \" \" + fmt.Sprint(n)\n\t\t\t}\n\t\t\tif r, ok := s.Kwargs[\"receiver\"]; ok {\n\t\t\t\tline += \" recv=\" + fmt.Sprint(r)\n\t\t\t}\n\t\t\tif p, ok := s.Kwargs[\"path\"]; ok {\n\t\t\t\tline += \" path=\" + fmt.Sprint(p)\n\t\t\t}\n\t\t}\n\t\tlines = append(lines, line)\n\t}\n\treturn strings.Join(lines, \"\\n\")\n}","file":"internal/adopt/goadopter_test.go","name":"goSummary"}}
// goSummary renders spans as one line each: kind, then action and name for units.
func goSummary(spans []Span) string {
	var lines []string
	for _, s := range spans {
		line := string(s.Kind)
		if s.Kind == KindUnit {
			line += " " + s.Action
			if n, ok := s.Kwargs["name"]; ok {
				line += " " + fmt.Sprint(n)
			}
			if r, ok := s.Kwargs["receiver"]; ok {
				line += " recv=" + fmt.Sprint(r)
			}
			if p, ok := s.Kwargs["path"]; ok {
				line += " path=" + fmt.Sprint(p)
			}
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:a782092fec4bcd6da1f8a0bbb80ea5d0f1e05bfedcdd28c25bffe608fa122e91","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoAdopterSplit(t *testing.T) {\n\tcases := []struct {\n\t\tname, src, want string\n\t}{\n\t\t{\"declarations\",\n\t\t\t\"package x\\n\\nimport \\\"b\\\"\\n\\n// F doc.\\nfunc F() {}\\n\\nfunc (t *T[K]) M() {}\\n\\ntype T struct{}\\n\",\n\t\t\t\"header\\nresidual\\nunit addImport path=b\\nresidual\\nunit setFunction F\\nresidual\\nunit setMethod M recv=T\\nresidual\\nunit setType T\"},\n\t\t{\"grouped import and block\",\n\t\t\t\"package x\\n\\nimport (\\n\\t\\\"a\\\"\\n\\t\\\"b\\\"\\n)\\n\\nconst (\\n\\tA = iota\\n\\tB\\n)\\n\",\n\t\t\t\"header\\nresidual\\nunit setImports imports#1\\nresidual\\nunit setConst A\"},\n\t\t{\"aliased import\",\n\t\t\t\"package x\\n\\nimport f \\\"fmt\\\"\\n\",\n\t\t\t\"header\\nresidual\\nunit setImports imports#1\"},\n\t\t{\"repeated names\",\n\t\t\t\"package x\\n\\nfunc init() {}\\n\\nfunc init() {}\\n\\nvar _ = 1\\nvar _ = 2\\n\",\n\t\t\t\"header\\nresidual\\nunit setFunction init#1\\nresidual\\nunit setFunction init#2\\nresidual\\nunit setVar _\\nunit setVar _#2\"},\n\t\t{\"package doc is a residual before the header\",\n\t\t\t\"// Package x does things.\\npackage x\\n\",\n\t\t\t\"residual\\nheader\"},\n\t\t{\"floating comment\",\n\t\t\t\"package x\\n\\n// floating\\n\\nfunc F() {}\\n\",\n\t\t\t\"header\\nresidual\\nunit setFunction F\"},\n\t\t{\"cgo import\",\n\t\t\t\"package x\\n\\n// #include <stdio.h>\\nimport \\\"C\\\"\\n\",\n\t\t\t\"header\\nresidual\\nresidual\"},\n\t}\n\tfor _, tc := range cases {\n\t\tt.Run(tc.name, func(t *testing.T) {\n\t\t\tspans, err := GoAdopter{}.Split(\"x.go\", []byte(tc.src))\n\t\t\tif err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tif err := Coverage([]byte(tc.src), spans); err != nil {\n\t\t\t\tt.Fatal(err)\n\t\t\t}\n\t\t\tif got := goSummary(spans); got != tc.want {\n\t\t\t\tt.Fatalf(\"spans:\\n%s\\nwant:\\n%s\", got, tc.want)\n\t\t\t}\n\t\t})\n\t}\n}","file":"internal/adopt/goadopter_test.go","name":"TestGoAdopterSplit"}}
func TestGoAdopterSplit(t *testing.T) {
	cases := []struct {
		name, src, want string
	}{
		{"declarations",
			"package x\n\nimport \"b\"\n\n// F doc.\nfunc F() {}\n\nfunc (t *T[K]) M() {}\n\ntype T struct{}\n",
			"header\nresidual\nunit addImport path=b\nresidual\nunit setFunction F\nresidual\nunit setMethod M recv=T\nresidual\nunit setType T"},
		{"grouped import and block",
			"package x\n\nimport (\n\t\"a\"\n\t\"b\"\n)\n\nconst (\n\tA = iota\n\tB\n)\n",
			"header\nresidual\nunit setImports imports#1\nresidual\nunit setConst A"},
		{"aliased import",
			"package x\n\nimport f \"fmt\"\n",
			"header\nresidual\nunit setImports imports#1"},
		{"repeated names",
			"package x\n\nfunc init() {}\n\nfunc init() {}\n\nvar _ = 1\nvar _ = 2\n",
			"header\nresidual\nunit setFunction init#1\nresidual\nunit setFunction init#2\nresidual\nunit setVar _\nunit setVar _#2"},
		{"package doc is a residual before the header",
			"// Package x does things.\npackage x\n",
			"residual\nheader"},
		{"floating comment",
			"package x\n\n// floating\n\nfunc F() {}\n",
			"header\nresidual\nunit setFunction F"},
		{"cgo import",
			"package x\n\n// #include <stdio.h>\nimport \"C\"\n",
			"header\nresidual\nresidual"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spans, err := GoAdopter{}.Split("x.go", []byte(tc.src))
			if err != nil {
				t.Fatal(err)
			}
			if err := Coverage([]byte(tc.src), spans); err != nil {
				t.Fatal(err)
			}
			if got := goSummary(spans); got != tc.want {
				t.Fatalf("spans:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:a782092fec4bcd6da1f8a0bbb80ea5d0f1e05bfedcdd28c25bffe608fa122e91","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoAdopterKwargs(t *testing.T) {\n\tsrc := \"package x\\n\\n// F doc.\\nfunc F() {}\\n\"\n\tspans, err := GoAdopter{}.Split(\"p/x.go\", []byte(src))\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif spans[0].Kind != KindHeader || spans[0].Text != \"package x\\n\" {\n\t\tt.Fatalf(\"header = %+v\", spans[0])\n\t}\n\tu := spans[2]\n\tif u.Kwargs[\"file\"] != \"p/x.go\" || u.Kwargs[\"code\"] != \"// F doc.\\nfunc F() {}\" || u.Text != \"// F doc.\\nfunc F() {}\\n\" {\n\t\tt.Fatalf(\"unit = %+v\", u)\n\t}\n}","file":"internal/adopt/goadopter_test.go","name":"TestGoAdopterKwargs"}}
func TestGoAdopterKwargs(t *testing.T) {
	src := "package x\n\n// F doc.\nfunc F() {}\n"
	spans, err := GoAdopter{}.Split("p/x.go", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if spans[0].Kind != KindHeader || spans[0].Text != "package x\n" {
		t.Fatalf("header = %+v", spans[0])
	}
	u := spans[2]
	if u.Kwargs["file"] != "p/x.go" || u.Kwargs["code"] != "// F doc.\nfunc F() {}" || u.Text != "// F doc.\nfunc F() {}\n" {
		t.Fatalf("unit = %+v", u)
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:a782092fec4bcd6da1f8a0bbb80ea5d0f1e05bfedcdd28c25bffe608fa122e91","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGoAdopterRejects(t *testing.T) {\n\tfor name, src := range map[string]string{\n\t\t\"not gofmt-clean\": \"package x\\nfunc  F() {}\\n\",\n\t\t\"parse error\":     \"package x\\nfunc (\",\n\t} {\n\t\tt.Run(name, func(t *testing.T) {\n\t\t\tif _, err := (GoAdopter{}).Split(\"x.go\", []byte(src)); err == nil {\n\t\t\t\tt.Fatal(\"want an error\")\n\t\t\t}\n\t\t})\n\t}\n\t_, err := GoAdopter{}.Split(\"x.go\", []byte(\"package x\\nfunc  F() {}\\n\"))\n\tif err == nil || !strings.Contains(err.Error(), \"not gofmt-clean\") {\n\t\tt.Fatalf(\"err = %v, want not gofmt-clean\", err)\n\t}\n}","file":"internal/adopt/goadopter_test.go","name":"TestGoAdopterRejects"}}
func TestGoAdopterRejects(t *testing.T) {
	for name, src := range map[string]string{
		"not gofmt-clean": "package x\nfunc  F() {}\n",
		"parse error":     "package x\nfunc (",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := (GoAdopter{}).Split("x.go", []byte(src)); err == nil {
				t.Fatal("want an error")
			}
		})
	}
	_, err := GoAdopter{}.Split("x.go", []byte("package x\nfunc  F() {}\n"))
	if err == nil || !strings.Contains(err.Error(), "not gofmt-clean") {
		t.Fatalf("err = %v, want not gofmt-clean", err)
	}
}

// /sedum:setFunction
