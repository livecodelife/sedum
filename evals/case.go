package evals

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Case is one cell of the matrix: an application, generated a particular way, by
// a particular model.
//
// Every axis is declared even where only one value of it exists today. Adding a
// model or a framework should be a new file rather than a schema change, because
// a schema that grows once results exist means rewriting the results too.
type Case struct {
	ID string `yaml:"id"`

	// Application is what is being built, and how hard it is. Complexity is an
	// ordinal tier rather than a score: a todo list, a URL shortener and a user
	// management service are 1, 2 and 3 because they are ordered, not because
	// the third is three times the first.
	Application struct {
		Name       string `yaml:"name"`
		Complexity int    `yaml:"complexity"`
	} `yaml:"application"`

	// Framework and Language name the target stack. Both are labels for
	// grouping results and Sedum reads neither - it has no language knowledge
	// and resolves packages by file extension.
	//
	// Language is carried separately because it groups differently: Rails and
	// Sinatra are one language and two frameworks, and a result that held
	// across both would be saying something about Ruby rather than about
	// either.
	Framework string `yaml:"framework"`

	// Variables are the run values this case's generator packages declare and
	// cannot know - a Go module path, a root namespace (prov-2026-6fc3d13d).
	//
	// On the case rather than on the runner, because a variable's value is a
	// property of the project being generated and for the eval the project is
	// the case. A flag would make it a property of the invocation, so two runs
	// of one case could differ in a way no stored entry records - and the entry
	// is what makes a measurement re-runnable (prov-2026-1c33a50b).
	Variables map[string]string `yaml:"variables,omitempty"`
	Language  string            `yaml:"language"`

	// Tightness is how much the generator package constrains the model:
	// "defined" for a package whose actions and variants cover the intended
	// output closely, "loose" for one that leaves more to a fallback.
	//
	// This is the axis Generators swaps along. Two cases naming one record and
	// two package sets differ in exactly one thing, which is what makes the
	// comparison mean anything - and it is why generators/ holds package sets
	// rather than packages.
	Tightness string `yaml:"tightness"`

	// Arm is "sedum", "baseline" or "intent" - a ladder of what the model was
	// given. sedum has a record and a catalog; baseline has the record and no
	// catalog; intent has the record's intent and nothing else, not its
	// constraints and not its file list.
	//
	// The two lower rungs exist because a selection rate has no meaning without
	// one, and because a baseline handed six constraints and a file list is a
	// model given a precise specification rather than a model given no tooling
	// (prov-2026-a4dbe65c, prov-2026-672c6471).
	Arm string `yaml:"arm"`

	// Generators and Records are directories, resolved against the root passed
	// to Load - evals/fixtures by default, so a case names a vendored snapshot
	// rather than a path into someone's workspace. An absolute path still
	// works, for pointing a case at a live project deliberately.
	//
	// Generators names a package *set*: the directory --generators is given,
	// whose subdirectories are the packages. Swapping it is how one record is
	// measured against two packages.
	Generators string `yaml:"generators"`
	Records    string `yaml:"records"`

	// Only narrows the run to named records, as --only does. A fixture holding
	// several records costs one model call each per sample, so a case that
	// cares about one of them says so rather than paying for the rest.
	Only []string `yaml:"only"`

	// Models are the models to measure, each run independently.
	Models []Model `yaml:"models"`

	// Check is the target's own parser, keyed by the extension it reads, run
	// over what Sedum wrote. Optional: a case that declares none is measured
	// on every other signal and reports no syntactic validity.
	//
	// It lives on the case rather than in the harness because it is target
	// knowledge, and a harness carrying a switch on language is one that has to
	// be edited to add a framework - which is the cost the matrix exists to
	// avoid (prov-2026-d61010a4).
	Check Check `yaml:"check,omitempty"`

	// PerSample is how long one of this case's samples takes, and it sizes both
	// halves of the plan. Optional: a case that declares none is planned with
	// the harness's own constant, exactly as it was before this field existed.
	//
	// It is a property of the case rather than of the model row. The chi record
	// drives more files and a larger catalog than the rails one, so a sample is
	// a longer prompt and a longer answer: on the same 14B row, rails samples
	// run 64-136s and chi samples run 154-494s. A single constant spends one
	// case's number on another (prov-2026-59ed14d5).
	//
	// Declare a typical sample rather than the slowest one. Headroom makes the
	// ceiling safe, and an estimate that overstates is one nobody starts a run
	// from.
	PerSample time.Duration `yaml:"per_sample,omitempty"`

	Expect Expectations `yaml:"expect"`
}

