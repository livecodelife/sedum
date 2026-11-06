// sedum:anchor:header

package verify

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-ab7a8e5f","gen":"sha256:2ebebcd17e9bb43789097822407d2ed76947c000ea826a212296e3571ec149e5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_windows.go","path":"context"},"kwargs":"sha256:7968766e7442655a210685f584563172355d0426da608cd75dcc22545d9ef77e","body":"sha256:0f7443fb25d7ad6a970617d4d0df6123283e24801b94c1d84fbe431230f31e1e"}
import "context"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-ab7a8e5f","gen":"sha256:2ebebcd17e9bb43789097822407d2ed76947c000ea826a212296e3571ec149e5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_windows.go","path":"os/exec"},"kwargs":"sha256:e4d55a8b1643a727714a325e0cf3167204c0939299d08b5e0feb687417ec8db7","body":"sha256:5970e9ca793f55e2beb08de623e98bdbc0ced178860fedd953eef7e15026ee19"}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-ab7a8e5f","gen":"sha256:2ebebcd17e9bb43789097822407d2ed76947c000ea826a212296e3571ec149e5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_windows.go","path":"strconv"},"kwargs":"sha256:22adc3c9bb8571c0d3fdeb93984a5aed63fe64ce72f8d7caf161f50ad46faeda","body":"sha256:f33f981c64a42d78a5daa7ff3fca15b10045080cccd282fef544b59eaf7f6dc8"}
import "strconv"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-ab7a8e5f","gen":"sha256:2ebebcd17e9bb43789097822407d2ed76947c000ea826a212296e3571ec149e5","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/verify/verify_windows.go","path":"syscall"},"kwargs":"sha256:7fb96f69a226c2b550a33f9365caab14fe6686155c39247727f9a15070dbe727","body":"sha256:1bede81c0dadcc701b7d096e89e517dcc5093c7f22a750533e538bf2882846ba"}
import "syscall"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-ab7a8e5f","gen":"sha256:2ebebcd17e9bb43789097822407d2ed76947c000ea826a212296e3571ec149e5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_windows.go","name":"shellCommand"},"kwargs":"sha256:27faf30df16ab69ed1155171137670c45d7a1aaef7760e36556165d9b1478ea0","body":"sha256:a972339dab49ebda1a070f372afbe0d89b0779aeca604461d2d78c755901700f"}
// shellCommand returns the unstarted command that runs one shell command line:
// through sh -c when an sh is on PATH, as Git for Windows provides, otherwise
// through cmd /C with the line passed to the process unquoted.
func shellCommand(ctx context.Context, command string) *exec.Cmd {
	if sh, err := exec.LookPath("sh"); err == nil {
		return exec.Command(sh, "-c", command)
	}
	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "cmd.exe /C " + command}
	return cmd
}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-ab7a8e5f","gen":"sha256:2ebebcd17e9bb43789097822407d2ed76947c000ea826a212296e3571ec149e5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_windows.go","name":"configureProcess"},"kwargs":"sha256:42e4796bdb26cb1cd6aa3dd1a332b049d25bc5f20da6b00c64770346cae961be","body":"sha256:b91ac1c54dd18fbc6e2c8aac048547762ab05e8bf8c5b12bd812c86b20042910"}
// configureProcess prepares cmd before it starts. Windows needs no process
// group, and it must not replace SysProcAttr, which shellCommand may have set.
func configureProcess(cmd *exec.Cmd) {}

// /sedum:setFunction
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-ab7a8e5f","gen":"sha256:2ebebcd17e9bb43789097822407d2ed76947c000ea826a212296e3571ec149e5","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/verify/verify_windows.go","name":"killTree"},"kwargs":"sha256:38629315ac6c143588186cbadb17a4bddcff11c623c32b523cc893c3972d739b","body":"sha256:29fad632a596a0d02030b8cf8b9383270af99805f070fa7f85e57caa559066dc"}
// killTree terminates the started command and every process it started.
func killTree(cmd *exec.Cmd) {
	_ = exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
	_ = cmd.Process.Kill()
}

// /sedum:setFunction
