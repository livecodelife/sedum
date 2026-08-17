#!/usr/bin/env python3
"""An OpenAI-compatible chat completions endpoint that always returns the same answer.

Sedum's only non-deterministic phase is the model call, and a behaviour harness
that has to make one cannot tell "the generated code is wrong" apart from "the
model chose badly". Pointing OPENAI_BASE_URL here replays a fixed selection
through the real `sedum grow` binary, so everything downstream of Phase 4 is
exercised exactly as it ships. It is a stand-in for `grow --execute`, which is
not built yet (M7).

    python3 stub_model.py <answer.json> <port>
"""
import json
import sys
from http.server import BaseHTTPRequestHandler, HTTPServer

ANSWER = ""


class Handler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", 0))
        self.rfile.read(length)
        body = json.dumps(
            {
                "id": "stub",
                "object": "chat.completion",
                "model": "canned",
                "choices": [
                    {
                        "index": 0,
                        "message": {"role": "assistant", "content": ANSWER},
                        "finish_reason": "stop",
                    }
                ],
                "usage": {"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
            }
        ).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        body = json.dumps({"object": "list", "data": [{"id": "canned"}]}).encode()
        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *args):
        pass


if __name__ == "__main__":
    with open(sys.argv[1]) as fh:
        ANSWER = fh.read()
    HTTPServer(("127.0.0.1", int(sys.argv[2])), Handler).serve_forever()
