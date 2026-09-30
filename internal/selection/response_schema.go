package selection

import "context"

// --response-schema (prov-2026-91c54941): opting a Client into structured
// output without changing the Client interface itself.
//
// The Client interface stays exactly Complete - selection_test.go's and
// pipeline_test.go's stub Clients, and every other record's specs against
// them, keep working unchanged. What decides whether a given call uses
// structured output is a type assertion Select performs internally: a
// caller wraps its Client with WithResponseSchema, and Select recognizes the
// wrapper and, for an underlying Client that also implements SchemaCapable,
// compiles that record's own (already-filtered) catalog into a schema and
// calls CompleteWithSchema instead of Complete.
//
// This is also how a client with no structured-output support - *Local, the
// bundled goinfer-serve backend, today - is refused rather than silently run
// unconstrained: it simply does not implement SchemaCapable, so the type
// assertion inside Select fails and Select reports a clear error before
// making any call. Whether goinfer-serve itself would accept response_format,
// or reproduce prov-2026-4bcabb2f's collapse if it did, is unverified and out
// of this record's scope to resolve - this is why that combination fails
// clearly instead of being wired through to it.

// SchemaCapable is satisfied by a Client that can send a response bound to a
// JSON Schema alongside a plain completion. It is a second, optional
// interface rather than a method Client itself declares, so that adding it
// can never change what every existing Client - real or a test's stub - has
// to implement.
type SchemaCapable interface {
	// CompleteWithSchema behaves like Complete, additively: schemaName and
	// schema (already-compiled JSON Schema bytes) select structured output
	// on top of the same request Complete would otherwise send.
	CompleteWithSchema(ctx context.Context, messages []Message, schemaName string, schema []byte) (Completion, error)
}

// responseSchemaClient marks a Client as opted into --response-schema.
// Select's own type assertion against this exact type is the only place that
// reads the marker - the wrapper otherwise behaves exactly as the Client it
// wraps, so any caller that does not know about --response-schema at all
// (selection_test.go's stub-driven tests among them) sees nothing different.
type responseSchemaClient struct {
	Client
}

// WithResponseSchema opts a Client into --response-schema. A caller (grow.go)
// wraps whichever Client it built, once, before running the pipeline; Select
// compiles a schema per record rather than once per run, because the schema
// is scoped to what survives that record's own authorized-path filter, which
// differs record to record.
func WithResponseSchema(c Client) Client {
	return &responseSchemaClient{Client: c}
}
