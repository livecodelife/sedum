package genpkg

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/livecodelife/sedum/internal/render"
)

// ConvParam is one stub parameter.
type ConvParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ConvStub is one function or method whose body is left open.
type ConvStub struct {
	Name         string      `json:"name"`
	ReceiverName string      `json:"receiver_name,omitempty"`
	ReceiverType string      `json:"receiver_type,omitempty"`
	Params       []ConvParam `json:"params,omitempty"`
	Results      []string    `json:"results,omitempty"`
	Doc          string      `json:"doc,omitempty"`
}

// ConvFile is one source file whose stubs delegate to helpers.
type ConvFile struct {
	Path    string     `json:"path"`
	Package string     `json:"package,omitempty"`
	Imports []string   `json:"imports,omitempty"`
	Stubs   []ConvStub `json:"stubs,omitempty"`
}

// ConvRenderInput is the input of RenderConvergence.
type ConvRenderInput struct {
	Files []ConvFile `json:"files"`
}

// ConvRenderedStub is one stub as rendered: the stub text to inject, the helper
// it calls, and the helper's signature without a body.
type ConvRenderedStub struct {
	Name string `json:"name"`
	// Action and Kwargs are the catalog invocation that writes Stub into the
	// file (kwargs carry file, name, code and, for setMethod, receiver).
	Action  string            `json:"action"`
	Kwargs  map[string]string `json:"kwargs"`
	Helper  string            `json:"helper"`
	Stub    string            `json:"stub"`
	Declare string            `json:"declare"`
}

// ConvRenderedFile is one file as rendered. Seed is the synth file's whole
// initial body (empty when the file has no stubs); SynthImport, when set, is
// the module a stub file must import to reach its helpers.
type ConvRenderedFile struct {
	Path      string `json:"path"`
	SynthPath string `json:"synth_path"`
	// Imports are the kwargs of each addImport invocation the file needs (the
	// file's own imports and, when set, its synth import), sorted by path.
	Imports     []map[string]string `json:"imports"`
	SynthImport string              `json:"synth_import,omitempty"`
	Seed        string              `json:"seed"`
	Stubs       []ConvRenderedStub  `json:"stubs"`
}

// ConvRenderOutput is the output of RenderConvergence.
type ConvRenderOutput struct {
	Package     string             `json:"package"`
	SynthSuffix string             `json:"synth_suffix"`
	Placeholder string             `json:"placeholder"`
	Files       []ConvRenderedFile `json:"files"`
}

var blankRuns = regexp.MustCompile(`\n{3,}`)

// RenderConvergence renders stubs, helper signatures and seeds for in using the
// package's convergence block. The core knows no language: it joins lists,
// picks the func or method variant and hands every template the values it is
// allowed to read.
func (p *Package) RenderConvergence(in ConvRenderInput) (ConvRenderOutput, error) {
	c := p.Convergence
	if c == nil {
		return ConvRenderOutput{}, fmt.Errorf("package %s declares no convergence block", p.Name)
	}
	r := &convRenderer{p: p, c: c, cache: map[string]*render.Template{}}
	out := ConvRenderOutput{Package: p.Name, SynthSuffix: c.SynthSuffix, Placeholder: ConvergencePlaceholder, Files: []ConvRenderedFile{}}
	for _, f := range in.Files {
		rf, err := r.file(f)
		if err != nil {
			return ConvRenderOutput{}, fmt.Errorf("%s: %w", f.Path, err)
		}
		out.Files = append(out.Files, rf)
	}
	return out, nil
}

type convRenderer struct {
	p     *Package
	c     *Convergence
	cache map[string]*render.Template
}

// do renders one template source with vals, which must bind every value the
// template references (the engine refuses an unbound one).
func (r *convRenderer) do(label, src string, vals map[string]string) (string, error) {
	t, ok := r.cache[src]
	if !ok {
		var err error
		t, err = render.Compile(r.p.Engine, src)
		if err != nil {
			return "", fmt.Errorf("%s: %w", label, err)
		}
		r.cache[src] = t
	}
	bound := make(map[string]any, len(vals))
	for k, v := range vals {
		bound[k] = v
	}
	s, err := t.Render(bound)
	if err != nil {
		return "", fmt.Errorf("%s: %w", label, err)
	}
	return s, nil
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	rs[0] = []rune(strings.ToUpper(string(rs[0])))[0]
	return string(rs)
}

var nonIdent = regexp.MustCompile(`[^A-Za-z0-9_]`)

// importAlias derives the name an import is bound to from its path.
func importAlias(p string) string {
	if i := strings.LastIndexAny(p, "/:"); i >= 0 && i < len(p)-1 {
		p = p[i+1:]
	}
	p = strings.TrimSuffix(p, path.Ext(p))
	p = nonIdent.ReplaceAllString(p, "_")
	if p != "" && p[0] >= '0' && p[0] <= '9' {
		p = "_" + p
	}
	return p
}

