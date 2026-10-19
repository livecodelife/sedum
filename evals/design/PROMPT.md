# Context for growing package design (advisory; not part of gen_hash)

Module: github.com/livecodelife/sedum. Package `design` (evals/design) is the
behaviour-only-intent eval harness. It is NOT part of sedum: nothing under
internal/ or cmd/ imports it. Standard library plus gopkg.in/yaml.v3 and the
internal growlog package. Every exported identifier has a doc comment. All code
must be gofmt-clean. Tests never call a model. Your job is TRANSCRIPTION: the
record's intent contains the exact code; write it unchanged, one invocation per
declaration (imports -> addImport, funcs -> setFunction, consts -> setConst,
types -> setType, methods -> setMethod).
