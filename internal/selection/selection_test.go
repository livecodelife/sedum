package selection

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/calebcowen/sedum/internal/catalog"
	"github.com/calebcowen/sedum/internal/expand"
	"github.com/calebcowen/sedum/internal/genpkg"
	"github.com/calebcowen/sedum/internal/recording"
	"github.com/calebcowen/sedum/internal/resolve"
)

// Phase 4 is the only non-deterministic step in Sedum, and Phase 5 is what
// makes that survivable: every way a response can be wrong is caught by a check
// specific enough to re-prompt with. These tests are that claim, one rule at a
// time, plus the loop those rules feed.
//
// No model is called. The client is an interface precisely so that the part
// most worth testing is the part that needs no server.

func generators() map[string]string {
	return map[string]string{
		"rails/sedum.yaml": `name: rails
extensions: [".rb"]
comment_prefix: "#"
transforms:
  constantize: [singular, pascal]
  instantize: [plural, "prefix:@"]
`,
		"rails/files/app/controllers/{name}_controller.rb": "class {{name|constantize}}Controller\n" +
			"  # sedum:anchor:class_body\nend\n",
		"rails/files/app/models/{name}.rb": "class {{name|constantize}}\n" +
			"  # sedum:anchor:class_body\nend\n",
		"rails/actions/actions.yaml": `actions:
  createControllerMethod:
    kwargs:
      controller: { type: string, required: true }
      name: { type: string, required: true }
      collection: { type: string, required: false }
      cached: { type: bool, required: false }
      limit: { type: int, required: false }
      only: { type: list, required: false }
    discriminator: name
    variants: [index, show]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body

  createModelClass:
    kwargs:
      name: { type: string, required: true }
    injects_into: "app/models/{{name|snake}}.rb"
    anchor: class_body

  # A literal kwarg: the template emits it into a Ruby file verbatim, so what is
  # bound is source rather than prose. Its own action rather than a kwarg added
  # to one above, so that the type has a fixture no other test binds against
  # (prov-2026-e3d7b9ac).
  declareDefault:
    kwargs:
      name: { type: string, required: true }
      value: { type: literal, required: true }
    injects_into: "app/models/{{name|snake}}.rb"
    anchor: class_body

  # An optional kwarg the template always renders, carrying a default. Without
  # the default this is exactly the shape checkDerived exists to reject, so the
  # two together are what says the default is what makes the omission safe
  # (prov-2026-f03916ba).
  declareDefaulted:
    kwargs:
      name: { type: string, required: true }
      value: { type: literal, required: false, default: "nil" }
    injects_into: "app/models/{{name|snake}}.rb"
    anchor: class_body

  # A required list kwarg the template does not render, which is what every
  # list kwarg in every real package looks like. It exists so that "an empty
  # array is a legitimate value" is a case with a test rather than an argument.
  tagResource:
    kwargs:
      controller: { type: string, required: true }
      tags: { type: list, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body

  # No _default template, so a value outside the variant list is a hard error
  # rather than a stub. It is the case has_default exists to distinguish.
  addCallback:
    kwargs:
      controller: { type: string, required: true }
      hook: { type: string, required: true }
    discriminator: hook
    variants: [before, after]
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body

  hiddenHelper:
    kwargs:
      controller: { type: string, required: true }
    injects_into: "app/controllers/{{controller|snake}}_controller.rb"
    anchor: class_body
    exposed: false
`,
		"rails/actions/createControllerMethod/index.rb": "def index\n" +
			"  {{collection|instantize}} = {{collection|constantize}}.all\nend\n",
		"rails/actions/createControllerMethod/show.rb":     "def show\nend\n",
		"rails/actions/createControllerMethod/_default.rb": "def {{name|snake}}\nend\n",
		"rails/actions/addCallback/before.rb":              "before_action :x\n",
		"rails/actions/addCallback/after.rb":               "after_action :x\n",
		"rails/actions/createModelClass.rb":                "# {{name}}\n",
		"rails/actions/declareDefault.rb":                  "# default {{value}}\n",
		"rails/actions/declareDefaulted.rb":                "# defaulted {{value}}\n",
		"rails/actions/tagResource.rb":                     "# tagged\n",
		"rails/actions/hiddenHelper.rb":                    "# helper\n",

		// A second package whose one exposed action is a composite spanning
		// two files. It is what the target check has to cover per child.
		"cairn/sedum.yaml": `name: cairn
extensions: [".crn"]
comment_prefix: ";;"
transforms:
  slug: [plural, kebab]
`,
		"cairn/files/Units/{name}/Manifest.crn": ";; sedum:anchor:steps\n",
		"cairn/files/Shared/{name}.crn":         ";; shared\n",
		"cairn/actions/actions.yaml": `actions:
  provisionStep:
    composes: [declareConstant, addStep]

  addStep:
    kwargs:
      unit: { type: string, required: true }
      step: { type: string, required: true }
    injects_into: "Units/{{unit|slug}}/Manifest.crn"
    anchor: steps
    exposed: false

  declareConstant:
    kwargs:
      unit: { type: string, required: true }
      name: { type: string, required: true }
    injects_into: "Shared/{{name|slug}}.crn"
    anchor: end_of_file
    exposed: false
`,
		"cairn/actions/addStep.crn":         ";; step {{step}}\n",
		"cairn/actions/declareConstant.crn": ";; const {{name}}\n",
	}
}

// freeGenerators adds a package whose action's target is a kwarg rather than a
// convention, plus three templates that differ in what they plant: one carrying
// the anchor that action needs, one carrying a different anchor, and a fallback
// carrying none. Those are the three answers Phase 5's applicability check has
// to tell apart (prov-2026-14c832bf).
//
// It is a separate package rather than an action added to rails, so that
// planting a new anchor cannot change what the completeness pass reports about
// the fixtures every other test in this file is built on.
func freeGenerators() map[string]string {
	files := generators()
	files["free/sedum.yaml"] = `name: free
extensions: [".ts"]
comment_prefix: "//"
`
	files["free/files/src/{name}.ts"] = "// sedum:anchor:imports\n"
	files["free/files/config/{name}.ts"] = "// sedum:anchor:settings\n"
	files["free/files/_default.ts"] = "// generated\n"
	files["free/actions/actions.yaml"] = `actions:
  addImport:
    kwargs:
      file: { type: string, required: true }
      symbol: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: imports

  setOption:
    kwargs:
      file: { type: string, required: true }
      key: { type: string, required: true }
    injects_into: "{{file}}"
    anchor: settings
`
	files["free/actions/addImport.ts"] = "import { {{symbol}} }\n"
	files["free/actions/setOption.ts"] = "// {{key}}\n"
	return files
}

// rendered builds the Phase 3 output for a path, carrying what the file was
// created with. created leaves that empty, which is the right default for the
// tests that predate it - a file whose content this run did not produce is one
// Sedum has observed nothing about.
func rendered(t *testing.T, set *genpkg.Set, path, pkgName, content string) resolve.File {
	t.Helper()
	pkg, ok := set.Lookup(pkgName)
	if !ok {
		t.Fatalf("package %q did not load", pkgName)
	}
	return resolve.File{
		Resolution: resolve.Resolution{RecordID: "PR-014", Path: path, Package: pkg},
		Rendered:   content,
	}
}

