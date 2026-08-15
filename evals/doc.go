// Package evals measures how well a model selects actions, and how well what it
// selected works.
//
// It is not part of Sedum. Nothing in internal/ or cmd/ imports it, it is not
// compiled into the binary, and its runner is behind a build tag so that
// go test ./... never reaches it. That separation is load-bearing rather than
// tidy: Sedum does not run or grade the code it generates, and a harness that
// does both would make that non-goal untrue by adjacency if it shipped as part
// of the tool. This grades Sedum from outside.
//
// # Why it exists
//
// The measurement that produced prov-2026-6d87dc11's correction was nineteen
// shell loops with grep -c. One run's output was lost before it could be read,
// two more were mangled by a bad sed, and every number in that record was hand
// tallied. The finding was worth having and the method was not repeatable, which
// is the whole argument for this package.
//
// # What it measures
//
// Two things, at different costs.
//
// Selection is cheap and needs only a model endpoint: run a case N times and
// report, per action, how often it was selected and how many invocations it
// drew. Rates rather than assertions, because the thing being measured is a
// sampled distribution - an eval that failed a build on a coin flip would be a
// flaky test wearing a suit.
//
// Behavior is expensive and needs the generated application to run: apply what
// was selected, start the target, and report the fraction of its linespec
// contracts that pass. Reserved in the schema and not yet implemented.
//
// # The matrix
//
// A single measurement says almost nothing, because every interesting question
// is comparative: does this model do better than that one, does a tighter
// generator package beat a looser one, does an application of this shape hold up
// where a simpler one did. So a case carries the axes it varies along even where
// only one value of each exists today - model, framework, application and its
// complexity tier, how tightly the generator package constrains the model, and
// whether the arm used Sedum at all.
//
// The baseline arm is the one worth explaining. A rate is only meaningful next
// to what the same model produces without Sedum on the same application, and
// without that column a good number is unfalsifiable.
package evals
