package verify

// sedum:anchor:imports
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"bytes"}}
import "bytes"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"context"}}
import "context"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"errors"}}
import "errors"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"fmt"}}
import "fmt"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"io"}}
import "io"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"io/fs"}}
import "io/fs"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"os"}}
import "os"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"os/exec"}}
import "os/exec"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"path/filepath"}}
import "path/filepath"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"strings"}}
import "strings"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"syscall"}}
import "syscall"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"time"}}
import "time"

// /sedum:addImport
// sedum:addImport {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","kwargs":{"file":"internal/verify/verify.go","path":"gopkg.in/yaml.v3"}}
import "gopkg.in/yaml.v3"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setConst {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// ConfigPath is where a project names its verification commands, relative to\n// the output directory (prov-2026-b5c21813).\nconst ConfigPath = \".sedum/config.yml\"","file":"internal/verify/verify.go","name":"ConfigPath"}}
// ConfigPath is where a project names its verification commands, relative to
// the output directory (prov-2026-b5c21813).
const ConfigPath = ".sedum/config.yml"

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// DefaultTimeout bounds each command when the project names none.\nconst DefaultTimeout = 120 * time.Second","file":"internal/verify/verify.go","name":"DefaultTimeout"}}
// DefaultTimeout bounds each command when the project names none.
const DefaultTimeout = 120 * time.Second

// /sedum:setConst
// sedum:setConst {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:51be1c18998d63cc558c253e211a6bd188d96cb251dad68c6ea88e03c2d5cc4c","kwargs":{"code":"// MaxOutput is how much of a failing command's output is fed back to the\n// model: the last bytes, where the error usually is.\nconst MaxOutput = 4000","file":"internal/verify/verify.go","name":"MaxOutput"}}
// MaxOutput is how much of a failing command's output is fed back to the
// model: the last bytes, where the error usually is.
const MaxOutput = 4000

// /sedum:setConst
// sedum:setType {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Config is a project's verification settings: shell commands run, in order,\n// against the tree a model's selection would produce. Commands come from the\n// project's own config file or the command line, never from a generator\n// package.\ntype Config struct {\n\tCommands []string\n\tTimeout  time.Duration\n\t// Skip are directory names not copied into the scratch tree, besides every\n\t// hidden directory.\n\tSkip []string\n}","file":"internal/verify/verify.go","name":"Config"}}
// Config is a project's verification settings: shell commands run, in order,
// against the tree a model's selection would produce. Commands come from the
// project's own config file or the command line, never from a generator
// package.
type Config struct {
	Commands []string
	Timeout  time.Duration
	// Skip are directory names not copied into the scratch tree, besides every
	// hidden directory.
	Skip []string
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"// Failure is one command that did not pass.\ntype Failure struct {\n\tCommand string\n\tStatus  string\n\tOutput  string\n}","file":"internal/verify/verify.go","name":"Failure"}}
// Failure is one command that did not pass.
type Failure struct {
	Command string
	Status  string
	Output  string
}

// /sedum:setType
// sedum:setType {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","kwargs":{"code":"type fileConfig struct {\n\tVerify  []string `yaml:\"verify\"`\n\tTimeout string   `yaml:\"verify_timeout\"`\n}","file":"internal/verify/verify.go","name":"fileConfig"}}
type fileConfig struct {
	Verify  []string `yaml:"verify"`
	Timeout string   `yaml:"verify_timeout"`
}

// /sedum:setType
// sedum:setFunction {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// LoadConfig reads ConfigPath under output. A missing file is not an error: it\n// is a project that verifies nothing.\nfunc LoadConfig(output string) (Config, error) {\n\tcfg := Config{Timeout: DefaultTimeout}\n\tdata, err := os.ReadFile(filepath.Join(output, filepath.FromSlash(ConfigPath)))\n\tif errors.Is(err, fs.ErrNotExist) {\n\t\treturn cfg, nil\n\t}\n\tif err != nil {\n\t\treturn Config{}, err\n\t}\n\tvar file fileConfig\n\tif err := yaml.Unmarshal(data, &file); err != nil {\n\t\treturn Config{}, fmt.Errorf(\"%s: %w\", ConfigPath, err)\n\t}\n\tcfg.Commands = file.Verify\n\tif file.Timeout != \"\" {\n\t\ttimeout, err := time.ParseDuration(file.Timeout)\n\t\tif err != nil || timeout <= 0 {\n\t\t\treturn Config{}, fmt.Errorf(\"%s: verify_timeout %q is not a positive duration\", ConfigPath, file.Timeout)\n\t\t}\n\t\tcfg.Timeout = timeout\n\t}\n\treturn cfg, nil\n}","file":"internal/verify/verify.go","name":"LoadConfig"}}
// LoadConfig reads ConfigPath under output. A missing file is not an error: it
// is a project that verifies nothing.
func LoadConfig(output string) (Config, error) {
	cfg := Config{Timeout: DefaultTimeout}
	data, err := os.ReadFile(filepath.Join(output, filepath.FromSlash(ConfigPath)))
	if errors.Is(err, fs.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return Config{}, err
	}
	var file fileConfig
	if err := yaml.Unmarshal(data, &file); err != nil {
		return Config{}, fmt.Errorf("%s: %w", ConfigPath, err)
	}
	cfg.Commands = file.Verify
	if file.Timeout != "" {
		timeout, err := time.ParseDuration(file.Timeout)
		if err != nil || timeout <= 0 {
			return Config{}, fmt.Errorf("%s: verify_timeout %q is not a positive duration", ConfigPath, file.Timeout)
		}
		cfg.Timeout = timeout
	}
	return cfg, nil
}

