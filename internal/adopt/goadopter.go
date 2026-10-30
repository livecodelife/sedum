package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter.go","path":"bytes"},"kwargs":"sha256:10d4b3aa994ac760576c84b7675867b74dd86228dfc3aa9e05ac03d61863f266","body":"sha256:468731f6dd513df0a66690d4649dc9e7244c96de227d8e6ead98cc5e9747c57a"}
import "bytes"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter.go","path":"fmt"},"kwargs":"sha256:0d58f82fc1ed49a368692bb0c8054ba6701be2624373a4b3a32caf26ff0db9a8","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter.go","path":"go/ast"},"kwargs":"sha256:5c15e8c88614f6b49242537f41f8cdf6e8f7e032d94c9c1967400d70dc93c562","body":"sha256:87ecade8212d1fdaaf2d64e42d4ca60dcb760b315cef84fec0e0ea33c456f9bd"}
import "go/ast"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter.go","path":"go/format"},"kwargs":"sha256:210700bb9d0670faec087c0386389a45ddd965e55ea9bcd644459f9efa45d176","body":"sha256:beeb535979eb17fc4afc70c9568d421f75255020e7b8349095a3acbcb7c60ca8"}
import "go/format"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter.go","path":"go/parser"},"kwargs":"sha256:18ac952efab261702ae4d7cc6abb905d86992be3be4dcb55c256fa4e2ee1b8ef","body":"sha256:0ac92d8b32b3e16ba33377f8c2360a9e5c3de365754dcdc80d2b394f274d3af3"}
import "go/parser"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter.go","path":"go/token"},"kwargs":"sha256:d8228c8c2b931521127457105f48512831686c43c99709e2375118e234f1b34c","body":"sha256:5f01267812de3043b1217a858c3f33ce10a6b446ef5481ed844862caa5f6fc56"}
import "go/token"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter.go","path":"strconv"},"kwargs":"sha256:a898843b3e251fcf57c0dbe5114a165a92507cfc2f3ee6f456648288947945c8","body":"sha256:f33f981c64a42d78a5daa7ff3fca15b10045080cccd282fef544b59eaf7f6dc8"}
import "strconv"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/goadopter.go","path":"strings"},"kwargs":"sha256:1cf8eac22ff09bafae2147a63b49928986ee2cf7bb3d0df4ea37b34657162383","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/goadopter.go","name":"GoAdopter"},"kwargs":"sha256:c9b9d46d7f856ba22098043377d6061c71117c9af1f6b77f3a9cd3c6d9dce8ac","body":"sha256:9eb7b505b423713a6a34c33ceb76345f14f51b2fb65b9e3bb9e7c5b013a04c4a"}
// GoAdopter is the built-in adopter for Go source (adopt: {builtin: go}). It
// needs gofmt-clean input and splits a file with go/parser into a header, one
// unit per declaration, and residual spans for everything else.
type GoAdopter struct{}

