package genpkg

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func loadConvPackage(t *testing.T, name string) *Package {
	t.Helper()
	set, findings, err := Load(filepath.Join("..", "..", "generators"), Options{Only: []string{name}})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range findings {
		if f.Kind == KindError {
			t.Fatalf("%s: %s", name, f)
		}
	}
	pkg, ok := set.Lookup(name)
	if !ok {
		t.Fatalf("no package %s", name)
	}
	return pkg
}

func TestTypeScriptConvergence(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("testdata", "convergence", "typescript", "input.json"))
	if err != nil {
		t.Fatal(err)
	}
	var in ConvRenderInput
	if err := json.Unmarshal(b, &in); err != nil {
		t.Fatal(err)
	}
	out, err := loadConvPackage(t, "typescript").RenderConvergence(in)
	if err != nil {
		t.Fatal(err)
	}
	f := out.Files[0]
	if f.SynthPath != "src/server_synth.ts" || f.SynthImport != "./server_synth.ts" {
		t.Errorf("synth path/import = %q, %q", f.SynthPath, f.SynthImport)
	}
	wantSeed := "// Helpers for src/server.ts, written by graft's convergence step.\n" +
		"import * as net from \"node:net\";\n\n" +
		"export function serverServeBody(s: Server, c: net.Socket): void {\n  throw new Error(\"graft: not converged\");\n}\n\n" +
		"export function serverParseBody(line: string): [string, number] {\n  throw new Error(\"graft: not converged\");\n}\n\n" +
		"export function serverTickBody() {\n  throw new Error(\"graft: not converged\");\n}\n\n"
	if f.Seed != wantSeed {
		t.Errorf("seed:\n got: %q\nwant: %q", f.Seed, wantSeed)
	}
	wantStub := "// Serve.\n// serve is rendered by graft; its body is a convergence region (serverServeBody).\n" +
		"export function Server_serve(s: Server, c: net.Socket): void {\n  return server_synth.serverServeBody(s, c);\n}"
	st := f.Stubs[0]
	if st.Stub != wantStub || st.Action != "setFunction" || st.Kwargs["name"] != "Server_serve" || st.Kwargs["receiver"] != "" {
		t.Errorf("method stub = %+v", st)
	}
	if got := f.Stubs[2].Stub; !strings.Contains(got, "  server_synth.serverTickBody();\n}") {
		t.Errorf("void stub = %q", got)
	}
	if len(f.Imports) != 2 || f.Imports[0]["path"] != "./server_synth.ts" || f.Imports[0]["alias"] != "server_synth" || f.Imports[1]["alias"] != "net" {
		t.Errorf("imports = %v", f.Imports)
	}

	// The output is real TypeScript: with the declarations a skeleton would
	// supply, Node's strip-types loads both modules.
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not installed")
	}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, text string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("package.json", `{"type":"module"}`)
	var main strings.Builder
	for _, im := range f.Imports {
		main.WriteString("import * as " + im["alias"] + " from \"" + im["path"] + "\";\n")
	}
	main.WriteString("export type Server = { n: number };\n")
	for _, s := range f.Stubs {
		main.WriteString(s.Stub + "\n")
	}
	write(f.Path, main.String())
	write(f.SynthPath, f.Seed)
	cmd := exec.Command(node, "--experimental-strip-types", "--no-warnings", "-e",
		`const m = await import("./src/server.ts"); try { m.parse("x") } catch (e) { if (!String(e).includes("graft: not converged")) throw e; process.exit(0) } process.exit(1)`)
	cmd.Dir = dir
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("node could not load and run the rendered files: %v\n%s", err, outb)
	}
}

// tsc finds a TypeScript compiler and Node's types: tsc on PATH with
// NODE_TYPES_DIR unset is not enough, so it installs both into a temp dir.
func tscCommand(t *testing.T) (dir, tsc string) {
	t.Helper()
	dir = t.TempDir()
	cmd := exec.Command("npm", "install", "--silent", "--no-audit", "--no-fund", "--prefix", dir, "typescript@5.9.3", "@types/node@24")
	if out, err := cmd.CombinedOutput(); err != nil {
		if os.Getenv("CI") != "" {
			t.Fatalf("npm install typescript: %v\n%s", err, out)
		}
		t.Skipf("cannot install typescript (%v)", err)
	}
	return dir, filepath.Join(dir, "node_modules", ".bin", "tsc")
}

// The rendered files are TypeScript that type-checks, not only code Node can
// strip: the seed imports the types its signatures use, hyphenated file names
// give valid identifiers, and an optional parameter is called without its "?".
func TestTypeScriptConvergenceTypechecks(t *testing.T) {
	in := ConvRenderInput{Files: []ConvFile{{
		Path:    "my-mod.ts",
		Imports: []string{"node:net"},
		Decls:   []ConvDecl{{Kind: "type", Name: "Server"}, {Kind: "type", Name: "Client"}, {Kind: "const", Name: "serverName"}},
		Stubs: []ConvStub{
			{Name: "serve", ReceiverName: "s", ReceiverType: "Server", Params: []ConvParam{{Name: "sock", Type: "net.Socket"}, {Name: "n?", Type: "number"}}, Results: []string{"string"}},
			{Name: "join", Params: []ConvParam{{Name: "c", Type: "Client"}, {Name: "...parts", Type: "string[]"}}},
		},
	}}}
	out, err := loadConvPackage(t, "typescript").RenderConvergence(in)
	if err != nil {
		t.Fatal(err)
	}
	f := out.Files[0]
	if strings.Contains(f.Stubs[0].Stub, "n?)") || !strings.Contains(f.Stubs[0].Stub, "myModServeBody(s, sock, n)") {
		t.Errorf("optional parameter: %q", f.Stubs[0].Stub)
	}
	if !strings.Contains(f.Stubs[1].Stub, "myModJoinBody(c, ...parts)") {
		t.Errorf("rest parameter: %q", f.Stubs[1].Stub)
	}
	if !strings.Contains(f.Seed, `import type { Client, Server } from "./my-mod.ts";`) {
		t.Errorf("seed lacks its type import:\n%s", f.Seed)
	}
	dir, tsc := tscCommand(t)
	work := filepath.Join(dir, "work")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatal(err)
	}
	var main strings.Builder
	for _, im := range f.Imports {
		main.WriteString("import * as " + im["alias"] + " from \"" + im["path"] + "\";\n")
	}
	main.WriteString("export class Server { n = 0; }\nexport class Client { id = 1; }\nexport const serverName = \"x\";\n")
	for _, s := range f.Stubs {
		main.WriteString(s.Stub + "\n")
	}
	for name, text := range map[string]string{
		f.Path: main.String(), f.SynthPath: f.Seed,
		"tsconfig.json": `{"compilerOptions":{"noEmit":true,"strict":true,"target":"es2022","module":"nodenext","allowImportingTsExtensions":true,"erasableSyntaxOnly":true,"types":["node"],"typeRoots":["../node_modules/@types"]},"files":["my-mod.ts","my-mod_synth.ts"]}`,
	} {
		if err := os.WriteFile(filepath.Join(work, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command(tsc, "-p", ".")
	cmd.Dir = work
	if outb, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("tsc --noEmit failed: %v\n%s\n--- stub file\n%s\n--- seed\n%s", err, outb, main.String(), f.Seed)
	}
}