func (r *convRenderer) file(f ConvFile) (ConvRenderedFile, error) {
	c := r.c
	dir, base := path.Dir(f.Path), path.Base(f.Path)
	ext := path.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	synthStem := stem + c.SynthSuffix
	rf := ConvRenderedFile{Path: f.Path, SynthPath: path.Join(dir, synthStem+ext), Stubs: []ConvRenderedStub{}, Imports: []map[string]string{}}

	if c.SynthImport != "" && len(f.Stubs) > 0 {
		s, err := r.do("synth_import", c.SynthImport, map[string]string{"synth_stem": synthStem, "stem": stem, "dir": dir, "ext": ext})
		if err != nil {
			return rf, err
		}
		rf.SynthImport = s
	}
	allImports := append([]string(nil), f.Imports...)
	if rf.SynthImport != "" {
		allImports = append(allImports, rf.SynthImport)
	}
	sort.Strings(allImports)
	for _, i := range allImports {
		kw := map[string]string{"file": f.Path, "path": i}
		if c.ImportAlias {
			kw["alias"] = importAlias(i)
		}
		rf.Imports = append(rf.Imports, kw)
	}
	pkg, err := r.packageName(f, dir)
	if err != nil {
		return rf, err
	}
	var declares, decls []string
	for _, s := range f.Stubs {
		rs, err := r.stub(f, s, stem, dir, ext, synthStem)
		if err != nil {
			return rf, fmt.Errorf("stub %s: %w", s.Name, err)
		}
		rf.Stubs = append(rf.Stubs, rs)
		declares = append(declares, rs.Declare)
		d, err := r.do("seed_decl", c.SeedDecl, map[string]string{"declare": rs.Declare, "placeholder": ConvergencePlaceholder})
		if err != nil {
			return rf, err
		}
		decls = append(decls, d)
	}
	if len(f.Stubs) == 0 {
		return rf, nil
	}
	all := strings.Join(declares, "\n")
	var imps []string
	for _, i := range f.Imports {
		if c.ImportUse != "" {
			use, err := r.do("import_use", c.ImportUse, map[string]string{"base": path.Base(i), "alias": importAlias(i)})
			if err != nil {
				return rf, err
			}
			if !strings.Contains(all, use) {
				continue
			}
		}
		imps = append(imps, i)
	}
	sort.Strings(imps)
	var lines []string
	for _, i := range imps {
		l, err := r.do("import", c.Import, map[string]string{"path": i, "alias": importAlias(i)})
		if err != nil {
			return rf, err
		}
		lines = append(lines, l)
	}
	seed, err := r.do("seed", c.Seed, map[string]string{
		"package": pkg, "imports": strings.Join(lines, "\n"), "decls": strings.Join(decls, "\n\n"),
		"file": f.Path, "stem": stem, "dir": dir, "ext": ext,
	})
	if err != nil {
		return rf, err
	}
	rf.Seed = blankRuns.ReplaceAllString(seed, "\n\n")
	return rf, nil
}

func (r *convRenderer) packageName(f ConvFile, dir string) (string, error) {
	if f.Package != "" {
		return f.Package, nil
	}
	pk := r.c.Package
	if pk.Entry != "" {
		for _, pat := range pk.EntryDirs {
			if ok, _ := path.Match(pat, dir); ok {
				return r.do("package.entry", pk.Entry, map[string]string{})
			}
		}
	}
	if pk.Default == "" {
		return "", nil
	}
	return r.do("package.default", pk.Default, map[string]string{"dirname": path.Base(dir)})
}