func loadSet(t *testing.T, files map[string]string) *genpkg.Set {
	t.Helper()

	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}

	set, findings, err := genpkg.Load(root, genpkg.Options{})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, f := range findings {
		if f.Kind == genpkg.KindError {
			t.Fatalf("fixture package does not load: %s", f)
		}
	}
	return set
}

// created builds the Phase 3 output a record would have produced.
func created(t *testing.T, set *genpkg.Set, paths ...string) []resolve.File {
	t.Helper()

	var out []resolve.File
	for _, spec := range paths {
		path, pkgName, _ := strings.Cut(spec, "@")
		pkg, ok := set.Lookup(pkgName)
		if !ok {
			t.Fatalf("package %q did not load", pkgName)
		}
		out = append(out, resolve.File{
			Resolution: resolve.Resolution{RecordID: "PR-014", Path: path, Package: pkg},
		})
	}
	return out
}

func railsRequest(t *testing.T) Request {
	t.Helper()
	set := loadSet(t, generators())
	return Request{
		RecordID:    "PR-014",
		Intent:      "Add a read-only users controller.",
		Constraints: []string{"No writes."},
		Files:       created(t, set, "app/controllers/users_controller.rb@rails"),
	}
}

// modelRequest authorizes the model file declareDefault injects into, which is
// the action carrying the fixture's literal kwarg.
func modelRequest(t *testing.T) Request {
	t.Helper()
	set := loadSet(t, generators())
	return Request{
		RecordID:    "PR-014",
		Intent:      "Give the user model a default.",
		Constraints: []string{"One default and nothing else."},
		Files:       created(t, set, "app/models/user.rb@rails"),
	}
}

// stub is a Client that returns canned responses and remembers what it was
// asked. It is what makes the retry loop testable without a server.
type stub struct {
	responses []string
	err       error
	seen      [][]Message
}

func (s *stub) Complete(_ context.Context, messages []Message) (Completion, error) {
	s.seen = append(s.seen, slices.Clone(messages))
	if s.err != nil {
		return Completion{}, s.err
	}
	if len(s.seen) > len(s.responses) {
		return Completion{}, fmt.Errorf("called %d times, but only %d responses were staged", len(s.seen), len(s.responses))
	}
	// Token counts a server would have reported, so the loop's accounting is
	// exercised rather than assumed. They are deliberately unequal per call.
	return Completion{
		Content:          s.responses[len(s.seen)-1],
		PromptTokens:     100 * len(s.seen),
		CompletionTokens: 10 * len(s.seen),
	}, nil
}

func selectWith(t *testing.T, req Request, retries int, responses ...string) ([]recording.Invocation, *stub, error) {
	t.Helper()
	answer, client, err := selectAnswer(t, req, retries, responses...)
	return answer.Invocations, client, err
}

// selectAnswer keeps what the answer cost, which selectWith drops for the many
// tests that only care about what was chosen.
func selectAnswer(t *testing.T, req Request, retries int, responses ...string) (Answer, *stub, error) {
	t.Helper()
	client := &stub{responses: responses}
	got, err := Select(context.Background(), client, req, Options{Retries: retries})
	return got, client, err
}

const validResponse = `{"invocations":[{"action":"createControllerMethod",` +
	`"kwargs":{"controller":"users","name":"index","collection":"users"}}]}`

// wantViolation asserts a response was rejected and that the diagnostic names
// the rule and everything the model would need to correct it. A check that does
// not say what was wrong is a defect: the retry loop's whole value is the
// specificity of what it re-prompts with.
func wantViolation(t *testing.T, req Request, response, rule string, needles ...string) {
	t.Helper()

	_, client, err := selectWith(t, req, 0, response)
	if err == nil {
		t.Fatalf("expected %s to be rejected", rule)
	}
	msg := err.Error()
	if !strings.Contains(msg, rule) {
		t.Errorf("rejection does not name the rule %q:\n%s", rule, msg)
	}
	for _, n := range needles {
		if !strings.Contains(msg, n) {
			t.Errorf("rejection does not mention %q:\n%s", n, msg)
		}
	}
	if len(client.seen) != 1 {
		t.Errorf("with no retries the model should be called once, got %d calls", len(client.seen))
	}
}

func TestValidOutputPasses(t *testing.T) {
	got, _, err := selectWith(t, railsRequest(t), 0, validResponse)
	if err != nil {
		t.Fatalf("a valid response was rejected: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d invocations, want 1", len(got))
	}
	if got[0].Action != "createControllerMethod" || got[0].Kwargs["controller"] != "users" {
		t.Errorf("invocation did not survive decoding: %+v", got[0])
	}
}

// A record whose intent needs no action is legitimate, and an empty array is a
// considered answer rather than a failure to answer.
func TestEmptyInvocationListIsValid(t *testing.T) {
	got, _, err := selectWith(t, railsRequest(t), 0, `{"invocations": []}`)
	if err != nil {
		t.Fatalf("an empty invocation list was rejected: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d invocations, want none", len(got))
	}
}

// The PRD describes a bare array and the wire format wraps it, because JSON
// object mode requires an object at the root. Both are accepted, and both are
// validated identically (prov-2026-abd43bb4).
func TestBothEnvelopeShapesAreAccepted(t *testing.T) {
	// collection is bound because the index template renders it. This test is
	// about the envelope, and an invocation that also tripped the derived
	// requirement check would stop isolating it.
	bare := `[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"index","collection":"users"}}]`

	got, _, err := selectWith(t, railsRequest(t), 0, bare)
	if err != nil {
		t.Fatalf("a bare array was rejected: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d invocations, want 1", len(got))
	}

	// Fenced, which is the formatting habit a model asked for JSON falls into.
	fenced := "```json\n" + validResponse + "\n```"
	if _, _, err := selectWith(t, railsRequest(t), 0, fenced); err != nil {
		t.Fatalf("a fenced response was rejected: %v", err)
	}
}

// The tolerance now includes a key nobody asked for. abd43bb4 rejected one, on
// the reasoning that ignoring it would discard the clearest signal that the
// model had misread the task; measured, it discarded the opposite - five of
// fifty samples on the 4B described row never reached selection for it, while
// the other forty-five produced 900 of 900 assertions (prov-2026-ddb398f3).
//
// What still rejects is an answer that cannot be read, which is a different
// thing from one that says more than it was asked to.
func TestAKeyBeyondTheSchemaIsIgnored(t *testing.T) {
	req := railsRequest(t)

	// validResponse with an annotation added, so the only difference from an
	// answer that is accepted is the key nobody asked for.
	annotated := `{"invocations":[{"action":"createControllerMethod",` +
		`"kwargs":{"controller":"users","name":"index","collection":"users"},` +
		`"reason":"the intent asks for a list"}]}`
	if _, _, err := selectWith(t, req, 0, annotated); err != nil {
		t.Fatalf("an annotated invocation was rejected: %v", err)
	}

	// The same rule at the envelope, or it would depend on nesting depth.
	if _, _, err := selectWith(t, req, 0, `{"invocations":[],"notes":"nothing to do"}`); err != nil {
		t.Fatalf("an annotated envelope was rejected: %v", err)
	}
}

