package pipeline

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/restamp.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/restamp.go","path":"github.com/livecodelife/sedum/internal/inject"}}
import "github.com/livecodelife/sedum/internal/inject"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/restamp.go","path":"github.com/livecodelife/sedum/internal/record"}}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/pipeline/restamp.go","path":"github.com/livecodelife/sedum/internal/recording"}}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Restamp is one record whose gen a replay took from the current record\n// instead of from the recording (prov-2026-c637c262).\ntype Restamp struct {\n\tRecordID string\n\tOld      string\n\tNew      string\n}","file":"internal/pipeline/restamp.go","name":"Restamp"}}
// Restamp is one record whose gen a replay took from the current record
// instead of from the recording (prov-2026-c637c262).
type Restamp struct {
	RecordID string
	Old      string
	New      string
}

// /sedum:setType
// sedum:setMethod {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// String is the one-line report: the record and the first eight characters of\n// each gen, without the algorithm prefix.\nfunc (r Restamp) String() string {\n\treturn fmt.Sprintf(\"restamped %s: gen %s -> %s\", r.RecordID, shortGen(r.Old), shortGen(r.New))\n}","file":"internal/pipeline/restamp.go","name":"String","receiver":"Restamp"}}
// String is the one-line report: the record and the first eight characters of
// each gen, without the algorithm prefix.
func (r Restamp) String() string {
	return fmt.Sprintf("restamped %s: gen %s -> %s", r.RecordID, shortGen(r.Old), shortGen(r.New))
}

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// shortGen is the first eight hex characters of a gen, \"none\" when it is empty.\nfunc shortGen(gen string) string {\n\tconst prefix = \"sha256:\"\n\tif len(gen) > len(prefix) && gen[:len(prefix)] == prefix {\n\t\tgen = gen[len(prefix):]\n\t}\n\tif gen == \"\" {\n\t\treturn \"none\"\n\t}\n\tif len(gen) > 8 {\n\t\treturn gen[:8]\n\t}\n\treturn gen\n}","file":"internal/pipeline/restamp.go","name":"shortGen"}}
// shortGen is the first eight hex characters of a gen, "none" when it is empty.
func shortGen(gen string) string {
	const prefix = "sha256:"
	if len(gen) > len(prefix) && gen[:len(prefix)] == prefix {
		gen = gen[len(prefix):]
	}
	if gen == "" {
		return "none"
	}
	if len(gen) > 8 {
		return gen[:8]
	}
	return gen
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// replayGenFor is the gen a replay stamps on a record's regions. Without\n// restamp it is what replayGen decides. With restamp it is the current\n// record's gen, and the second result says what changed; a record that is\n// missing or deprecated has no current gen and keeps the recorded one, and a\n// record whose gen did not change reports nothing.\nfunc replayGenFor(rec recording.Recording, records *record.Set, recordID string, restamp bool) (string, *Restamp) {\n\tgen := replayGen(rec, records, recordID)\n\tif !restamp {\n\t\treturn gen, nil\n\t}\n\tcurrent := genOf(records, recordID)\n\tif current == \"\" || current == gen {\n\t\treturn gen, nil\n\t}\n\treturn current, &Restamp{RecordID: recordID, Old: gen, New: current}\n}","file":"internal/pipeline/restamp.go","name":"replayGenFor"}}
// replayGenFor is the gen a replay stamps on a record's regions. Without
// restamp it is what replayGen decides. With restamp it is the current
// record's gen, and the second result says what changed; a record that is
// missing or deprecated has no current gen and keeps the recorded one, and a
// record whose gen did not change reports nothing.
func replayGenFor(rec recording.Recording, records *record.Set, recordID string, restamp bool) (string, *Restamp) {
	gen := replayGen(rec, records, recordID)
	if !restamp {
		return gen, nil
	}
	current := genOf(records, recordID)
	if current == "" || current == gen {
		return gen, nil
	}
	return current, &Restamp{RecordID: recordID, Old: gen, New: current}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// replayReconcile names the records a replay reconciles: every record the\n// recording selected for, over its managed files. Whether that prunes or only\n// reports is the recording's to say (recording.Recording.Complete).\nfunc replayReconcile(selections []Selection) []inject.Reconcile {\n\tvar out []inject.Reconcile\n\tfor _, s := range selections {\n\t\tif !anyManaged(s.Files) {\n\t\t\tcontinue\n\t\t}\n\t\tout = append(out, inject.Reconcile{RecordID: s.RecordID, Targets: targetsOf(s.Files)})\n\t}\n\treturn out\n}","file":"internal/pipeline/restamp.go","name":"replayReconcile"}}
// replayReconcile names the records a replay reconciles: every record the
// recording selected for, over its managed files. Whether that prunes or only
// reports is the recording's to say (recording.Recording.Complete).
func replayReconcile(selections []Selection) []inject.Reconcile {
	var out []inject.Reconcile
	for _, s := range selections {
		if !anyManaged(s.Files) {
			continue
		}
		out = append(out, inject.Reconcile{RecordID: s.RecordID, Targets: targetsOf(s.Files)})
	}
	return out
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// ReplayWouldChange reports that a replay's result differs from the tree: a\n// region it writes with different bytes, or a region of a named record it no\n// longer selects. It is the verdict of a replay run with --check.\nfunc ReplayWouldChange(applied []inject.Result) bool {\n\tfor _, r := range applied {\n\t\tif r.Changed || r.Unselected {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}","file":"internal/pipeline/restamp.go","name":"ReplayWouldChange"}}
// ReplayWouldChange reports that a replay's result differs from the tree: a
// region it writes with different bytes, or a region of a named record it no
// longer selects. It is the verdict of a replay run with --check.
func ReplayWouldChange(applied []inject.Result) bool {
	for _, r := range applied {
		if r.Changed || r.Unselected {
			return true
		}
	}
	return false
}

// /sedum:setFunction
