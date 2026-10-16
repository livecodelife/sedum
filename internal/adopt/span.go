package adopt

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/span.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/adopt/span.go","path":"strings"}}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type Kind string","file":"internal/adopt/span.go","name":"Kind"}}
type Kind string

// /sedum:setType
// sedum:setConst {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"const KindHeader Kind = \"header\"","file":"internal/adopt/span.go","name":"KindHeader"}}
const KindHeader Kind = "header"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"const KindUnit Kind = \"unit\"","file":"internal/adopt/span.go","name":"KindUnit"}}
const KindUnit Kind = "unit"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"const KindResidual Kind = \"residual\"","file":"internal/adopt/span.go","name":"KindResidual"}}
const KindResidual Kind = "residual"

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type Span struct {\n\tKind   Kind           `json:\"kind\"`\n\tText   string         `json:\"text\"`\n\tAction string         `json:\"action,omitempty\"`\n\tKwargs map[string]any `json:\"kwargs,omitempty\"`\n}","file":"internal/adopt/span.go","name":"Span"}}
type Span struct {
	Kind   Kind           `json:"kind"`
	Text   string         `json:"text"`
	Action string         `json:"action,omitempty"`
	Kwargs map[string]any `json:"kwargs,omitempty"`
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type Adopter interface {\n\tSplit(path string, src []byte) ([]Span, error)\n}","file":"internal/adopt/span.go","name":"Adopter"}}
type Adopter interface {
	Split(path string, src []byte) ([]Span, error)
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type CoverageError struct {\n\tOffset int\n\tLine   int\n\tWant   string\n\tGot    string\n}","file":"internal/adopt/span.go","name":"CoverageError"}}
type CoverageError struct {
	Offset int
	Line   int
	Want   string
	Got    string
}

// /sedum:setType
// sedum:setMethod {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (e CoverageError) Error() string {\n\treturn fmt.Sprintf(\"coverage: spans differ from the file at byte %d (line %d): want %q got %q\", e.Offset, e.Line, e.Want, e.Got)\n}","file":"internal/adopt/span.go","name":"Error","receiver":"CoverageError"}}
func (e CoverageError) Error() string {
	return fmt.Sprintf("coverage: spans differ from the file at byte %d (line %d): want %q got %q", e.Offset, e.Line, e.Want, e.Got)
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func Coverage(src []byte, spans []Span) error {\n\tvar got []byte\n\tfor _, s := range spans {\n\t\tgot = append(got, s.Text...)\n\t}\n\tif string(got) == string(src) {\n\t\treturn nil\n\t}\n\toff := 0\n\tfor off < len(got) && off < len(src) && got[off] == src[off] {\n\t\toff++\n\t}\n\tline := 1\n\tfor _, c := range src[:off] {\n\t\tif c == '\\n' {\n\t\t\tline++\n\t\t}\n\t}\n\tend := off + 40\n\tif end > len(src) {\n\t\tend = len(src)\n\t}\n\tgend := off + 40\n\tif gend > len(got) {\n\t\tgend = len(got)\n\t}\n\treturn &CoverageError{Offset: off, Line: line, Want: string(src[off:end]), Got: string(got[off:gend])}\n}","file":"internal/adopt/span.go","name":"Coverage"}}
func Coverage(src []byte, spans []Span) error {
	var got []byte
	for _, s := range spans {
		got = append(got, s.Text...)
	}
	if string(got) == string(src) {
		return nil
	}
	off := 0
	for off < len(got) && off < len(src) && got[off] == src[off] {
		off++
	}
	line := 1
	for _, c := range src[:off] {
		if c == '\n' {
			line++
		}
	}
	end := off + 40
	if end > len(src) {
		end = len(src)
	}
	gend := off + 40
	if gend > len(got) {
		gend = len(got)
	}
	return &CoverageError{Offset: off, Line: line, Want: string(src[off:end]), Got: string(got[off:gend])}
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"func (s Span) Valid() error {\n\tswitch s.Kind {\n\tcase KindUnit:\n\t\tif s.Action == \"\" || s.Text == \"\" {\n\t\t\treturn fmt.Errorf(\"unit span needs an action and text\")\n\t\t}\n\tcase KindHeader:\n\t\tif s.Action != \"\" || s.Kwargs != nil {\n\t\t\treturn fmt.Errorf(\"%s span must not carry an action or kwargs\", s.Kind)\n\t\t}\n\tcase KindResidual:\n\t\tif s.Action != \"\" || s.Kwargs != nil {\n\t\t\treturn fmt.Errorf(\"%s span must not carry an action or kwargs\", s.Kind)\n\t\t}\n\tdefault:\n\t\treturn fmt.Errorf(\"unknown span kind %q\", s.Kind)\n\t}\n\treturn nil\n}","file":"internal/adopt/span.go","name":"Valid","receiver":"Span"}}
func (s Span) Valid() error {
	switch s.Kind {
	case KindUnit:
		if s.Action == "" || s.Text == "" {
			return fmt.Errorf("unit span needs an action and text")
		}
	case KindHeader:
		if s.Action != "" || s.Kwargs != nil {
			return fmt.Errorf("%s span must not carry an action or kwargs", s.Kind)
		}
	case KindResidual:
		if s.Action != "" || s.Kwargs != nil {
			return fmt.Errorf("%s span must not carry an action or kwargs", s.Kind)
		}
	default:
		return fmt.Errorf("unknown span kind %q", s.Kind)
	}
	return nil
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func CheckSpans(src []byte, spans []Span) error {\n\tif len(spans) == 0 && len(src) > 0 {\n\t\treturn fmt.Errorf(\"no spans for a non-empty file\")\n\t}\n\tfor i, s := range spans {\n\t\tif err := s.Valid(); err != nil {\n\t\t\treturn fmt.Errorf(\"span %d: %w\", i, err)\n\t\t}\n\t}\n\treturn Coverage(src, spans)\n}","file":"internal/adopt/span.go","name":"CheckSpans"}}
func CheckSpans(src []byte, spans []Span) error {
	if len(spans) == 0 && len(src) > 0 {
		return fmt.Errorf("no spans for a non-empty file")
	}
	for i, s := range spans {
		if err := s.Valid(); err != nil {
			return fmt.Errorf("span %d: %w", i, err)
		}
	}
	return Coverage(src, spans)
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func ResidualIsBlank(s Span) bool {\n\treturn s.Kind == KindResidual && strings.TrimSpace(s.Text) == \"\"\n}","file":"internal/adopt/span.go","name":"ResidualIsBlank"}}
func ResidualIsBlank(s Span) bool {
	return s.Kind == KindResidual && strings.TrimSpace(s.Text) == ""
}

// /sedum:setFunction
