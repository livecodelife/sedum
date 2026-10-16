package genpkg

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/genpkg/adopt.go","path":"regexp"}}
import "regexp"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Adopt is the optional adopt: block of a package's sedum.yaml. It takes\n// one of three forms: a command (Command), a declarative description of\n// units (UnitStart, UnitEnd, Header, AttachAbove, Actions), or the name of\n// a built-in adopter (Builtin). The forms do not mix.\ntype Adopt struct {\n\tBuiltin     string            `yaml:\"builtin\"`\n\tCommand     []string          `yaml:\"command\"`\n\tUnitStart   string            `yaml:\"unit_start\"`\n\tUnitEnd     string            `yaml:\"unit_end\"`\n\tHeader      string            `yaml:\"header\"`\n\tAttachAbove []string          `yaml:\"attach_above\"`\n\tActions     map[string]string `yaml:\"actions\"`\n}","file":"internal/genpkg/adopt.go","name":"Adopt"}}
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
// sedum:setConst {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// RuleAdoptInvalid is reported when the adopt: block of sedum.yaml is\n// malformed: a bad unit pattern, an unknown or ill-fitting action, or a\n// command mixed with declarative fields.\nconst RuleAdoptInvalid = \"adopt_invalid\"","file":"internal/genpkg/adopt.go","name":"RuleAdoptInvalid"}}
// RuleAdoptInvalid is reported when the adopt: block of sedum.yaml is
// malformed: a bad unit pattern, an unknown or ill-fitting action, or a
// command mixed with declarative fields.
const RuleAdoptInvalid = "adopt_invalid"

// /sedum:setConst
// sedum:setFunction {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func checkAdopt(pkg *Package, r *reporter) {\n\ta := pkg.Adopt\n\tif a == nil {\n\t\treturn\n\t}\n\thasCommand := len(a.Command) > 0\n\tforms := 0\n\tfor _, set := range []bool{hasCommand, a.UnitStart != \"\", a.Builtin != \"\"} {\n\t\tif set {\n\t\t\tforms++\n\t\t}\n\t}\n\tif forms != 1 {\n\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: set exactly one of command, unit_start or builtin\")\n\t}\n\tif a.Builtin != \"\" {\n\t\tif !BuiltinAdopters[a.Builtin] {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: builtin %q is not a built-in adopter\", a.Builtin)\n\t\t}\n\t\tif hasCommand || a.UnitStart != \"\" || a.Header != \"\" || a.UnitEnd != \"\" || len(a.AttachAbove) > 0 || len(a.Actions) > 0 {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: builtin cannot be combined with command, unit_start, header, unit_end, attach_above or actions\")\n\t\t}\n\t\treturn\n\t}\n\tif hasCommand {\n\t\tif a.Command[0] == \"\" {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: command must not start with an empty element\")\n\t\t}\n\t\tif a.Header != \"\" || a.UnitEnd != \"\" || len(a.AttachAbove) > 0 || len(a.Actions) > 0 {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: command cannot be combined with header, unit_end, attach_above or actions\")\n\t\t}\n\t\treturn\n\t}\n\tif a.UnitStart == \"\" {\n\t\treturn\n\t}\n\tif re, err := regexp.Compile(a.UnitStart); err != nil {\n\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: unit_start %q does not compile: %v\", a.UnitStart, err)\n\t} else {\n\t\tgroups := make(map[string]bool)\n\t\tfor _, n := range re.SubexpNames() {\n\t\t\tgroups[n] = true\n\t\t}\n\t\tfor _, want := range []string{\"kind\", \"name\"} {\n\t\t\tif !groups[want] {\n\t\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: unit_start %q lacks the named group %q\", a.UnitStart, want)\n\t\t\t}\n\t\t}\n\t}\n\tswitch a.UnitEnd {\n\tcase \"\":\n\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: unit_end is required with unit_start; use dedent, blank, or a regexp\")\n\tcase \"dedent\", \"blank\":\n\tdefault:\n\t\tif _, err := regexp.Compile(a.UnitEnd); err != nil {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: unit_end %q is not dedent, blank, or a regexp that compiles: %v\", a.UnitEnd, err)\n\t\t}\n\t}\n\tif a.Header != \"\" {\n\t\tif _, err := regexp.Compile(a.Header); err != nil {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: header %q does not compile: %v\", a.Header, err)\n\t\t}\n\t}\n\tfor _, p := range a.AttachAbove {\n\t\tif _, err := regexp.Compile(p); err != nil {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: attach_above %q does not compile: %v\", p, err)\n\t\t}\n\t}\n\tif len(a.Actions) == 0 {\n\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: actions must map at least one unit kind to an action\")\n\t}\n\tfor kind, name := range a.Actions {\n\t\tact, ok := pkg.Actions[name]\n\t\tif !ok {\n\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: actions maps %q to %q, which the package does not define\", kind, name)\n\t\t\tcontinue\n\t\t}\n\t\tfor kw, k := range act.Kwargs {\n\t\t\tif k.Required && kw != \"file\" && kw != \"name\" && kw != \"code\" {\n\t\t\t\tr.errorf(manifestFile, RuleAdoptInvalid, \"adopt: action %q requires kwarg %q, which the declarative adopter does not supply\", name, kw)\n\t\t\t}\n\t\t}\n\t}\n}","file":"internal/genpkg/adopt.go","name":"checkAdopt"}}
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
// sedum:setVar {"tier":"owned","record":"prov-2026-3d0321d4","gen":"sha256:fb10210cd9c7db2ab7b1e5ec7e7d74d1bce1f7fdf73fc802b53480149b8dd849","pkg":"sha256:c7b9d9ac1655e8211b8b2ca8c49cf34b8d510d8d823d4f5eb60ec7e1018aa2eb","kwargs":{"code":"// BuiltinAdopters names the adopters compiled into sedum (internal/adopt\n// holds their implementations).\nvar BuiltinAdopters = map[string]bool{\"go\": true}","file":"internal/genpkg/adopt.go","name":"BuiltinAdopters"}}
// BuiltinAdopters names the adopters compiled into sedum (internal/adopt
// holds their implementations).
var BuiltinAdopters = map[string]bool{"go": true}

// /sedum:setVar
