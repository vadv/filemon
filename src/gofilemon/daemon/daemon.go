package daemon

// https://github.com/golang/go/issues/227
// https://habrahabr.ru/post/187668/

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"gofilemon/logger"
	"gofilemon/protocol"
)

const (
	envDaemonName  = "_runned_is_daemon"
	envDaemonValue = "1"
)

func Daemonize(pwd string) error {
	return reborn(0002, pwd)
}

func reborn(umask uint32, workDir string) (err error) {
	if !IsDaemon() {
		var path string
		if path, err = filepath.Abs(os.Args[0]); err != nil {
			return
		}
		cmd := exec.Command(path) // daemonize without parametrs
		envVar := fmt.Sprintf("%s=%s", envDaemonName, envDaemonValue)
		cmd.Env = append(os.Environ(), envVar)
		if fd, err := logger.GetLogFD(); err == nil {
			cmd.Stdout = fd
			cmd.Stderr = fd
		}
		if err = cmd.Start(); err != nil {
			return
		}
		fmt.Printf(protocol.ERROR_ONLY_DAEMON_START)
		os.Exit(protocol.ERROR_COMMON_CODE)
	}
	syscall.Umask(int(umask))
	if len(workDir) != 0 {
		if err = os.Chdir(workDir); err != nil {
			return
		}
	}
	_, err = syscall.Setsid()
	return
}

func IsDaemon() bool {
	return os.Getenv(envDaemonName) == envDaemonValue
}
