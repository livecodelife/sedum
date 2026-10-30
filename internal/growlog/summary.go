package growlog

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/summary.go","path":"sort"},"kwargs":"sha256:bb05d4e14452d4efb84f463c79b499f2269523d25f246501e6182fca900f1854","body":"sha256:f197172617ae21ef23b1855b09aafb14571020ae81bff877d63d3ff81f4bf63c"}
import "sort"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/growlog/summary.go","path":"time"},"kwargs":"sha256:e855050b466d925b297ef01f1b75cae8e603761b10541fb0608e843e5df2b733","body":"sha256:e99128ebf4e36998f40b33adf2eef57f4d6c315f1a4d1bdd3b7fc09a0d7bc82d"}
import "time"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/growlog/summary.go","name":"Filter"},"kwargs":"sha256:eb3775d630d999d352cb9031331bc1c23fc532b79f34bed7087d9d3985a7cae0","body":"sha256:c864dfb104422ca50c77557315c8f86b854dbb5e8d1acafcc3bb419eeeacd18d"}
// Filter selects which entries Summarize counts. Zero fields match everything.
type Filter struct {
	Since  time.Time
	Record string
	Model  string
}

// /sedum:setType
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/growlog/summary.go","name":"match","receiver":"Filter"},"kwargs":"sha256:39c25703bf236745b9c75a06883c7f8b0a65af420a6b08d6cc5635644886d2d2","body":"sha256:f2082f3eca00721c7338d0c197680231573d6440d6d67366a1c78d4c0e5a2866"}
// match reports whether e passes the filter. An entry whose ts does not parse
// fails a non-zero Since.
func (f Filter) match(e Entry) bool {
	if f.Record != "" && e.Record != f.Record {
		return false
	}
	if f.Model != "" && e.Model != f.Model {
		return false
	}
	if !f.Since.IsZero() {
		ts, err := time.Parse(time.RFC3339, e.TS)
		if err != nil || ts.Before(f.Since) {
			return false
		}
	}
	return true
}

// /sedum:setMethod
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/growlog/summary.go","name":"ModelStats"},"kwargs":"sha256:e6dfedb6fbd4df59f83eb50121dcbe8528c31eda6d9941ea8f51b253bd05ccf8","body":"sha256:95c329d5f15611d01beb7cd72f09923399755a5c0c3819ade9fec370bf404e1e"}
// ModelStats is the totals for one model. Retries are the answers the model
// had refused (Entry.Rejected); Repairs are verification send-backs.
type ModelStats struct {
	Model            string `json:"model"`
	Grows            int    `json:"grows"`
	Failed           int    `json:"failed"`
	Calls            int    `json:"calls"`
	Retries          int    `json:"retries"`
	Repairs          int    `json:"repairs"`
	Completeness     int    `json:"completeness"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	Invocations      int    `json:"invocations"`
	Literal          int    `json:"literal"`
}

// /sedum:setType
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/growlog/summary.go","name":"LiteralShare","receiver":"ModelStats"},"kwargs":"sha256:64ed066202e9a3de62801d02aa46d468122d6760c28ea4e0347eddfb1f3fdee1","body":"sha256:9ab403208e25c8dec5014756ec61ddc2f87a1edbc82f09079c3cf04233498634"}
// LiteralShare is the fraction of invocations that carried model-written code,
// or 0 when there were none.
func (m ModelStats) LiteralShare() float64 {
	if m.Invocations == 0 {
		return 0
	}
	return float64(m.Literal) / float64(m.Invocations)
}

// /sedum:setMethod
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/growlog/summary.go","name":"FailureRate","receiver":"ModelStats"},"kwargs":"sha256:2b718da1155c3a91fab0bd10809fc6c72d93bb5bc4c269245b13614052807d84","body":"sha256:357478fd550b3f359b0c8845c1439b8ac98229ba608972bc6819a79fc5fddba3"}
// FailureRate is the fraction of grows that failed, or 0 when there were none.
func (m ModelStats) FailureRate() float64 {
	if m.Grows == 0 {
		return 0
	}
	return float64(m.Failed) / float64(m.Grows)
}

// /sedum:setMethod
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/growlog/summary.go","name":"Summary"},"kwargs":"sha256:1bce523b4417fe42d6edc431551feb5bbcd637cb0630807771d615c45d21a0f4","body":"sha256:c56e72390d267666e433798322f3f5b9ebf03d5198c1d860f49a4413e215c359"}
// Summary is Summarize's result: per-model stats sorted by model name, and the
// total over all of them.
type Summary struct {
	Entries   int          `json:"entries"`
	Malformed int          `json:"malformed"`
	Models    []ModelStats `json:"models"`
	Total     ModelStats   `json:"total"`
}

// /sedum:setType
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/growlog/summary.go","name":"add","receiver":"ModelStats"},"kwargs":"sha256:4fcf67653e9448be115dd3263590f02d7a524c0bb05d7f2ce0da63579be42d32","body":"sha256:204fc1d14f833c14061baf0fe69ef5c135f514a087faf7a931c9c77fd3156b78"}
// add counts e into m. A skip line is not a grow.
func (m *ModelStats) add(e Entry) {
	if e.Mode != "skip" {
		m.Grows++
		if e.Outcome == "failed" {
			m.Failed++
		}
	}
	m.Calls += e.Calls
	m.Retries += e.Rejected
	m.Repairs += e.Repairs
	m.Completeness += e.Completeness
	m.PromptTokens += e.PromptTokens
	m.CompletionTokens += e.CompletionTokens
	m.Invocations += e.Invocations
	m.Literal += e.Literal
}

// /sedum:setMethod
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-78be21c1","gen":"sha256:6dc1f98c8c3071a5eda80455e83d39a557cba593d6fb1dabc9bccf5eeefb6b2a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/growlog/summary.go","name":"Summarize"},"kwargs":"sha256:de223267b4e0becac6d575cd52b9d000bc75192fc045c0b6014de739f91107b2","body":"sha256:ab2a7a90e02d9b9ced501a42ac9172ff9feabb4fdf71c8a0e9b4fdc6f8d89057"}
// Summarize totals the entries that pass f, per model and overall.
func Summarize(entries []Entry, f Filter) Summary {
	by := map[string]*ModelStats{}
	var s Summary
	for _, e := range entries {
		if !f.match(e) {
			continue
		}
		s.Entries++
		m := by[e.Model]
		if m == nil {
			m = &ModelStats{Model: e.Model}
			by[e.Model] = m
		}
		m.add(e)
		s.Total.add(e)
	}
	for _, m := range by {
		s.Models = append(s.Models, *m)
	}
	sort.Slice(s.Models, func(i, j int) bool { return s.Models[i].Model < s.Models[j].Model })
	if s.Models == nil {
		s.Models = []ModelStats{}
	}
	return s
}

// /sedum:setFunction