// /sedum:setFunction
// sedum:setMethod {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// With returns the config with extra commands run after its own.\nfunc (c Config) With(extra []string) Config {\n\tc.Commands = append(append([]string(nil), c.Commands...), extra...)\n\treturn c\n}","file":"internal/verify/verify.go","name":"With","receiver":"Config"}}
// With returns the config with extra commands run after its own.
func (c Config) With(extra []string) Config {
	c.Commands = append(append([]string(nil), c.Commands...), extra...)
	return c
}

// /sedum:setMethod
// sedum:setMethod {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","kwargs":{"code":"// Enabled reports that there is something to run.\nfunc (c Config) Enabled() bool { return len(c.Commands) > 0 }","file":"internal/verify/verify.go","name":"Enabled","receiver":"Config"}}
// Enabled reports that there is something to run.
func (c Config) Enabled() bool { return len(c.Commands) > 0 }

// /sedum:setMethod
// sedum:setFunction {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Run applies candidate (paths relative to output, with their new content) over\n// a scratch copy of output and runs each command there, printing it to log\n// first. It returns the commands that failed, none when all passed. The\n// scratch directory is removed on every exit path.\nfunc Run(ctx context.Context, cfg Config, output string, candidate map[string]string, log io.Writer) ([]Failure, error) {\n\tscratch, err := os.MkdirTemp(\"\", \"sedum-verify-\")\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tdefer os.RemoveAll(scratch)\n\n\tif err := copyTree(output, scratch, cfg.Skip); err != nil {\n\t\treturn nil, err\n\t}\n\tfor rel, content := range candidate {\n\t\tpath := filepath.Join(scratch, filepath.FromSlash(rel))\n\t\tif err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\tif err := os.WriteFile(path, []byte(content), 0o644); err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t}\n\n\ttimeout := cfg.Timeout\n\tif timeout <= 0 {\n\t\ttimeout = DefaultTimeout\n\t}\n\tvar failures []Failure\n\tfor _, command := range cfg.Commands {\n\t\tfmt.Fprintf(log, \"verify: %s\\n\", command)\n\t\tstatus, text, err := runCommand(ctx, scratch, command, timeout)\n\t\tif err != nil {\n\t\t\treturn nil, err\n\t\t}\n\t\tif status != \"\" {\n\t\t\tfailures = append(failures, Failure{Command: command, Status: status, Output: text})\n\t\t}\n\t}\n\treturn failures, nil\n}","file":"internal/verify/verify.go","name":"Run"}}
// Run applies candidate (paths relative to output, with their new content) over
// a scratch copy of output and runs each command there, printing it to log
// first. It returns the commands that failed, none when all passed. The
// scratch directory is removed on every exit path.
func Run(ctx context.Context, cfg Config, output string, candidate map[string]string, log io.Writer) ([]Failure, error) {
	scratch, err := os.MkdirTemp("", "sedum-verify-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(scratch)

	if err := copyTree(output, scratch, cfg.Skip); err != nil {
		return nil, err
	}
	for rel, content := range candidate {
		path := filepath.Join(scratch, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return nil, err
		}
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	var failures []Failure
	for _, command := range cfg.Commands {
		fmt.Fprintf(log, "verify: %s\n", command)
		status, text, err := runCommand(ctx, scratch, command, timeout)
		if err != nil {
			return nil, err
		}
		if status != "" {
			failures = append(failures, Failure{Command: command, Status: status, Output: text})
		}
	}
	return failures, nil
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// runCommand runs one shell command in dir, in its own process group so a\n// timeout kills everything it started. status is empty when it exited 0.\nfunc runCommand(ctx context.Context, dir, command string, timeout time.Duration) (status, output string, err error) {\n\tcmd := exec.Command(\"sh\", \"-c\", command)\n\tcmd.Dir = dir\n\tcmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}\n\tvar buf bytes.Buffer\n\tcmd.Stdout, cmd.Stderr = &buf, &buf\n\tif err := cmd.Start(); err != nil {\n\t\treturn \"\", \"\", err\n\t}\n\tdone := make(chan error, 1)\n\tgo func() { done <- cmd.Wait() }()\n\n\ttimer := time.NewTimer(timeout)\n\tdefer timer.Stop()\n\tselect {\n\tcase waitErr := <-done:\n\t\tif waitErr == nil {\n\t\t\treturn \"\", \"\", nil\n\t\t}\n\t\tvar exit *exec.ExitError\n\t\tif errors.As(waitErr, &exit) {\n\t\t\treturn exit.String(), tail(buf.String()), nil\n\t\t}\n\t\treturn \"\", \"\", waitErr\n\tcase <-timer.C:\n\t\t_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)\n\t\t<-done\n\t\treturn fmt.Sprintf(\"timed out after %s\", timeout), tail(buf.String()), nil\n\tcase <-ctx.Done():\n\t\t_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)\n\t\t<-done\n\t\treturn \"\", \"\", ctx.Err()\n\t}\n}","file":"internal/verify/verify.go","name":"runCommand"}}
// runCommand runs one shell command in dir, in its own process group so a
// timeout kills everything it started. status is empty when it exited 0.
func runCommand(ctx context.Context, dir, command string, timeout time.Duration) (status, output string, err error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var buf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &buf, &buf
	if err := cmd.Start(); err != nil {
		return "", "", err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case waitErr := <-done:
		if waitErr == nil {
			return "", "", nil
		}
		var exit *exec.ExitError
		if errors.As(waitErr, &exit) {
			return exit.String(), tail(buf.String()), nil
		}
		return "", "", waitErr
	case <-timer.C:
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		return fmt.Sprintf("timed out after %s", timeout), tail(buf.String()), nil
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		<-done
		return "", "", ctx.Err()
	}
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// tail keeps the last MaxOutput bytes of s.\nfunc tail(s string) string {\n\tif len(s) <= MaxOutput {\n\t\treturn s\n\t}\n\treturn s[len(s)-MaxOutput:]\n}","file":"internal/verify/verify.go","name":"tail"}}
// tail keeps the last MaxOutput bytes of s.
func tail(s string) string {
	if len(s) <= MaxOutput {
		return s
	}
	return s[len(s)-MaxOutput:]
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// copyTree copies the regular files and directories under src into dst,\n// skipping hidden directories and the named top-level directories.\nfunc copyTree(src, dst string, skip []string) error {\n\tskipped := map[string]bool{}\n\tfor _, name := range skip {\n\t\tskipped[name] = true\n\t}\n\treturn filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\trel, err := filepath.Rel(src, path)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tif rel == \".\" {\n\t\t\treturn nil\n\t\t}\n\t\tif d.IsDir() {\n\t\t\tif strings.HasPrefix(d.Name(), \".\") || (!strings.Contains(rel, string(filepath.Separator)) && skipped[d.Name()]) {\n\t\t\t\treturn filepath.SkipDir\n\t\t\t}\n\t\t\treturn os.MkdirAll(filepath.Join(dst, rel), 0o755)\n\t\t}\n\t\tif !d.Type().IsRegular() {\n\t\t\treturn nil\n\t\t}\n\t\tdata, err := os.ReadFile(path)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tinfo, err := d.Info()\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\treturn os.WriteFile(filepath.Join(dst, rel), data, info.Mode().Perm()|0o200)\n\t})\n}","file":"internal/verify/verify.go","name":"copyTree"}}
// copyTree copies the regular files and directories under src into dst,
// skipping hidden directories and the named top-level directories.
func copyTree(src, dst string, skip []string) error {
	skipped := map[string]bool{}
	for _, name := range skip {
		skipped[name] = true
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") || (!strings.Contains(rel, string(filepath.Separator)) && skipped[d.Name()]) {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), data, info.Mode().Perm()|0o200)
	})
}

// /sedum:setFunction
// sedum:setFunction {"tier":"owned","record":"prov-2026-aabeac10","gen":"sha256:101d811fc4a470dbb6ac18b477cc46a0afe98d566820267f4bfdd1b8526beb0b","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","kwargs":{"code":"// Message is what the model is told when a selection did not verify.\nfunc Message(failures []Failure) string {\n\tvar b strings.Builder\n\tb.WriteString(\"These checks failed after applying your selection:\\n\")\n\tfor _, f := range failures {\n\t\tfmt.Fprintf(&b, \"\\n$ %s\\n%s\\n%s\", f.Command, f.Status, strings.TrimRight(f.Output, \"\\n\"))\n\t\tb.WriteString(\"\\n\")\n\t}\n\tb.WriteString(\"\\nFix the selection so every one of them passes.\")\n\treturn b.String()\n}","file":"internal/verify/verify.go","name":"Message"}}
// Message is what the model is told when a selection did not verify.
func Message(failures []Failure) string {
	var b strings.Builder
	b.WriteString("These checks failed after applying your selection:\n")
	for _, f := range failures {
		fmt.Fprintf(&b, "\n$ %s\n%s\n%s", f.Command, f.Status, strings.TrimRight(f.Output, "\n"))
		b.WriteString("\n")
	}
	b.WriteString("\nFix the selection so every one of them passes.")
	return b.String()
}

// /sedum:setFunction
