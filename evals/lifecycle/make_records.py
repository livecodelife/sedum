#!/usr/bin/env python3
"""Write the testbed's v1 records into fixtures/stepN/records/, one group per step that adds them.

Fixed ids (prov-2026-7e5700NN) keep every check readable; created_at times are
spread so ownership and anchor order are deterministic.

    python3 -I make_records.py
"""
import os

ROOT = os.path.dirname(os.path.abspath(__file__))

BRIEF, BP_API, NOTES_OLD, HEALTH, LEGACY, NOTES, BUG1, DUP = (
    f"prov-2026-7e5700{n:02d}" for n in range(1, 9))

RECORDS = [
    # step 1: the starting v1 set
    ("step1", dict(id=BRIEF, type="brief", title="A notes service", created_at="2026-10-01",
                   intent="A small HTTP service for keeping notes, with a health check for its load balancer.",
                   constraints=["Every endpoint lives under internal/handlers and is registered with chi."],
                   affected_scope=["internal/**"])),
    ("step1", dict(id=BP_API, type="blueprint", implements=BRIEF, title="Notes HTTP API", created_at="2026-10-02",
                   intent="An HTTP API for notes plus a health endpoint, one handler file per resource.",
                   constraints=["One handler file per resource, named after it."],
                   affected_scope=["internal/handlers/**"])),
    ("step1", dict(id=NOTES_OLD, type="imprint", implements=BP_API, title="List and create notes", created_at="2026-10-03",
                   intent="Serve the note resource: GET /note lists notes and POST /note creates one.",
                   constraints=["Handlers live in internal/handlers/note.go."],
                   affected_scope=["internal/handlers/note.go"])),
    ("step1", dict(id=HEALTH, type="imprint", implements=BP_API, title="Health endpoint", created_at="2026-10-04",
                   intent="Serve the health resource: GET /health answers 200 when the service is up.",
                   constraints=["Handlers live in internal/handlers/health.go."],
                   affected_scope=["internal/handlers/health.go"])),
    ("step1", dict(id=LEGACY, type="imprint", implements=BP_API, title="Legacy endpoint", created_at="2026-10-05",
                   intent="Serve the legacy resource: GET /legacy keeps old clients working until they move.",
                   constraints=["Handlers live in internal/handlers/legacy.go."],
                   affected_scope=["internal/handlers/legacy.go"])),
    # step 3: changes made the v1 way (supersession, a bug)
    ("step3", dict(id=NOTES, type="imprint", implements=BP_API, supersedes=NOTES_OLD,
                   title="List, create and delete notes", created_at="2026-10-06",
                   intent="Serve the note resource: GET /note lists notes, POST /note creates one, and DELETE /note removes one.",
                   constraints=["Handlers live in internal/handlers/note.go."],
                   affected_scope=["internal/handlers/note.go"])),
    ("step3", dict(id=BUG1, type="bug", extends=BP_API, title="Handlers return JSON errors", created_at="2026-10-07",
                   intent="Errors come back as plain text, which clients cannot parse. Every handler must return errors as JSON.",
                   constraints=["Every handler error response is JSON with an error field."],
                   affected_scope=["internal/handlers/errors.go"])),
    # step 11: a second record claiming a region another record owns
    ("step11", dict(id=DUP, type="imprint", implements=BP_API, title="Read-only note mirror", created_at="2026-10-08",
                    intent="A read-only mirror of the note resource for reporting: GET /note.",
                    constraints=["Handlers live in internal/handlers/note.go."],
                    affected_scope=["internal/handlers/note.go"])),
]


def q(s):
    return '"' + s.replace("\\", "\\\\").replace('"', '\\"') + '"'


def lst(items):
    return " []" if not items else "\n" + "".join(f"  - {q(i)}\n" for i in items).rstrip("\n")


def render(r, superseded_by=""):
    lines = [
        f"id: {r['id']}",
        f"title: {q(r['title'])}",
        f"status: {r.get('status', 'implemented')}",
        f"created_at: {q(r['created_at'])}",
        "author: testbed@example.invalid",
        f"intent: {q(r['intent'])}",
        f"constraints:{lst(r['constraints'])}",
        f"affected_scope:{lst(r['affected_scope'])}",
        "forbidden_scope: []",
        f"type: {r['type']}",
    ]
    if r.get("implements"):
        lines.append(f"implements: {r['implements']}")
    if r.get("extends"):
        lines.append(f"extends: {r['extends']}")
    lines += [
        f"supersedes: {q(r.get('supersedes', ''))}",
        f"superseded_by: {q(superseded_by)}",
        "related: []",
        'sealed_at_sha: ""',
        "associated_specs: []",
        "associated_traces: []",
        "monitors: []",
        "tags: [testbed]",
    ]
    return "\n".join(lines) + "\n"


def fixture_name(rid):
    # Stored without the prov-YYYY- prefix: a repo's provenance pre-commit hook lints every
    # staged prov-YYYY-*.yml as one of its own records. run.sh restores the name in work/.
    return rid.split("-", 2)[2] + ".yml"


def main():
    for step, r in RECORDS:
        d = os.path.join(ROOT, "fixtures", step, "records")
        os.makedirs(d, exist_ok=True)
        with open(os.path.join(d, fixture_name(r["id"])), "w") as fh:
            fh.write(render(r))
    # step 3 also flips the old imprint to superseded, as `create --supersedes` would
    old = next(r for _, r in RECORDS if r["id"] == NOTES_OLD)
    d = os.path.join(ROOT, "fixtures", "step3", "records")
    with open(os.path.join(d, fixture_name(NOTES_OLD)), "w") as fh:
        fh.write(render(dict(old, status="superseded"), superseded_by=NOTES))


if __name__ == "__main__":
    main()