func (r *convRenderer) stub(f ConvFile, s ConvStub, stem, dir, ext, synthStem string) (ConvRenderedStub, error) {
	c := r.c
	method := s.ReceiverType != ""
	if (s.ReceiverType == "") != (s.ReceiverName == "") {
		return ConvRenderedStub{}, fmt.Errorf("receiver_name and receiver_type go together")
	}
	if s.Name == "" {
		return ConvRenderedStub{}, fmt.Errorf("a stub needs a name")
	}
	helper, err := r.do("helper", c.Helper, map[string]string{"stem": stem, "name": s.Name, "Name": upperFirst(s.Name), "dir": dir, "ext": ext})
	if err != nil {
		return ConvRenderedStub{}, err
	}
	var params, args []string
	for _, pa := range s.Params {
		if pa.Name == "" || pa.Type == "" {
			return ConvRenderedStub{}, fmt.Errorf("parameter without a name or type")
		}
		t, err := r.do("param", c.Param, map[string]string{"name": pa.Name, "type": pa.Type})
		if err != nil {
			return ConvRenderedStub{}, err
		}
		params = append(params, t)
		tmpl := c.Arg
		if c.SpreadPrefix != "" && strings.HasPrefix(pa.Type, c.SpreadPrefix) {
			tmpl = c.SpreadArg
		}
		a, err := r.do("arg", tmpl, map[string]string{"name": pa.Name})
		if err != nil {
			return ConvRenderedStub{}, err
		}
		args = append(args, a)
	}
	var res string
	switch len(s.Results) {
	case 0:
		res, err = r.do("result.none", c.Result.None, map[string]string{})
	case 1:
		res, err = r.do("result.one", c.Result.One, map[string]string{"result": s.Results[0]})
	default:
		res, err = r.do("result.many", c.Result.Many, map[string]string{"results": strings.Join(s.Results, ", ")})
	}
	if err != nil {
		return ConvRenderedStub{}, err
	}

	pj, aj := strings.Join(params, ", "), strings.Join(args, ", ")
	var call, declare, recvParam string
	if method {
		ra, err := r.do("receiver_arg", c.ReceiverArg, map[string]string{"receiver_name": s.ReceiverName})
		if err != nil {
			return ConvRenderedStub{}, err
		}
		rp, err := r.do("receiver_param", c.ReceiverParam, map[string]string{"receiver_name": s.ReceiverName, "receiver_type": s.ReceiverType})
		if err != nil {
			return ConvRenderedStub{}, err
		}
		recvParam = rp
		call, err = r.do("call.method", c.Call.Method, map[string]string{
			"helper": helper, "args": aj, "synth_stem": synthStem, "receiver_name": s.ReceiverName,
			"recv_args": strings.Join(append([]string{ra}, args...), ", ")})
		if err != nil {
			return ConvRenderedStub{}, err
		}
		declare, err = r.do("declare.method", c.Declare.Method, map[string]string{
			"helper": helper, "params": pj, "result": res, "receiver_name": s.ReceiverName, "receiver_type": s.ReceiverType,
			"recv_params": strings.Join(append([]string{rp}, params...), ", ")})
		if err != nil {
			return ConvRenderedStub{}, err
		}
	} else {
		call, err = r.do("call.func", c.Call.Func, map[string]string{"helper": helper, "args": aj, "synth_stem": synthStem})
		if err != nil {
			return ConvRenderedStub{}, err
		}
		declare, err = r.do("declare.func", c.Declare.Func, map[string]string{"helper": helper, "params": pj, "result": res})
		if err != nil {
			return ConvRenderedStub{}, err
		}
	}
	bodyTmpl, label := c.Body.Void, "body.void"
	if len(s.Results) > 0 {
		bodyTmpl, label = c.Body.Value, "body.value"
	}
	body, err := r.do(label, bodyTmpl, map[string]string{"call": call})
	if err != nil {
		return ConvRenderedStub{}, err
	}

	note, err := r.do("note", c.Note, map[string]string{"name": s.Name, "helper": helper})
	if err != nil {
		return ConvRenderedStub{}, err
	}
	cp, cs := r.p.CommentPrefix+" ", r.p.CommentSuffix
	var doc []string
	if strings.TrimSpace(s.Doc) != "" {
		for _, l := range strings.Split(strings.TrimSpace(s.Doc), "\n") {
			doc = append(doc, cp+l+cs)
		}
	}
	doc = append(doc, cp+note+cs)

	vals := map[string]string{
		"name": s.Name, "params": pj, "result": res, "body": body,
		"doc": strings.Join(doc, "\n"), "helper": helper,
	}
	tmpl, label := c.Stub.Func, "stub.func"
	if method {
		vals["receiver_name"], vals["receiver_type"] = s.ReceiverName, s.ReceiverType
		vals["recv_params"] = strings.Join(append([]string{recvParam}, params...), ", ")
		tmpl, label = c.Stub.Method, "stub.method"
	}
	stub, err := r.do(label, tmpl, vals)
	if err != nil {
		return ConvRenderedStub{}, err
	}
	action, declName := "setFunction", s.Name
	kw := map[string]string{"file": f.Path}
	if method {
		action = c.MethodAction
		if action == "" {
			action = "setMethod"
		}
		if c.MethodName != "" {
			declName, err = r.do("method_name", c.MethodName, map[string]string{"receiver_type": s.ReceiverType, "name": s.Name})
			if err != nil {
				return ConvRenderedStub{}, err
			}
		}
		if action == "setMethod" {
			kw["receiver"] = strings.TrimPrefix(s.ReceiverType, "*")
		}
	}
	kw["name"], kw["code"] = declName, stub
	return ConvRenderedStub{Name: s.Name, Action: action, Kwargs: kw, Helper: helper, Stub: stub, Declare: declare}, nil
}
