package cli

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-ef240c78","gen":"sha256:d649ee55e43c11f0d256d42a6575e7a6988a88592de41be511d9bba9e8614b28","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/cli/growlog_wiring_test.go","path":"testing"}}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"tier":"owned","record":"prov-2026-ef240c78","gen":"sha256:d649ee55e43c11f0d256d42a6575e7a6988a88592de41be511d9bba9e8614b28","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"func TestGrowLogCommandRegistered(t *testing.T) {\n\tfindCommand(t, \"grow-log\")\n}","file":"internal/cli/growlog_wiring_test.go","name":"TestGrowLogCommandRegistered"}}
func TestGrowLogCommandRegistered(t *testing.T) {
	findCommand(t, "grow-log")
}

// /sedum:setFunction
