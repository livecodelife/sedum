package genpkg

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/genpkg/adopt.go","path":"regexp"},"kwargs":"sha256:565fcd04637c6ea760d56c66ec6fed241bb7df4bdd300dbf0a151e7066a1178d","body":"sha256:b437bd748876cfcdfa00102314b52d7ce19d023c66492336bdfd27288b1329d2"}
import "regexp"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/genpkg/adopt.go","name":"Adopt"},"kwargs":"sha256:0e7b22e6cc0b9ec35b98a20027814c026a8ca378dd30999753aa224d57ab87b7","body":"sha256:acdb72e4f327089bd23641b8583a67e3bb5d08c0b7f24f99b831c4a8900cdd2a"}
// Adopt is the optional adopt: block of a package's sedum.yaml. It takes
// one of three forms: a command (Command), a declarative description of
// units (UnitStart, UnitEnd, Header, AttachAbove, Actions), or the name of
// a built-in adopter (Builtin). The forms do not mix.
type Adopt struct {
	Builtin     string            `yaml:"builtin"`
	Command     []string          `yaml:"command"`
	UnitStart   string            `yaml:"unit_start"`
	UnitEnd     string            `yaml:"unit_end"`
	Header      string            `yaml:"header"`
	AttachAbove []string          `yaml:"attach_above"`
	Actions     map[string]string `yaml:"actions"`
}

// /sedum:setType
// sedum:setConst {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","id":{"file":"internal/genpkg/adopt.go","name":"RuleAdoptInvalid"},"kwargs":"sha256:4406e6644fc64755419b6dc985481cbd251d707090f2b1745bb5c6bac4a376ce","body":"sha256:0e7cf89f5cbdd8d314ff765b6a5efe77dc62958ab9aa8ad7a0370fd335fec2bf"}
// RuleAdoptInvalid is reported when the adopt: block of sedum.yaml is
// malformed: a bad unit pattern, an unknown or ill-fitting action, or a
// command mixed with declarative fields.
const RuleAdoptInvalid = "adopt_invalid"

// /sedum:setConst
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/genpkg/adopt.go","name":"checkAdopt"},"kwargs":"sha256:a24e3d8f60a8a73129c406a0efa85e0a500a2334eb505864178efc238bd225a0","body":"sha256:4f180d6dcba109970c512486355eb1c67073f3427559069d4562c23f4c60cfbd"}
func checkAdopt(pkg *Package, r *reporter) {
	a := pkg.Adopt
	if a == nil {
		return
	}
	hasCommand := len(a.Command) > 0
	forms := 0
	for _, set := range []bool{hasCommand, a.UnitStart != "", a.Builtin != ""} {
		if set {
			forms++
		}
	}
	if forms != 1 {
		r.errorf(manifestFile, RuleAdoptInvalid, "adopt: set exactly one of command, unit_start or builtin")
	}
	if a.Builtin != "" {
		if !BuiltinAdopters[a.Builtin] {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: builtin %q is not a built-in adopter", a.Builtin)
		}
		if hasCommand || a.UnitStart != "" || a.Header != "" || a.UnitEnd != "" || len(a.AttachAbove) > 0 || len(a.Actions) > 0 {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: builtin cannot be combined with command, unit_start, header, unit_end, attach_above or actions")
		}
		return
	}
	if hasCommand {
		if a.Command[0] == "" {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: command must not start with an empty element")
		}
		if a.Header != "" || a.UnitEnd != "" || len(a.AttachAbove) > 0 || len(a.Actions) > 0 {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: command cannot be combined with header, unit_end, attach_above or actions")
		}
		return
	}
	if a.UnitStart == "" {
		return
	}
	if re, err := regexp.Compile(a.UnitStart); err != nil {
		r.errorf(manifestFile, RuleAdoptInvalid, "adopt: unit_start %q does not compile: %v", a.UnitStart, err)
	} else {
		groups := make(map[string]bool)
		for _, n := range re.SubexpNames() {
			groups[n] = true
		}
		for _, want := range []string{"kind", "name"} {
			if !groups[want] {
				r.errorf(manifestFile, RuleAdoptInvalid, "adopt: unit_start %q lacks the named group %q", a.UnitStart, want)
			}
		}
	}
	switch a.UnitEnd {
	case "":
		r.errorf(manifestFile, RuleAdoptInvalid, "adopt: unit_end is required with unit_start; use dedent, blank, or a regexp")
	case "dedent", "blank":
	default:
		if _, err := regexp.Compile(a.UnitEnd); err != nil {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: unit_end %q is not dedent, blank, or a regexp that compiles: %v", a.UnitEnd, err)
		}
	}
	if a.Header != "" {
		if _, err := regexp.Compile(a.Header); err != nil {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: header %q does not compile: %v", a.Header, err)
		}
	}
	for _, p := range a.AttachAbove {
		if _, err := regexp.Compile(p); err != nil {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: attach_above %q does not compile: %v", p, err)
		}
	}
	if len(a.Actions) == 0 {
		r.errorf(manifestFile, RuleAdoptInvalid, "adopt: actions must map at least one unit kind to an action")
	}
	for kind, name := range a.Actions {
		act, ok := pkg.Actions[name]
		if !ok {
			r.errorf(manifestFile, RuleAdoptInvalid, "adopt: actions maps %q to %q, which the package does not define", kind, name)
			continue
		}
		for kw, k := range act.Kwargs {
			if k.Required && kw != "file" && kw != "name" && kw != "code" {
				r.errorf(manifestFile, RuleAdoptInvalid, "adopt: action %q requires kwarg %q, which the declarative adopter does not supply", name, kw)
			}
		}
	}
}

// /sedum:setFunction
// sedum:setVar {"v":2,"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:c7b9d9ac1655e8211b8b2ca8c49cf34b8d510d8d823d4f5eb60ec7e1018aa2eb","id":{"file":"internal/genpkg/adopt.go","name":"BuiltinAdopters"},"kwargs":"sha256:06f80552361fa4d7aa9b1c96b6022880baca01e887a534d0a3a1600b0f686117","body":"sha256:11090f3a3888812d397cfad01ff62200b442b8b2af33c245e9afa35a6140552e"}
// BuiltinAdopters names the adopters compiled into sedum (internal/adopt
// holds their implementations).
var BuiltinAdopters = map[string]bool{"go": true}

// /sedum:setVar
