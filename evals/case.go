package evals

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	Language  string `yaml:"language"`

	// Tightness is how much the generator package constrains the model:
	// "defined" for a package whose actions and variants cover the intended
	// output closely, "loose" for one that leaves more to a fallback.
	//
	// This is the axis Generators swaps along. Two cases naming one record and
	// two package sets differ in exactly one thing, which is what makes the
	// comparison mean anything - and it is why generators/ holds package sets
	// rather than packages.
	Tightness string `yaml:"tightness"`

	// Arm is "sedum" or "baseline". A baseline arm asks the same model for the
	// same application with no generator package and no action vocabulary, and
	// exists because a selection rate has no meaning without one.
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

	Expect Expectations `yaml:"expect"`
}

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

	// Behavior is the fraction of the application's linespec contracts a
	// complete answer would satisfy. Reserved: measuring it means applying the
	// selection, starting the target and running its suite, which is a
	// different order of cost from a model call and is not implemented.
	Behavior *BehaviorExpectation `yaml:"behavior,omitempty"`
}

// BehaviorExpectation reserves the shape of the expensive half.
type BehaviorExpectation struct {
	// Specs is the directory holding the application's linespec contracts.
	Specs string `yaml:"specs"`
	// PassRate is the fraction expected to pass, between 0 and 1.
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
		if c.Expect.Behavior != nil {
			c.Expect.Behavior.Specs = resolve(root, c.Expect.Behavior.Specs)
		}

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
	case "baseline":
		// A baseline arm has no generator package by construction. It is
		// declared here so the matrix can hold it; running one is not
		// implemented.
	default:
		return fmt.Errorf("%s: case %s has arm %q, want \"sedum\" or \"baseline\"", file, c.ID, c.Arm)
	}
	return nil
}
