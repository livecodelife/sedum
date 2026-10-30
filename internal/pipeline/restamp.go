package pipeline

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/restamp.go","path":"fmt"},"kwargs":"sha256:0b7a5f8645c9030627fd98d0c0b100e78b54d834e1ff2d69ce1f5fa9e83e659a","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/restamp.go","path":"github.com/livecodelife/sedum/internal/inject"},"kwargs":"sha256:f576aaed08719652313b3e7dc17ab8837606997af1b6f5e213a1d0626871f38f","body":"sha256:c32cfeba25229dcfd3a66390c7c5b1c1c3cb162181c8360a4df0ff1bdfad701f"}
import "github.com/livecodelife/sedum/internal/inject"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/restamp.go","path":"github.com/livecodelife/sedum/internal/record"},"kwargs":"sha256:9eb5d77e4801594bc8f88d10467fea9e35c58743ed45e56d2236801797c948a0","body":"sha256:527c0befadfcad873cb5f9da4ab560755bc3040424ba6241011c3abf87fcc5e0"}
import "github.com/livecodelife/sedum/internal/record"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/restamp.go","path":"github.com/livecodelife/sedum/internal/recording"},"kwargs":"sha256:964df166f915fdc1ee0fab80b6eaf9bdef4285e0b81e096a0a4ebc5de008089d","body":"sha256:77eda6b9872215a21291b42d59b9ec4b021705be869f51a3d61adc6dcb120e03"}
import "github.com/livecodelife/sedum/internal/recording"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/pipeline/restamp.go","name":"Restamp"},"kwargs":"sha256:074550799e1a8b56f0adbb718716acfc98af6f41b9981d3047d60d364d10f9cc","body":"sha256:d6a6627aacbc0756c21ae5e6aea1be4aff3a04caee512f476d4d8a16f74dc10e"}
// Restamp is one record whose gen a replay took from the current record
// instead of from the recording (prov-2026-c637c262).
type Restamp struct {
	RecordID string
	Old      string
	New      string
}

// /sedum:setType
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/pipeline/restamp.go","name":"String","receiver":"Restamp"},"kwargs":"sha256:618b799de52e1b59c2770b6d91b0ff09837bd880575aac16be3a2e3583a63577","body":"sha256:5aea58e1ec483b9ceb4f680a49e95b744f4b01f814e1d5cde64c9b2df85894c6"}
// String is the one-line report: the record and the first eight characters of
// each gen, without the algorithm prefix.
func (r Restamp) String() string {
	return fmt.Sprintf("restamped %s: gen %s -> %s", r.RecordID, shortGen(r.Old), shortGen(r.New))
}

// /sedum:setMethod
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp.go","name":"shortGen"},"kwargs":"sha256:9c48bd387429654c89f118dd7d90891543473d88e3937db9a0acb190b59659db","body":"sha256:4cd0305ef8976d27a11612b8072f7b53fcc7c8d2aa27c7b292a0a1aa94e154aa"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp.go","name":"replayGenFor"},"kwargs":"sha256:05170e0d67aec9645aa406fab7ac5a1297f2c115181b558bf7d8954aa9b92097","body":"sha256:7099cfd62830f3c955834af4c69fd04a6302c0846b2a0116e877406d41fa72ac"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp.go","name":"replayReconcile"},"kwargs":"sha256:b5f23e2af16861858de7ac57767fb0b03768043d80caff76fcd296aa0c616a78","body":"sha256:6273a32733ca9ef517e3eb46d0e63c151e27fce61d6805b74196d78ed1ce6ea8"}
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
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-010fd5b2","gen":"sha256:b75409b540af694c1bb43f58f412e0fa8c9ac878e8d5f64dfb6ed63753e9d10f","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/restamp.go","name":"ReplayWouldChange"},"kwargs":"sha256:0668c038cfb8733b8f56c5472988f992c08be0ff61daa18c4555ab073a7b48c5","body":"sha256:e8d3341c65f75992e795ad7c2e0854f2e3f2fe16cfa95395a1b9455317ab4ee5"}
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