// WithoutPackage reports whether the arm has no generator package, and so no
// selection, binding, anchor fill or idempotency to score - which makes
// behaviour the only rung that can measure it.
//
// A predicate rather than a comparison at each site, because the sites that ask
// are asking that question and not which arm it is: a third arm that is also
// package-free would otherwise have to be remembered in four places.
func (c Case) WithoutPackage() bool {
	return c.Arm == "baseline" || c.Arm == "intent"
}

// HeldToPaths reports whether an answer from this arm is scored against the
// paths the record authorizes.
//
// The intent arm is not. It is told no paths, so filtering against a list it
// never saw would score whether it guessed the ones this standard happens to
// use (prov-2026-672c6471) - and a rate over a list that admits everything is
// wrote/wrote on every run, which is arithmetic rather than a measurement.
//
// A predicate for the same reason WithoutPackage is one, and after the same
// mistake: the answer is consumed by the parser's allowed list, by the report,
// and by the published page, and three copies of `arm == "intent"` is how the
// next arm gets one of them wrong (prov-2026-d773a705).
func (c Case) HeldToPaths() bool { return heldToPaths(c.Arm) }

// heldToPaths is the arm-level form, for the call that has the arm and not the
// case it came from.
func heldToPaths(arm string) bool { return arm != "intent" }

// Model is one row of the model axis: what to ask for, and what is actually
// answering.
//
// Engine and Quant are carried because the same weights under two runtimes are
// not the same model for measurement purposes. MLX 4-bit and a llama.cpp
// Q4_K_M build of one checkpoint use different quantization schemes and do not
// produce identical output, so folding them into a single row would let a rate
// measured on one be read as a claim about the other.
//
// It also matters practically: only the llama.cpp engine supports continuous
// batching today, so the row that runs fast and the row that runs on the
// interactive default are different rows.
type Model struct {
	// ID is what the endpoint is asked for, verbatim.
	ID string `yaml:"id"`
	// Engine is the runtime serving it - "mlx", "llama.cpp", or a hosted
	// provider's name. Free-form: nothing dispatches on it, it is recorded so
	// that a result says what produced it.
	Engine string `yaml:"engine"`
	// Quant is the quantization, as the build names it: "4bit", "q4_k_m", or
	// empty for a hosted model whose weights are not ours to describe.
	Quant string `yaml:"quant"`
}

// Label identifies a model in a report, and is what distinguishes two rows over
// one checkpoint.
func (m Model) Label() string {
	out := m.ID
	if m.Engine != "" {
		out += "/" + m.Engine
	}
	if m.Quant != "" {
		out += "-" + m.Quant
	}
	return out
}

// Expectations is what a correct run would have produced.
//
// These are the target, not a threshold. Nothing here fails a build: the output
// is the observed rate beside the expectation, because a sampled rate is a
// measurement and a measurement does not pass or fail.
type Expectations struct {
	// Actions maps an action name to how many invocations a complete answer
	// contains. An action listed with count 0 is one that should not appear.
	Actions map[string]int `yaml:"actions"`

	// Bindings is what a complete answer binds those invocations to, keyed by
	// action. Optional: a case that names none measures exactly what it
	// measured before this field existed.
	//
	// Counts and bindings are authored under opposite disciplines, and which
	// applies turns on whether the correct answer exists independently of the
	// model. A count is a property of the package and the record together that
	// nobody has computed, so it is established from a run that produced a
	// complete answer. A binding is stated in the record and settled by the
	// target framework, so establishing it from a run would be asking the thing
	// under test to grade itself (prov-2026-2b121b62).
	Bindings map[string]ActionBinding `yaml:"bindings,omitempty"`

	// Behavior is the fraction of the application's linespec contracts a
	// complete answer would satisfy. Reserved: measuring it means applying the
	// selection, starting the target and running its suite, which is a
	// different order of cost from a model call and is not implemented.
	Behavior *BehaviorExpectation `yaml:"behavior,omitempty"`
}

