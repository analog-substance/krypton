package run

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	static "github.com/analog-substance/arsenic-static"
	"github.com/analog-substance/krypton/internal/bin"
)

type ScriptExecMode int

const (
	ExecFromDisk ScriptExecMode = iota
	ExecFromMem
)

var nmapPath = ""
var execMode ScriptExecMode = ExecFromMem

func init() {
	var err error
	nmapPath, err = bin.Locate("nmap")
	if err != nil {
		nmapPath = "./nmap"

		fmt.Println("[-] nmap either not found or an error occurred. Falling back writing to disk")

		err = bin.WriteAs(fmt.Sprintf("nmap_%s", runtime.GOARCH), "nmap")
		if err != nil {
			panic(fmt.Errorf("error occurred while writing nmap to disk: %v", err))
		}
	}
}

// SetExecMode sets how the scripts are executed, whether from memory or from disk
func SetExecMode(mode ScriptExecMode) {
	execMode = mode
	if execMode == ExecFromDisk {
		ensureScripts()
	}
}

func ensureScripts() error {
	scripts := []string{
		"bin/as-recon-discover-hosts",
		"bin/as-recon-discover-services",
	}
	for _, script := range scripts {
		err := static.Write(script)
		if err != nil {
			return err
		}
	}
	return nil
}

func diskCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	return cmd
}

func DiscoverHosts(networks string) (string, error) {
	var cmd *exec.Cmd
	var err error

	script := "bin/as-recon-discover-hosts"
	switch execMode {
	default:
		cmd, err = static.Command(script, networks)
	case ExecFromDisk:
		cmd = diskCommand(script, networks)
	}

	if err != nil {
		return "", err
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("NMAP=%s", nmapPath))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error occurred while discovering hosts: %v", err)
	}
	return string(output), nil
}

func discoverServices(hosts string, isUDP bool) error {
	var cmd *exec.Cmd
	var err error

	args := ""
	if isUDP {
		args = "--udp"
	}

	env := []string{
		fmt.Sprintf("NMAP=%s", nmapPath),
	}
	script := "bin/as-recon-discover-services"
	switch execMode {
	default:
		cmd, err = static.Command(script, args)
		env = append(env, fmt.Sprintf("SCRIPT_STDIN=%s", hosts))
	case ExecFromDisk:
		cmd = diskCommand(script, args)
		cmd.Stdin = strings.NewReader(hosts)
	}

	cmd.Env = append(cmd.Env, env...)
	cmd.Stdout = os.Stdout

	if err != nil {
		return err
	}

	return cmd.Run()
}

func DiscoverTCPServices(hosts string) error {
	return discoverServices(hosts, false)
}

func DiscoverUDPServices(hosts string) error {
	return discoverServices(hosts, true)
}
