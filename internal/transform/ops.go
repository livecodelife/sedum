package transform

import (
	"errors"
	"fmt"
	"reflect"
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

	// join is the only operation whose input is a list, so it reads the value
	// before the flattening below refuses one. It is also the only way a list
	// reaches a file at all: every other path refuses it.
	if r.Name == "join" {
		return joinList(r.Arg, v)
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

// joinList renders every member of a list and joins them with a literal
// separator.
//
// A member is rendered by the same extraction every other value goes through,
// so a list of numbers is a list of numbers and a nested list is refused where
// it sits rather than formatted.
//
// The order is the order bound. A joined list is data — reordering an
// alternation changes which alternative a regular expression prefers — so
// nothing here sorts.
func joinList(sep string, v any) (string, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		// The mistake is in the package rather than the binding: join was
		// written against a kwarg that is not a list.
		return "", fmt.Errorf("join takes a list, but this value is %s", describeKind(rv))
	}

	// An empty list joins to the empty string, and the empty string looks like
	// a value in every place a joined list is going: an empty alternation
	// matches everything, an empty selector list matches nothing. Nothing
	// downstream reports either, so it is reported here.
	if rv.Len() == 0 {
		return "", errors.New("join received an empty list, and an empty list has nothing to join")
	}

	parts := make([]string, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		part, err := scalar(rv.Index(i).Interface())
		if err != nil {
			return "", fmt.Errorf("join member %d: %w", i, err)
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, sep), nil
}

// describeKind names what arrived, so the diagnostic tells a package author
// what they wrote join against.
func describeKind(rv reflect.Value) string {
	switch rv.Kind() {
	case reflect.Invalid:
		return "empty"
	case reflect.String:
		return "a single value"
	default:
		return "a " + rv.Kind().String()
	}
}
