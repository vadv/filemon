package socket

import (
	"fmt"
	"os"
	"path/filepath"

	"gofilemon/protocol"
)

// return filepath of socket, based on binary name
func FilePath() string {
	fullBinaryName := os.Args[0]
	binaryName := filepath.Base(fullBinaryName)
	return filepath.Join("/tmp", fmt.Sprintf("%s-mon.%s", binaryName, "sock"))
}

func Exists() bool {
	if _, err := os.Stat(FilePath()); err == nil {
		return true
	}
	return false
}

func Alive() bool {

	if !Exists() {
		return false
	}

	client, err := NewClient()
	if err != nil {
		return false
	}
	defer client.Close()

	if err := client.Write([]byte(protocol.PING)); err != nil {
		return false
	}

	if _, err := client.Read(); err != nil {
		return false
	}

	return true
}