// A mistyped kwargs key is the case abd43bb4 wanted caught, and it still is -
// by the rule whose message names what is missing rather than one that says the
// entry was unreadable.
func TestAMistypedKwargsKeyIsRejectedAsMissingArguments(t *testing.T) {
	wantViolation(t, railsRequest(t),
		`{"invocations":[{"action":"createControllerMethod",`+
			`"kwarg":{"controller":"users","name":"index","collection":"users"}}]}`,
		"missing_kwarg")
}

// Leniency is not acceptance of an unreadable answer.
func TestAnAnswerThatCannotBeReadIsStillRejected(t *testing.T) {
	req := railsRequest(t)

	wantViolation(t, req, `not json at all`, "response_shape")
	wantViolation(t, req, ``, "empty_response")
	wantViolation(t, req, `{"actions":[]}`, "response_shape", "invocations")
	wantViolation(t, req,
		`{"invocations":[{"action":"createControllerMethod","kwargs":"users"}]}`,
		"invocation_shape")
	wantViolation(t, req, `{"invocations":[]} trailing`, "response_shape")
}

// Rule one: the action exists and is exposed. An unexposed action is absent
// from the catalog, so naming it is indistinguishable from naming something
// that does not exist - which is the point of the tier.
func TestUnknownAndUnexposedActionsAreRejected(t *testing.T) {
	req := railsRequest(t)

	wantViolation(t, req,
		`{"invocations":[{"action":"createSerializer","kwargs":{}}]}`,
		"unknown_action", "createSerializer", "createControllerMethod")

	wantViolation(t, req,
		`{"invocations":[{"action":"hiddenHelper","kwargs":{"controller":"users"}}]}`,
		"unknown_action", "hiddenHelper")
}

// Rule two: every required kwarg is bound, and the diagnostic names which.
func TestMissingRequiredKwargIsRejected(t *testing.T) {
	wantViolation(t, railsRequest(t),
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users"}}]}`,
		"missing_kwarg", `"name"`)
}

// A required kwarg the model mentioned and left empty is not bound. It passed
// every check before this: it is a string, so the type is right, and it is
// present, so nothing called it missing - and the rails standard's addColumn
// rendered "t.string :title, null: false, default:", which is Ruby that does
// not parse (prov-2026-9a491128).
func TestARequiredKwargBoundToNothingIsRejected(t *testing.T) {
	req := railsRequest(t)

	wantViolation(t, req,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":""}}]}`,
		RuleEmptyKwarg, `"name"`)

	// Whitespace is a value nobody chose that renders as one.
	wantViolation(t, req,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"   "}}]}`,
		RuleEmptyKwarg, `"name"`)

	// Distinct from missing_kwarg, because a kwarg nobody wrote and one
	// deliberately emptied are different mistakes - and a shared slug would
	// make them one row in every count drawn from a results file.
	_, _, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":""}}]}`)
	if err == nil {
		t.Fatal("an empty required kwarg was accepted")
	}
	if strings.Contains(err.Error(), RuleMissingKwarg) {
		t.Errorf("an empty kwarg was reported as missing:\n%s", err)
	}
}

// Rule three: no kwarg the action does not declare. The diagnostic lists what
// it does declare, so the correction does not need another round trip.
func TestUnknownKwargIsRejected(t *testing.T) {
	wantViolation(t, railsRequest(t),
		`{"invocations":[{"action":"createModelClass","kwargs":{"name":"user","table":"users"}}]}`,
		"unknown_kwarg", `"table"`, `"name"`)
}

// Rule four: every value matches its declared type, and the diagnostic says
// what it actually was - "bound to the string \"3\"" is actionable where "wrong
// type" is not.
func TestKwargTypeMismatchIsRejected(t *testing.T) {
	req := railsRequest(t)

	wantViolation(t, req,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"index","limit":"3"}}]}`,
		"kwarg_type", "limit", "int", `"3"`)

	wantViolation(t, req,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"index","cached":"yes"}}]}`,
		"kwarg_type", "cached", "bool")

	wantViolation(t, req,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"index","only":"show"}}]}`,
		"kwarg_type", "only", "list")

	// JSON has one number type, so int means integral rather than differently
	// encoded.
	wantViolation(t, req,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"index","limit":2.5}}]}`,
		"kwarg_type", "limit")
}

// A literal is carried on the wire as a JSON string, so a string satisfies it
// and every other wire type does not. This is the one place a declared type and
// a wire type differ, which is why it is asserted rather than left to the
// equality that covers the other four (prov-2026-e3d7b9ac).
func TestLiteralIsBoundAsAString(t *testing.T) {
	req := modelRequest(t)

	// Ruby source, bound as a string, emitted verbatim. Accepted.
	got, _, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"declareDefault","kwargs":{"name":"user","value":"nil"}}]}`)
	if err != nil {
		t.Fatalf("a literal bound to a string should be accepted: %v", err)
	}
	if len(got) != 1 || got[0].Kwargs["value"] != "nil" {
		t.Errorf("the bound literal did not survive validation: %#v", got)
	}

	// JSON has no way to carry source, so every non-string is a mistake the
	// model can act on - including null, which is what a model reaches for when
	// the literal it wants to write is the target language's absence.
	wantViolation(t, req,
		`{"invocations":[{"action":"declareDefault","kwargs":{"name":"user","value":null}}]}`,
		"kwarg_type", "value", "literal")

	wantViolation(t, req,
		`{"invocations":[{"action":"declareDefault","kwargs":{"name":"user","value":3}}]}`,
		"kwarg_type", "value", "literal")

	wantViolation(t, req,
		`{"invocations":[{"action":"declareDefault","kwargs":{"name":"user","value":true}}]}`,
		"kwarg_type", "value", "literal")
}

// The catalog carries the declared type through untouched. Rewriting literal to
// string on the way out would put the model back where this type exists to move
// it from, and nothing downstream would notice (prov-2026-e3d7b9ac).
func TestLiteralReachesTheModelAsLiteral(t *testing.T) {
	req := modelRequest(t)

	_, client, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"declareDefault","kwargs":{"name":"user","value":"nil"}}]}`)
	if err != nil {
		t.Fatalf("selection failed: %v", err)
	}

	var prompt string
	for _, m := range client.seen[0] {
		prompt += m.Content
	}
	if !strings.Contains(prompt, "literal") {
		t.Errorf("the prompt never says literal, so the model cannot act on the type:\n%s", prompt)
	}

	// And it rules the wrong shape out rather than only describing the right
	// one. A model bound {"literal": "false"} to a literal kwarg once - it made
	// an envelope out of the type's own name, and the whole object rendered
	// into the file (prov-2026-32b0f3e9).
	if !strings.Contains(prompt, "not an object naming the type") {
		t.Errorf("the prompt does not rule out wrapping a literal in the type's name:\n%s", prompt)
	}
}

