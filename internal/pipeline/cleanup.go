package pipeline

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup.go","path":"os"},"kwargs":"sha256:d9d201f9c642c133f11666688957530e5b5ebab08d251b77adc0f8e462694e8a","body":"sha256:7c6e309eb2575d4aa451ec2ab83161d13e9ae03bbca3a4d2ecf30ce12adf0905"}
import "os"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup.go","path":"path"},"kwargs":"sha256:c9ffcc5c2d3e572db8b841e12011f23937f13beb1030cab0006a7f0c01bec7bf","body":"sha256:c8c46e2b4136d27afcfddd73e765db239e056b6cba04d989689431e735645f13"}
import "path"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup.go","path":"path/filepath"},"kwargs":"sha256:54255374ee2975bc8b2f5cc4ddf11f058029f312061a820f237b40316e227fc7","body":"sha256:a6d196b558491e646be93f56a321f339068c64c76aa9d49bb18046e7f9c777a1"}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup.go","path":"sort"},"kwargs":"sha256:620644543575e1f6fe9a40f73c36dd3a39c3e22220db1ad31bbcac1c522f5301","body":"sha256:f197172617ae21ef23b1855b09aafb14571020ae81bff877d63d3ff81f4bf63c"}
import "sort"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup.go","path":"strings"},"kwargs":"sha256:190b2a14a63d6f3ff4bc5d8fbfb8efba718d643573d1b96d77755eca2ec77274","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup.go","path":"github.com/livecodelife/sedum/internal/resolve"},"kwargs":"sha256:32319f7c64bc8558fa7eb54a6a5c516f13e8ea0d0609da3f4269518916fb1b98","body":"sha256:e12df022982be7cfcdb6de1e42d5d1480f2b3d08431010742e4e2d9d5f13e140"}
import "github.com/livecodelife/sedum/internal/resolve"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/pipeline/cleanup.go","path":"github.com/livecodelife/sedum/internal/runlog"},"kwargs":"sha256:6dd5492b4ceb440307f5498916b37e38fdeca13e0c973a1b7f06a431d3bd23e9","body":"sha256:b60feee19062baf6c10593bd005c135bbc0f06f358a2a5b0d0c592bdad5e726e"}
import "github.com/livecodelife/sedum/internal/runlog"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/pipeline/cleanup.go","name":"runCleanup"},"kwargs":"sha256:5077ef1dc9bb56790910ae8e3a90588ea4fffd78fe1b3f62482e46c2f7ddeda2","body":"sha256:a8697c22278ef8db7e85466face119d25ba6e258cf0cf6358fd11c75dbd6aa23"}
// runCleanup records what a run created so that a failed run can remove it.
type runCleanup struct {
	output string
	dirs   []string
	files  []string
}

// /sedum:setType
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/pipeline/cleanup.go","name":"newRunCleanup"},"kwargs":"sha256:165bec51aeac0f642877ffcfe38278604f986d028f00c961f04a91be5e9bb3bc","body":"sha256:be63587d60a4944a887aff3e1b2f2a5f7ea9ab983c83eec55ec645d7eb99eb0b"}
// newRunCleanup records the ancestor directories of paths that do not exist yet, deepest first.
func newRunCleanup(output string, paths []string) *runCleanup {
	c := &runCleanup{output: output}
	seen := map[string]bool{}
	for _, p := range paths {
		for dir := path.Dir(p); dir != "." && dir != "/" && !seen[dir]; dir = path.Dir(dir) {
			seen[dir] = true
			if _, err := os.Stat(filepath.Join(output, filepath.FromSlash(dir))); os.IsNotExist(err) {
				c.dirs = append(c.dirs, dir)
			}
		}
	}
	sort.SliceStable(c.dirs, func(i, j int) bool {
		return strings.Count(c.dirs[i], "/") > strings.Count(c.dirs[j], "/")
	})
	return c
}

// /sedum:setFunction
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/pipeline/cleanup.go","name":"created","receiver":"runCleanup"},"kwargs":"sha256:6ac9f280db3e044908be8f06511ff8a5231a1bfe576016d474afed794437e1bb","body":"sha256:90eb78bcd96a666773f0281c9df2a25826f7c194c2ba390a107bcd8362a57334"}
// created records the files this run wrote that did not exist before it.
func (c *runCleanup) created(files []resolve.File) {
	for _, f := range files {
		if !f.Existed && !f.Unmanaged {
			c.files = append(c.files, f.Path)
		}
	}
}

// /sedum:setMethod
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-2558649d","gen":"sha256:e272e990fe4924d53d2b169d8a093cdb7f927c8235aa079ab23c06cb072a803a","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/pipeline/cleanup.go","name":"undo","receiver":"runCleanup"},"kwargs":"sha256:07261dd682c90dd585779d4b9539fa62239a08f97048ec875408b0a613c58a2d","body":"sha256:b8c5a8cc99a0f0390503894589cb020e9dcf52442e6456d8464801ec5fcf6261"}
// undo removes the recorded files, then the recorded directories deepest first.
func (c *runCleanup) undo(log *runlog.Log) {
	if log == nil {
		log = runlog.Discard()
	}
	for _, f := range c.files {
		if err := os.Remove(filepath.Join(c.output, filepath.FromSlash(f))); err == nil {
			log.Info("removed file created by a failed run", "path", f)
		}
	}
	for _, d := range c.dirs {
		_ = os.Remove(filepath.Join(c.output, filepath.FromSlash(d)))
	}
}

// /sedum:setMethod
