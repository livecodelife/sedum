# Context for growing package growlog (advisory; not part of gen_hash)

Module: github.com/livecodelife/sedum. Package `growlog` (internal/growlog) is the
persistent grow log: an append-only JSON Lines file, one line per record per
`sedum grow`. Standard library only. Every exported identifier has a doc comment.
All code must be gofmt-clean. Tests are table-driven or small and use t.TempDir().
Never run git in tests. Your job is TRANSCRIPTION: the record's intent contains
the exact code; write it unchanged, one invocation per declaration.