// An optional kwarg the template renders is normally rejected when unbound, so
// that Phase 6 never halts on something the retry loop could have fixed. A
// declared default is what makes the same shape safe, and Phase 5 has to know
// that or the default is unreachable through the only path to it
// (prov-2026-f03916ba).
func TestADefaultMakesAnOmissionLegal(t *testing.T) {
	req := modelRequest(t)

	// The same omission, without a default: still rejected.
	wantViolation(t, req,
		`{"invocations":[{"action":"declareDefault","kwargs":{"name":"user"}}]}`,
		"missing_kwarg", "value")

	got, _, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"declareDefaulted","kwargs":{"name":"user"}}]}`)
	if err != nil {
		t.Fatalf("omitting a kwarg that declares a default should be accepted: %v", err)
	}

	// And the recording still says the model bound nothing. The default is
	// applied at resolution, so an answer that omitted the kwarg reads as
	// having omitted it however the file eventually renders.
	if _, bound := got[0].Kwargs["value"]; bound {
		t.Errorf("validation wrote the default into the recorded answer: %#v", got[0].Kwargs)
	}
}

// The model is told what omission resolves to. A default it cannot see makes
// silence a gamble rather than a choice.
func TestTheDefaultReachesTheModel(t *testing.T) {
	req := modelRequest(t)

	_, client, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"declareDefaulted","kwargs":{"name":"user"}}]}`)
	if err != nil {
		t.Fatalf("selection failed: %v", err)
	}

	var prompt string
	for _, m := range client.seen[0] {
		prompt += m.Content
	}
	if !strings.Contains(prompt, `"default"`) {
		t.Errorf("the catalog never carries the declared default:\n%s", prompt)
	}
}

// Rule five: a discriminator value with no template is legal when the package
// ships a fallback and an error when it does not. Both halves are asserted,
// because the difference between them is the whole reason the catalog reports
// which is the case (prov-2026-21031113).
func TestVariantIsCheckedAgainstTheFallback(t *testing.T) {
	req := railsRequest(t)

	// createControllerMethod ships _default.rb, so an uncovered value is a
	// knowing fallback rather than a violation.
	if _, _, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"search"}}]}`); err != nil {
		t.Errorf("an uncovered variant with a fallback should be allowed: %v", err)
	}

	// addCallback ships no _default, so the same shape of answer is an error
	// naming what it does cover.
	wantViolation(t, req,
		`{"invocations":[{"action":"addCallback","kwargs":{"controller":"users","hook":"around"}}]}`,
		"variant", "around", `"before"`, `"after"`)
}

// Rule six: the file an action resolves to is one this record authorized. This
// is where a record that named an implementation file but forgot its companion
// lands - before a retry is spent, naming the file that is missing.
func TestActionResolvingToAnUnauthorizedPathIsRejected(t *testing.T) {
	set := loadSet(t, generators())
	req := Request{
		RecordID: "PR-014",
		Intent:   "Add a users controller.",
		Files:    created(t, set, "app/controllers/users_controller.rb@rails"),
	}

	// createModelClass renders app/models/user.rb, which this record did not
	// authorize.
	wantViolation(t, req,
		`{"invocations":[{"action":"createModelClass","kwargs":{"name":"user"}}]}`,
		"unauthorized_path", "app/models/user.rb")
}

// A composite is checked against every file its children touch. Omitting one of
// them from affected_scope fails here rather than halfway through Phase 7.
func TestCompositeIsCheckedAgainstEveryChildsTarget(t *testing.T) {
	set := loadSet(t, generators())
	partial := Request{
		RecordID: "PR-020",
		Intent:   "Provision a step.",
		Files:    created(t, set, "Units/billing-runs/Manifest.crn@cairn"),
	}

	wantViolation(t, partial,
		`{"invocations":[{"action":"provisionStep","kwargs":{"unit":"billing-run","step":"charge","name":"retry-limit"}}]}`,
		"unauthorized_path", "Shared/retry-limits.crn")

	// With both files authorized, the same invocation validates.
	whole := partial
	whole.Files = created(t, set,
		"Units/billing-runs/Manifest.crn@cairn", "Shared/retry-limits.crn@cairn")
	if _, _, err := selectWith(t, whole, 0,
		`{"invocations":[{"action":"provisionStep","kwargs":{"unit":"billing-run","step":"charge","name":"retry-limit"}}]}`); err != nil {
		t.Errorf("a composite whose children are both authorized was rejected: %v", err)
	}
}

// A response with three faults re-prompts with three. Spending one model call
// per fault is precisely the cost deterministic validation exists to avoid.
func TestEveryInvocationIsChecked(t *testing.T) {
	_, _, err := selectWith(t, railsRequest(t), 0, `{"invocations":[
		{"action":"nope","kwargs":{}},
		{"action":"createControllerMethod","kwargs":{"controller":"users"}},
		{"action":"createModelClass","kwargs":{"name":"user","table":"users"}}
	]}`)
	if err == nil {
		t.Fatal("expected the response to be rejected")
	}

	msg := err.Error()
	for _, rule := range []string{"unknown_action", "missing_kwarg", "unknown_kwarg"} {
		if !strings.Contains(msg, rule) {
			t.Errorf("only some faults were reported; %q is missing:\n%s", rule, msg)
		}
	}
}

// A diagnostic must never blame a path for a kwarg's fault. The schema checks
// run first, so an invocation missing its required argument is told that, not
// told about the file the unrendered pattern produced (prov-2026-9dcf2658).
func TestSchemaFaultsAreReportedBeforePathFaults(t *testing.T) {
	set := loadSet(t, generators())
	req := Request{
		RecordID: "PR-014",
		Intent:   "Add a users controller.",
		Files:    created(t, set, "app/controllers/users_controller.rb@rails"),
	}

	_, _, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"createModelClass","kwargs":{}}]}`)
	if err == nil {
		t.Fatal("expected the response to be rejected")
	}
	if !strings.Contains(err.Error(), "missing_kwarg") {
		t.Errorf("the missing kwarg was not reported:\n%s", err)
	}
	if strings.Contains(err.Error(), "unauthorized_path") {
		t.Errorf("a missing kwarg was reported as a path problem:\n%s", err)
	}
}

// The loop: a rejected response is re-prompted with what was wrong, and the
// rejected response stays in the conversation so the correction refers to
// something the model can still see.
func TestRejectedResponseIsRePromptedWithItsViolations(t *testing.T) {
	bad := `{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users"}}]}`

	got, client, err := selectWith(t, railsRequest(t), 3, bad, validResponse)
	if err != nil {
		t.Fatalf("the retry should have succeeded: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d invocations, want 1", len(got))
	}
	if len(client.seen) != 2 {
		t.Fatalf("the model was called %d times, want 2", len(client.seen))
	}

	second := client.seen[1]
	if len(second) != 4 {
		t.Fatalf("the retry conversation has %d messages, want the prompt plus the rejected answer and the correction", len(second))
	}
	if second[2].Role != RoleAssistant || second[2].Content != bad {
		t.Errorf("the rejected response was not kept in the conversation: %+v", second[2])
	}
	if !strings.Contains(second[3].Content, "name") {
		t.Errorf("the correction does not say what was missing:\n%s", second[3].Content)
	}
}

