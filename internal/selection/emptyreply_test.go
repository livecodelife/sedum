package selection

import (
	"strings"
	"testing"
)

// A model reply with no text must not be sent back as an assistant turn: the
// API refuses an empty message with HTTP 400 (prov-2026-db703ddb). The loop
// still sends its follow-up user message and still counts the attempt.

const emptyReplyAnchored = "class UsersController\n  # sedum:anchor:class_body\nend\n"

// wantNoEmptyMessage fails if any request the stub saw carried a message with
// empty or whitespace-only content.
func wantNoEmptyMessage(t *testing.T, client *stub) {
	t.Helper()
	for i, req := range client.seen {
		for j, m := range req {
			if strings.TrimSpace(m.Content) == "" {
				t.Errorf("request %d message %d (role %q) has empty content", i+1, j+1, m.Role)
			}
		}
	}
}

func wantEndsWithUser(t *testing.T, client *stub, request int, contains string) {
	t.Helper()
	msgs := client.seen[request]
	last := msgs[len(msgs)-1]
	if last.Role != RoleUser || !strings.Contains(last.Content, contains) {
		t.Errorf("request %d should end with a user message containing %q, got role %q: %q", request+1, contains, last.Role, last.Content)
	}
}

// Rejection site: the reply is empty, so it is rejected, and the retry carries
// the rejection text as the user turn with no assistant turn before it.
func TestEmptyReplyRejectionSendsNoAssistantTurn(t *testing.T) {
	for name, reply := range map[string]string{"empty": "", "whitespace": " \n\t "} {
		t.Run(name, func(t *testing.T) {
			got, client, err := selectAnswer(t, railsRequest(t), 1, reply, validResponse)
			if err != nil {
				t.Fatalf("Select: %v", err)
			}
			if got.Rejected != 1 || got.Calls != 2 || got.Repairs != 0 || got.Completeness != 0 {
				t.Errorf("Rejected=%d Calls=%d Repairs=%d Completeness=%d, want 1, 2, 0, 0", got.Rejected, got.Calls, got.Repairs, got.Completeness)
			}
			wantNoEmptyMessage(t, client)
			wantEndsWithUser(t, client, 1, "empty")
			for _, m := range client.seen[1] {
				if m.Role == RoleAssistant {
					t.Errorf("the retry carries an assistant turn for an empty reply: %q", m.Content)
				}
			}
		})
	}
}

// An empty reply exhausting the retries is still a Rejection with the same
// counts, and no request carried an empty message on the way there.
func TestEmptyReplyRejectionKeepsRetryLimit(t *testing.T) {
	_, client, err := selectAnswer(t, railsRequest(t), 2, "", "  ", "")
	rej, ok := err.(*Rejection)
	if !ok {
		t.Fatalf("want *Rejection, got %T: %v", err, err)
	}
	if rej.Rejected != 3 || rej.Calls != 3 {
		t.Errorf("Rejected=%d Calls=%d, want 3 and 3", rej.Rejected, rej.Calls)
	}
	wantNoEmptyMessage(t, client)
}

// Completeness site: an empty reply earlier in the conversation must not leave
// an empty message in the request that carries the completeness note.
func TestEmptyReplyThenCompletenessNote(t *testing.T) {
	req := railsRequest(t)
	req.Files[0].Rendered = emptyReplyAnchored
	got, client, err := selectAnswer(t, req, 2, "", `{"invocations":[]}`, validResponse)
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if got.Rejected != 1 || got.Completeness != 1 || got.Calls != 3 {
		t.Errorf("Rejected=%d Completeness=%d Calls=%d, want 1, 1, 3", got.Rejected, got.Completeness, got.Calls)
	}
	wantNoEmptyMessage(t, client)
	wantEndsWithUser(t, client, 2, "class_body")
	// The non-empty reply that earned the note is carried unchanged.
	msgs := client.seen[2]
	a := msgs[len(msgs)-2]
	if a.Role != RoleAssistant || a.Content != `{"invocations":[]}` {
		t.Errorf("want the non-empty reply carried as an assistant turn, got role %q: %q", a.Role, a.Content)
	}
}

// Verification-repair site, same shape.
func TestEmptyReplyThenVerificationRepair(t *testing.T) {
	check := &failing{n: 1, message: "These checks failed: go build ./..."}
	got, client, err := selectRepairing(t, 1, check.verify, "", validResponse, validResponse)
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if got.Rejected != 1 || got.Repairs != 1 || got.Calls != 3 {
		t.Errorf("Rejected=%d Repairs=%d Calls=%d, want 1, 1, 3", got.Rejected, got.Repairs, got.Calls)
	}
	wantNoEmptyMessage(t, client)
	wantEndsWithUser(t, client, 2, check.message)
	msgs := client.seen[2]
	if a := msgs[len(msgs)-2]; a.Role != RoleAssistant || a.Content != validResponse {
		t.Errorf("want the non-empty reply carried as an assistant turn, got role %q: %q", a.Role, a.Content)
	}
}

// A non-empty rejected reply is still carried as an assistant turn, unchanged,
// followed by the user turn.
func TestEmptyReplyDoesNotDropNonEmptyReply(t *testing.T) {
	bad := `{"invocations":[{"action":"noSuchAction","kwargs":{}}]}`
	_, client, err := selectAnswer(t, railsRequest(t), 1, bad, validResponse)
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	msgs := client.seen[1]
	if len(msgs) < 3 {
		t.Fatalf("want at least 3 messages, got %d", len(msgs))
	}
	a, u := msgs[len(msgs)-2], msgs[len(msgs)-1]
	if a.Role != RoleAssistant || a.Content != bad || u.Role != RoleUser {
		t.Errorf("want assistant(%q) then user, got %q(%q) then %q", bad, a.Role, a.Content, u.Role)
	}
	wantNoEmptyMessage(t, client)
}
