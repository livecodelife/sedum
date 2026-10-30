package cli

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-ef240c78","gen":"sha256:d649ee55e43c11f0d256d42a6575e7a6988a88592de41be511d9bba9e8614b28","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/cli/growlog_wiring_test.go","path":"testing"},"kwargs":"sha256:3dd2cea9ff99f800259ad980647772c21ab6d15140ff43806bb9ae42e31a71ce","body":"sha256:14a4a5f3bfbf9325de04706efa8c15a9446f01d2a26624ce9c66bd1d69dc212e"}
import "testing"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-ef240c78","gen":"sha256:d649ee55e43c11f0d256d42a6575e7a6988a88592de41be511d9bba9e8614b28","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/cli/growlog_wiring_test.go","name":"TestGrowLogCommandRegistered"},"kwargs":"sha256:5584ecccd372d562c54af30b20c3f38502fc149b14675980a768d591e97516c6","body":"sha256:907a93077e2955d44dac32003c19eaa52929bda5e8e59583d19604e5e4b5b1bc"}
func TestGrowLogCommandRegistered(t *testing.T) {
	findCommand(t, "grow-log")
}

// /sedum:setFunction
