package builtin

import "embed"

// packagesFS carries the built-in generator packages, one zip each, produced by
// hack/sync-builtin.sh. It is declared by hand: a go:embed directive must sit
// directly above its var, and the ownership marker a grown region plants there
// makes gofmt rewrite the comment block, so Sedum cannot grow it.
//
//go:embed packages/*.zip
var packagesFS embed.FS