// Exhaustion reports every attempt, not just the last. A model making a
// different mistake each time and one making the same mistake three times are
// different problems with different fixes.
func TestRetryExhaustionReportsEveryAttempt(t *testing.T) {
	_, client, err := selectWith(t, railsRequest(t), 2,
		`{"invocations":[{"action":"nope","kwargs":{}}]}`,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users"}}]}`,
		`{"invocations":[{"action":"createModelClass","kwargs":{"name":"user","table":"users"}}]}`,
	)
	if err == nil {
		t.Fatal("expected exhaustion to halt the run")
	}
	if len(client.seen) != 3 {
		t.Fatalf("the model was called %d times, want 3 (the first attempt plus two retries)", len(client.seen))
	}

	msg := err.Error()
	for _, want := range []string{"attempt 1", "attempt 2", "attempt 3", "unknown_action", "missing_kwarg", "unknown_kwarg"} {
		if !strings.Contains(msg, want) {
			t.Errorf("exhaustion does not report %q:\n%s", want, msg)
		}
	}
}

// A transport failure is not the model's to correct, and retrying it would
// consume the budget reserved for the mistakes that are.
func TestTransportFailureIsNotRetried(t *testing.T) {
	client := &stub{err: errors.New("connection refused")}

	_, err := Select(context.Background(), client, railsRequest(t), Options{Retries: 3})
	if err == nil {
		t.Fatal("expected the model call to fail")
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("the transport error was not reported: %v", err)
	}
	if len(client.seen) != 1 {
		t.Errorf("a transport failure was retried %d times", len(client.seen)-1)
	}
}

// The prompt carries the record's own words, the files the run created, and the
// catalog. The catalog is embedded as exactly the bytes sedum actions --json
// prints, which is what keeps that command evidence rather than a second
// rendering.
func TestPromptCarriesTheRecordAndTheCatalogVerbatim(t *testing.T) {
	req := railsRequest(t)
	packages := expand.Packages(req.Files)
	cat := catalog.Build(packages, catalog.Options{})

	messages, err := Prompt(req, cat)
	if err != nil {
		t.Fatalf("Prompt: %v", err)
	}
	if len(messages) != 2 || messages[0].Role != RoleSystem || messages[1].Role != RoleUser {
		t.Fatalf("prompt is %d messages in roles %v", len(messages), messages)
	}

	user := messages[1].Content
	for _, want := range []string{
		"Add a read-only users controller.",
		"No writes.",
		"app/controllers/users_controller.rb",
	} {
		if !strings.Contains(user, want) {
			t.Errorf("prompt does not carry %q:\n%s", want, user)
		}
	}

	payload, err := cat.JSON()
	if err != nil {
		t.Fatalf("catalog JSON: %v", err)
	}
	if !strings.Contains(user, string(payload)) {
		t.Errorf("prompt does not embed the catalog payload verbatim:\n%s", user)
	}

	// An unexposed action must not be visible anywhere in the prompt.
	if strings.Contains(user, "hiddenHelper") {
		t.Errorf("an unexposed action reached the prompt:\n%s", user)
	}
}

// Sedum's core carries no target knowledge, and the prompt is the place that
// constraint is easiest to break: one helpful sentence about controllers would
// put a framework's vocabulary in the binary. Every target-specific word the
// model sees has to come from the record or the package.
func TestSystemPromptNamesNoTarget(t *testing.T) {
	lowered := strings.ToLower(systemPrompt)
	for _, forbidden := range []string{
		"rails", "ruby", "golang", "python", "java", "react",
		"controller", "model class", "header file", "import statement",
	} {
		if strings.Contains(lowered, forbidden) {
			t.Errorf("the system prompt names a target concept (%q), which puts target knowledge in Sedum's core", forbidden)
		}
	}
}

// A path a package declares unmanaged is named rather than hidden. The intent
// may refer to it, and a model that could not see it would keep reaching for it
// with an action that cannot get there.
func TestPromptNamesUnmanagedPathsAsUnreachable(t *testing.T) {
	set := loadSet(t, generators())
	req := Request{
		RecordID: "PR-014",
		Intent:   "Add a users controller and the dependency it needs.",
		Files: append(created(t, set, "app/controllers/users_controller.rb@rails"),
			resolve.File{Resolution: resolve.Resolution{
				RecordID: "PR-014", Path: "Gemfile", Unmanaged: true, UnmanagedBy: "rails", UnmanagedAs: "Gemfile",
			}}),
	}

	messages, err := Prompt(req, catalog.Build(expand.Packages(req.Files), catalog.Options{}))
	if err != nil {
		t.Fatalf("Prompt: %v", err)
	}

	user := messages[1].Content
	if !strings.Contains(user, "Gemfile") {
		t.Errorf("prompt does not mention the unmanaged path:\n%s", user)
	}
	if !strings.Contains(user, "not written by this tool") {
		t.Errorf("prompt does not say the unmanaged path is unreachable:\n%s", user)
	}
}

// The failure this record was written from: a discriminated action shares one
// kwarg schema across every variant, so a value that index needs and show
// forbids can only be declared optional. The catalog said optional, the model
// omitted it, Phase 5 passed, and Phase 6 halted rendering a template - with
// the retry loop already skipped, because nothing was wrong with the selection
// (prov-2026-369544c1).
func TestAVariantsOwnRequirementIsCheckedInPhase5(t *testing.T) {
	req := railsRequest(t)

	// index renders {{collection}}, so omitting it is now re-promptable rather
	// than a Phase 6 halt. The diagnostic names the kwarg and the variant.
	wantViolation(t, req,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"index"}}]}`,
		"missing_derived_kwarg", "collection", "index")

	// show renders nothing, so the same omission is legal there. The schema
	// really does make it optional, and the derivation must not promote a
	// requirement one variant happens to have into one every variant carries.
	if _, _, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"show"}}]}`,
	); err != nil {
		t.Errorf("a variant whose template renders no optional kwarg was rejected: %v", err)
	}
}

// A discriminator value with no dedicated template inherits the fallback's
// requirements rather than none, because that is the template it will render.
func TestAnUncoveredVariantInheritsTheDefaultsRequirements(t *testing.T) {
	// _default.rb renders {{name|snake}}, and name is the discriminator, so it
	// is bound by construction - the fallback is satisfied here. What must not
	// happen is index's requirement leaking onto a value index does not serve.
	if _, _, err := selectWith(t, railsRequest(t), 0,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"archive"}}]}`,
	); err != nil {
		t.Errorf("an uncovered variant was held to another variant's template: %v", err)
	}
}

// The two halves are reported separately because they have different fixes. A
// declared requirement that is absent is a binding mistake; a derived one is
// usually a schema understating what a variant needs, and a diagnostic that
// conflated them would send an author to the wrong file.
func TestDeclaredAndDerivedRequirementsAreDistinctViolations(t *testing.T) {
	wantViolation(t, railsRequest(t),
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"name":"index"}}]}`,
		"missing_kwarg", "controller")
}

