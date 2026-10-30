package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_test.go","path":"fmt"},"kwargs":"sha256:09eade496c83201d5066d7b21fc1a88a3ae0b129e4a6512204ff0c75c9ca3c6b","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_test.go","path":"strings"},"kwargs":"sha256:49c71df4c9cfc9a40680649cf295958b8978bd7cade2770dd0b75f18616f9999","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter_test.go","path":"testing"},"kwargs":"sha256:35d298a2250f83bbe67b412901ee769c562c8153d07828436c2ef74f92537930","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter_test.go","name":"goSummary"},"kwargs":"sha256:94f3257a463dc953a7a669c80e0544bf71c24d2788aaa9f006834d585923625f","body":"sha256:56fd5bafb1d6f03562a09d2142a0486cf3ae994706864ac9ec8b410d6f83ceb8"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter_test.go","name":"TestGoAdopterSplit"},"kwargs":"sha256:2a2e56fd662b09f5c18366a2ffe8f7e6353ebd8e22bc57bde44805ed1a3f3c7b","body":"sha256:a51b7bbdf13fdcbc00e361f1eb1fe6a31101c68e3f1c97b3e53d16dcd9829eeb"}
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
			spans, err := GoAdopter{}.Split("p/x.go", []byte(tc.src))
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter_test.go","name":"TestGoAdopterKwargs"},"kwargs":"sha256:c1358fde40e55d477805e9ca112786e07428c09e14dd3a75d2b9000c11bfe59c","body":"sha256:67921f3bcb162e413ee4017fb9b3da769567ab443b4a34a2cb36beaf7cbdd439"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter_test.go","name":"TestGoAdopterRejects"},"kwargs":"sha256:2aa3729198ff0fcb4dfaf5f1b0c8ef958e7302126619f60ef41b384be9dbd1a1","body":"sha256:41a3543416050cc864e49c476d57da72349f0708ab6e2a13d0bb1082d939a453"}
func TestGoAdopterRejects(t *testing.T) {
	for name, src := range map[string]string{
		"not gofmt-clean": "package x\nfunc  F() {}\n",
		"parse error":     "package x\nfunc (",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := (GoAdopter{}).Split("p/x.go", []byte(src)); err == nil {
				t.Fatal("want an error")
			}
		})
	}
	_, err := GoAdopter{}.Split("p/x.go", []byte("package x\nfunc  F() {}\n"))
	if err == nil || !strings.Contains(err.Error(), "not gofmt-clean") {
		t.Fatalf("err = %v, want not gofmt-clean", err)
	}
}

// /sedum:setFunction
