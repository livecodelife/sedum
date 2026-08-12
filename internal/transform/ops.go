package transform

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

// The built-in operations.
//
// Case conversion normalizes its input before doing anything else, so that
// UserURL, user-id, and users_controller all arrive at the same answer. That
// word splitting is the fiddly part of the job and is what strcase is here for.

// operate applies one built-in operation.
func (e *Engine) operate(r Ref, v any) (string, error) {
	// Inflection is the only place a value's own forms matter, so it reads
	// the value before it is flattened to text.
	switch r.Name {
	case "plural":
		if val, ok := v.(Value); ok && val.Plural != "" {
			return val.Plural, nil
		}
	case "singular":
		if val, ok := v.(Value); ok && val.Singular != "" {
			return val.Singular, nil
		}
	}

	s, err := scalar(v)
	if err != nil {
		return "", err
	}

	switch r.Name {
	case "pascal":
		return e.joinWords(s, "pascal", false), nil
	case "camel":
		return e.joinWords(s, "camel", true), nil
	case "snake":
		return strcase.ToSnake(s), nil
	case "kebab":
		return strcase.ToKebab(s), nil
	case "upper":
		// upper and lower are whole-token case folds, not word-splitting
		// operations: a screaming constant is the composition
		// [snake, upper] rather than an eleventh built-in.
		return strings.ToUpper(s), nil
	case "lower":
		return strings.ToLower(s), nil
	case "plural":
		return e.inflector.plural(s), nil
	case "singular":
		return e.inflector.singular(s), nil
	case "prefix":
		return r.Arg + s, nil
	case "suffix":
		return s + r.Arg, nil
	}
	return "", checkOperation(r)
}

// joinWords is pascal and camel. It normalizes to snake form, splits on the
// underscore, renders each word, and joins.
//
// The exception table is consulted per word rather than over the whole token,
// so that declaring url -> URL once reaches every token containing that word:
// pascal over user_url yields UserURL, not UserUrl.
//
// camel lowercases the leading word and skips the table for it, since a leading
// acronym would otherwise render uRL.
func (e *Engine) joinWords(s, op string, lowerFirst bool) string {
	table := e.exceptions[op]

	var out strings.Builder
	for i, word := range strings.Split(strcase.ToSnake(s), "_") {
		if word == "" {
			continue
		}
		if i == 0 && lowerFirst {
			out.WriteString(strings.ToLower(word))
			continue
		}
		if replacement, ok := table[strings.ToLower(word)]; ok {
			out.WriteString(replacement)
			continue
		}
		out.WriteString(capitalize(word))
	}
	return out.String()
}

// capitalize uppercases the first rune and leaves the rest, which is already
// normalized by the time it gets here.
func capitalize(word string) string {
	runes := []rune(word)
	if len(runes) == 0 {
		return word
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
