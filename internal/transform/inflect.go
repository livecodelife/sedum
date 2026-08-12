package transform

import (
	"embed"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// Inflection.
//
// plural and singular are the only operations that cannot be expressed as
// pattern rules over the input alone, because morphology is irregular. The
// rules therefore ship as a data table per language rather than as code: adding
// a language means adding a file under inflections/, never a branch here.

//go:embed inflections/*.yaml
var inflectionTables embed.FS

// DefaultLanguage names the table loaded when a caller asks for none. It is the
// name of a data file, not a case in the code.
const DefaultLanguage = "en"

// inflectionTable is one language's rules, as the data file declares them.
type inflectionTable struct {
	// Uncountable words have one form. They are checked first because no
	// rule below should ever see them.
	Uncountable []string `yaml:"uncountable"`
	// Irregular maps singular to plural for words no rule can derive. It is
	// read in both directions, which is also what makes an already-inflected
	// irregular a fixed point.
	Irregular map[string]string `yaml:"irregular"`
	// Plural and Singular are ordered pattern rules, most specific first.
	// The first match wins.
	Plural   []inflectionRule `yaml:"plural"`
	Singular []inflectionRule `yaml:"singular"`
}

type inflectionRule struct {
	Match   string `yaml:"match"`
	Replace string `yaml:"replace"`
}

type compiledRule struct {
	match   *regexp.Regexp
	replace string
}

type inflector struct {
	uncountable   map[string]bool
	toPlural      map[string]string
	toSingular    map[string]string
	pluralRules   []compiledRule
	singularRules []compiledRule
}

// loadInflector reads and compiles a language's rule table. A malformed or
// missing table is a hard failure: an engine that cannot inflect would silently
// produce singular collection names.
func loadInflector(language string) (*inflector, error) {
	if language == "" {
		language = DefaultLanguage
	}

	data, err := inflectionTables.ReadFile("inflections/" + language + ".yaml")
	if err != nil {
		return nil, fmt.Errorf("no inflection table for %q; tables are data files under inflections/", language)
	}

	var table inflectionTable
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	dec.KnownFields(true)
	if err := dec.Decode(&table); err != nil {
		return nil, fmt.Errorf("inflection table %s: %w", language, err)
	}

	inf := &inflector{
		uncountable: map[string]bool{},
		toPlural:    map[string]string{},
		toSingular:  map[string]string{},
	}
	for _, word := range table.Uncountable {
		inf.uncountable[strings.ToLower(word)] = true
	}
	for singular, plural := range table.Irregular {
		inf.toPlural[strings.ToLower(singular)] = plural
		inf.toSingular[strings.ToLower(plural)] = singular
	}

	inf.pluralRules, err = compileRules(language, "plural", table.Plural)
	if err != nil {
		return nil, err
	}
	inf.singularRules, err = compileRules(language, "singular", table.Singular)
	if err != nil {
		return nil, err
	}
	return inf, nil
}

func compileRules(language, direction string, rules []inflectionRule) ([]compiledRule, error) {
	out := make([]compiledRule, 0, len(rules))
	for _, rule := range rules {
		re, err := regexp.Compile(rule.Match)
		if err != nil {
			return nil, fmt.Errorf("inflection table %s: %s rule %q: %w", language, direction, rule.Match, err)
		}
		out = append(out, compiledRule{match: re, replace: rule.Replace})
	}
	return out, nil
}

func (i *inflector) plural(s string) string {
	return i.inflect(s, i.toPlural, i.toSingular, i.pluralRules)
}

func (i *inflector) singular(s string) string {
	return i.inflect(s, i.toSingular, i.toPlural, i.singularRules)
}

// inflect resolves one direction.
//
// known maps a word to its form in the target direction; already holds the
// words that are already in the target direction, so that inflecting twice
// changes nothing. Idempotency is not decoration: a run may inflect a value the
// model already supplied in the form asked for.
func (i *inflector) inflect(s string, known, already map[string]string, rules []compiledRule) string {
	if s == "" {
		return s
	}

	word := strings.ToLower(s)
	if i.uncountable[word] {
		return s
	}
	if _, ok := already[word]; ok {
		return s
	}
	if form, ok := known[word]; ok {
		return matchLeadingCase(s, form)
	}

	for _, rule := range rules {
		if rule.match.MatchString(s) {
			return rule.match.ReplaceAllString(s, rule.replace)
		}
	}
	return s
}

// matchLeadingCase carries the input's leading case onto a table lookup, since
// a table entry is written in one case and the value it is applied to may not
// be. Pattern rules need no equivalent: they rewrite a suffix and leave the
// rest of the word as it was.
func matchLeadingCase(in, out string) string {
	if in == "" || out == "" {
		return out
	}
	first := []rune(in)[0]
	if !unicode.IsUpper(first) {
		return out
	}
	runes := []rune(out)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
