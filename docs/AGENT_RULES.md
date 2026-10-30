# Writing records that Sedum will regenerate

This is guidance for whoever writes provenance records (a person or an agent) for code that Sedum generates and regenerates. The general rules for good records (one requirement per line, no code in a record, behaviour rather than history) apply here too and live with the provenance tooling. What follows are the rules that exist only because a model, not a person, writes the code from the record, and why each exists.

The aim: a record is good enough when the code can be regenerated from it by a model and still pass the repository's tests.

## Say the interface, not only the behaviour

Regeneration fails where a constraint names a behaviour but not the interface. A model that is told a function "returns the user's slug" will invent a name, a parameter list and a return shape, and the invention will differ from the code the tests call. State each name, parameter and result the code must have, and the shape of each value (its fields, or its form), in the constraint.

## Name every collaborator

A model cannot see another record's unit. Anything the record does not state about a collaborator, such as the function it calls, what that returns or which type it takes, will be invented. State each collaborator by name with the shape the unit relies on.

## State type-level steps as behaviour, in words

Type-level steps (a conversion, a narrowing, a cast, a generic parameter) are written as behaviour in words, and each is stated explicitly. A model follows them far less reliably than ordinary behaviour, so an unstated type step is the usual cause of a generated file that does not compile.

## Own each import once per file

Exactly one record in a file states each imported name, with its module and whether it is a value or a type. Every other record in that file that uses the name says that it adds no imports and where the name comes from. If two records name the same import, the generated file declares it twice.

Review check: each imported name is stated by exactly one record per file.

## One top-level function per imprint

Keep one imprint per top-level function. A function together with its imports, a type, or a group of constants is fine as one unit. Never merge two functions into one imprint.

Why: measured on 2026-10-21, Haiku selects only one top-level function per record, even when the intent says two. The second function is silently not generated. This is a measurement of one model on one date; re-measure when the model changes and update the date here.

Review check: an imprint with two or more top-level functions is flagged. A blueprint may still cover several behaviours, one imprint for each function.

## Review before generating

Review a record before any generation or regeneration, so that its faults cost a review and not a model run. A reviewer reads the records and never runs `sedum`. The checks above (collaborators, type-level steps, imports owned once, one function per imprint, and the size of a unit) are the reviewer's Sedum-specific checks.

## Keeping this current

When the owner or the maintainers give a new rule of this kind, add it here with its reason and, if it comes from a measurement, the date.
