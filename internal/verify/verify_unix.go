// sedum:anchor:header
// sedum:setBuildConstraint:unix {"v":2,"tier":"owned","record":"prov-2026-521fe542","gen":"sha256:8287731096042d9b924342ec733f960b6104e040ed90980ae7668091e4dd65a8","pkg":"sha256:e2efa6e46e6ce7d19a2a47c83c100bd0cd1267447a0338c145217cdfd6429112","id":{"file":"internal/verify/verify_unix.go"},"kwargs":"sha256:f8734a10c7e6529b8fb6be2ed7814818fc0089d4c9f47c23fb1b4996d06a8f0a","body":"sha256:ea2c2ccba781f3ab2d8664fbfc2936bce438c3004178e81464b93bdb88d18749"}
//go:build unix

// /sedum:setBuildConstraint:unix

package verify

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-521fe542","gen":"sha256:8287731096042d9b924342ec733f960b6104e040ed90980ae7668091e4dd65a8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_unix.go","path":"context"},"kwargs":"sha256:423edbe27133a567b396896ec28f49920e97ecdadb81c2dffdec0cbb789ecbd0","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-521fe542","gen":"sha256:8287731096042d9b924342ec733f960b6104e040ed90980ae7668091e4dd65a8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_unix.go","path":"os/exec"},"kwargs":"sha256:8579c28a879b978be67771f562c424f7291ace75c934046735e630f459bba60f","body":"sha256:5970e9ca793f55e2beb08de623e98bdbc0ced178860fedd953eef7e15026ee19"}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-521fe542","gen":"sha256:8287731096042d9b924342ec733f960b6104e040ed90980ae7668091e4dd65a8","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_unix.go","path":"syscall"},"kwargs":"sha256:410ceb074254042c31cfc5b9c3a0f0dba24c91c0175081c3d573db262920c853","body":"sha256:1bede81c0dadcc701b7d096e89e517dcc5093c7f22a750533e538bf2882846ba"}
import "syscall"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-521fe542","gen":"sha256:8287731096042d9b924342ec733f960b6104e040ed90980ae7668091e4dd65a8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_unix.go","name":"shellCommand"},"kwargs":"sha256:b80c7d4ba613aac82af5523369478d9c9dd61e27148935828edc840224151d2f","body":"sha256:6e8b7039f0cfb247d39a8c34847f9b54592e8bfb994212a29d9fa683402a7bb5"}
// shellCommand returns the unstarted command that runs one shell command line
// through sh -c. ctx is unused here; the parameter keeps the signature the same
// on every platform.
func shellCommand(ctx context.Context, command string) *exec.Cmd {
	return exec.Command("sh", "-c", command)
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-521fe542","gen":"sha256:8287731096042d9b924342ec733f960b6104e040ed90980ae7668091e4dd65a8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_unix.go","name":"configureProcess"},"kwargs":"sha256:2cf1b91a64fb7eeaebf59799d2f2c7b602a1608128694a307745405b7dfe900c","body":"sha256:af61c0ce654492dff0acff1aefe0f9a9141fa51126e0a559c6ce1aeaca7c53ed"}
// configureProcess prepares cmd before it starts: the command leads its own
// process group, so killTree can signal everything it started.
func configureProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-521fe542","gen":"sha256:8287731096042d9b924342ec733f960b6104e040ed90980ae7668091e4dd65a8","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_unix.go","name":"killTree"},"kwargs":"sha256:a51c9d43fbb998e81f861b0ebe3e3b1e9117178ac904b17fbd60c3f1631aec1e","body":"sha256:7463d7bd9329e0f93de2cee5e4b279fad4ff216e1b9038425c3ef03024307690"}
// killTree terminates the started command and every process in its group.
func killTree(cmd *exec.Cmd) {
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

// /sedum:setFunction
