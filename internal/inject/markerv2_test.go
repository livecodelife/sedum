package inject

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/livecodelife/sedum/internal/genpkg"
)

// Marker v2 (prov-2026-a3771fbf): a region's marker carries what identifies it
// (the id object) and what detects drift (a kwargs hash and a body hash). These
// cases pin the shape of the types and the reading and writing of the marker
// line. They are written against the records, before the code.

// fieldSet describes a struct as a map from field name to type text.
func fieldSet(typ reflect.Type) map[string]string {
	out := map[string]string{}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		out[f.Name] = f.Type.String()
	}
	return out
}

func checkFieldSet(t *testing.T, typ reflect.Type, want map[string]string) {
	t.Helper()
	got := fieldSet(typ)
	for name, wantType := range want {
		gotType, ok := got[name]
		if !ok {
			t.Errorf("%s has no field %s", typ.Name(), name)
			continue
		}
		if gotType != wantType {
			t.Errorf("%s.%s has type %s, want %s", typ.Name(), name, gotType, wantType)
		}
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			t.Errorf("%s has an unexpected field %s", typ.Name(), name)
		}
	}
}

// prov-2026-1dd3221d
func TestMarkerFields(t *testing.T) {
	checkFieldSet(t, reflect.TypeOf(Marker{}), map[string]string{
		"Action":     "string",
		"Variant":    "string",
		"Record":     "string",
		"Gen":        "string",
		"Pkg":        "string",
		"Writer":     "string",
		"Tier":       "inject.Tier",
		"Kwargs":     "map[string]interface {}",
		"Extra":      "map[string]json.RawMessage",
		"Version":    "int",
		"Foreign":    "bool",
		"KwargsHash": "string",
		"Body":       "string",
	})
}

// prov-2026-a2f0e140
func TestRegionFields(t *testing.T) {
	checkFieldSet(t, reflect.TypeOf(Region{}), map[string]string{
		"Marker":   "inject.Marker",
		"Start":    "int",
		"End":      "int",
		"BodyHash": "string",
	})
}

// prov-2026-727ffa6c
func TestResultFields(t *testing.T) {
	checkFieldSet(t, reflect.TypeOf(Result{}), map[string]string{
		"Path":       "string",
		"Action":     "string",
		"Variant":    "string",
		"Record":     "string",
		"Replaced":   "bool",
		"Skipped":    "bool",
		"Pruned":     "bool",
		"Unselected": "bool",
		"Held":       "bool",
		"Orphan":     "bool",
		"Changed":    "bool",
		"Marker":     "inject.Marker",
		"Kwargs":     "map[string]interface {}",
	})
}

