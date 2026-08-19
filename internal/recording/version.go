package recording

// Version is the semver string sedum --version emits and the value a recording
// carries in its sedum_version field.
//
// It lives here, in the package that defines the field, because
// prov-2026-b5465dfa requires the two to be incapable of drifting. A constant
// in the command layer with the writer copying it would satisfy that on the day
// it was written and not afterwards; a caller pinning a version floor is
// relying on the binary's answer and the artifact's agreeing forever.
//
// It is a constant rather than a build-time stamp. A caller told to invoke a
// binary needs to know which one it has, and a version that is absent from a
// `go build` or `go run` invocation would be absent exactly when someone is
// developing against it (prov-2026-b5465dfa).
const Version = "0.1.0"
