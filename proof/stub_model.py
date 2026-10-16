"""A local stand-in for a chat-completions endpoint, for proof/capture.sh.

It answers with canned selections chosen by text in the prompt, and appends
one line per call to the log file, so the proof can count model calls. It binds
to 127.0.0.1 only and never forwards anything: no real model is involved.

Usage: python3 stub_model.py <port> <routes.json> <calls.log>
routes.json maps a prompt substring to the response content to return.
"""
import json
import sys
from http.server import BaseHTTPRequestHandler, HTTPServer

port, routes_path, log_path = int(sys.argv[1]), sys.argv[2], sys.argv[3]


class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        body = json.loads(self.rfile.read(int(self.headers["Content-Length"])))
        prompt = "\n".join(m.get("content", "") for m in body.get("messages", []))
        with open(routes_path) as f:
            routes = json.load(f)
        content = next((answer for needle, answer in routes.items() if needle in prompt), '{"invocations": []}')
        with open(log_path, "a") as log:
            log.write(json.dumps({"matched": content[:60]}) + "\n")
        out = json.dumps({
            "id": "stub", "object": "chat.completion", "model": "stub",
            "choices": [{"index": 0, "finish_reason": "stop",
                         "message": {"role": "assistant", "content": content}}],
            "usage": {"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
        }).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(out)))
        self.end_headers()
        self.wfile.write(out)

    def do_GET(self):
        self.send_response(200)
        self.end_headers()

    def log_message(self, *args):
        pass


HTTPServer(("127.0.0.1", port), Handler).serve_forever()