// ActionBinding is what one action's invocations should carry.
type ActionBinding struct {
	// Because is the constraint this expectation was authored from, quoted or
	// paraphrased from the record the case runs.
	//
	// Required, and loading refuses a binding block without one. Nothing
	// mechanical can check that an expectation is right - type comparison can
	// say the string "false" is not the boolean false, and nothing can say nil
	// was the correct default except the record saying so. What protects the
	// judgment is that it is legible where someone can disagree with it, and a
	// reader deciding whether an expectation is correct should not have to find
	// the provenance record to do it (prov-2026-2b121b62).
	Because string `yaml:"because"`

	// Key names the kwargs that identify an invocation, and pairing an answer
	// against Invocations is by their values.
	//
	// Declared rather than inferred because the alternatives are both wrong.
	// Position carries no meaning - the model emits invocations in whatever
	// order it produced them. Best-match pairing would silently pair an
	// invocation whose identifying kwarg is wrong with the expectation it least
	// resembles, reporting a near-miss where the model addressed the wrong
	// thing entirely.
	Key []string `yaml:"key"`

	// Invocations is one entry per invocation a complete answer contains, each
	// holding the kwargs being expected of it.
	//
	// Only the kwargs named here are scored. An unnamed one reports as not
	// scored rather than as passing, so a fixture can adopt this one kwarg at a
	// time instead of having to be fully specified before any of it can be
	// measured.
	Invocations []map[string]any `yaml:"invocations"`
}

// BehaviorExpectation names the harness that says whether a selection produces
// a working application, and what fraction of samples should.
//
// It reserved `specs`, naming a linespec suite, before there was anything to
// point one at. What landed instead is a target: a directory under
// evals/behavior/targets holding the scaffold, build, boot and assertion steps
// for one stack. The reservation was early rather than wrong - a suite needs
// the application running before it can assert anything, and standing one up
// was the larger half. A linespec assertion phase would go behind `target`
// rather than beside it (prov-2026-83340ba0).
type BehaviorExpectation struct {
	// Target is a directory name under evals/behavior/targets, not a path. What
	// it asserts is its own business, which is what keeps adding a stack from
	// being a change to this schema.
	Target string `yaml:"target"`

	// PassRate is the fraction of measured samples expected to produce a
	// working application, between 0 and 1.
	PassRate float64 `yaml:"pass_rate"`
}

// Load reads every case file in dir.
//
// Paths inside a case are resolved against root rather than against the case
// file, so that one checkout can point at fixture applications wherever they
// live without every case file carrying the same ../.. prefix.
func Load(dir, root string) ([]Case, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading case directory: %w", err)
	}

	var cases []Case
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".yaml" {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}

		var c Case
		dec := yaml.NewDecoder(strings.NewReader(string(raw)))
		dec.KnownFields(true)
		if err := dec.Decode(&c); err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}

		c.Generators = resolve(root, c.Generators)
		c.Records = resolve(root, c.Records)
		// Behavior.Target is deliberately not resolved. It is a name the
		// harness looks up under its own targets directory, and turning it into
		// a path here would put the harness's layout in the case schema.

		if err := c.validate(e.Name()); err != nil {
			return nil, err
		}
		cases = append(cases, c)
	}
	return cases, nil
}

