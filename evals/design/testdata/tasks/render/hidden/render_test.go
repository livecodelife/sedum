package growlog

import (
	"strings"
	"testing"
)

func TestSummaryText(t *testing.T) {
	s := Summarize(summaryEntries(), Filter{})
	s.Malformed = 1
	text := s.Text()
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	if lines[0] != "grow log: 4 line(s), 1 malformed" {
		t.Fatalf("first line = %q", lines[0])
	}
	if len(lines) != 6 {
		t.Fatalf("got %d lines, want 6:\n%s", len(lines), text)
	}
	want := map[string][]string{
		"haiku": {"2", "1", "5", "4", "1", "150/25", "75%"},
		"stub":  {"1", "0", "1", "0", "0", "0/0", "0%"},
		"total": {"3", "1", "6", "4", "1", "150/25", "50%"},
	}
	for _, l := range lines[2:] {
		f := strings.Fields(l)
		if w, ok := want[f[0]]; ok {
			if strings.Join(f[1:], " ") != strings.Join(w, " ") {
				t.Errorf("row %s = %v, want %v", f[0], f[1:], w)
			}
			delete(want, f[0])
		}
	}
	if len(want) != 0 {
		t.Errorf("rows missing: %v\n%s", want, text)
	}
	if !strings.Contains(text, "(none)") {
		t.Errorf("a run with no model must show as (none):\n%s", text)
	}
}

func TestSummaryJSON(t *testing.T) {
	b, err := Summarize(nil, Filter{}).JSON()
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !strings.HasSuffix(s, "\n") || !strings.Contains(s, `"models": []`) {
		t.Fatalf("JSON = %q", s)
	}
	if strings.Contains(s, "timestamp") {
		t.Fatalf("a summary carries no timestamp: %q", s)
	}
}
