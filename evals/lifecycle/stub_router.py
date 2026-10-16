#!/usr/bin/env python3
"""A deterministic OpenAI-compatible chat endpoint for the lifecycle testbed.

Extends sedum's evals/behavior/stub_model.py: instead of one fixed answer, it
picks a canned answer by the first route whose `match` text appears in the
request (sedum's selection prompt carries the record's intent and authorized
files, but not its id). Every request is appended to the call log, so "grow
made no model calls" is checked by the log's line count.

    python3 stub_router.py <routes.json> <port> <call-log>

routes.json: [{"match": "text in the prompt", "answer": "file.json"}, ...]
Answer paths are relative to routes.json.
"""
import json
import os
import sys
from http.server import BaseHTTPRequestHandler, HTTPServer

ROUTES = []
BASE = ""
LOG = ""


def pick(prompt):
    for r in ROUTES:
        if r["match"] in prompt:
            with open(os.path.join(BASE, r["answer"])) as fh:
                return r["answer"], fh.read()
    return None, json.dumps({"invocations": []})


class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        req = json.loads(self.rfile.read(length) or b"{}")
        prompt = "\n".join(str(m.get("content", "")) for m in req.get("messages", []))
        name, answer = pick(prompt)
        with open(LOG, "a") as fh:
            fh.write(json.dumps({"route": name}) + "\n")
        body = json.dumps({
            "id": "stub", "object": "chat.completion", "model": "stub",
            "choices": [{"index": 0, "message": {"role": "assistant", "content": answer},
                         "finish_reason": "stop"}],
            "usage": {"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
        }).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *args):
        pass


def main():
    global ROUTES, BASE, LOG
    routes, port, LOG = sys.argv[1], int(sys.argv[2]), sys.argv[3]
    BASE = os.path.dirname(os.path.abspath(routes))
    with open(routes) as fh:
        ROUTES = json.load(fh)
    HTTPServer(("127.0.0.1", port), Handler).serve_forever()


if __name__ == "__main__":
    main()