func resolve(root, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

// validate rejects a case that cannot be run, naming the file. A case pointing
// at a fixture that is not checked out is the common one, and it should say so
// rather than failing later inside a package load.
func (c Case) validate(file string) error {
	if c.ID == "" {
		return fmt.Errorf("%s: a case needs an id", file)
	}
	if len(c.Models) == 0 {
		return fmt.Errorf("%s: case %s names no models", file, c.ID)
	}
	for i, m := range c.Models {
		if m.ID == "" {
			return fmt.Errorf("%s: case %s model %d has no id", file, c.ID, i)
		}
		// Engine is required because its absence is what silently merges two
		// rows over one checkpoint, which is the mistake this field exists to
		// prevent. Quant is optional: a hosted model's weights are not ours to
		// describe.
		if m.Engine == "" {
			return fmt.Errorf("%s: case %s model %q declares no engine; two runtimes over one checkpoint are two rows, not one",
				file, c.ID, m.ID)
		}
	}

	for action, b := range c.Expect.Bindings {
		if err := b.validate(file, c.ID, action); err != nil {
			return err
		}
	}

	switch c.Arm {
	case "sedum":
		for label, dir := range map[string]string{"generators": c.Generators, "records": c.Records} {
			if dir == "" {
				return fmt.Errorf("%s: case %s declares no %s", file, c.ID, label)
			}
			if _, err := os.Stat(dir); err != nil {
				return fmt.Errorf("%s: case %s points at a %s directory that is not there: %s",
					file, c.ID, label, dir)
			}
		}
	case "baseline", "intent":
		// Neither arm has a generator package by construction - that
		// absence is the arm - so the directory checks above do not apply. What
		// it does need is somewhere to boot what the model writes, because
		// behaviour is the only rung that can score it (prov-2026-a4dbe65c).
		if c.Expect.Behavior == nil {
			return fmt.Errorf("%s: case %s has arm baseline and declares no expect.behavior target; only behaviour can score a baseline",
				file, c.ID)
		}
		if c.Records == "" {
			return fmt.Errorf("%s: case %s has arm baseline and names no records; the record is the prompt", file, c.ID)
		}
		if _, err := os.Stat(c.Records); err != nil {
			return fmt.Errorf("%s: case %s points at a records directory that is not there: %s", file, c.ID, c.Records)
		}
	default:
		return fmt.Errorf("%s: case %s has arm %q, want \"sedum\" or \"baseline\"", file, c.ID, c.Arm)
	}
	return nil
}

// validate refuses a binding block that cannot be scored or cannot be argued
// with.
func (b ActionBinding) validate(file, caseID, action string) error {
	if strings.TrimSpace(b.Because) == "" {
		return fmt.Errorf("%s: case %s binding for %s states no `because`; a binding expectation is a judgment about what is correct, and one nobody can check against the record is worse than none",
			file, caseID, action)
	}
	if len(b.Key) == 0 {
		return fmt.Errorf("%s: case %s binding for %s declares no key; an answer's invocations arrive in no particular order and something has to say which expectation each one answers",
			file, caseID, action)
	}
	if len(b.Invocations) == 0 {
		return fmt.Errorf("%s: case %s binding for %s expects no invocations", file, caseID, action)
	}

	seen := map[string]bool{}
	for i, inv := range b.Invocations {
		for _, k := range b.Key {
			v, ok := inv[k]
			if !ok {
				return fmt.Errorf("%s: case %s binding for %s invocation %d carries no %q, which is one of its key kwargs",
					file, caseID, action, i+1, k)
			}
			// A kwarg may expect any of several literals, but not one that
			// identifies the invocation: pairing looks a value up, and a key
			// that could be two things cannot be looked up as either.
			if _, isList := v.([]any); isList {
				return fmt.Errorf("%s: case %s binding for %s invocation %d expects several values for %q, which is one of its key kwargs; a key is looked up rather than compared",
					file, caseID, action, i+1, k)
			}
		}
		// Two expected invocations sharing a key would each match the same
		// answer, so one of them could never be missed however wrong the
		// answer was.
		id := keyOf(inv, b.Key)
		if seen[id] {
			return fmt.Errorf("%s: case %s binding for %s has two invocations keyed %s; a key that does not distinguish them cannot pair either one",
				file, caseID, action, id)
		}
		seen[id] = true
	}
	return nil
}

// keyOf is an invocation's identity under a key, as a comparable string.
//
// This is the one place values are rendered to text, and it is safe here
// because it decides *which* expectation an invocation answers rather than
// whether it answers it correctly. Scoring never stringifies (prov-2026-2b121b62).
func keyOf(kwargs map[string]any, key []string) string {
	parts := make([]string, 0, len(key))
	for _, k := range key {
		parts = append(parts, fmt.Sprintf("%s=%v", k, kwargs[k]))
	}
	return strings.Join(parts, " ")
}
