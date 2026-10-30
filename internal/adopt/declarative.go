package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative.go","path":"fmt"},"kwargs":"sha256:79da9bebb30be3c7919ae968752e0e410d15315f2d4e78f5223b070b5d05c9cb","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative.go","path":"regexp"},"kwargs":"sha256:320ea0fa006e2bbbdb5869cb95fc5eb73213834667802b9f9d694dae91353318","body":"sha256:b437bd748876cfcdfa00102314b52d7ce19d023c66492336bdfd27288b1329d2"}
import "regexp"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative.go","path":"strings"},"kwargs":"sha256:1bd1e834bc95c41a81b6f917147965c11e9372513203be3334f5f33809298642","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/declarative.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:f5d41a6e241d7cf587c87cb707771f5a91c5aac70c557a209218a7d93cc49b34","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/declarative.go","name":"Declarative"},"kwargs":"sha256:974d12e103e4de87463609d9c82e385bdcd689003de64d3e79679d1118fb2072","body":"sha256:2c1064614787e8fadcb039c8fca5436a02bd19630aaa1b4e745294a7e24e16d6"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative.go","name":"NewDeclarative"},"kwargs":"sha256:717d9bd168487a370360933a9861122c0889042ea9803092987f01211ffc77fb","body":"sha256:413cbb6c47f52eee2ce702261c0a93e7b1e5bff6115b85ad4f0513c1e3b58120"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative.go","name":"splitLines"},"kwargs":"sha256:1af4da0aea8c09c4b062aa2d17b6f77312cdf73fac14a7d6915d746d529118d2","body":"sha256:00143b7dc4b76648106f5258392d39592e1c6bf16f377eb24e8b2ab161bee593"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/declarative.go","name":"Split","receiver":"Declarative"},"kwargs":"sha256:2dffe1843396c33a0fb8757f42edf127ab414fccf66a726c4f2a5e5c22fb348f","body":"sha256:b46849b49f902a79e36afdda508d3ab2bb537ce9318c20d9121d4577607008ed"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/declarative.go","name":"endOf","receiver":"Declarative"},"kwargs":"sha256:cb75ff24ba5d92fa149983f18dfe067aaa852fc74cb744e5af6ecc7c678e72ea","body":"sha256:4bfb36266ad25c66ea2a5dcd899445e8f1c51a67fa726bcf78bf69a28a78a409"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/declarative.go","name":"matchStart","receiver":"Declarative"},"kwargs":"sha256:cb0ebe509523b8d98e7e4a582b6d3b638d56a82cfd19dd1d117366b7807f39e0","body":"sha256:a16c35929ed6f6cc0538293498333d1ac8b028802f2e915b86b780a7281b883c"}
// matchStart reports the kind and name of a line that starts a unit.
func (d *Declarative) matchStart(line string) (kind, name string, ok bool) {
	m := d.unitStart.FindStringSubmatch(stripEOL(line))
	if m == nil {
		return "", "", false
	}
	return namedGroup(d.unitStart, m, "kind"), namedGroup(d.unitStart, m, "name"), true
}

// /sedum:setMethod
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/declarative.go","name":"attaches","receiver":"Declarative"},"kwargs":"sha256:e3ba63dbeb7e31a1dd9bd96ea33a726ba2c972072ccf2251ecb2ed367bb38bfe","body":"sha256:b097b29963d446336323feb85445ea8ee38753f2ad2620d23517a23b9e94b93a"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/declarative.go","name":"unitSpan","receiver":"Declarative"},"kwargs":"sha256:cdc6764a4e47cf8ce3972b08e1427cbda2ede29df8820e5b677acde4c1c779ff","body":"sha256:029313d39d4e089a5745d0a20ad5bcea2a2d6ba04f7146181f433277363db9e5"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative.go","name":"residualSpan"},"kwargs":"sha256:27491e688928bd4976a806ff977cd51b1ecb31f241c4ebc4b2c94e88b252f7a2","body":"sha256:b2d36af21f34af114e92b5a7215579843c7decbfe650a65b54aa8b2395b860bc"}
// residualSpan builds a residual span from consecutive leftover lines.
func residualSpan(lines []string) Span {
	return Span{Kind: KindResidual, Text: strings.Join(lines, "")}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative.go","name":"stripEOL"},"kwargs":"sha256:eb8a7e0a9d9a58ef3a47e45d73051757cdddb987c54baf2480fe0b6c8bacb321","body":"sha256:bb07cd5f9ed5d810a21594e553851a6fb57e53f1970d6e00ca1854f6c920f9f1"}
// stripEOL removes a trailing line ending.
func stripEOL(line string) string {
	return strings.TrimRight(line, "\r\n")
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative.go","name":"namedGroup"},"kwargs":"sha256:8438e430e74efb760f9757e8e18d940d7409d87b6e3c7c243ac13295e2e02c82","body":"sha256:dd4c921546ba90dd22a1cc30687edc84550c7e38172e81975279797e4abca412"}
// namedGroup returns the text of a named submatch, or empty when absent.
func namedGroup(re *regexp.Regexp, m []string, name string) string {
	i := re.SubexpIndex(name)
	if i < 0 || i >= len(m) {
		return ""
	}
	return m[i]
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-856f85eb","gen":"sha256:e6570a4fd6ea725d919e68451b96892a2f182bf97c341497b36f9681d7725f3a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/declarative.go","name":"isBlank"},"kwargs":"sha256:c04f407a65b988ebe3cf3e081e095641d6ff1a793277519ab6fc0e2cb25290e5","body":"sha256:687075b4717ba26a201ebcc9696b4124987f2b9066e377ea4f3ed2a93ce74a80"}
// isBlank reports whether a line is empty or whitespace only.
func isBlank(line string) bool {
	return strings.TrimSpace(line) == ""
}

// /sedum:setFunction
