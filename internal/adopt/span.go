package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/span.go","path":"fmt"},"kwargs":"sha256:a511fead2ad746f881b30316e314362a81081e1368f6bee25dce1e82bba7637a","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/span.go","path":"strings"},"kwargs":"sha256:537e91dc804ab1fea6670bb8d9a2d0273a1bcf2cbd10980f6ed564470ddb4a2e","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/span.go","name":"Kind"},"kwargs":"sha256:1d3fd9c43039035de381a38d7d10f78f7f27502c78ba7b67ebdb137816ca4f7f","body":"sha256:3758224b6a2cac6ae53f6d3dad4627aaf0d22e62b3b46a15ddb134aae8f8dd7a"}
type Kind string

// /sedum:setType
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/adopt/span.go","name":"KindHeader"},"kwargs":"sha256:ecaf69c6f5df6364e891ccb182e8bb085178ef6f747afc0c26d43bcf4dc86ebd","body":"sha256:0f726fdb3940e4e5d72e661141630ac2c2f07c4a5c7ece37911e24691c3425a7"}
const KindHeader Kind = "header"

// /sedum:setConst
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/adopt/span.go","name":"KindUnit"},"kwargs":"sha256:c4093ec6cdd2ad6167b3eed34efda79347809f5d6deae8cd21507f20e8c45e0c","body":"sha256:c6d70636874902b973eedd0f2e0e58954cf27f54ab9e15c50645ab232bba7641"}
const KindUnit Kind = "unit"

// /sedum:setConst
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/adopt/span.go","name":"KindResidual"},"kwargs":"sha256:b01592e9393633d3aa7cd2ecd231b16e39cbb3d997cd5ac7da35cd38f41aad59","body":"sha256:14a81aeadabcfccb501e0e91b2bf7bb5ee5f2dd483418370cdb53ad318539770"}
const KindResidual Kind = "residual"

// /sedum:setConst
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/span.go","name":"Span"},"kwargs":"sha256:d09c9bcb4d5a1af2f5ea2569fc3d45ca215ec08395b1a4ea5e22d54b1b665e92","body":"sha256:3c62b48fea21b3cb45b745550203c24182a4ebc744bac8091f304f6426addaf8"}
type Span struct {
	Kind   Kind           `json:"kind"`
	Text   string         `json:"text"`
	Action string         `json:"action,omitempty"`
	Kwargs map[string]any `json:"kwargs,omitempty"`
}

// /sedum:setType
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/span.go","name":"Adopter"},"kwargs":"sha256:9987c437b00799b7c974214c09a85e194375ea3224139454f84bdbc57f6ed9f2","body":"sha256:74e9d1b6d169faf234808d70ecc265f0d06562e72546a230df80b05f1f052d2c"}
type Adopter interface {
	Split(path string, src []byte) ([]Span, error)
}

// /sedum:setType
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/span.go","name":"CoverageError"},"kwargs":"sha256:e3ab8c44d149de907a6c676ccb24d8a9ec01e2130cbcca0d33dad542b1bcf8b1","body":"sha256:3eec0e038600c7c9f1a4209f50eeeb45a33ba59703ad31a0c85660689a16c6ce"}
type CoverageError struct {
	Offset int
	Line   int
	Want   string
	Got    string
}

// /sedum:setType
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/span.go","name":"Error","receiver":"CoverageError"},"kwargs":"sha256:93edfd02762d4f6f6dbb072f05548127dd0b2ec2d711af66b720f7afe4a65386","body":"sha256:6a70ad0b766af6d5e6b8e77f0c0011405dd24b078182e121477faca791b828f4"}
func (e CoverageError) Error() string {
	return fmt.Sprintf("coverage: spans differ from the file at byte %d (line %d): want %q got %q", e.Offset, e.Line, e.Want, e.Got)
}

// /sedum:setMethod
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/span.go","name":"Coverage"},"kwargs":"sha256:476440da31bf6b33d6a41a9b18df80bcfd9a6d3f25ead9aeebbfc4132dfe37af","body":"sha256:9f29d866ba61b1aa56195aea57660b6df845a93d91a1a72113023a6e10f82d97"}
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
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/span.go","name":"Valid","receiver":"Span"},"kwargs":"sha256:7ae3dda69265bab622435e896b95ce8a84f7d9ecafe1449fcce8af0f17c4cb04","body":"sha256:d4f2ea104e9587c2ea880cb2b8544287d2a6b0f6258a3243028d48dea258d39c"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/span.go","name":"CheckSpans"},"kwargs":"sha256:a46bffa7e6f31ecae21a7ae88877550718444a7791e39a4c8133849f6daecb99","body":"sha256:d20677e8b95efff0470ef10a7c4bcf687324449e7ebeac65b54b54c5cd790ba7"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-d87e7f92","gen":"sha256:be079f96fee77ea810342a75f22d0f5b007cec1f10482b632a21550cc79041fc","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/span.go","name":"ResidualIsBlank"},"kwargs":"sha256:21bb0baabdc53d22b4047cf0ebc45ed09320d4b5d2074bc2f1a99bf2c49859a4","body":"sha256:2e7608a373418fefec44ed86a7434f2632c80aebfd0521a80f4bad2704dce003"}
func ResidualIsBlank(s Span) bool {
	return s.Kind == KindResidual && strings.TrimSpace(s.Text) == ""
}

// /sedum:setFunction
