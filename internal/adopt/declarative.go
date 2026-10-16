package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative.go","path":"regexp"}}
import "regexp"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type Declarative struct {\n\tactions   map[string]string\n\tdeclared  map[string]map[string]bool\n\tunitStart *regexp.Regexp\n\theader    *regexp.Regexp\n\tattach    []*regexp.Regexp\n\tunitEnd   *regexp.Regexp\n\tendMode   string\n}","file":"internal/adopt/declarative.go","name":"Declarative"}}
type Declarative struct {
	actions   map[string]string
	declared  map[string]map[string]bool
	unitStart *regexp.Regexp
	header    *regexp.Regexp
	attach    []*regexp.Regexp
	unitEnd   *regexp.Regexp
	endMode   string
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type unitSpan struct {\n\tfrom, end int\n\taction    string\n\tname      string\n}","file":"internal/adopt/declarative.go","name":"unitSpan"}}
type unitSpan struct {
	from, end int
	action    string
	name      string
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func NewDeclarative(pkg *genpkg.Package) (*Declarative, error) {\n\tif pkg == nil || pkg.Adopt == nil {\n\t\treturn nil, fmt.Errorf(\"declarative adopter: package has no adopt section\")\n\t}\n\ta := pkg.Adopt\n\tif len(a.Command) > 0 {\n\t\treturn nil, fmt.Errorf(\"declarative adopter: package uses a command adopter\")\n\t}\n\td := &Declarative{actions: a.Actions, declared: map[string]map[string]bool{}, endMode: a.UnitEnd}\n\tfor name, act := range pkg.Actions {\n\t\tif act == nil {\n\t\t\tcontinue\n\t\t}\n\t\tset := map[string]bool{}\n\t\tfor k := range act.Kwargs {\n\t\t\tset[k] = true\n\t\t}\n\t\td.declared[name] = set\n\t}\n\tvar err error\n\tif d.unitStart, err = compileRe(\"unit start\", a.UnitStart); err != nil {\n\t\treturn nil, err\n\t}\n\tif a.Header != \"\" {\n\t\tif d.header, err = compileRe(\"header\", a.Header); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t}\n\tfor _, p := range a.AttachAbove {\n\t\tre, err := compileRe(\"attach above\", p)\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\td.attach = append(d.attach, re)\n\t}\n\tif a.UnitEnd != \"dedent\" && a.UnitEnd != \"blank\" {\n\t\tif d.unitEnd, err = compileRe(\"unit end\", a.UnitEnd); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t}\n\treturn d, nil\n}","file":"internal/adopt/declarative.go","name":"NewDeclarative"}}
func NewDeclarative(pkg *genpkg.Package) (*Declarative, error) {
	if pkg == nil || pkg.Adopt == nil {
		return nil, fmt.Errorf("declarative adopter: package has no adopt section")
	}
	a := pkg.Adopt
	if len(a.Command) > 0 {
		return nil, fmt.Errorf("declarative adopter: package uses a command adopter")
	}
	d := &Declarative{actions: a.Actions, declared: map[string]map[string]bool{}, endMode: a.UnitEnd}
	for name, act := range pkg.Actions {
		if act == nil {
			continue
		}
		set := map[string]bool{}
		for k := range act.Kwargs {
			set[k] = true
		}
		d.declared[name] = set
	}
	var err error
	if d.unitStart, err = compileRe("unit start", a.UnitStart); err != nil {
		return nil, err
	}
	if a.Header != "" {
		if d.header, err = compileRe("header", a.Header); err != nil {
			return nil, err
		}
	}
	for _, p := range a.AttachAbove {
		re, err := compileRe("attach above", p)
		if err != nil {
			return nil, err
		}
		d.attach = append(d.attach, re)
	}
	if a.UnitEnd != "dedent" && a.UnitEnd != "blank" {
		if d.unitEnd, err = compileRe("unit end", a.UnitEnd); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func compileRe(what, expr string) (*regexp.Regexp, error) {\n\tre, err := regexp.Compile(expr)\n\tif err != nil {\n\t\treturn nil, fmt.Errorf(\"declarative adopter: %s regexp: %w\", what, err)\n\t}\n\treturn re, nil\n}","file":"internal/adopt/declarative.go","name":"compileRe"}}
func compileRe(what, expr string) (*regexp.Regexp, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("declarative adopter: %s regexp: %w", what, err)
	}
	return re, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func splitLines(s string) []string {\n\tparts := strings.SplitAfter(s, \"\\n\")\n\tif len(parts) > 0 && parts[len(parts)-1] == \"\" {\n\t\tparts = parts[:len(parts)-1]\n\t}\n\treturn parts\n}","file":"internal/adopt/declarative.go","name":"splitLines"}}
func splitLines(s string) []string {
	parts := strings.SplitAfter(s, "\n")
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (d *Declarative) Split(path string, src []byte) ([]Span, error) {\n\ttext := string(src)\n\tspans := []Span{}\n\tstart := 0\n\tif d.header != nil {\n\t\tif loc := d.header.FindStringIndex(text); loc != nil && loc[0] == 0 {\n\t\t\tspans = append(spans, Span{Kind: KindHeader, Text: text[:loc[1]]})\n\t\t\tstart = loc[1]\n\t\t}\n\t}\n\tlines := splitLines(text[start:])\n\tvar units []unitSpan\n\tprevEnd := 0\n\tfor i := 0; i < len(lines); i++ {\n\t\tkind, name, ok := d.startOf(lines[i])\n\t\tif !ok {\n\t\t\tcontinue\n\t\t}\n\t\taction, has := d.actions[kind]\n\t\tif !has {\n\t\t\tcontinue\n\t\t}\n\t\tfrom := i\n\t\tfor from > prevEnd && d.attachable(lines[from-1]) {\n\t\t\tfrom--\n\t\t}\n\t\tend := d.endOf(lines, i)\n\t\tunits = append(units, unitSpan{from: from, end: end, action: action, name: name})\n\t\tprevEnd = end\n\t\ti = end - 1\n\t}\n\tvar resid strings.Builder\n\tflush := func() {\n\t\tif resid.Len() > 0 {\n\t\t\tspans = append(spans, Span{Kind: KindResidual, Text: resid.String()})\n\t\t\tresid.Reset()\n\t\t}\n\t}\n\tpos := 0\n\tfor _, u := range units {\n\t\tfor ; pos < u.from; pos++ {\n\t\t\tresid.WriteString(lines[pos])\n\t\t}\n\t\tflush()\n\t\tbody := strings.Join(lines[u.from:u.end], \"\")\n\t\tspans = append(spans, Span{Kind: KindUnit, Text: body, Action: u.action, Kwargs: d.kwargs(u.action, path, u.name, strings.TrimSuffix(body, \"\\n\"))})\n\t\tpos = u.end\n\t}\n\tfor ; pos < len(lines); pos++ {\n\t\tresid.WriteString(lines[pos])\n\t}\n\tflush()\n\treturn spans, nil\n}","file":"internal/adopt/declarative.go","name":"Split","receiver":"Declarative"}}
func (d *Declarative) Split(path string, src []byte) ([]Span, error) {
	text := string(src)
	spans := []Span{}
	start := 0
	if d.header != nil {
		if loc := d.header.FindStringIndex(text); loc != nil && loc[0] == 0 {
			spans = append(spans, Span{Kind: KindHeader, Text: text[:loc[1]]})
			start = loc[1]
		}
	}
	lines := splitLines(text[start:])
	var units []unitSpan
	prevEnd := 0
	for i := 0; i < len(lines); i++ {
		kind, name, ok := d.startOf(lines[i])
		if !ok {
			continue
		}
		action, has := d.actions[kind]
		if !has {
			continue
		}
		from := i
		for from > prevEnd && d.attachable(lines[from-1]) {
			from--
		}
		end := d.endOf(lines, i)
		units = append(units, unitSpan{from: from, end: end, action: action, name: name})
		prevEnd = end
		i = end - 1
	}
	var resid strings.Builder
	flush := func() {
		if resid.Len() > 0 {
			spans = append(spans, Span{Kind: KindResidual, Text: resid.String()})
			resid.Reset()
		}
	}
	pos := 0
	for _, u := range units {
		for ; pos < u.from; pos++ {
			resid.WriteString(lines[pos])
		}
		flush()
		body := strings.Join(lines[u.from:u.end], "")
		spans = append(spans, Span{Kind: KindUnit, Text: body, Action: u.action, Kwargs: d.kwargs(u.action, path, u.name, strings.TrimSuffix(body, "\n"))})
		pos = u.end
	}
	for ; pos < len(lines); pos++ {
		resid.WriteString(lines[pos])
	}
	flush()
	return spans, nil
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (d *Declarative) startOf(line string) (string, string, bool) {\n\tm := d.unitStart.FindStringSubmatch(strings.TrimRight(line, \"\\r\\n\"))\n\tif m == nil {\n\t\treturn \"\", \"\", false\n\t}\n\tki := d.unitStart.SubexpIndex(\"kind\")\n\tif ki < 0 {\n\t\treturn \"\", \"\", false\n\t}\n\tname := \"\"\n\tif ni := d.unitStart.SubexpIndex(\"name\"); ni >= 0 {\n\t\tname = m[ni]\n\t}\n\treturn m[ki], name, true\n}","file":"internal/adopt/declarative.go","name":"startOf","receiver":"Declarative"}}
func (d *Declarative) startOf(line string) (string, string, bool) {
	m := d.unitStart.FindStringSubmatch(strings.TrimRight(line, "\r\n"))
	if m == nil {
		return "", "", false
	}
	ki := d.unitStart.SubexpIndex("kind")
	if ki < 0 {
		return "", "", false
	}
	name := ""
	if ni := d.unitStart.SubexpIndex("name"); ni >= 0 {
		name = m[ni]
	}
	return m[ki], name, true
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (d *Declarative) attachable(line string) bool {\n\tl := strings.TrimRight(line, \"\\r\\n\")\n\tfor _, re := range d.attach {\n\t\tif re.MatchString(l) {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}","file":"internal/adopt/declarative.go","name":"attachable","receiver":"Declarative"}}
func (d *Declarative) attachable(line string) bool {
	l := strings.TrimRight(line, "\r\n")
	for _, re := range d.attach {
		if re.MatchString(l) {
			return true
		}
	}
	return false
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (d *Declarative) endOf(lines []string, i int) int {\n\tswitch d.endMode {\n\tcase \"dedent\":\n\t\tend := i + 1\n\t\tfor j := i + 1; j < len(lines); j++ {\n\t\t\tl := strings.TrimRight(lines[j], \"\\r\\n\")\n\t\t\tif strings.TrimSpace(l) == \"\" {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tif l[0] != ' ' && l[0] != '\\t' {\n\t\t\t\tbreak\n\t\t\t}\n\t\t\tend = j + 1\n\t\t}\n\t\treturn end\n\tcase \"blank\":\n\t\tj := i + 1\n\t\tfor j < len(lines) && strings.TrimSpace(lines[j]) != \"\" {\n\t\t\tj++\n\t\t}\n\t\treturn j\n\t}\n\tfor j := i; j < len(lines); j++ {\n\t\tif d.unitEnd.MatchString(strings.TrimRight(lines[j], \"\\r\\n\")) {\n\t\t\treturn j + 1\n\t\t}\n\t}\n\treturn len(lines)\n}","file":"internal/adopt/declarative.go","name":"endOf","receiver":"Declarative"}}
func (d *Declarative) endOf(lines []string, i int) int {
	switch d.endMode {
	case "dedent":
		end := i + 1
		for j := i + 1; j < len(lines); j++ {
			l := strings.TrimRight(lines[j], "\r\n")
			if strings.TrimSpace(l) == "" {
				continue
			}
			if l[0] != ' ' && l[0] != '\t' {
				break
			}
			end = j + 1
		}
		return end
	case "blank":
		j := i + 1
		for j < len(lines) && strings.TrimSpace(lines[j]) != "" {
			j++
		}
		return j
	}
	for j := i; j < len(lines); j++ {
		if d.unitEnd.MatchString(strings.TrimRight(lines[j], "\r\n")) {
			return j + 1
		}
	}
	return len(lines)
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:12e7118d99ed19c20430ed17e0d4592c3ec126f7a93f9e496f349a39e5e89a4d","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (d *Declarative) kwargs(action, path, name, code string) map[string]any {\n\tout := map[string]any{}\n\tdecl := d.declared[action]\n\tif decl[\"file\"] {\n\t\tout[\"file\"] = path\n\t}\n\tif decl[\"name\"] {\n\t\tout[\"name\"] = name\n\t}\n\tif decl[\"code\"] {\n\t\tout[\"code\"] = code\n\t}\n\treturn out\n}","file":"internal/adopt/declarative.go","name":"kwargs","receiver":"Declarative"}}
func (d *Declarative) kwargs(action, path, name, code string) map[string]any {
	out := map[string]any{}
	decl := d.declared[action]
	if decl["file"] {
		out["file"] = path
	}
	if decl["name"] {
		out["name"] = name
	}
	if decl["code"] {
		out["code"] = code
	}
	return out
}

// /sedum:setMethod