// Under-selection is the one failure class Phase 5 is structurally blind to: a
// short list is valid output, so a run that selects thirteen of fourteen correct
// invocations passes every rule, first try, with no retry. The observation is
// fed back once and the model keeps the judgment (prov-2026-6d87dc11).
func TestAnIncompleteSelectionIsFedBackOnce(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:class_body\nend\n"

	// The first answer selects nothing, leaving class_body unfilled. The second
	// fills it, and the run accepts that.
	got, client, err := selectWith(t, req, 0,
		`{"invocations":[]}`,
		validResponse,
	)
	if err != nil {
		t.Fatalf("a completeness re-prompt was treated as a failure: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d invocations, want the second answer's 1", len(got))
	}
	if len(client.seen) != 2 {
		t.Fatalf("the model was called %d time(s), want 2 - one answer and one completeness re-prompt", len(client.seen))
	}

	// The observation names the file and the anchor, which is the bar every
	// Phase 5 diagnostic is held to.
	followUp := client.seen[1][len(client.seen[1])-1].Content
	for _, want := range []string{"users_controller.rb", "class_body"} {
		if !strings.Contains(followUp, want) {
			t.Errorf("the completeness note does not mention %q:\n%s", want, followUp)
		}
	}
}

// The model keeps the judgment. An anchor with nothing in it is not necessarily
// a mistake - a template may plant a region a given change does not need - so a
// second declining answer stands and the run continues.
func TestADecliningAnswerStandsAndTheRunContinues(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:class_body\nend\n"

	got, client, err := selectWith(t, req, 0,
		`{"invocations":[]}`,
		`{"invocations":[]}`,
	)
	if err != nil {
		t.Fatalf("a twice-empty selection became an error; an unfilled anchor is never a hard error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d invocations, want the model's own answer of none", len(got))
	}
	if len(client.seen) != 2 {
		t.Errorf("the model was called %d time(s), want exactly 2 - the observation is fed back at most once", len(client.seen))
	}
}

// A response that leaves nothing unfilled is never re-prompted, so the common
// case costs nothing.
func TestACompleteSelectionIsNotReprompted(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:class_body\nend\n"

	_, client, err := selectWith(t, req, 0, validResponse)
	if err != nil {
		t.Fatalf("a complete selection was rejected: %v", err)
	}
	if len(client.seen) != 1 {
		t.Errorf("the model was called %d time(s), want 1 - nothing was left unfilled", len(client.seen))
	}
}

// A difference the model could not have closed earns no call at all. The chi
// eval fixture plants an extensions anchor no action targets - a deliberate
// extension point for a later record - and paid a re-prompt for it on every
// sample, declined every time (prov-2026-206fa618).
func TestAnUnfillableAnchorEarnsNoRePrompt(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:audit_log\nend\n"

	got, client, err := selectWith(t, req, 0, `{"invocations":[]}`)
	if err != nil {
		t.Fatalf("a selection leaving only an unfillable anchor was rejected: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d invocations, want the model's own answer of none", len(got))
	}
	if len(client.seen) != 1 {
		t.Errorf("the model was called %d time(s), want 1 - nothing it selected could have filled audit_log", len(client.seen))
	}
}

// The narrowing is to which anchors are asked about, not to whether they are.
// A file planting both kinds still earns the observation, and the note names
// only the anchor an action could have filled.
func TestTheObservationNamesOnlyTheFillableAnchors(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:audit_log\n" +
		"  # sedum:anchor:class_body\nend\n"

	_, client, err := selectWith(t, req, 0, `{"invocations":[]}`, validResponse)
	if err != nil {
		t.Fatalf("a completeness re-prompt was treated as a failure: %v", err)
	}
	if len(client.seen) != 2 {
		t.Fatalf("the model was called %d time(s), want 2 - class_body was fillable and unfilled", len(client.seen))
	}

	followUp := client.seen[1][len(client.seen[1])-1].Content
	if !strings.Contains(followUp, "class_body") {
		t.Errorf("the completeness note does not mention the fillable anchor:\n%s", followUp)
	}
	if strings.Contains(followUp, "audit_log") {
		t.Errorf("the completeness note asks about an anchor nothing can fill:\n%s", followUp)
	}
}

// What a record cost is the loop's to report, because the loop is the only
// thing that can tell a rejected answer from a completeness observation after
// the fact. A caller left to infer it has nothing but a clock, which is calls
// multiplied by an unknown per-call cost (prov-2026-0811425c).
func TestAnAnswerReportsWhatItCost(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:class_body\nend\n"

	// One rejected answer, then an empty one that earns the observation, then
	// a valid and complete one: three calls, one rejection, one observation.
	bad := `{"invocations":[{"action":"noSuchAction","kwargs":{}}]}`
	got, _, err := selectAnswer(t, req, 1, bad, `{"invocations":[]}`, validResponse)
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if got.Calls != 3 {
		t.Errorf("reported %d calls, want 3", got.Calls)
	}
	if got.Rejected != 1 {
		t.Errorf("reported %d rejections, want 1", got.Rejected)
	}
	// Counted apart from the rejection: the answer that earned it was valid,
	// and it draws from its own budget.
	if got.Completeness != 1 {
		t.Errorf("reported %d completeness calls, want 1", got.Completeness)
	}
}

// Token counts leave with the calls that made them, summed across a record's
// attempts. The stub reports a different count per call, so a sum that only
// kept the last one would not add up.
func TestAnAnswerCarriesWhatTheCallsCostInTokens(t *testing.T) {
	bad := `{"invocations":[{"action":"noSuchAction","kwargs":{}}]}`
	got, _, err := selectAnswer(t, railsRequest(t), 1, bad, validResponse)
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	// Two calls, at 100/10 then 200/20.
	if got.PromptTokens != 300 || got.CompletionTokens != 30 {
		t.Errorf("reported %d prompt and %d completion tokens, want 300 and 30",
			got.PromptTokens, got.CompletionTokens)
	}
}

// A rejected answer carries its token cost too. A sample that spent its whole
// budget being refused is the expensive one, so dropping its counts would lose
// the cost most worth knowing.
func TestARejectionCarriesItsTokenCost(t *testing.T) {
	bad := `{"invocations":[{"action":"noSuchAction","kwargs":{}}]}`
	_, _, err := selectAnswer(t, railsRequest(t), 1, bad, bad)

	var rejected *Rejection
	if !errors.As(err, &rejected) {
		t.Fatalf("want *Rejection, got %T", err)
	}
	if rejected.PromptTokens != 300 || rejected.CompletionTokens != 30 {
		t.Errorf("rejection reports %d prompt and %d completion tokens, want 300 and 30",
			rejected.PromptTokens, rejected.CompletionTokens)
	}
}

// First-call validity survives any retry budget, which is the measurement the
// budget used to destroy: an answer with no rejections validated first try,
// whatever the budget would have allowed.
func TestNoRejectionsMeansItValidatedFirstTry(t *testing.T) {
	got, _, err := selectAnswer(t, railsRequest(t), 2, validResponse)
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if got.Rejected != 0 || got.Calls != 1 {
		t.Errorf("a first-try answer reported %d call(s) and %d rejection(s), want 1 and 0",
			got.Calls, got.Rejected)
	}
}