// prov-2026-0f92d922
func TestAttrsV2(t *testing.T) {
	typ := reflect.TypeOf(attrsV2{})
	wantFields := []struct {
		name string
		typ  string
		tag  string
	}{
		{"V", "int", "v"},
		{"Tier", "inject.Tier", "tier"},
		{"Record", "string", "record,omitempty"},
		{"Gen", "string", "gen,omitempty"},
		{"Pkg", "string", "pkg,omitempty"},
		{"Writer", "string", "writer,omitempty"},
		{"ID", "map[string]interface {}", "id,omitempty"},
		{"Kwargs", "string", "kwargs,omitempty"},
		{"Body", "string", "body,omitempty"},
	}
	if typ.NumField() != len(wantFields) {
		t.Fatalf("attrsV2 has %d fields, want %d", typ.NumField(), len(wantFields))
	}
	for i, want := range wantFields {
		f := typ.Field(i)
		if f.Name != want.name || f.Type.String() != want.typ {
			t.Errorf("field %d is %s %s, want %s %s", i, f.Name, f.Type, want.name, want.typ)
		}
		if tag := f.Tag.Get("json"); tag != want.tag {
			t.Errorf("field %s has json tag %q, want %q", f.Name, tag, want.tag)
		}
	}

	tests := []struct {
		name string
		in   attrsV2
		want string
	}{
		{"v and tier are always written", attrsV2{V: 2, Tier: TierOwned}, `{"v":2,"tier":"owned"}`},
		{"empty tier is still written", attrsV2{V: 2}, `{"v":2,"tier":""}`},
		{"every key, in field order",
			attrsV2{V: 2, Tier: TierSeeded, Record: "R", Gen: "g", Pkg: "p", Writer: "w",
				ID: map[string]any{"b": 1, "a": "x"}, Kwargs: "sha256:k", Body: "sha256:b"},
			`{"v":2,"tier":"seeded","record":"R","gen":"g","pkg":"p","writer":"w","id":{"a":"x","b":1},"kwargs":"sha256:k","body":"sha256:b"}`},
		{"empty id is omitted", attrsV2{V: 2, Tier: TierOwned, ID: map[string]any{}}, `{"v":2,"tier":"owned"}`},
		{"empty strings are omitted", attrsV2{V: 2, Tier: TierOwned, Record: "R"}, `{"v":2,"tier":"owned","record":"R"}`},
		{"body alone", attrsV2{V: 2, Tier: TierOwned, Body: "sha256:b"}, `{"v":2,"tier":"owned","body":"sha256:b"}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.in)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("attrsV2 encodes as\n  %s\nwant\n  %s", got, tc.want)
			}
		})
	}
}

// prov-2026-a73496ad
func TestEncodeAttrs(t *testing.T) {
	v2 := func(mutate func(*Marker)) Marker {
		m := Marker{
			Action: "addPolicy", Version: 2, Tier: TierOwned,
			Record: "PR-1", Gen: "sha256:g", Pkg: "sha256:p", Writer: DefaultWriter,
			Kwargs:     map[string]any{"name": "read", "controller": "users"},
			KwargsHash: "sha256:k", Body: "sha256:b",
		}
		if mutate != nil {
			mutate(&m)
		}
		return m
	}
	raw := func(s string) json.RawMessage { return json.RawMessage(s) }

	tests := []struct {
		name string
		in   Marker
		want string
	}{
		{"v2 key order is v, tier, record, gen, pkg, writer, id, kwargs, body", v2(nil),
			`{"v":2,"tier":"owned","record":"PR-1","gen":"sha256:g","pkg":"sha256:p","id":{"controller":"users","name":"read"},"kwargs":"sha256:k","body":"sha256:b"}`},
		{"a named writer sits between pkg and id",
			v2(func(m *Marker) { m.Writer = "harness" }),
			`{"v":2,"tier":"owned","record":"PR-1","gen":"sha256:g","pkg":"sha256:p","writer":"harness","id":{"controller":"users","name":"read"},"kwargs":"sha256:k","body":"sha256:b"}`},
		{"an empty tier is written as owned",
			v2(func(m *Marker) { m.Tier = "" }),
			`{"v":2,"tier":"owned","record":"PR-1","gen":"sha256:g","pkg":"sha256:p","id":{"controller":"users","name":"read"},"kwargs":"sha256:k","body":"sha256:b"}`},
		{"the seeded tier is written",
			v2(func(m *Marker) { m.Tier = TierSeeded }),
			`{"v":2,"tier":"seeded","record":"PR-1","gen":"sha256:g","pkg":"sha256:p","id":{"controller":"users","name":"read"},"kwargs":"sha256:k","body":"sha256:b"}`},
		{"empty fields are omitted",
			Marker{Action: "a", Version: 2},
			`{"v":2,"tier":"owned"}`},
		{"HTML characters in id are written as themselves",
			Marker{Action: "a", Version: 2, Kwargs: map[string]any{"t": "&t.ID <T> >"}},
			`{"v":2,"tier":"owned","id":{"t":"&t.ID <T> >"}}`},
		{"extra keys follow the declared keys, sorted",
			v2(func(m *Marker) {
				m.Extra = map[string]json.RawMessage{"zeta": raw(`1`), "alpha": raw(`"x"`), "mid": raw(`{"a":[1,2]}`)}
			}),
			`{"v":2,"tier":"owned","record":"PR-1","gen":"sha256:g","pkg":"sha256:p","id":{"controller":"users","name":"read"},"kwargs":"sha256:k","body":"sha256:b","alpha":"x","mid":{"a":[1,2]},"zeta":1}`},
		{"extra values are compacted",
			Marker{Action: "a", Version: 2, Extra: map[string]json.RawMessage{"n": raw(` { "a" : [ 1 , 2 ] } `)}},
			`{"v":2,"tier":"owned","n":{"a":[1,2]}}`},
		{"extra key names are encoded as strings",
			Marker{Action: "a", Version: 2, Extra: map[string]json.RawMessage{`a"b`: raw(`true`)}},
			`{"v":2,"tier":"owned","a\"b":true}`},
		{"on v2 an extra key named one of the nine v2 keys is not appended",
			Marker{Action: "a", Version: 2, Extra: map[string]json.RawMessage{
				"v": raw(`3`), "tier": raw(`"seeded"`), "record": raw(`"X"`), "gen": raw(`"X"`), "pkg": raw(`"X"`),
				"writer": raw(`"X"`), "id": raw(`{"x":1}`), "kwargs": raw(`{"x":1}`), "body": raw(`"X"`), "keep": raw(`1`)}},
			`{"v":2,"tier":"owned","keep":1}`},
		{"v1 marker has the kwargs object and no v",
			Marker{Action: "a", Tier: TierOwned, Record: "R", Kwargs: map[string]any{"b": 1, "a": "x"}},
			`{"tier":"owned","record":"R","kwargs":{"a":"x","b":1}}`},
		{"v1 marker never writes the v2 hashes",
			Marker{Action: "a", Tier: TierOwned, KwargsHash: "sha256:k", Body: "sha256:b"},
			`{"tier":"owned"}`},
		{"on v1 an extra key named id or body is appended",
			Marker{Action: "a", Tier: TierOwned, Extra: map[string]json.RawMessage{"id": raw(`{"x":1}`), "body": raw(`"b"`)}},
			`{"tier":"owned","body":"b","id":{"x":1}}`},
		{"on v1 an extra key named one of the six v1 keys is not appended",
			Marker{Action: "a", Tier: TierOwned, Extra: map[string]json.RawMessage{
				"tier": raw(`"seeded"`), "record": raw(`"X"`), "gen": raw(`"X"`), "pkg": raw(`"X"`),
				"writer": raw(`"X"`), "kwargs": raw(`{"x":1}`), "keep": raw(`1`)}},
			`{"tier":"owned","keep":1}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.in.encodeAttrs()
			if err != nil {
				t.Fatalf("encodeAttrs: %v", err)
			}
			if got != tc.want {
				t.Errorf("encodeAttrs =\n  %s\nwant\n  %s", got, tc.want)
			}
			if strings.Contains(got, "\n") {
				t.Errorf("encodeAttrs wrote a newline: %q", got)
			}
		})
	}
}

func TestEncodeAttrsErrors(t *testing.T) {
	tests := []struct {
		name string
		in   Marker
		want string
	}{
		{"v2 with an id that cannot be encoded",
			Marker{Action: "addStep", Variant: "run", Version: 2, Kwargs: map[string]any{"f": func() {}}},
			"action addStep:run: kwargs cannot be recorded on a marker: "},
		{"v1 with kwargs that cannot be encoded",
			Marker{Action: "addStep", Kwargs: map[string]any{"f": func() {}}},
			"action addStep: kwargs cannot be recorded on a marker: "},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.in.encodeAttrs()
			if err == nil {
				t.Fatalf("encodeAttrs = %q, want an error", got)
			}
			if !strings.HasPrefix(err.Error(), tc.want) {
				t.Errorf("error = %q, want it to start with %q", err.Error(), tc.want)
			}
			if !strings.Contains(err.Error(), "unsupported type") {
				t.Errorf("error does not carry the cause: %q", err.Error())
			}
		})
	}
}

// The opening line is the comment-wrapped label and the attribute object.
func TestOpenWritesAV2Marker(t *testing.T) {
	m := Marker{
		Action: "addPolicy", Variant: "", Version: 2, Tier: TierOwned, Record: "PR-1",
		Kwargs: map[string]any{"name": "read"}, KwargsHash: "sha256:k", Body: "sha256:b",
	}
	got, err := m.Open(lineComment("#"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	want := `# sedum:addPolicy {"v":2,"tier":"owned","record":"PR-1","id":{"name":"read"},"kwargs":"sha256:k","body":"sha256:b"}`
	if got != want {
		t.Errorf("Open =\n  %s\nwant\n  %s", got, want)
	}

	delimited, err := m.Open(genpkg.Comment{Prefix: "<!--", Suffix: "-->"})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	want = `<!-- sedum:addPolicy {"v":2,"tier":"owned","record":"PR-1","id":{"name":"read"},"kwargs":"sha256:k","body":"sha256:b"} -->`
	if delimited != want {
		t.Errorf("Open =\n  %s\nwant\n  %s", delimited, want)
	}
}

// normalized drops the difference between a nil map and an empty one, which the
// records do not distinguish ("empty").
func normalized(m Marker) Marker {
	if len(m.Kwargs) == 0 {
		m.Kwargs = nil
	}
	if len(m.Extra) == 0 {
		m.Extra = nil
	}
	return m
}

func TestParseOpenNotAMarker(t *testing.T) {
	for _, line := range []string{
		"",
		"x := 1",
		"# a plain comment",
		"# sedum:anchor:class_body",
		"# /sedum:addStep",
		"sedum:addStep {}",
	} {
		got, ok, err := parseOpen(lineComment("#"), line)
		if ok || err != nil {
			t.Errorf("parseOpen(%q) = %v, %v; want false and no error", line, ok, err)
		}
		if !reflect.DeepEqual(got, Marker{}) {
			t.Errorf("parseOpen(%q) returned %+v, want the zero Marker", line, got)
		}
	}
}

func TestParseOpenWithNoAttributes(t *testing.T) {
	got, ok, err := parseOpen(lineComment("#"), "# sedum:addStep:run")
	if err != nil || !ok {
		t.Fatalf("parseOpen = %v, %v", ok, err)
	}
	want := Marker{Action: "addStep", Variant: "run", Tier: TierOwned, Writer: "sedum"}
	if !reflect.DeepEqual(normalized(got), want) {
		t.Errorf("parseOpen = %+v, want %+v", got, want)
	}
	if got.Version != 0 || got.Foreign {
		t.Errorf("a marker with no attributes has Version %d, Foreign %v; want 0, false", got.Version, got.Foreign)
	}
}

func TestParseOpenV2(t *testing.T) {
	const comment = "#"
	tests := []struct {
		name string
		line string
		want Marker
	}{
		{"baseline",
			`# sedum:addPolicy:x {"v":2,"tier":"owned","record":"PR-1","gen":"sha256:g","pkg":"sha256:p","id":{"controller":"users","name":"read"},"kwargs":"sha256:k","body":"sha256:b"}`,
			Marker{Action: "addPolicy", Variant: "x", Version: 2, Tier: TierOwned, Writer: "sedum",
				Record: "PR-1", Gen: "sha256:g", Pkg: "sha256:p",
				Kwargs:     map[string]any{"controller": "users", "name": "read"},
				KwargsHash: "sha256:k", Body: "sha256:b"}},
		{"seeded tier and a named writer",
			`# sedum:addPolicy {"v":2,"tier":"seeded","writer":"harness","id":{"n":1}}`,
			Marker{Action: "addPolicy", Version: 2, Tier: TierSeeded, Writer: "harness",
				Kwargs: map[string]any{"n": float64(1)}}},
		{"no id, kwargs or body",
			`# sedum:addPolicy {"v":2,"tier":"owned","record":"R"}`,
			Marker{Action: "addPolicy", Version: 2, Tier: TierOwned, Writer: "sedum", Record: "R"}},
		{"only v",
			`# sedum:addPolicy {"v":2}`,
			Marker{Action: "addPolicy", Version: 2, Tier: TierOwned, Writer: "sedum"}},
		{"unknown keys are kept in extra, the nine are not",
			`# sedum:addPolicy {"v":2,"tier":"owned","id":{"a":1},"kwargs":"sha256:k","body":"sha256:b","verified_by":"spec/x","attempts":3}`,
			Marker{Action: "addPolicy", Version: 2, Tier: TierOwned, Writer: "sedum",
				Kwargs: map[string]any{"a": float64(1)}, KwargsHash: "sha256:k", Body: "sha256:b",
				Extra: map[string]json.RawMessage{"verified_by": json.RawMessage(`"spec/x"`), "attempts": json.RawMessage(`3`)}}},
		{"v written after the other keys",
			`# sedum:addPolicy {"tier":"owned","kwargs":"sha256:k","v":2}`,
			Marker{Action: "addPolicy", Version: 2, Tier: TierOwned, Writer: "sedum", KwargsHash: "sha256:k"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok, err := parseOpen(lineComment(comment), tc.line)
			if err != nil || !ok {
				t.Fatalf("parseOpen = %v, %v", ok, err)
			}
			if !reflect.DeepEqual(normalized(got), normalized(tc.want)) {
				t.Errorf("parseOpen =\n  %+v\nwant\n  %+v", got, tc.want)
			}
		})
	}
}

func TestParseOpenV2WithACommentSuffix(t *testing.T) {
	got, ok, err := parseOpen(genpkg.Comment{Prefix: "<!--", Suffix: "-->"},
		`  <!-- sedum:addCol:a {"v":2,"tier":"owned","id":{"c":"a"},"kwargs":"sha256:k","body":"sha256:b"} -->`)
	if err != nil || !ok {
		t.Fatalf("parseOpen = %v, %v", ok, err)
	}
	want := Marker{Action: "addCol", Variant: "a", Version: 2, Tier: TierOwned, Writer: "sedum",
		Kwargs: map[string]any{"c": "a"}, KwargsHash: "sha256:k", Body: "sha256:b"}
	if !reflect.DeepEqual(normalized(got), want) {
		t.Errorf("parseOpen = %+v, want %+v", got, want)
	}
}

func TestParseOpenV1(t *testing.T) {
	tests := []struct {
		name string
		line string
		want Marker
	}{
		{"kwargs object, no version, no body",
			`# sedum:addStep:run {"tier":"owned","record":"R","gen":"g","pkg":"p","kwargs":{"a":"x"}}`,
			Marker{Action: "addStep", Variant: "run", Tier: TierOwned, Writer: "sedum", Record: "R", Gen: "g", Pkg: "p",
				Kwargs: map[string]any{"a": "x"}}},
		{"id and body keys are kept in extra",
			`# sedum:addStep {"tier":"owned","kwargs":{"a":"x"},"id":{"a":1},"body":"sha256:b"}`,
			Marker{Action: "addStep", Tier: TierOwned, Writer: "sedum", Kwargs: map[string]any{"a": "x"},
				Extra: map[string]json.RawMessage{"id": json.RawMessage(`{"a":1}`), "body": json.RawMessage(`"sha256:b"`)}}},
		{"unknown key is kept in extra",
			`# sedum:addStep {"tier":"seeded","writer":"w","zz":[1,2]}`,
			Marker{Action: "addStep", Tier: TierSeeded, Writer: "w",
				Extra: map[string]json.RawMessage{"zz": json.RawMessage(`[1,2]`)}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok, err := parseOpen(lineComment("#"), tc.line)
			if err != nil || !ok {
				t.Fatalf("parseOpen = %v, %v", ok, err)
			}
			if !reflect.DeepEqual(normalized(got), normalized(tc.want)) {
				t.Errorf("parseOpen =\n  %+v\nwant\n  %+v", got, tc.want)
			}
			if got.Version != 0 || got.Foreign || got.KwargsHash != "" || got.Body != "" {
				t.Errorf("a v1 marker has Version %d, Foreign %v, KwargsHash %q, Body %q; want none of them",
					got.Version, got.Foreign, got.KwargsHash, got.Body)
			}
		})
	}
}

func TestParseOpenForeign(t *testing.T) {
	tests := []struct {
		name  string
		attrs string
		extra map[string]string
	}{
		{"v is 3", `{"v":3,"tier":"seeded","record":"R","kwargs":{"a":1},"x":true}`,
			map[string]string{"v": `3`, "tier": `"seeded"`, "record": `"R"`, "kwargs": `{"a":1}`, "x": `true`}},
		{"v is 1", `{"v":1,"tier":"owned","kwargs":{"a":1}}`,
			map[string]string{"v": `1`, "tier": `"owned"`, "kwargs": `{"a":1}`}},
		{"v is a string", `{"v":"2","tier":"owned","id":{"a":1},"kwargs":"sha256:k","body":"sha256:b"}`,
			map[string]string{"v": `"2"`, "tier": `"owned"`, "id": `{"a":1}`, "kwargs": `"sha256:k"`, "body": `"sha256:b"`}},
		{"v is fractional", `{"v":2.5}`, map[string]string{"v": `2.5`}},
		{"v is zero", `{"v":0}`, map[string]string{"v": `0`}},
		{"kwargs of any shape is not an error", `{"v":9,"kwargs":5}`, map[string]string{"v": `9`, "kwargs": `5`}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok, err := parseOpen(lineComment("#"), `# sedum:addStep:run `+tc.attrs)
			if err != nil {
				t.Fatalf("parseOpen: %v", err)
			}
			if !ok {
				t.Fatal("a foreign marker is returned with true")
			}
			if !got.Foreign {
				t.Error("Foreign = false, want true")
			}
			if got.Action != "addStep" || got.Variant != "run" {
				t.Errorf("label = %s:%s, want addStep:run", got.Action, got.Variant)
			}
			if got.Tier != TierOwned || got.Writer != "sedum" || got.Version != 0 {
				t.Errorf("Tier %q, Writer %q, Version %d; want owned, sedum, 0", got.Tier, got.Writer, got.Version)
			}
			if got.Record != "" || got.Gen != "" || got.Pkg != "" || len(got.Kwargs) != 0 || got.KwargsHash != "" || got.Body != "" {
				t.Errorf("a foreign marker has fields set: %+v", got)
			}
			extra := map[string]string{}
			for k, v := range got.Extra {
				extra[k] = string(v)
			}
			if !reflect.DeepEqual(extra, tc.extra) {
				t.Errorf("Extra = %v, want every attribute key as read: %v", extra, tc.extra)
			}
		})
	}
}

func TestParseOpenErrors(t *testing.T) {
	const tail = " carries attributes that are not readable as JSON: "
	tests := []struct {
		name string
		line string
		// exact is the whole message where the cause is fixed by the record;
		// otherwise prefix is what the message starts with.
		exact  string
		prefix string
	}{
		{"not JSON", `# sedum:addStep:run {not json}`, "", "marker for addStep:run" + tail},
		{"truncated JSON", `# sedum:addStep {"tier":"owned"`, "", "marker for addStep" + tail},
		{"an array", `# sedum:addStep [1]`, "", "marker for addStep" + tail},
		{"a bare string", `# sedum:addStep "x"`, "", "marker for addStep" + tail},
		{"a bare word", `# sedum:addStep:run owned`, "", "marker for addStep:run" + tail},
		{"v1 kwargs is a string", `# sedum:addStep:run {"tier":"owned","kwargs":"sha256:k"}`,
			"marker for addStep:run" + tail + "kwargs must be an object", ""},
		{"v1 kwargs is an array", `# sedum:addStep {"kwargs":[1]}`,
			"marker for addStep" + tail + "kwargs must be an object", ""},
		{"v1 kwargs is a number", `# sedum:addStep {"tier":"owned","kwargs":7}`,
			"marker for addStep" + tail + "kwargs must be an object", ""},
		{"v2 kwargs is an object", `# sedum:addStep:run {"v":2,"tier":"owned","kwargs":{"a":1}}`,
			"marker for addStep:run" + tail + "kwargs must be a string", ""},
		{"v2 kwargs is a number", `# sedum:addStep {"v":2,"kwargs":7}`,
			"marker for addStep" + tail + "kwargs must be a string", ""},
		{"v2 kwargs is an array", `# sedum:addStep {"v":2,"kwargs":["sha256:k"]}`,
			"marker for addStep" + tail + "kwargs must be a string", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok, err := parseOpen(lineComment("#"), tc.line)
			if err == nil {
				t.Fatalf("parseOpen = %+v, %v; want an error", got, ok)
			}
			if ok {
				t.Error("a marker with unreadable attributes was reported as read")
			}
			if tc.exact != "" && err.Error() != tc.exact {
				t.Errorf("error = %q, want %q", err.Error(), tc.exact)
			}
			if tc.prefix != "" && !strings.HasPrefix(err.Error(), tc.prefix) {
				t.Errorf("error = %q, want it to start with %q", err.Error(), tc.prefix)
			}
		})
	}
}

// A v2 marker read and written back is the same bytes, whatever the unknown
// keys were, once they are in sorted order.
func TestV2MarkerRoundTripsByteForByte(t *testing.T) {
	line := `# sedum:addPolicy:x {"v":2,"tier":"owned","record":"PR-1","gen":"sha256:g","pkg":"sha256:p","writer":"harness","id":{"controller":"users","name":"read"},"kwargs":"sha256:k","body":"sha256:b","attempts":3,"verified_by":"spec/x"}`
	m, ok, err := parseOpen(lineComment("#"), line)
	if err != nil || !ok {
		t.Fatalf("parseOpen = %v, %v", ok, err)
	}
	got, err := m.Open(lineComment("#"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if got != line {
		t.Errorf("round trip changed the marker:\n  %s\n  %s", line, got)
	}
}