// /sedum:setType
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter.go","name":"goRecvName"},"kwargs":"sha256:e80fac2f6bcbdf42e9450707b15e13c148be94cd1c20520aeb5273dc6afb0277","body":"sha256:a97d85799ca7cac92bd75971401a59cb9aeb3168187d46914777cc937574b60c"}
// goRecvName returns the receiver's base type name: no star, parentheses or type parameters.
func goRecvName(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.StarExpr:
		return goRecvName(t.X)
	case *ast.ParenExpr:
		return goRecvName(t.X)
	case *ast.IndexExpr:
		return goRecvName(t.X)
	case *ast.IndexListExpr:
		return goRecvName(t.X)
	case *ast.Ident:
		return t.Name
	}
	return ""
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter.go","name":"goGenName"},"kwargs":"sha256:08ee4e9af33ccacab8c02181b552e4230b149f51d7eeea826997cdec0aca8e3d","body":"sha256:c3893992faff9fe4075df1680218332ac16c6d8c35181a584cc107df1421c5ff"}
// goGenName returns the name of a type, var or const declaration: the first
// name of its first spec. It returns "" for a declaration with no specs.
func goGenName(d *ast.GenDecl) string {
	if len(d.Specs) == 0 {
		return ""
	}
	switch s := d.Specs[0].(type) {
	case *ast.TypeSpec:
		return s.Name.Name
	case *ast.ValueSpec:
		return s.Names[0].Name
	}
	return ""
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter.go","name":"goLineEnd"},"kwargs":"sha256:0117e3dcd77550eff08e8de6f82dfd43b91399e356eefc7130022efac82d9794","body":"sha256:a3a629349b5b73fa8bc7be3104a4c00ee0b56235db9076d5aca873427742bf66"}
// goLineEnd returns the offset just past the line that contains offset off.
func goLineEnd(src []byte, off int) int {
	if i := bytes.IndexByte(src[off:], '\n'); i >= 0 {
		return off + i + 1
	}
	return len(src)
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/goadopter.go","name":"goPlainImport"},"kwargs":"sha256:6036c1cb49d8827628affba036d9403f879063ac3a9cf8b84894c68244a4addc","body":"sha256:900ae0c1c8b828ab15fb21549b2cb966929c081c9c2d4dfa96927735666d878a"}
// goPlainImport reports whether d is one ungrouped, unaliased, uncommented
// import that is not cgo, and returns its path.
func goPlainImport(d *ast.GenDecl) (string, bool) {
	if d.Lparen.IsValid() || len(d.Specs) != 1 || d.Doc != nil {
		return "", false
	}
	s := d.Specs[0].(*ast.ImportSpec)
	if s.Name != nil || s.Doc != nil || s.Comment != nil {
		return "", false
	}
	p, err := strconv.Unquote(s.Path.Value)
	if err != nil || p == "C" {
		return "", false
	}
	return p, true
}

// /sedum:setFunction
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-9bc22cf7","gen":"sha256:bc2ac7f452b6ca292315e0250983a1bc6ea42ba5db0787d556851e7dba11d8e7","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/goadopter.go","name":"Split","receiver":"GoAdopter"},"kwargs":"sha256:8b87e6d4fec84319b4f28e025b9b464e5b344c8b1704a0fae628846ed2128eb0","body":"sha256:ea186de96947b8188faa3823faa2697ae64c86f246505cae957ceaf318dc29f8"}
// Split implements Adopter.
func (GoAdopter) Split(path string, src []byte) ([]Span, error) {
	formatted, err := format.Source(src)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if !bytes.Equal(formatted, src) {
		return nil, fmt.Errorf("%s: not gofmt-clean", path)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	off := func(p token.Pos) int { return fset.Position(p).Offset }
	var spans []Span
	pos := 0
	gap := func(to int) {
		if to > pos {
			spans = append(spans, Span{Kind: KindResidual, Text: string(src[pos:to])})
			pos = to
		}
	}
	gap(off(file.Package))
	headerEnd := goLineEnd(src, off(file.Name.End()))
	headerKind := KindHeader
	if !strings.Contains(path, "/") && file.Name.Name != "main" {
		headerKind = KindResidual
	}
	spans = append(spans, Span{Kind: headerKind, Text: string(src[pos:headerEnd])})
	pos = headerEnd
	seen := map[string]int{}
	imports := 0
	for _, decl := range file.Decls {
		start := off(decl.Pos())
		var doc *ast.CommentGroup
		var action, name string
		kwargs := map[string]any{"file": path}
		switch d := decl.(type) {
		case *ast.FuncDecl:
			doc = d.Doc
			name = d.Name.Name
			action = "setFunction"
			if d.Recv != nil && len(d.Recv.List) == 1 {
				action = "setMethod"
				kwargs["receiver"] = goRecvName(d.Recv.List[0].Type)
			}
			if d.Recv == nil && name == "init" {
				seen["init"]++
				name = fmt.Sprintf("init#%d", seen["init"])
			}
		case *ast.GenDecl:
			doc = d.Doc
			switch d.Tok {
			case token.IMPORT:
				if p, ok := goPlainImport(d); ok {
					action = "addImport"
					kwargs["path"] = p
				} else if len(d.Specs) == 1 && d.Specs[0].(*ast.ImportSpec).Path.Value == `"C"` {
					action = ""
				} else {
					imports++
					action = "setImports"
					name = fmt.Sprintf("imports#%d", imports)
				}
			case token.TYPE:
				action, name = "setType", goGenName(d)
			case token.VAR:
				action, name = "setVar", goGenName(d)
			case token.CONST:
				action, name = "setConst", goGenName(d)
			}
			if d.Tok != token.IMPORT && name == "" {
				action = ""
			}
		}
		if doc != nil {
			start = off(doc.Pos())
		}
		end := goLineEnd(src, off(decl.End()))
		gap(start)
		text := string(src[start:end])
		if action == "" {
			spans = append(spans, Span{Kind: KindResidual, Text: text})
			pos = end
			continue
		}
		if action != "addImport" && !(action == "setFunction" && len(name) > 5 && name[:5] == "init#") {
			key := action + "/" + name
			if r, ok := kwargs["receiver"].(string); ok {
				key += "/" + r
			}
			seen[key]++
			if seen[key] > 1 {
				name = fmt.Sprintf("%s#%d", name, seen[key])
			}
		}
		if action != "addImport" {
			kwargs["name"] = name
		}
		if action != "addImport" {
			kwargs["code"] = strings.TrimSuffix(text, "\n")
		}
		spans = append(spans, Span{Kind: KindUnit, Text: text, Action: action, Kwargs: kwargs})
		pos = end
	}
	gap(len(src))
	return spans, nil
}

// /sedum:setMethod