// A rejected answer is a type, so a caller can tell it from a transport failure
// without matching on text. They are different measurements - a model that
// chose badly and a server that was not there - and the harness classified them
// by matching "did not validate" until this existed.
func TestARejectedAnswerIsATypedError(t *testing.T) {
	bad := `{"invocations":[{"action":"noSuchAction","kwargs":{}}]}`
	_, _, err := selectAnswer(t, railsRequest(t), 1, bad, bad)

	var rejected *Rejection
	if !errors.As(err, &rejected) {
		t.Fatalf("a rejected answer is not a *Rejection: %T %v", err, err)
	}
	if rejected.Calls != 2 || rejected.Rejected != 2 {
		t.Errorf("rejection reports %d call(s) and %d rejection(s), want 2 and 2",
			rejected.Calls, rejected.Rejected)
	}
	if len(rejected.Attempts) != 2 {
		t.Fatalf("rejection kept %d attempt(s), want both", len(rejected.Attempts))
	}
	// The text stays exactly as informative as it was: every attempt's
	// violations, not only the last response's.
	msg := rejected.Error()
	for _, want := range []string{"attempt 1", "attempt 2", "noSuchAction", "did not validate"} {
		if !strings.Contains(msg, want) {
			t.Errorf("rejection text lost %q:\n%s", want, msg)
		}
	}
}

// A transport failure is not a Rejection. Nothing about a connection error is
// the model's to correct, and a caller that classified it as a bad answer would
// be counting an unreachable server as a model that chose badly.
func TestATransportFailureIsNotARejection(t *testing.T) {
	client := &stub{err: errors.New("connection refused")}
	_, err := Select(context.Background(), client, railsRequest(t), Options{Retries: 2})

	var rejected *Rejection
	if errors.As(err, &rejected) {
		t.Error("a transport failure was reported as a rejected answer")
	}
	if err == nil || !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("the transport failure was not reported as itself: %v", err)
	}
}

// A completeness re-prompt is not a validation retry. Exhausting one must not
// consume the other, so an invalid answer still gets its full retry budget after
// the observation has been spent.
func TestCompletenessAndValidationDrawFromSeparateBudgets(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = "class UsersController\n  # sedum:anchor:class_body\nend\n"

	// Empty (triggers the observation), then invalid (consumes the one retry),
	// then valid. With Retries=1 this must succeed: the completeness call is
	// not drawn from the retry budget.
	got, client, err := selectWith(t, req, 1,
		`{"invocations":[]}`,
		`{"invocations":[{"action":"noSuchAction","kwargs":{}}]}`,
		validResponse,
	)
	if err != nil {
		t.Fatalf("the completeness re-prompt consumed a validation retry: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d invocations, want 1", len(got))
	}
	if len(client.seen) != 3 {
		t.Errorf("the model was called %d time(s), want 3", len(client.seen))
	}
}

// Strings only. The string case is grounded in what rendering does to it - an
// empty one renders to nothing where a token is required - and nothing else
// carries that property, so nothing else is read (prov-2026-9a491128).
func TestOnlyAnEmptyStringCountsAsEmpty(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value any
		empty bool
	}{
		{"the empty string", "", true},
		{"whitespace", "  \t ", true},
		{"a value", "nil", false},
		// An empty list is NOT empty for this purpose. A list has no usable
		// rendering either way - ["title","completed"] emits "[title completed]"
		// and [] emits "[]", both unusable as source - so emptiness is not what
		// would make a rendered list wrong, and an action for which an empty
		// array is a meaningful value must not be refused for it.
		{"an empty list", []any{}, false},
		{"a list with a member", []any{"title"}, false},
		// Zero and false are values an author may legitimately want.
		{"zero", float64(0), false},
		{"false", false, false},
	} {
		if got := isEmpty(tc.value); got != tc.empty {
			t.Errorf("%s: isEmpty(%#v) = %v, want %v", tc.name, tc.value, got, tc.empty)
		}
	}
}

// A kwarg nothing needs is untouched. Absence is what optional means, and an
// author who made one optional has already said the template can do without it.
func TestAnEmptyKwargNothingNeedsIsNotRejected(t *testing.T) {
	// The show variant renders nothing, so collection is neither schema- nor
	// derived-required here, and its emptiness is the only thing under test.
	got, _, err := selectWith(t, railsRequest(t), 0,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"show","collection":""}}]}`)
	if err != nil {
		t.Fatalf("an empty optional kwarg was rejected: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d invocations, want 1", len(got))
	}
}

// A kwarg the schema calls optional and the selected template renders is not
// optional, and that holds for what it contains as well as for whether it is
// there. Checking only the schema flag missed exactly the shape this matters
// most for: a discriminated action whose variants need different arguments has
// to declare every one of them optional (prov-2026-9a491128).
func TestAnEmptyDerivedRequiredKwargIsRejected(t *testing.T) {
	// collection is required:false, and the index template renders
	// {{collection|instantize}} - so index needs it and show does not.
	wantViolation(t, railsRequest(t),
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"index","collection":""}}]}`,
		RuleEmptyKwarg, `"collection"`)

	// Absence stays checkDerived's to report. One mistake, one violation.
	_, _, err := selectWith(t, railsRequest(t), 0,
		`{"invocations":[{"action":"createControllerMethod","kwargs":{"controller":"users","name":"index"}}]}`)
	if err == nil {
		t.Fatal("an absent derived-required kwarg was accepted")
	}
	if !strings.Contains(err.Error(), RuleMissingDerivedKwarg) {
		t.Errorf("absence is not reported as a derived requirement:\n%s", err)
	}
	if strings.Contains(err.Error(), RuleEmptyKwarg) {
		t.Errorf("an absent kwarg was also reported as empty:\n%s", err)
	}
}

// The case an empty-list clause would have broken: an action for which an empty
// array is a meaningful value. A list has no usable rendering either way, so
// emptiness is not what would make one wrong, and refusing it would have been a
// rule with no failure behind it (prov-2026-9a491128).
func TestARequiredListBoundToAnEmptyArrayIsAccepted(t *testing.T) {
	got, _, err := selectWith(t, railsRequest(t), 0,
		`{"invocations":[{"action":"tagResource","kwargs":{"controller":"users","tags":[]}}]}`)
	if err != nil {
		t.Fatalf("an empty array was refused for a required list kwarg: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d invocations, want 1", len(got))
	}

	// And the string rule still bites on the same invocation, so this is a
	// narrowing rather than a hole.
	wantViolation(t, railsRequest(t),
		`{"invocations":[{"action":"tagResource","kwargs":{"controller":"","tags":[]}}]}`,
		RuleEmptyKwarg, `"controller"`)
}

