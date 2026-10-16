package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative.go","path":"regexp"}}
import "regexp"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/declarative.go","path":"github.com/livecodelife/sedum/internal/genpkg"}}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Declarative splits files by the rules in a package's adopt section.\ntype Declarative struct {\n\theader    *regexp.Regexp\n\tunitStart *regexp.Regexp\n\tattach    []*regexp.Regexp\n\tunitEnd   *regexp.Regexp\n\tendMode   string\n\tkinds     map[string]string\n\tdeclared  map[string]map[string]bool\n}","file":"internal/adopt/declarative.go","name":"Declarative"}}
// Declarative splits files by the rules in a package's adopt section.
type Declarative struct {
	header    *regexp.Regexp
	unitStart *regexp.Regexp
	attach    []*regexp.Regexp
	unitEnd   *regexp.Regexp
	endMode   string
	kinds     map[string]string
	declared  map[string]map[string]bool
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// NewDeclarative builds a Declarative adopter from a declarative package.\nfunc NewDeclarative(pkg *genpkg.Package) (*Declarative, error) {\n\tif pkg.Adopt == nil {\n\t\treturn nil, fmt.Errorf(\"package %s: no adopt section\", pkg.Name)\n\t}\n\tad := pkg.Adopt\n\tif len(ad.Command) > 0 {\n\t\treturn nil, fmt.Errorf(\"package %s: adopt command is not declarative\", pkg.Name)\n\t}\n\tif ad.UnitEnd == \"\" {\n\t\treturn nil, fmt.Errorf(\"package %s: unit_end is required\", pkg.Name)\n\t}\n\td := &Declarative{\n\t\tendMode:  ad.UnitEnd,\n\t\tkinds:    map[string]string{},\n\t\tdeclared: map[string]map[string]bool{},\n\t}\n\tvar err error\n\tif d.unitStart, err = regexp.Compile(ad.UnitStart); err != nil {\n\t\treturn nil, fmt.Errorf(\"package %s: unit_start: %w\", pkg.Name, err)\n\t}\n\tif ad.Header != \"\" {\n\t\tif d.header, err = regexp.Compile(ad.Header); err != nil {\n\t\t\treturn nil, fmt.Errorf(\"package %s: header: %w\", pkg.Name, err)\n\t\t}\n\t}\n\tfor _, pat := range ad.AttachAbove {\n\t\tre, err := regexp.Compile(pat)\n\t\tif err != nil {\n\t\t\treturn nil, fmt.Errorf(\"package %s: attach_above: %w\", pkg.Name, err)\n\t\t}\n\t\td.attach = append(d.attach, re)\n\t}\n\tif ad.UnitEnd != \"dedent\" && ad.UnitEnd != \"blank\" {\n\t\tif d.unitEnd, err = regexp.Compile(ad.UnitEnd); err != nil {\n\t\t\treturn nil, fmt.Errorf(\"package %s: unit_end: %w\", pkg.Name, err)\n\t\t}\n\t}\n\tfor kind, action := range ad.Actions {\n\t\td.kinds[kind] = action\n\t\tset := map[string]bool{}\n\t\tif a, ok := pkg.Actions[action]; ok && a != nil {\n\t\t\tfor _, k := range []string{\"file\", \"name\", \"code\"} {\n\t\t\t\tif _, ok := a.Kwargs[k]; ok {\n\t\t\t\t\tset[k] = true\n\t\t\t\t}\n\t\t\t}\n\t\t}\n\t\td.declared[action] = set\n\t}\n\treturn d, nil\n}","file":"internal/adopt/declarative.go","name":"NewDeclarative"}}
// NewDeclarative builds a Declarative adopter from a declarative package.
func NewDeclarative(pkg *genpkg.Package) (*Declarative, error) {
	if pkg.Adopt == nil {
		return nil, fmt.Errorf("package %s: no adopt section", pkg.Name)
	}
	ad := pkg.Adopt
	if len(ad.Command) > 0 {
		return nil, fmt.Errorf("package %s: adopt command is not declarative", pkg.Name)
	}
	if ad.UnitEnd == "" {
		return nil, fmt.Errorf("package %s: unit_end is required", pkg.Name)
	}
	d := &Declarative{
		endMode:  ad.UnitEnd,
		kinds:    map[string]string{},
		declared: map[string]map[string]bool{},
	}
	var err error
	if d.unitStart, err = regexp.Compile(ad.UnitStart); err != nil {
		return nil, fmt.Errorf("package %s: unit_start: %w", pkg.Name, err)
	}
	if ad.Header != "" {
		if d.header, err = regexp.Compile(ad.Header); err != nil {
			return nil, fmt.Errorf("package %s: header: %w", pkg.Name, err)
		}
	}
	for _, pat := range ad.AttachAbove {
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("package %s: attach_above: %w", pkg.Name, err)
		}
		d.attach = append(d.attach, re)
	}
	if ad.UnitEnd != "dedent" && ad.UnitEnd != "blank" {
		if d.unitEnd, err = regexp.Compile(ad.UnitEnd); err != nil {
			return nil, fmt.Errorf("package %s: unit_end: %w", pkg.Name, err)
		}
	}
	for kind, action := range ad.Actions {
		d.kinds[kind] = action
		set := map[string]bool{}
		if a, ok := pkg.Actions[action]; ok && a != nil {
			for _, k := range []string{"file", "name", "code"} {
				if _, ok := a.Kwargs[k]; ok {
					set[k] = true
				}
			}
		}
		d.declared[action] = set
	}
	return d, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// splitLines splits s into lines, keeping each line's newline.\nfunc splitLines(s string) []string {\n\tvar out []string\n\tfor len(s) > 0 {\n\t\ti := strings.IndexByte(s, '\\n')\n\t\tif i < 0 {\n\t\t\tout = append(out, s)\n\t\t\tbreak\n\t\t}\n\t\tout = append(out, s[:i+1])\n\t\ts = s[i+1:]\n\t}\n\treturn out\n}","file":"internal/adopt/declarative.go","name":"splitLines"}}
// splitLines splits s into lines, keeping each line's newline.
func splitLines(s string) []string {
	var out []string
	for len(s) > 0 {
		i := strings.IndexByte(s, '\n')
		if i < 0 {
			out = append(out, s)
			break
		}
		out = append(out, s[:i+1])
		s = s[i+1:]
	}
	return out
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Split divides src into header, unit and residual spans in file order.\nfunc (d *Declarative) Split(path string, src []byte) ([]Span, error) {\n\tout := []Span{}\n\ttext := string(src)\n\tif d.header != nil {\n\t\tif loc := d.header.FindStringIndex(text); loc != nil && loc[0] == 0 && loc[1] > 0 {\n\t\t\tout = append(out, Span{Kind: KindHeader, Text: text[:loc[1]]})\n\t\t\ttext = text[loc[1]:]\n\t\t}\n\t}\n\tlines := splitLines(text)\n\tpending := 0\n\tfor i := 0; i < len(lines); {\n\t\tkind, name, ok := d.matchStart(lines[i])\n\t\taction, known := d.kinds[kind]\n\t\tif !ok || !known {\n\t\t\ti++\n\t\t\tcontinue\n\t\t}\n\t\tb := i\n\t\tfor b > pending && d.attaches(lines[b-1]) {\n\t\t\tb--\n\t\t}\n\t\te := d.endOf(lines, i)\n\t\tif b > pending {\n\t\t\tout = append(out, residualSpan(lines[pending:b]))\n\t\t}\n\t\tout = append(out, d.unitSpan(path, action, name, lines[b:e]))\n\t\tpending = e\n\t\ti = e\n\t}\n\tif pending < len(lines) {\n\t\tout = append(out, residualSpan(lines[pending:]))\n\t}\n\treturn out, nil\n}","file":"internal/adopt/declarative.go","name":"Split","receiver":"Declarative"}}
// Split divides src into header, unit and residual spans in file order.
func (d *Declarative) Split(path string, src []byte) ([]Span, error) {
	out := []Span{}
	text := string(src)
	if d.header != nil {
		if loc := d.header.FindStringIndex(text); loc != nil && loc[0] == 0 && loc[1] > 0 {
			out = append(out, Span{Kind: KindHeader, Text: text[:loc[1]]})
			text = text[loc[1]:]
		}
	}
	lines := splitLines(text)
	pending := 0
	for i := 0; i < len(lines); {
		kind, name, ok := d.matchStart(lines[i])
		action, known := d.kinds[kind]
		if !ok || !known {
			i++
			continue
		}
		b := i
		for b > pending && d.attaches(lines[b-1]) {
			b--
		}
		e := d.endOf(lines, i)
		if b > pending {
			out = append(out, residualSpan(lines[pending:b]))
		}
		out = append(out, d.unitSpan(path, action, name, lines[b:e]))
		pending = e
		i = e
	}
	if pending < len(lines) {
		out = append(out, residualSpan(lines[pending:]))
	}
	return out, nil
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// endOf returns the exclusive index where the unit starting at start ends.\nfunc (d *Declarative) endOf(lines []string, start int) int {\n\tswitch d.endMode {\n\tcase \"dedent\":\n\t\te := start + 1\n\t\tfor j := start + 1; j < len(lines); j++ {\n\t\t\tline := lines[j]\n\t\t\tif isBlank(line) {\n\t\t\t\tcontinue\n\t\t\t}\n\t\t\tif line[0] != ' ' && line[0] != '\\t' {\n\t\t\t\tbreak\n\t\t\t}\n\t\t\te = j + 1\n\t\t}\n\t\treturn e\n\tcase \"blank\":\n\t\te := start + 1\n\t\tfor e < len(lines) && !isBlank(lines[e]) {\n\t\t\te++\n\t\t}\n\t\treturn e\n\t}\n\tfor j := start; j < len(lines); j++ {\n\t\tif d.unitEnd.MatchString(stripEOL(lines[j])) {\n\t\t\treturn j + 1\n\t\t}\n\t}\n\treturn len(lines)\n}","file":"internal/adopt/declarative.go","name":"endOf","receiver":"Declarative"}}
// endOf returns the exclusive index where the unit starting at start ends.
func (d *Declarative) endOf(lines []string, start int) int {
	switch d.endMode {
	case "dedent":
		e := start + 1
		for j := start + 1; j < len(lines); j++ {
			line := lines[j]
			if isBlank(line) {
				continue
			}
			if line[0] != ' ' && line[0] != '\t' {
				break
			}
			e = j + 1
		}
		return e
	case "blank":
		e := start + 1
		for e < len(lines) && !isBlank(lines[e]) {
			e++
		}
		return e
	}
	for j := start; j < len(lines); j++ {
		if d.unitEnd.MatchString(stripEOL(lines[j])) {
			return j + 1
		}
	}
	return len(lines)
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// matchStart reports the kind and name of a line that starts a unit.\nfunc (d *Declarative) matchStart(line string) (kind, name string, ok bool) {\n\tm := d.unitStart.FindStringSubmatch(stripEOL(line))\n\tif m == nil {\n\t\treturn \"\", \"\", false\n\t}\n\treturn namedGroup(d.unitStart, m, \"kind\"), namedGroup(d.unitStart, m, \"name\"), true\n}","file":"internal/adopt/declarative.go","name":"matchStart","receiver":"Declarative"}}
// matchStart reports the kind and name of a line that starts a unit.
func (d *Declarative) matchStart(line string) (kind, name string, ok bool) {
	m := d.unitStart.FindStringSubmatch(stripEOL(line))
	if m == nil {
		return "", "", false
	}
	return namedGroup(d.unitStart, m, "kind"), namedGroup(d.unitStart, m, "name"), true
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// attaches reports whether a line matches any AttachAbove regexp.\nfunc (d *Declarative) attaches(line string) bool {\n\tbody := stripEOL(line)\n\tfor _, re := range d.attach {\n\t\tif re.MatchString(body) {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}","file":"internal/adopt/declarative.go","name":"attaches","receiver":"Declarative"}}
// attaches reports whether a line matches any AttachAbove regexp.
func (d *Declarative) attaches(line string) bool {
	body := stripEOL(line)
	for _, re := range d.attach {
		if re.MatchString(body) {
			return true
		}
	}
	return false
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// unitSpan builds a unit span carrying only the kwargs its action declares.\nfunc (d *Declarative) unitSpan(path, action, name string, lines []string) Span {\n\ttext := strings.Join(lines, \"\")\n\tkwargs := map[string]any{}\n\tfor k := range d.declared[action] {\n\t\tswitch k {\n\t\tcase \"file\":\n\t\t\tkwargs[k] = path\n\t\tcase \"name\":\n\t\t\tkwargs[k] = name\n\t\tcase \"code\":\n\t\t\tkwargs[k] = strings.TrimSuffix(text, \"\\n\")\n\t\t}\n\t}\n\treturn Span{Kind: KindUnit, Text: text, Action: action, Kwargs: kwargs}\n}","file":"internal/adopt/declarative.go","name":"unitSpan","receiver":"Declarative"}}
// unitSpan builds a unit span carrying only the kwargs its action declares.
func (d *Declarative) unitSpan(path, action, name string, lines []string) Span {
	text := strings.Join(lines, "")
	kwargs := map[string]any{}
	for k := range d.declared[action] {
		switch k {
		case "file":
			kwargs[k] = path
		case "name":
			kwargs[k] = name
		case "code":
			kwargs[k] = strings.TrimSuffix(text, "\n")
		}
	}
	return Span{Kind: KindUnit, Text: text, Action: action, Kwargs: kwargs}
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// residualSpan builds a residual span from consecutive leftover lines.\nfunc residualSpan(lines []string) Span {\n\treturn Span{Kind: KindResidual, Text: strings.Join(lines, \"\")}\n}","file":"internal/adopt/declarative.go","name":"residualSpan"}}
// residualSpan builds a residual span from consecutive leftover lines.
func residualSpan(lines []string) Span {
	return Span{Kind: KindResidual, Text: strings.Join(lines, "")}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// stripEOL removes a trailing line ending.\nfunc stripEOL(line string) string {\n\treturn strings.TrimRight(line, \"\\r\\n\")\n}","file":"internal/adopt/declarative.go","name":"stripEOL"}}
// stripEOL removes a trailing line ending.
func stripEOL(line string) string {
	return strings.TrimRight(line, "\r\n")
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// namedGroup returns the text of a named submatch, or empty when absent.\nfunc namedGroup(re *regexp.Regexp, m []string, name string) string {\n\ti := re.SubexpIndex(name)\n\tif i < 0 || i >= len(m) {\n\t\treturn \"\"\n\t}\n\treturn m[i]\n}","file":"internal/adopt/declarative.go","name":"namedGroup"}}
// namedGroup returns the text of a named submatch, or empty when absent.
func namedGroup(re *regexp.Regexp, m []string, name string) string {
	i := re.SubexpIndex(name)
	if i < 0 || i >= len(m) {
		return ""
	}
	return m[i]
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// isBlank reports whether a line is empty or whitespace only.\nfunc isBlank(line string) bool {\n\treturn strings.TrimSpace(line) == \"\"\n}","file":"internal/adopt/declarative.go","name":"isBlank"}}
// isBlank reports whether a line is empty or whitespace only.
func isBlank(line string) bool {
	return strings.TrimSpace(line) == ""
}

// /sedum:setFunction
