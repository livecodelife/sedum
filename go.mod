module github.com/livecodelife/sedum

go 1.25.7

require (
	github.com/iancoleman/strcase v0.3.0
	github.com/sashabaranov/go-openai v1.42.0
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.9
	gopkg.in/yaml.v3 v3.0.1
)

// The shared provenance record library, which owns gen_hash. It has no
// fetchable version yet, so it resolves to a sibling checkout; this branch must
// not reach main until it does (prov-2026-dac49f78).
require github.com/livecodelife/provenance v0.0.0-00010101000000-000000000000

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	golang.org/x/text v0.34.0 // indirect
)

replace github.com/livecodelife/provenance => ../provenance