// Nothing capped a completion, so a model that does not stop generated until the
// context ran out - 40,076 tokens on a prompt that returns 173 when replayed
// (prov-2026-7f79ba36).
func TestACompletionIsCapped(t *testing.T) {
	if got := maxCompletionTokens(); got != 16384 {
		t.Errorf("default cap is %d, want 16384", got)
	}

	// Above every completion this harness has legitimately received: the
	// largest on record is 9,759, a baseline arm writing six whole files.
	if maxCompletionTokens() <= 9759 {
		t.Errorf("cap %d is not above the largest answer on record", maxCompletionTokens())
	}

	t.Setenv(EnvMaxTokens, "2048")
	if got := maxCompletionTokens(); got != 2048 {
		t.Errorf("override gave %d, want 2048", got)
	}

	// A value nobody can generate against is ignored rather than obeyed.
	t.Setenv(EnvMaxTokens, "0")
	if got := maxCompletionTokens(); got != 16384 {
		t.Errorf("zero gave %d, want the default", got)
	}
	t.Setenv(EnvMaxTokens, "not a number")
	if got := maxCompletionTokens(); got != 16384 {
		t.Errorf("junk gave %d, want the default", got)
	}
}

// Validate is the entry point replay uses. It exists so that there is one
// validator rather than two that agree on the day they were written, so what is
// asserted here is that reaching it from outside produces the same verdicts the
// model path gets (prov-2026-6903db40).

func TestValidateRejectsAnUnknownActionForANonModelCaller(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, "app/controllers/users_controller.rb@rails")

	violations := Validate(files, []recording.Invocation{
		{Action: "createControllerMethd", Kwargs: map[string]any{"controller": "users", "name": "index"}},
	})

	if len(violations) == 0 {
		t.Fatal("an action that does not exist was accepted")
	}
	if violations[0].Rule != RuleUnknownAction {
		t.Errorf("rule = %v, want the same rule the model path reports", violations[0].Rule)
	}
}

func TestValidateRejectsAMissingRequiredKwarg(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, "app/controllers/users_controller.rb@rails")

	violations := Validate(files, []recording.Invocation{
		{Action: "createControllerMethod", Kwargs: map[string]any{"controller": "users"}},
	})

	if len(violations) == 0 {
		t.Fatal("an invocation omitting a required kwarg was accepted")
	}
}

// A correct invocation passes, so the checks above are rejecting the fault
// rather than rejecting everything.
func TestValidateAcceptsACorrectInvocation(t *testing.T) {
	set := loadSet(t, generators())
	files := created(t, set, "app/controllers/users_controller.rb@rails")

	if violations := Validate(files, []recording.Invocation{
		{Action: "createControllerMethod", Kwargs: map[string]any{"controller": "users", "name": "index", "collection": "users"}},
	}); len(violations) > 0 {
		t.Errorf("a correct invocation was rejected: %v", violations)
	}
}

// An action's anchor declares the region kind it needs, and a file template
// declares which regions the files it creates carry. Applicability is the
// relation between them, checked per invocation now that a target can be
// aimed rather than derived (prov-2026-14c832bf).
func TestAnActionAimedAtAFileThatCarriesADifferentAnchorIsRejected(t *testing.T) {
	set := loadSet(t, freeGenerators())
	req := Request{
		RecordID: "PR-014",
		Intent:   "Import the settings helper.",
		Files: []resolve.File{
			rendered(t, set, "config/app.ts", "free", "// sedum:anchor:settings\n"),
		},
	}

	wantViolation(t, req,
		`{"invocations":[{"action":"addImport","kwargs":{"file":"config/app.ts","symbol":"Helper"}}]}`,
		RuleAnchorUnplanted, "config/app.ts", "imports", "settings")
}

// The two cases have different fixes, so the diagnostic distinguishes them. A
// file whose template plants nothing is a legitimate shape - a fallback
// template is boilerplate for paths no action targets - and telling an author
// it is missing one anchor would send them looking for the wrong thing.
func TestAFileThatPlantsNoAnchorsAtAllSaysSo(t *testing.T) {
	set := loadSet(t, freeGenerators())
	req := Request{
		RecordID: "PR-014",
		Intent:   "Import the settings helper.",
		Files: []resolve.File{
			rendered(t, set, "notes.ts", "free", "// generated\n"),
		},
	}

	wantViolation(t, req,
		`{"invocations":[{"action":"addImport","kwargs":{"file":"notes.ts","symbol":"Helper"}}]}`,
		RuleAnchorUnplanted, "no injection points at all")
}

func TestAnActionAimedAtAFileCarryingItsAnchorIsAccepted(t *testing.T) {
	set := loadSet(t, freeGenerators())
	req := Request{
		RecordID: "PR-014",
		Intent:   "Import the settings helper.",
		Files: []resolve.File{
			rendered(t, set, "src/app.ts", "free", "// sedum:anchor:imports\n"),
		},
	}

	invocations, _, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"addImport","kwargs":{"file":"src/app.ts","symbol":"Helper"}}]}`)
	if err != nil {
		t.Fatalf("a free-target action aimed at a file carrying its anchor was rejected: %v", err)
	}
	if len(invocations) != 1 {
		t.Fatalf("got %d invocations, want 1", len(invocations))
	}
}

// Only files this run rendered can be spoken about. A path Sedum produced no
// content for is one it has observed nothing about, and asserting an anchor is
// missing from it would be a claim rather than a reading. Phase 7 stays the
// backstop there, which is where this case landed before the check existed.
func TestAFileThisRunDidNotRenderIsNotJudged(t *testing.T) {
	set := loadSet(t, freeGenerators())
	req := Request{
		RecordID: "PR-014",
		Intent:   "Import the settings helper.",
		Files:    created(t, set, "src/app.ts@free"),
	}

	if _, _, err := selectWith(t, req, 0,
		`{"invocations":[{"action":"addImport","kwargs":{"file":"src/app.ts","symbol":"Helper"}}]}`); err != nil {
		t.Fatalf("a file with no rendered content was judged anyway: %v", err)
	}
}

// A free-target action's pattern says nothing about which file to aim it at, so
// the file list has to. Without this the model is told an anchor and given no
// way to find a file that carries one.
func TestThePromptNamesEachFilesInjectionPoints(t *testing.T) {
	set := loadSet(t, freeGenerators())
	req := Request{
		RecordID: "PR-014",
		Intent:   "Import the settings helper.",
		Files: []resolve.File{
			rendered(t, set, "src/app.ts", "free", "// sedum:anchor:imports\n"),
			rendered(t, set, "notes.ts", "free", "// generated\n"),
		},
	}

	packages := expand.Packages(req.Files)
	messages, err := Prompt(req, catalog.Build(packages, catalog.Options{}))
	if err != nil {
		t.Fatalf("Prompt: %v", err)
	}
	user := messages[len(messages)-1].Content

	if !strings.Contains(user, "src/app.ts (package free; injection points: imports)") {
		t.Errorf("the file list does not name the anchors a file carries:\n%s", user)
	}
	// A file carrying none is still listed, and says nothing it does not know.
	if !strings.Contains(user, "notes.ts (package free)") {
		t.Errorf("a file carrying no anchors was described as if it did:\n%s", user)
	}

	// The bare pattern is suppressed for the model and the anchor is not.
	if strings.Contains(user, `"{{file}}"`) {
		t.Errorf("the prompt shows a bare target pattern:\n%s", user)
	}
	if !strings.Contains(user, `"anchor"`) {
		t.Errorf("the prompt does not carry the anchor:\n%s", user)
	}
}
