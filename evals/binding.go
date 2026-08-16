package evals

import (
	"encoding/json"
	"sort"

	"github.com/calebcowen/sedum/internal/recording"
)

// Binding scoring: what the model bound, against what the record says is
// correct.
//
// Selection asks whether the model reached for the right actions. This asks
// whether it reached for them correctly, and the two are different measurements
// that a combined rate would hide. The case for keeping them apart is not
// theoretical: after invocations began being stored, four smoke samples scored
// a perfect six-of-six on selection while every one of them bound `default`
// wrong on at least one column - one of them to an empty string that renders
// `default: ` and does not parse (prov-2026-2b121b62).

// KwargResult is one kwarg's binding rate across a measurement's samples.
//
// Two denominators, and only one of them may carry an interval. A sample is one
// prompt and one act of reasoning; the rails fixture asks for two columns, so a
// model that has decided on `"false"` binds it twice and those two comparisons
// are one observation seen twice. An interval over invocations would not move
// the rate - it would shrink the uncertainty around it, reporting more
// confidence than was bought (prov-2026-7cb96bf0).
type KwargResult struct {
	Action string
	Kwarg  string

	// CleanSamples is how many samples bound this kwarg correctly in every
	// invocation of it, out of Samples that bound it at all. This is the
	// citable rate and the only one an interval goes on.
	CleanSamples int
	Samples      int

	// Correct and Scored are the same comparison counted per invocation. Kept
	// because a sample rate cannot say whether a failing sample got one
	// invocation wrong or all of them - a slip and a decision applied
	// consistently are different findings - and reported as a bare ratio
	// because an interval here is the error this record exists to remove.
	Correct int
	Scored  int
}

// Rate is the fraction of samples that bound this kwarg correctly throughout.
func (k KwargResult) Rate() float64 {
	if k.Samples == 0 {
		return 0
	}
	return float64(k.CleanSamples) / float64(k.Samples)
}

// BindingResult is one action's binding outcome across a measurement.
type BindingResult struct {
	Action string
	// Kwargs are the per-kwarg rates, sorted by name.
	Kwargs []KwargResult
	// Missing is how many expected invocations no answer accounted for, and
	// Unexpected is how many answered invocations no expectation named. Both
	// are counted across scored samples rather than per sample.
	//
	// An invocation whose key kwarg is bound wrong lands here as one of each,
	// which is the correct reading: the model did not fumble an argument, it
	// wrote about something nobody asked for.
	Missing    int
	Unexpected int
}

// ScoredBindings scores every action a case declares bindings for.
//
// Only valid samples are scored, for the same reason only valid samples are
// counted: a rejected answer has no invocation list, and scoring it would
// report a rejection as a wrong argument.
//
// Returns nil when the case declares no bindings, which is what makes the field
// optional - a case that names none measures exactly what it measured before
// bindings existed.
func (m Measurement) ScoredBindings() []BindingResult {
	if len(m.Case.Expect.Bindings) == 0 {
		return nil
	}

	actions := make([]string, 0, len(m.Case.Expect.Bindings))
	for action := range m.Case.Expect.Bindings {
		actions = append(actions, action)
	}
	sort.Strings(actions)

	var out []BindingResult
	for _, action := range actions {
		want := m.Case.Expect.Bindings[action]
		res := BindingResult{Action: action}
		byKwarg := map[string]*KwargResult{}

		for _, s := range m.Samples {
			if s.Invalid || s.Err != nil {
				continue
			}
			scoreSample(action, want, s.Invocations, &res, byKwarg)
		}

		for _, k := range byKwarg {
			res.Kwargs = append(res.Kwargs, *k)
		}
		sort.Slice(res.Kwargs, func(i, j int) bool { return res.Kwargs[i].Kwarg < res.Kwargs[j].Kwarg })
		out = append(out, res)
	}
	return out
}

// sampleTally is one kwarg's comparisons within a single sample, which is the
// grain the clustering lives at.
type sampleTally struct{ correct, scored int }

// scoreSample pairs one sample's invocations of an action against the
// expectation and tallies the comparison.
//
// The per-invocation counts accumulate straight into byKwarg. The per-sample
// verdict is decided here and folded in once at the end, because "bound
// correctly in this sample" means every invocation of it in this sample was
// correct.
func scoreSample(action string, want ActionBinding, invocations []recording.Invocation, res *BindingResult, byKwarg map[string]*KwargResult) {
	// Pairing is by the declared key, never by position and never by
	// closest-match. An invocation the key does not match is unexpected rather
	// than a near-miss.
	actual := map[string]map[string]any{}
	for _, inv := range invocations {
		if inv.Action != action {
			continue
		}
		actual[keyOf(inv.Kwargs, want.Key)] = inv.Kwargs
	}

	matched := map[string]bool{}
	within := map[string]*sampleTally{}

	for _, expected := range want.Invocations {
		id := keyOf(expected, want.Key)
		got, ok := actual[id]
		if !ok {
			res.Missing++
			continue
		}
		matched[id] = true

		for kwarg, wantValue := range expected {
			k := byKwarg[kwarg]
			if k == nil {
				k = &KwargResult{Action: action, Kwarg: kwarg}
				byKwarg[kwarg] = k
			}
			t := within[kwarg]
			if t == nil {
				t = &sampleTally{}
				within[kwarg] = t
			}

			k.Scored++
			t.scored++
			if equalBinding(wantValue, got[kwarg]) {
				k.Correct++
				t.correct++
			}
		}
	}
	for id := range actual {
		if !matched[id] {
			res.Unexpected++
		}
	}

	// One observation per kwarg per sample. A kwarg the sample never bound -
	// because none of its invocations paired - contributes no observation
	// rather than a failed one, which is the same rule that keeps a rejected
	// sample out of the selection denominator.
	for kwarg, t := range within {
		k := byKwarg[kwarg]
		k.Samples++
		if t.correct == t.scored {
			k.CleanSamples++
		}
	}
}

// equalBinding compares an expected value against a bound one without
// stringifying either.
//
// Rendering both to text before comparing would make the string "false" equal
// the boolean false, which is exactly the distinction that motivated scoring
// bindings at all: the described package set bound the string "false" to a
// boolean column that wanted the literal false, and to a string column that
// wanted nil. A scoring rule that erases the finding it was built for is not a
// scoring rule.
//
// Numbers are the one normalisation. A kwarg arrives through encoding/json,
// where every number decodes to float64, while the expectation arrives through
// YAML, where 5 decodes to int - so an int kwarg would never match its own
// expectation on representation alone. That is a difference between two
// decoders rather than between two values.
func equalBinding(want, got any) bool {
	if wantNum, ok := numeric(want); ok {
		gotNum, ok := numeric(got)
		return ok && wantNum == gotNum
	}
	return want == got
}

func numeric(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}
